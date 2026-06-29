// Package parser provides a Lua table parser for nvim-web-devicons icon files.
// It uses regex-based extraction without requiring a Lua runtime.
package parser

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// Icon represents a single icon entry parsed from Lua source.
type Icon struct {
	// Name is the key used to look up this icon (filename or extension).
	Name string `json:"name"`
	// Glyph is the Nerd Font Unicode character for this icon.
	Glyph string `json:"glyph"`
	// Color is the hex colour string, e.g. "#F05340".
	Color string `json:"color"`
}

// entryRe matches a single nvim-web-devicons table entry on one line, e.g.:
//
//	["go"] = { icon = "", color = "#00ADD8", cterm_color = "38", name = "Go" },
var entryRe = regexp.MustCompile(`\["([^"]+)"\]\s*=\s*\{\s*icon\s*=\s*"([^"]+)"[^}]*color\s*=\s*"(#[0-9A-Fa-f]{6})"`)

// ParseLua reads a nvim-web-devicons Lua table file and returns a map of
// key → Icon. The key is either a filename (e.g. "go.mod") or a file
// extension (e.g. "go").
//
// Lines that do not match the expected pattern (comments, blank lines, the
// opening/closing braces) are silently skipped, making the parser robust
// against minor formatting variations.
func ParseLua(r io.Reader) (map[string]Icon, error) {
	result := make(map[string]Icon)

	scanner := bufio.NewScanner(r)
	// Allocate a large buffer so that wide Unicode lines are not split.
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()

		// Fast-path: skip blank lines and comment/structural lines.
		trimmed := strings.TrimSpace(line)
		if trimmed == "" ||
			strings.HasPrefix(trimmed, "--") ||
			trimmed == "return {" ||
			trimmed == "}" {
			continue
		}

		matches := entryRe.FindStringSubmatch(line)
		if len(matches) < 4 {
			continue
		}

		key := matches[1]
		glyph := matches[2]
		color := matches[3]

		result[key] = Icon{
			Name:  key,
			Glyph: glyph,
			Color: color,
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("parser: scan error: %w", err)
	}

	return result, nil
}
