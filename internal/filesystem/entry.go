// Package filesystem provides directory reading and file entry construction.
package filesystem

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Indicator represents a special file type indicator character.
type Indicator string

const (
	IndicatorDir    Indicator = "/"
	IndicatorPipe   Indicator = "|"
	IndicatorSymlink Indicator = "@"
	IndicatorSocket  Indicator = "="
	IndicatorExec    Indicator = "*"
	IndicatorNone    Indicator = ""
)

// Entry represents a single file or directory entry with all metadata needed
// for rendering. It is intentionally immutable after construction.
type Entry struct {
	Name      string    // base name without extension (for plain files)
	Ext       string    // extension including dot (e.g. ".go")
	FullName  string    // full name = Name + Ext
	Indicator Indicator // "/" for dir, "@" for symlink, etc.
	Size      int64     // file size in bytes
	Blocks    int64     // allocated blocks (Unix only)
	ModTime   time.Time
	Mode      fs.FileMode
	ModeStr   string
	Owner     string
	Group     string
	GitStatus string
	Icon      string
	IconColor string // ANSI escape for icon color; empty if colors disabled
	IsHidden  bool   // true if name starts with "."
}

// GetName implements sorting.Entry.
func (e *Entry) GetName() string { return e.FullName }

// GetExt implements sorting.Entry.
func (e *Entry) GetExt() string { return e.Ext }

// GetSize implements sorting.Entry.
func (e *Entry) GetSize() int64 { return e.Size }

// GetModTime implements sorting.Entry.
func (e *Entry) GetModTime() time.Time { return e.ModTime }

// Options configures what metadata is collected during directory reading.
type Options struct {
	ShowHidden bool // include entries starting with "."
	ShowAll    bool // include "." and ".." (implies ShowHidden)
	AlmostAll  bool // include entries starting with "." but not "." or ".."
	LongFormat bool // collect mode, owner, group
	Blocks     bool // collect block count
	Icons      bool // resolve icon glyph + color
	Colors     bool // include ANSI color codes
	Git        bool // collect git status
	DirSelf    bool // list the directory itself, not its contents
}

// ModeIndicator converts a file mode to an Indicator.
func ModeIndicator(m fs.FileMode) Indicator {
	switch {
	case m&fs.ModeDir != 0:
		return IndicatorDir
	case m&fs.ModeNamedPipe != 0:
		return IndicatorPipe
	case m&fs.ModeSymlink != 0:
		return IndicatorSymlink
	case m&fs.ModeSocket != 0:
		return IndicatorSocket
	case m&0o111 != 0:
		return IndicatorExec
	default:
		return IndicatorNone
	}
}

// SplitName splits a filename into (name, ext) where ext includes the dot.
// Hidden files like ".gitignore" return (".gitignore", "").
// Regular files like "main.go" return ("main", ".go").
func SplitName(full string) (name, ext string) {
	// Don't treat leading dot as extension separator
	if strings.HasPrefix(full, ".") {
		// Could still have an extension, e.g. ".bashrc.bak"
		rest := full[1:]
		if i := strings.LastIndex(rest, "."); i >= 0 {
			return full[:1+i+1], full[1+i+1:]
		}
		return full, ""
	}
	ext = filepath.Ext(full)
	name = full[:len(full)-len(ext)]
	return name, ext
}

// buildEntry creates an Entry from an os.FileInfo, path, and options.
// gitStatus is the pre-computed git status for this entry (may be empty).
func buildEntry(fi fs.FileInfo, gitStatus string, opts Options, iconFn func(name, ext string, ind Indicator) (glyph, color string)) *Entry {
	full := fi.Name()
	name, ext := SplitName(full)
	ind := ModeIndicator(fi.Mode())

	e := &Entry{
		Name:      name,
		Ext:       ext,
		FullName:  full,
		Indicator: ind,
		Size:      fi.Size(),
		ModTime:   fi.ModTime(),
		Mode:      fi.Mode(),
		GitStatus: gitStatus,
		IsHidden:  strings.HasPrefix(full, "."),
	}

	if opts.LongFormat {
		e.ModeStr = fi.Mode().String()
		e.Owner, e.Group = ownerGroup(fi)
	}
	if opts.Blocks {
		e.Blocks = fileBlocks(fi)
	}
	if opts.Icons && iconFn != nil {
		glyph, color := iconFn(name, ext, ind)
		e.Icon = glyph
		if opts.Colors {
			e.IconColor = color
		}
	}
	return e
}

// ReadDir reads the given open directory and returns a slice of Entries.
// The iconFn parameter is injected to avoid a direct dependency on the icons package.
func ReadDir(d *os.File, gitStatus map[string]string, opts Options, iconFn func(name, ext string, ind Indicator) (glyph, color string)) ([]*Entry, []string, error) {
	infos, err := d.Readdir(0)

	var entries []*Entry
	var subdirs []string

	showHidden := opts.ShowHidden || opts.ShowAll || opts.AlmostAll

	for _, fi := range infos {
		full := fi.Name()

		// Filter hidden files
		if !showHidden && strings.HasPrefix(full, ".") {
			continue
		}

		gs := gitStatus[full]
		if fi.IsDir() {
			gs = gitStatus[full+"/"]
		}

		e := buildEntry(fi, gs, opts, iconFn)
		entries = append(entries, e)

		if fi.IsDir() {
			subdirs = append(subdirs, full+"/")
		}
	}

	return entries, subdirs, err
}
