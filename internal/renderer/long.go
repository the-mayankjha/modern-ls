package renderer

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/the-mayankjha/modern-ls/internal/filesystem"
)

// Long implements the long-format listing (-l, -o, -g flags).
type Long struct{}

// Render writes one file per line in the classic long format:
//
//	[blocks] mode owner group size date icon name git
func (l *Long) Render(w io.Writer, entries []*filesystem.Entry, _ int, opts Options) error {
	if len(entries) == 0 {
		return nil
	}

	// First pass: collect column widths.
	maxMode, maxOwner, maxGroup, maxSize, maxDate, maxBlocks := 0, 0, 0, 0, 0, 0
	for _, e := range entries {
		if n := len(e.ModeStr); n > maxMode {
			maxMode = n
		}
		if n := len(e.Owner); n > maxOwner {
			maxOwner = n
		}
		if n := len(e.Group); n > maxGroup {
			maxGroup = n
		}
		sz := formatSizeHuman(e.Size, opts)
		if n := len(sz); n > maxSize {
			maxSize = n
		}
		dt := e.ModTime.Format(timeFormat(opts.TimeFormat))
		if n := len(dt); n > maxDate {
			maxDate = n
		}
		if opts.ShowBlocks {
			bs := formatSize(e.Blocks*512, false)
			if n := len(bs); n > maxBlocks {
				maxBlocks = n
			}
		}
	}

	reset := ""
	if opts.Colors {
		reset = noColor
	}

	for _, e := range entries {
		var sb strings.Builder

		// Block size
		if opts.ShowBlocks {
			bs := formatSize(e.Blocks*512, false)
			fmt.Fprintf(&sb, "%-*s%s", maxBlocks, bs, brailSep)
		}

		// Mode
		if opts.ShowMode && e.ModeStr != "" {
			fmt.Fprintf(&sb, "%-*s%s", maxMode, e.ModeStr, brailSep)
		}

		// Owner
		if opts.ShowOwner && e.Owner != "" {
			fmt.Fprintf(&sb, "%-*s%s", maxOwner, e.Owner, brailSep)
		}

		// Group
		if opts.ShowGroup && e.Group != "" {
			fmt.Fprintf(&sb, "%-*s%s", maxGroup, e.Group, brailSep)
		}

		// Size
		sz := formatSizeHuman(e.Size, opts)
		fmt.Fprintf(&sb, "%*s%s", maxSize, sz, brailSep)

		// Date
		dt := e.ModTime.Format(timeFormat(opts.TimeFormat))
		fmt.Fprintf(&sb, "%-*s%s", maxDate, dt, brailSep)

		// Icon
		if opts.Icons {
			fmt.Fprintf(&sb, "%s%s%s%s", e.IconColor, e.Icon, reset, brailSep)
		}

		// Name
		name := e.FullName + string(e.Indicator)
		gc := ""
		gcReset := ""
		if opts.Colors && e.GitStatus != "" {
			gc = gitColorStr(e.GitStatus, opts)
			gcReset = reset
		}
		fmt.Fprintf(&sb, "%s%s%s", gc, name, gcReset)

		// Git status
		if opts.ShowGit && e.GitStatus != "" {
			fmt.Fprintf(&sb, " %s%s%s", gc, e.GitStatus, gcReset)
		}

		fmt.Fprintln(w, sb.String())
	}
	return nil
}

func formatSizeHuman(b int64, opts Options) string {
	// Check if human-readable mode is set via a convention: we use the
	// ShowBlocks field as a proxy; a real implementation would pass HumanReadable
	// in Options. For now, always use bytes in long mode.
	return fmt.Sprintf("%d", b)
}

func timeFormat(tf string) string {
	switch tf {
	case "StampMilli":
		return time.StampMilli
	case "Kitchen":
		return time.Kitchen
	case "ANSIC":
		return time.ANSIC
	case "UnixDate":
		return time.UnixDate
	case "RubyDate":
		return time.RubyDate
	case "RFC1123":
		return time.RFC1123
	case "RFC1123Z":
		return time.RFC1123Z
	case "RFC3339":
		return time.RFC3339
	case "RFC822":
		return time.RFC822
	case "RFC822Z":
		return time.RFC822Z
	case "RFC850":
		return time.RFC850
	default:
		return time.Stamp
	}
}
