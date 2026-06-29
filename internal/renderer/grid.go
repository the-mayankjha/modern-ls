package renderer

import (
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/the-mayankjha/modern-ls/internal/filesystem"
)

// brailSep is the invisible Braille Empty character used as a zero-width
// visual separator between icon and filename columns.
// This preserves the original logo-ls column alignment trick.
const brailSep = "\u2800"

// cell holds pre-computed display data for one entry.
type cell struct {
	size      string
	icon      string
	iconColor string
	name      string // name + ext + indicator
	git       string
	// visual widths
	sizeW int
	nameW int
	gitW  int
}

// cellWidth returns the total visual column width of a cell + padding.
func (c *cell) colWidth(showIcon bool) int {
	w := 0
	if c.sizeW > 0 {
		w += c.sizeW + 1 // size + space
	}
	if showIcon {
		w += 2 // icon glyph + brail sep = 2 columns
	}
	w += c.nameW
	if c.gitW > 0 {
		w += 2 // brail sep + git char
	}
	return w
}

// Grid implements the multi-column grid layout.
type Grid struct{}

// Render writes a multi-column grid layout to w, fitting within termWidth columns.
func (g *Grid) Render(w io.Writer, entries []*filesystem.Entry, termWidth int, opts Options) error {
	if len(entries) == 0 {
		return nil
	}

	cells := buildCells(entries, opts)
	showIcon := opts.Icons && anyHasIcon(cells)
	const pad = 2

	// Find optimal column count using a greedy doubling approach.
	cols, widths := computeColumns(cells, showIcon, termWidth, pad)
	_ = cols

	rows := int(math.Ceil(float64(len(cells)) / float64(len(widths))))

	for row := 0; row < rows; row++ {
		for col := 0; col < len(widths); col++ {
			idx := row + col*rows
			if idx >= len(cells) {
				continue
			}
			isLast := col == len(widths)-1
			p := pad
			if isLast {
				p = 0
			}
			writeCell(w, cells[idx], widths[col], showIcon, opts)
			fmt.Fprintf(w, "%*s", p, "")
		}
		fmt.Fprintln(w)
	}
	return nil
}

// colWidths holds the max widths for [size, name, git] in one display column.
type colWidths [3]int

func computeColumns(cells []cell, showIcon bool, termW, pad int) (int, []colWidths) {
	n := len(cells)
	var best []colWidths

	for ncols := 1; ncols <= n; ncols++ {
		jump := int(math.Ceil(float64(n) / float64(ncols)))
		ws := make([]colWidths, ncols)
		for col := 0; col < ncols; col++ {
			start := col * jump
			end := start + jump
			if start >= n {
				start = n
			}
			if end > n {
				end = n
			}
			if start < end {
				ws[col] = maxWidths(cells[start:end])
			}
		}
		total := totalWidth(ws, showIcon, pad)
		if total <= termW {
			best = ws
			if best != nil && ncols == n {
				break // everything fits on one line
			}
		} else if best == nil {
			// Even 1 column doesn't fit: use 1 column anyway
			best = ws[:1]
			break
		} else {
			break
		}
	}
	if best == nil {
		ws := []colWidths{maxWidths(cells)}
		best = ws
	}
	return len(best), best
}

func maxWidths(cells []cell) colWidths {
	var cw colWidths
	for _, c := range cells {
		if c.sizeW > cw[0] {
			cw[0] = c.sizeW
		}
		if c.nameW > cw[1] {
			cw[1] = c.nameW
		}
		if c.gitW > cw[2] {
			cw[2] = c.gitW
		}
	}
	return cw
}

func totalWidth(ws []colWidths, showIcon bool, pad int) int {
	total := 0
	for i, cw := range ws {
		w := 0
		if cw[0] > 0 {
			w += cw[0] + 1
		}
		if showIcon {
			w += 2
		}
		w += cw[1]
		if cw[2] > 0 {
			w += 2
		}
		if i < len(ws)-1 {
			w += pad
		}
		total += w
	}
	return total
}

func writeCell(w io.Writer, c cell, cw colWidths, showIcon bool, opts Options) {
	reset := ""
	if opts.Colors {
		reset = noColor
	}

	// Size column
	if cw[0] > 0 {
		fmt.Fprintf(w, "%-*s%s", cw[0], c.size, brailSep)
	}

	// Icon column
	if showIcon {
		fmt.Fprintf(w, "%s%s%s%s", c.iconColor, c.icon, reset, brailSep)
	}

	// Name column (padded to max width in this display column)
	gc := ""
	gcReset := ""
	if opts.Colors {
		gc = gitColorStr(c.git, opts)
		if gc != "" {
			gcReset = reset
		}
	}
	padding := cw[1] - c.nameW
	fmt.Fprintf(w, "%s%s%s%s", gc, c.name, gcReset, strings.Repeat(" ", padding))

	// Git column
	if cw[2] > 0 && c.git != "" {
		fmt.Fprintf(w, "%s%s%s%s", brailSep, gc, c.git, gcReset)
	}
}

func gitColorStr(git string, opts Options) string {
	if !opts.Colors || opts.Theme == nil {
		return ""
	}
	return gitColor(git, opts.Theme)
}

func anyHasIcon(cells []cell) bool {
	for _, c := range cells {
		if c.icon != "" {
			return true
		}
	}
	return false
}

func buildCells(entries []*filesystem.Entry, opts Options) []cell {
	cells := make([]cell, len(entries))
	for i, e := range entries {
		name := e.FullName + string(e.Indicator)
		c := cell{
			icon:      e.Icon,
			iconColor: e.IconColor,
			name:      name,
			git:       e.GitStatus,
			nameW:     runewidth.StringWidth(name),
			gitW:      runewidth.StringWidth(e.GitStatus),
		}
		if opts.ShowBlocks {
			c.size = formatSize(e.Blocks*512, false)
			c.sizeW = len(c.size)
		}
		cells[i] = c
	}
	return cells
}

func formatSize(b int64, human bool) string {
	if !human {
		return fmt.Sprintf("%d", b)
	}
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%c", float64(b)/float64(div), "KMGTPE"[exp])
}
