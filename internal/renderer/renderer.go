// Package renderer provides output formatters for modern-ls directory listings.
package renderer

import (
	"io"

	"github.com/the-mayankjha/modern-ls/internal/filesystem"
	"github.com/the-mayankjha/modern-ls/internal/themes"
)

// Options controls what a Renderer emits.
type Options struct {
	Theme      *themes.Theme
	TimeFormat string
	Colors     bool
	Icons      bool
	ShowBlocks bool
	ShowGit    bool
	// Long-listing specific
	ShowOwner bool
	ShowGroup bool
	ShowMode  bool
	// Tree specific
	Tree  bool
	Depth int
}

// Renderer is the interface for all output formatters.
type Renderer interface {
	// Render writes entries to w. The termWidth is used for grid layout;
	// it is ignored by long and oneline renderers.
	Render(w io.Writer, entries []*filesystem.Entry, termWidth int, opts Options) error
}

// noColor is the ANSI reset.
const noColor = "\033[0m"

// gitColor returns the ANSI color for a git status string under the given theme.
func gitColor(status string, t *themes.Theme) string {
	if t == nil {
		return ""
	}
	switch status {
	case "U":
		return t.GitUntracked.ANSI()
	case " ", "":
		return ""
	default:
		return t.GitModified.ANSI()
	}
}
