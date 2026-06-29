// Package config loads and provides modern-ls configuration.
//
// Precedence order (highest to lowest):
//  1. CLI flags (applied by the caller, not this package)
//  2. MODERN_LS_CONFIG environment variable (explicit path)
//  3. $XDG_CONFIG_HOME/modern-ls/config.toml  (typically ~/.config/…)
//  4. Built-in defaults (see Defaults)
package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/adrg/xdg"
	"github.com/pelletier/go-toml/v2"
)

// Display controls what information is rendered.
type Display struct {
	Icons  bool   `toml:"icons"`
	Colors bool   `toml:"colors"`
	Git    bool   `toml:"git"`
	Theme  string `toml:"theme"`
	Human  bool   `toml:"human_readable"`
	Date   string `toml:"date_format"`
}

// Sort controls entry ordering.
type Sort struct {
	By      string `toml:"by"`      // alpha | size | time | ext | version | none
	Reverse bool   `toml:"reverse"` // reverse the chosen sort order
}

// Output controls the output layout.
type Output struct {
	Layout string `toml:"layout"` // grid | long | oneline
}

// Git controls git repository integration.
type Git struct {
	Enabled bool `toml:"enabled"`
}

// Config is the complete modern-ls configuration structure.
// It is intentionally flat to map 1-to-1 with a TOML config file.
type Config struct {
	Display Display `toml:"display"`
	Sort    Sort    `toml:"sort"`
	Output  Output  `toml:"output"`
	Git     Git     `toml:"git"`
}

// Defaults returns a Config populated with sensible built-in values.
// All callers should start from Defaults() and layer overrides on top.
func Defaults() Config {
	return Config{
		Display: Display{
			Icons:  true,
			Colors: true,
			Git:    false,
			Theme:  "default",
			Human:  false,
			Date:   "Stamp",
		},
		Sort: Sort{
			By:      "alpha",
			Reverse: false,
		},
		Output: Output{
			Layout: "grid",
		},
		Git: Git{
			Enabled: false,
		},
	}
}

// Load reads the config file resolved via the standard path hierarchy,
// merging its values over Defaults(). A missing config file is not an
// error — Load returns the defaults silently.
func Load() (Config, error) {
	cfg := Defaults()

	path, err := resolvePath()
	if err != nil || path == "" {
		// Path resolution failure (e.g. bad XDG state) is non-fatal.
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("config: read %q: %w", path, err)
	}
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("config: parse %q: %w", path, err)
	}
	return cfg, nil
}

// LoadFrom reads a config file from an explicit path, merging over Defaults().
// Unlike Load, an absent file is returned as an error because the caller has
// explicitly requested a specific location.
func LoadFrom(path string) (Config, error) {
	cfg := Defaults()
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("config: read %q: %w", path, err)
	}
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("config: parse %q: %w", path, err)
	}
	return cfg, nil
}

// resolvePath returns the config file path to load.
// Priority: MODERN_LS_CONFIG env var > XDG config directory.
// Returns ("", nil) when no path can be determined without error.
func resolvePath() (string, error) {
	if p := os.Getenv("MODERN_LS_CONFIG"); p != "" {
		return p, nil
	}
	// xdg.ConfigFile does NOT create the file; it merely resolves the path.
	p, err := xdg.ConfigFile("modern-ls/config.toml")
	if err != nil {
		return "", err
	}
	return p, nil
}
