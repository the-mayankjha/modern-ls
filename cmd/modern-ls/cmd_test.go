package main

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestModernLs_E2E(t *testing.T) {
	// Build the binary
	bin := filepath.Join(t.TempDir(), "modern-ls")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to build modern-ls: %v", err)
	}

	// Create test directory structure
	testDir := t.TempDir()
	
	// Create some dummy files
	files := []string{
		"Makefile",
		"main.go",
		"README.md",
		".gitignore",
		"script.sh",
	}
	
	for _, f := range files {
		cmd = exec.Command("touch", filepath.Join(testDir, f))
		if err := cmd.Run(); err != nil {
			t.Fatalf("failed to touch %s: %v", f, err)
		}
	}
	
	// Make script.sh executable
	cmd = exec.Command("chmod", "+x", filepath.Join(testDir, "script.sh"))
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to chmod script.sh: %v", err)
	}

	// Create some dummy directories
	dirs := []string{
		"admin",
		"src",
		".git",
		"regular_folder",
	}
	
	for _, d := range dirs {
		cmd = exec.Command("mkdir", filepath.Join(testDir, d))
		if err := cmd.Run(); err != nil {
			t.Fatalf("failed to mkdir %s: %v", d, err)
		}
	}

	// Run modern-ls -1 in the test directory
	// Note: We use -1 for one-per-line to make it easy to parse
	// We disable color so we can just check the raw text/icons
	cmd = exec.Command(bin, "-1", "-c", "-a")
	cmd.Dir = testDir
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to run modern-ls: %v", err)
	}

	output := out.String()
	
	expected := []string{
		"⠀.git", // Custom folder
		"⠀.gitignore", // ByExtension / ByFilename
		"⠀admin", // Custom semantic folder
		"󰂺⠀README.md", // ByExtension / ByFilename
		"⠀regular_folder", // Default folder fallback
		"⠀main.go", // ByExtension
		"⠀Makefile", // ByFilename
		"⠀script.sh", // Executable (.sh)
		"⠀src", // Custom folder
	}
	
	for _, exp := range expected {
		if !strings.Contains(output, exp) {
			t.Errorf("expected output to contain %q, got:\n%s", exp, output)
		}
	}
}
