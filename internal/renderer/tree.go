package renderer

import (
	"fmt"
	"io"
	"strings"

	"github.com/the-mayankjha/modern-ls/internal/filesystem"
)

// Tree renders entries in a tree-like format, prefixing each entry with branch characters.
type Tree struct {
	Prefix     string
	IsLast     bool
	Depth      int // not really needed if prefix is managed by caller
	TimeFormat string
}

// Render here is not used the same way as Grid/Long, but we can implement it for single items.
func (t *Tree) Render(w io.Writer, entries []*filesystem.Entry, termWidth int, opts Options) error {
	return nil
}

// RenderEntry formats a single entry in the tree.
func (t *Tree) RenderEntry(w io.Writer, e *filesystem.Entry, opts Options) {
	branch := "├── "
	if t.IsLast {
		branch = "└── "
	}

	// Create a long format string if ShowMode is true
	var details string
	// We'll skip the complex Long format in tree view for now to keep it clean.

	name := e.Name
	if e.IsHidden {
		name = "." + name
	}
	name += e.Ext

	icon := ""
	if opts.Icons {
		c := ""
		if opts.Colors {
			c = iconColorStr(e, opts)
		}
		icon = c + e.Icon + noColor + " "
	}

	nc := entryColorStr(e, opts)
	ncReset := ""
	if nc != "" {
		ncReset = noColor
	}

	gitStr := ""
	if opts.ShowGit && e.GitStatus != "" && e.GitStatus != " " {
		gc := ""
		gcReset := ""
		if opts.Colors {
			gc = gitColor(strings.TrimSpace(e.GitStatus), opts.Theme)
			gcReset = noColor
		}
		gitStr = " [" + gc + strings.TrimSpace(e.GitStatus) + gcReset + "]"
	}

	ind := ""
	if e.Indicator != filesystem.IndicatorNone {
		ind = string(e.Indicator)
	}

	fmt.Fprintf(w, "%s%s%s%s%s%s%s%s\n", t.Prefix, branch, details, icon, nc, name, ncReset, ind+gitStr)
}
