package renderer

import (
	"fmt"
	"io"

	"github.com/the-mayankjha/modern-ls/internal/filesystem"
)

// OneLine implements one-file-per-line listing (-1 flag).
type OneLine struct{}

// Render writes one entry per line: [size] icon name [git].
func (o *OneLine) Render(w io.Writer, entries []*filesystem.Entry, _ int, opts Options) error {
	reset := ""
	if opts.Colors {
		reset = noColor
	}

	for _, e := range entries {
		// Optional block size
		if opts.ShowBlocks {
			bs := formatSize(e.Blocks*512, false)
			fmt.Fprintf(w, "%s%s", bs, brailSep)
		}

		// Icon
		if opts.Icons && e.Icon != "" {
			fmt.Fprintf(w, "%s%s%s%s", e.IconColor, e.Icon, reset, brailSep)
		}

		// Name
		name := e.FullName + string(e.Indicator)
		gc := gitColorStr(e.GitStatus, opts)
		gcReset := ""
		if gc != "" {
			gcReset = reset
		}
		fmt.Fprintf(w, "%s%s%s", gc, name, gcReset)

		// Git
		if opts.ShowGit && e.GitStatus != "" {
			fmt.Fprintf(w, " %s%s%s", gc, e.GitStatus, gcReset)
		}

		fmt.Fprintln(w)
	}
	return nil
}
