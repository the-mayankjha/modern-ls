// Package icons provides icon lookup for filesystem entries.
// Icon data is sourced from nvim-web-devicons and generated at build time
// via go generate. Do not modify icons_generated.go manually.
package icons

import (
	"fmt"
	"strings"

	"github.com/the-mayankjha/modern-ls/internal/filesystem"
)


// hexToANSI converts a #RRGGBB hex string to an ANSI 24-bit foreground code.
func hexToANSI(hex string) string {
	if len(hex) != 7 || hex[0] != '#' {
		return ""
	}
	var r, g, b uint8
	if _, err := fmt.Sscanf(hex[1:], "%02x%02x%02x", &r, &g, &b); err != nil {
		return ""
	}
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
}

// executableColor is the ANSI color for executable files (green, matching logo-ls).
const executableColor = "\033[38;2;76;175;080m"

// Lookup returns the icon glyph and ANSI color string for a given filesystem entry.
// It follows this priority chain:
//  1. Exact filename match (ByFilename)
//  2. Special: .go files ending in _test
//  3. Sub-extension match (compound extension, e.g. "spec.ts")
//  4. Extension match (ByExtension)
//  5. Hidden file / hidden directory fallback
//  6. Directory / file default
//
// The returned color is already an ANSI escape string.
// If the entry is executable, the color is overridden with green.
func Lookup(name, ext string, ind filesystem.Indicator) (glyph, color string) {
	fullName := name + ext

	switch ind {
	case filesystem.IndicatorDir:
		return lookupDir(fullName)
	default:
		return lookupFile(name, ext, fullName, ind)
	}
}

func lookupDir(fullName string) (glyph, color string) {
	lower := strings.ToLower(fullName)
	
	// 1. Specific folder match
	if entry, ok := ByFolder[lower]; ok {
		return entry.Glyph, hexToANSI(entry.Color)
	}
	
	// 2. Exact filename match (rare for folders, but occasionally exists in icon maps)
	if entry, ok := ByFilename[lower]; ok {
		return entry.Glyph, hexToANSI(entry.Color)
	}
	
	// 3. Hidden directory (starts with dot)
	if strings.HasPrefix(fullName, ".") {
		d := Defaults["hiddendir"]
		return d.Glyph, hexToANSI(d.Color)
	}
	
	// 4. Generic directory fallback
	d := Defaults["dir"]
	return d.Glyph, hexToANSI(d.Color)
}

func lookupFile(name, ext, fullName string, ind filesystem.Indicator) (glyph, color string) {
	lower := strings.ToLower(fullName)

	// 1. Exact filename match
	if entry, ok := ByFilename[lower]; ok {
		glyph, color = entry.Glyph, hexToANSI(entry.Color)
		return applyExec(glyph, color, ind)
	}

	// 2. Special Go test file
	if ext == ".go" && strings.HasSuffix(name, "_test") {
		if entry, ok := ByFilename["_test.go"]; ok {
			glyph, color = entry.Glyph, hexToANSI(entry.Color)
			return applyExec(glyph, color, ind)
		}
	}

	// 3. Sub-extension: e.g. for "foo.spec.ts", try "spec.ts"
	parts := strings.SplitN(name, ".", 2)
	if len(parts) == 2 && parts[0] != "" {
		subExt := strings.ToLower(parts[1] + ext)
		if entry, ok := ByFilename[subExt]; ok {
			glyph, color = entry.Glyph, hexToANSI(entry.Color)
			return applyExec(glyph, color, ind)
		}
	}

	// 4. Extension match
	if ext != "" {
		extLower := strings.ToLower(strings.TrimPrefix(ext, "."))
		if entry, ok := ByExtension[extLower]; ok {
			glyph, color = entry.Glyph, hexToANSI(entry.Color)
			return applyExec(glyph, color, ind)
		}
	}

	// 5. Hidden file
	if strings.HasPrefix(fullName, ".") {
		if ind == filesystem.IndicatorExec {
			return Defaults["exe"].Glyph, executableColor
		}
		d := Defaults["hiddenfile"]
		glyph, color = d.Glyph, hexToANSI(d.Color)
		return applyExec(glyph, color, ind)
	}

	// 6. Default file
	if ind == filesystem.IndicatorExec {
		return Defaults["exe"].Glyph, executableColor
	}
	d := Defaults["file"]
	glyph, color = d.Glyph, hexToANSI(d.Color)
	return applyExec(glyph, color, ind)
}

// applyExec overrides the color with green if the entry is executable.
func applyExec(glyph, color string, ind filesystem.Indicator) (string, string) {
	if ind == filesystem.IndicatorExec {
		return glyph, executableColor
	}
	return glyph, color
}
