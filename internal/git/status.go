// Package git provides Git repository status information for modern-ls.
// It uses go-git for repository detection and status queries, with an
// optional fast-path through the git binary (git status --porcelain -z).
package git

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	gogit "github.com/go-git/go-git/v5"
)

// Status represents the git status of a single file in the working tree.
type Status rune

const (
	// Untracked indicates a new file that is not tracked by the index.
	Untracked Status = 'U'
	// Modified indicates a tracked file that has been changed in the worktree.
	Modified Status = 'M'
	// DirMod indicates a directory that contains at least one modified or
	// untracked descendant.
	DirMod Status = '●'
	// Clean indicates no change, or that the path is not inside a git repo.
	Clean Status = ' '
)

// String returns the single-character display string for a status.
// Clean returns an empty string so callers can test for non-empty output.
func (s Status) String() string {
	switch s {
	case Untracked:
		return "U"
	case Modified:
		return "M"
	case DirMod:
		return "●"
	default:
		return ""
	}
}

// RepoStatus is a mapping from relative path to its Status for all changed
// files within a repository. Paths use forward slashes and are relative to
// the queried directory.
type RepoStatus map[string]Status

// Get returns the Status for the given relative path. Returns Clean when the
// path has no status entry or the receiver is nil.
func (rs RepoStatus) Get(relPath string) Status {
	if rs == nil {
		return Clean
	}
	if s, ok := rs[relPath]; ok {
		return s
	}
	return Clean
}

// ForDir computes the git status for all entries in dirPath.
// It walks up the directory tree looking for a .git directory (via
// go-git's DetectDotGit option) so it works from any sub-directory.
//
// Returns (nil, nil) when dirPath is not inside any git repository —
// this is not considered an error condition for modern-ls.
func ForDir(dirPath string) (RepoStatus, error) {
	opts := &gogit.PlainOpenOptions{DetectDotGit: true}
	repo, err := gogit.PlainOpenWithOptions(dirPath, opts)
	if err != nil {
		// Not a git repo — normal for non-versioned directories.
		return nil, nil //nolint:nilerr
	}

	wt, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("git: worktree: %w", err)
	}

	repoRoot := wt.Filesystem.Root()

	raw, err := rawStatus(wt)
	if err != nil {
		return nil, err
	}
	if len(raw) == 0 {
		return nil, nil
	}

	absDir, err := filepath.Abs(dirPath)
	if err != nil {
		return nil, fmt.Errorf("git: abs path: %w", err)
	}
	// Ensure trailing separator so HasPrefix matches correctly.
	absDir = filepath.ToSlash(absDir) + "/"
	root := filepath.ToSlash(repoRoot) + "/"

	result := make(RepoStatus)
	for filePath, fs := range raw {
		// filePath from go-git is relative to the repo root; make it absolute.
		absFile := root + filePath
		// Strip the queried directory prefix to get a path relative to dirPath.
		relToDir := relPath(absFile, absDir)
		if relToDir == "" {
			continue
		}

		segments := strings.SplitAfter(relToDir, "/")
		dirKey := ""
		for i, seg := range segments {
			if i == len(segments)-1 {
				// Leaf file — record its specific status.
				st := Modified
				if fs.Worktree == '?' {
					st = Untracked
				}
				result[relToDir] = st
			} else {
				// Ancestor directory — mark as containing modifications.
				dirKey += seg
				result[dirKey] = DirMod
			}
		}
	}
	return result, nil
}

// rawStatus retrieves git status preferring the git binary (faster for large
// repos) and falling back to go-git's pure-Go implementation.
func rawStatus(wt *gogit.Worktree) (gogit.Status, error) {
	cmd := exec.Command("git", "status", "--porcelain", "-z")
	cmd.Dir = wt.Filesystem.Root()
	out, err := cmd.Output()
	if err == nil {
		return parsePorcelain(string(out)), nil
	}
	// git binary unavailable or returned non-zero — fall back to go-git.
	return wt.Status()
}

// parsePorcelain parses the NUL-delimited output of `git status --porcelain -z`
// into a gogit.Status map. Each record is "<XY> <path>" where XY is a two-letter
// status code and fields are separated by a space.
func parsePorcelain(output string) gogit.Status {
	parts := strings.Split(output, "\x00")
	result := make(gogit.Status, len(parts))
	for _, line := range parts {
		if len(line) < 3 {
			continue
		}
		// First two bytes are the XY status code; byte 3 is a space.
		xy := line[:2]
		path := strings.TrimSpace(line[3:])
		if path == "" {
			continue
		}
		// Use the worktree status code (second character of XY).
		wtCode := gogit.StatusCode(xy[1])
		result[path] = &gogit.FileStatus{
			Worktree: wtCode,
		}
	}
	return result
}

// relPath returns the portion of gitPath after the dirPath prefix.
// Returns "" if gitPath does not start with dirPath.
func relPath(gitPath, dirPath string) string {
	if strings.HasPrefix(gitPath, dirPath) {
		return strings.TrimPrefix(gitPath, dirPath)
	}
	return ""
}
