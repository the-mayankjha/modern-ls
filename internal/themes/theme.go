// Package themes defines the theme system for modern-ls.
// Each theme is a named palette of 24-bit RGB colors that is registered via
// an init() function. Callers obtain a theme with Get or Default.
package themes

import (
	"fmt"
	"sort"
)

// Color is a 24-bit terminal color stored as its R, G, B components.
type Color struct {
	R, G, B uint8
}

// ANSI returns the ANSI SGR escape sequence for this color as a foreground.
// Output: ESC[38;2;<R>;<G>;<B>m
func (c Color) ANSI() string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", c.R, c.G, c.B)
}

// Reset is the ANSI SGR escape sequence that clears all text attributes.
const Reset = "\033[0m"

// Theme defines the complete color palette used across a modern-ls output.
// Fields are grouped by their purpose: filesystem entry types, Git status
// indicators, and long-listing column labels.
type Theme struct {
	// Name is the unique identifier used for CLI / config lookup.
	Name string

	// Filesystem entry types.
	Dir        Color // regular directory
	DirOpen    Color // open/expanded directory (tree mode)
	HiddenDir  Color // dot-directory (e.g. .git)
	File       Color // regular file
	HiddenFile Color // dot-file (e.g. .bashrc)
	Executable Color // file with execute bit set
	Symlink    Color // symbolic link
	Pipe       Color // named pipe (FIFO)
	Socket     Color // Unix domain socket
	Special    Color // block/character device

	// Git working-tree status indicators.
	GitUntracked Color
	GitModified  Color
	GitAdded     Color
	GitDeleted   Color
	GitRenamed   Color
	GitConflict  Color

	// Long-listing column decorations.
	Permissions Color
	Owner       Color
	Group       Color
	SizeUnit    Color // unit suffix (K, M, G…)
	SizeNum     Color // numeric part of size
	Date        Color

	// Block-allocation size column.
	BlockSize Color
}

// registry holds all registered themes, keyed by name.
// It is populated by each theme file's init() function.
var registry = map[string]*Theme{}

// Register adds t to the global theme registry.
// Duplicate names silently overwrite the previous entry; register each theme
// exactly once inside an init() function to avoid races.
func Register(t *Theme) {
	registry[t.Name] = t
}

// Get returns the named theme. If name is not found in the registry,
// Get returns Default() so callers always receive a valid palette.
func Get(name string) *Theme {
	if t, ok := registry[name]; ok {
		return t
	}
	return Default()
}

// Names returns the names of all registered themes in sorted order.
func Names() []string {
	names := make([]string, 0, len(registry))
	for k := range registry {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}
