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
			fmt.Fprintf(w, "%s%s%s%s", iconColorStr(e, opts), e.Icon, reset, brailSep)
		}

		// Name
		name := e.FullName + string(e.Indicator)
		nc := entryColorStr(e, opts)
		ncReset := ""
		if nc != "" {
			ncReset = reset
		}
		fmt.Fprintf(w, "%s%s%s", nc, name, ncReset)

		// Git
		if opts.ShowGit && e.GitStatus != "" {
			gc := ""
			gcReset := ""
			if opts.Colors {
				gc = gitColorStr(e.GitStatus, opts)
				if gc != "" {
					gcReset = reset
				}
			}
			fmt.Fprintf(w, " %s%s%s", gc, e.GitStatus, gcReset)
		}

		fmt.Fprintln(w)
	}
	return nil
}
