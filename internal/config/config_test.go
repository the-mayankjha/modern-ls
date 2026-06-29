package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaults(t *testing.T) {
	cfg := Defaults()

	if !cfg.Display.Icons {
		t.Error("Defaults().Display.Icons should be true")
	}
	if !cfg.Display.Colors {
		t.Error("Defaults().Display.Colors should be true")
	}
	if cfg.Display.Theme != "default" {
		t.Errorf("Defaults().Display.Theme = %q, want %q", cfg.Display.Theme, "default")
	}
	if cfg.Sort.By != "alpha" {
		t.Errorf("Defaults().Sort.By = %q, want %q", cfg.Sort.By, "alpha")
	}
	if cfg.Sort.Reverse {
		t.Error("Defaults().Sort.Reverse should be false")
	}
	if cfg.Output.Layout != "grid" {
		t.Errorf("Defaults().Output.Layout = %q, want %q", cfg.Output.Layout, "grid")
	}
	if cfg.Git.Enabled {
		t.Error("Defaults().Git.Enabled should be false")
	}
}

func TestLoad_NoFile_ReturnsDefaults(t *testing.T) {
	// Point MODERN_LS_CONFIG to a path that does not exist.
	tmp := filepath.Join(t.TempDir(), "nonexistent.toml")
	t.Setenv("MODERN_LS_CONFIG", tmp)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}
	want := Defaults()
	if cfg != want {
		t.Errorf("Load() = %+v, want defaults %+v", cfg, want)
	}
}

func TestLoad_ValidFile(t *testing.T) {
	const tomlContent = `
[display]
icons  = false
colors = true
theme  = "dracula"

[sort]
by      = "size"
reverse = true

[output]
layout = "long"

[git]
enabled = true
`
	tmp := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(tmp, []byte(tomlContent), 0o644); err != nil {
		t.Fatalf("setup: write temp config: %v", err)
	}
	t.Setenv("MODERN_LS_CONFIG", tmp)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}

	if cfg.Display.Icons {
		t.Error("Display.Icons: got true, want false")
	}
	if cfg.Display.Theme != "dracula" {
		t.Errorf("Display.Theme = %q, want %q", cfg.Display.Theme, "dracula")
	}
	if cfg.Sort.By != "size" {
		t.Errorf("Sort.By = %q, want %q", cfg.Sort.By, "size")
	}
	if !cfg.Sort.Reverse {
		t.Error("Sort.Reverse: got false, want true")
	}
	if cfg.Output.Layout != "long" {
		t.Errorf("Output.Layout = %q, want %q", cfg.Output.Layout, "long")
	}
	if !cfg.Git.Enabled {
		t.Error("Git.Enabled: got false, want true")
	}
}

func TestLoadFrom_InvalidTOML(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "bad.toml")
	if err := os.WriteFile(tmp, []byte("[[[ not valid toml"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	_, err := LoadFrom(tmp)
	if err == nil {
		t.Fatal("LoadFrom() with invalid TOML: want error, got nil")
	}
}

func TestLoadFrom_MissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.toml")
	_, err := LoadFrom(path)
	if err == nil {
		t.Fatal("LoadFrom() with missing file: want error, got nil")
	}
}

func TestLoadFrom_PartialOverride(t *testing.T) {
	// Only override one field; remaining fields should retain defaults.
	const partial = `
[display]
theme = "nord"
`
	tmp := filepath.Join(t.TempDir(), "partial.toml")
	if err := os.WriteFile(tmp, []byte(partial), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	cfg, err := LoadFrom(tmp)
	if err != nil {
		t.Fatalf("LoadFrom() error = %v", err)
	}
	if cfg.Display.Theme != "nord" {
		t.Errorf("Display.Theme = %q, want %q", cfg.Display.Theme, "nord")
	}
	// Default for Icons should still hold.
	if !cfg.Display.Icons {
		t.Error("Display.Icons: got false after partial override, want true (default)")
	}
}
