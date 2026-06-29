package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// gitRun runs a git sub-command in dir. Returns an error if the git binary is
// not on PATH, allowing callers to skip git-dependent tests gracefully.
func gitRun(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

// hasGit reports whether the git binary is available on PATH.
func hasGit() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

// ---------------------------------------------------------------------------
// Status.String
// ---------------------------------------------------------------------------

func TestStatus_String(t *testing.T) {
	tests := []struct {
		status Status
		want   string
	}{
		{Untracked, "U"},
		{Modified, "M"},
		{DirMod, "●"},
		{Clean, ""},
		{Status(' '), ""},
	}
	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			if got := tc.status.String(); got != tc.want {
				t.Errorf("Status(%q).String() = %q, want %q", rune(tc.status), got, tc.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// RepoStatus.Get
// ---------------------------------------------------------------------------

func TestRepoStatus_Get(t *testing.T) {
	rs := RepoStatus{
		"main.go":      Modified,
		"internal/foo": Untracked,
	}

	tests := []struct {
		path string
		want Status
	}{
		{"main.go", Modified},
		{"internal/foo", Untracked},
		{"notexist.go", Clean},
	}
	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			got := rs.Get(tc.path)
			if got != tc.want {
				t.Errorf("RepoStatus.Get(%q) = %v, want %v", tc.path, got, tc.want)
			}
		})
	}
}

func TestRepoStatus_Get_Nil(t *testing.T) {
	var rs RepoStatus
	if got := rs.Get("any/path"); got != Clean {
		t.Errorf("nil RepoStatus.Get() = %v, want Clean", got)
	}
}

// ---------------------------------------------------------------------------
// relPath (unexported helper)
// ---------------------------------------------------------------------------

func TestRelPath(t *testing.T) {
	tests := []struct {
		gitPath string
		dirPath string
		want    string
	}{
		{"/repo/src/main.go", "/repo/src/", "main.go"},
		{"/repo/src/main.go", "/repo/", "src/main.go"},
		{"/other/path", "/repo/", ""},
		{"/repo/", "/repo/", ""},
	}
	for _, tc := range tests {
		t.Run(tc.gitPath, func(t *testing.T) {
			got := relPath(tc.gitPath, tc.dirPath)
			if got != tc.want {
				t.Errorf("relPath(%q, %q) = %q, want %q",
					tc.gitPath, tc.dirPath, got, tc.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// parsePorcelain (unexported helper)
// ---------------------------------------------------------------------------

func TestParsePorcelain(t *testing.T) {
	// Simulate: "M  main.go NUL ?? untracked.go NUL"
	// XY codes: "M " = modified in index, " " in worktree → worktree = ' '
	//           "??" = untracked → worktree = '?'
	input := "M  main.go\x00?? untracked.go\x00"
	result := parsePorcelain(input)

	if s, ok := result["main.go"]; !ok {
		t.Error("parsePorcelain: expected main.go in result")
	} else if s.Worktree != ' ' {
		t.Errorf("main.go Worktree = %q, want ' '", s.Worktree)
	}

	if s, ok := result["untracked.go"]; !ok {
		t.Error("parsePorcelain: expected untracked.go in result")
	} else if s.Worktree != '?' {
		t.Errorf("untracked.go Worktree = %q, want '?'", s.Worktree)
	}
}

// ---------------------------------------------------------------------------
// ForDir — integration-style tests
// ---------------------------------------------------------------------------

func TestForDir_NotARepo(t *testing.T) {
	dir := t.TempDir()
	rs, err := ForDir(dir)
	if err != nil {
		t.Fatalf("ForDir(non-repo) error = %v, want nil", err)
	}
	if rs != nil {
		t.Errorf("ForDir(non-repo) = %v, want nil", rs)
	}
}

func TestForDir_WithUntrackedFile(t *testing.T) {
	if !hasGit() {
		t.Skip("git binary not found on PATH")
	}

	dir := t.TempDir()
	gitRun(t, dir, "init")
	gitRun(t, dir, "config", "user.email", "ci@example.com")
	gitRun(t, dir, "config", "user.name", "CI")

	// Create an untracked file.
	newFile := filepath.Join(dir, "hello.go")
	if err := os.WriteFile(newFile, []byte("package main\n"), 0o644); err != nil {
		t.Fatalf("create test file: %v", err)
	}

	rs, err := ForDir(dir)
	if err != nil {
		t.Fatalf("ForDir error = %v", err)
	}
	if rs == nil {
		t.Fatal("ForDir returned nil, want non-nil RepoStatus")
	}
	got := rs.Get("hello.go")
	if got != Untracked {
		t.Errorf("hello.go status = %v, want Untracked", got)
	}
}

func TestForDir_CleanRepo(t *testing.T) {
	if !hasGit() {
		t.Skip("git binary not found on PATH")
	}

	dir := t.TempDir()
	gitRun(t, dir, "init")
	gitRun(t, dir, "config", "user.email", "ci@example.com")
	gitRun(t, dir, "config", "user.name", "CI")

	// Commit a file so the repo is clean.
	f := filepath.Join(dir, "README.md")
	if err := os.WriteFile(f, []byte("# test\n"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	gitRun(t, dir, "add", ".")
	gitRun(t, dir, "commit", "-m", "init")

	rs, err := ForDir(dir)
	if err != nil {
		t.Fatalf("ForDir error = %v", err)
	}
	// A clean repo should return nil (no changes).
	if rs != nil && len(rs) > 0 {
		t.Errorf("clean repo ForDir() = %v, want nil or empty", rs)
	}
}
