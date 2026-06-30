// Package cli implements the modern-ls command-line interface using spf13/pflag.
package cli

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"

	"github.com/spf13/pflag"
	"github.com/the-mayankjha/modern-ls/internal/config"
	"github.com/the-mayankjha/modern-ls/internal/filesystem"
	"github.com/the-mayankjha/modern-ls/internal/git"
	"github.com/the-mayankjha/modern-ls/internal/icons"
	"github.com/the-mayankjha/modern-ls/internal/renderer"
	"github.com/the-mayankjha/modern-ls/internal/sorting"
	"github.com/the-mayankjha/modern-ls/internal/terminal"
	"github.com/the-mayankjha/modern-ls/internal/themes"
)

// Version is injected at build time via -ldflags.
var Version = "dev"

// exitCode tracks the process exit code during a run.
type exitCode int

const (
	exitOK      exitCode = 0
	exitMinor   exitCode = 1
	exitSerious exitCode = 2
)

// Run is the main entry point for the CLI. It parses flags, loads config,
// and runs the listing. It returns the process exit code.
func Run(args []string, stdout io.Writer) exitCode {
	fs := pflag.NewFlagSet("modern-ls", pflag.ContinueOnError)
	fs.SetOutput(stdout)

	// ── Content flags ──────────────────────────────────────────────────────────
	fA := fs.BoolP("almost-all", "A", false, "do not list implied . and ..")
	fa := fs.BoolP("all", "a", false, "do not ignore entries starting with .")

	// ── Display flags ──────────────────────────────────────────────────────────
	f1 := fs.BoolP("one-per-line", "1", false, "list one file per line")
	fd := fs.BoolP("directory", "d", false, "list directories themselves, not their contents")
	fl := fs.BoolP("long", "l", false, "use a long listing format")
	fo := fs.Bool("no-group-info", false, "like -l, but do not list group information (use -o historically)")
	fg := fs.Bool("no-owner", false, "like -l, but do not list owner (use -g historically)")
	fG := fs.BoolP("no-group", "G", false, "in a long listing, don't print group names")
	_ = fs.BoolP("human-readable", "h", false, "with -l and -s, print sizes like 1K 234M 2G etc.")
	fs2 := fs.BoolP("size", "s", false, "print the allocated size of each file, in blocks")

	// ── Sort flags ─────────────────────────────────────────────────────────────
	fS := fs.BoolP("sort-size", "S", false, "sort by file size, largest first")
	fU := fs.BoolP("no-sort", "U", false, "do not sort; list entries in directory order")
	fX := fs.BoolP("sort-ext", "X", false, "sort alphabetically by entry extension")
	fv := fs.BoolP("version-sort", "v", false, "natural sort of (version) numbers within text")
	ft := fs.BoolP("sort-time", "t", false, "sort by modification time, newest first")
	fr := fs.BoolP("reverse", "r", false, "reverse order while sorting")
	fR := fs.BoolP("recursive", "R", false, "list subdirectories recursively")
	fTree := fs.Bool("tree", false, "recurse into directories as a tree")
	fDepth := fs.Int("depth", 0, "limit the depth of the tree (0 means unlimited)")

	// ── Time ───────────────────────────────────────────────────────────────────
	fT := fs.StringP("time-style", "T", "Stamp", "time/date format: Stamp|StampMilli|Kitchen|ANSIC|UnixDate|RubyDate|RFC1123|RFC1123Z|RFC3339|RFC822|RFC822Z|RFC850")

	// ── modern-ls specific ─────────────────────────────────────────────────────
	fD := fs.BoolP("git-status", "D", false, "print git status of files")
	fc := fs.BoolP("disable-color", "c", false, "don't color icons, filenames and git status")
	fi := fs.BoolP("disable-icon", "i", false, "don't print icons of the files")
	fTheme := fs.String("theme", "", "color theme: default|catppuccin|tokyonight|gruvbox|dracula|nord|rose-pine")
	fConfig := fs.String("config", "", "path to config file")

	// ── Meta flags ─────────────────────────────────────────────────────────────
	fVersion := fs.BoolP("version", "V", false, "output version information and exit")
	fHelp := fs.BoolP("help", "?", false, "display this help and exit")

	// Legacy short flags preserved for compatibility
	fs.Bool("o", false, "like -l, but do not list group information")
	fs.Bool("g", false, "like -l, but do not list owner")

	if err := fs.Parse(args); err != nil {
		if err == pflag.ErrHelp {
			return exitOK
		}
		fmt.Fprintf(stdout, "modern-ls: %v\nTry 'modern-ls --help' for more information.\n", err)
		return exitSerious
	}

	// ── Help ────────────────────────────────────────────────────────────────
	if *fHelp {
		printHelp(stdout, fs)
		return exitOK
	}

	// ── Version ─────────────────────────────────────────────────────────────
	if *fVersion {
		fmt.Fprintf(stdout, "modern-ls %s\nCopyright (c) 2024 Mayank Jha\nLicense MIT <https://opensource.org/licenses/MIT>\nThis is free software: you are free to change and redistribute it.\nThere is NO WARRANTY, to the extent permitted by law.\n\nWritten by Mayank Jha\n", Version)
		return exitOK
	}

	// ── Load config ─────────────────────────────────────────────────────────
	var cfg config.Config
	var cfgErr error
	if *fConfig != "" {
		cfg, cfgErr = config.LoadFrom(*fConfig)
	} else {
		cfg, cfgErr = config.Load()
	}
	if cfgErr != nil {
		log.Printf("modern-ls: config: %v", cfgErr)
	}

	// ── Apply flag overrides on top of config ────────────────────────────────
	// Colors: disabled if -c flag; enabled by default (auto TTY detection)
	useColors := cfg.Display.Colors
	if *fc {
		useColors = false
	}
	if useColors {
		// Disable colors if stdout is not a terminal
		useColors = terminal.IsTerminal()
	}

	useIcons := cfg.Display.Icons
	if *fi {
		useIcons = false
	}

	useGit := cfg.Display.Git || *fD

	themeName := cfg.Display.Theme
	if *fTheme != "" {
		themeName = *fTheme
	}
	theme := themes.Get(themeName)

	// ── Determine sort strategy ──────────────────────────────────────────────
	sortStrategy := pickSort(*fS, *ft, *fX, *fv, *fU, *fr)

	// ── Determine layout ─────────────────────────────────────────────────────
	longMode := *fl || *fo || *fg
	onePerLine := *f1

	renderOpts := renderer.Options{
		Theme:      theme,
		TimeFormat: *fT,
		Colors:     useColors,
		Icons:      useIcons,
		ShowBlocks: *fs2,
		ShowGit:    useGit,
		ShowOwner:  !*fo && !*fg,
		ShowGroup:  !*fG && !*fg,
		ShowMode:   longMode,
	}

	fsOpts := filesystem.Options{
		ShowAll:    *fa,
		AlmostAll:  *fA,
		LongFormat: longMode,
		Blocks:     *fs2,
		Icons:      useIcons,
		Colors:     useColors,
		Git:        useGit,
		DirSelf:    *fd,
	}

	renderOpts.Tree = *fTree
	renderOpts.Depth = *fDepth

	// ── Collect paths from arguments ─────────────────────────────────────────
	paths := fs.Args()
	if len(paths) == 0 {
		paths = []string{"."}
	}
	sort.Strings(paths)

	// ── Segregate files from directories ────────────────────────────────────
	var fileInfos []os.FileInfo
	var dirs []*os.File
	ec := exitOK

	for _, p := range paths {
		f, err := os.Open(p)
		if err != nil {
			log.Printf("modern-ls: cannot access %q: %v", p, err)
			ec = max(ec, exitSerious)
			continue
		}
		stat, err := f.Stat()
		if err != nil {
			log.Printf("modern-ls: cannot access %q: %v", p, err)
			f.Close()
			ec = max(ec, exitSerious)
			continue
		}
		if stat.IsDir() {
			dirs = append(dirs, f)
		} else {
			fileInfos = append(fileInfos, stat)
			f.Close()
		}
	}

	// ── Icon resolver (injected closure) ─────────────────────────────────────
	iconFn := func(name, ext string, ind filesystem.Indicator) (string, string) {
		if !useIcons {
			return "", ""
		}
		return icons.Lookup(name, ext, ind)
	}

	// ── Render standalone files ──────────────────────────────────────────────
	if len(fileInfos) > 0 {
		entries := buildEntries(fileInfos, nil, fsOpts, iconFn)
		sortEntries(entries, sortStrategy)
		r := pickRenderer(longMode, onePerLine)
		if err := r.Render(stdout, entries, terminal.Width(), renderOpts); err != nil {
			log.Printf("modern-ls: render: %v", err)
		}
		if len(dirs) > 0 {
			fmt.Fprintln(stdout)
		}
	}

	// ── Render directories ───────────────────────────────────────────────────
	if *fTree {
		// Tree mode
		for i, d := range dirs {
			if i > 0 {
				fmt.Fprintln(stdout)
			}
			openDirIcon := dirOpenIcon(useIcons, useColors, theme)
			fmt.Fprintf(stdout, "%s%s\n", openDirIcon, d.Name())
			runTree(stdout, d, 1, *fDepth, "", fsOpts, renderOpts, sortStrategy, iconFn, useIcons, useColors, theme)
		}
	} else if *fR {
		// Recursive mode
		for i, d := range dirs {
			if i > 0 {
				fmt.Fprintln(stdout)
			}
			openDirIcon := dirOpenIcon(useIcons, useColors, theme)
			fmt.Fprintf(stdout, "%s%s:\n", openDirIcon, d.Name())
			runRecursive(stdout, d, fsOpts, renderOpts, sortStrategy, iconFn, useIcons, useColors, theme)
		}
	} else {
		// Non-recursive: print dir header only when multiple dirs given
		printHeader := len(dirs) > 1 || (len(fileInfos) > 0 && len(dirs) > 0)
		for i, d := range dirs {
			if printHeader {
				openDirIcon := dirOpenIcon(useIcons, useColors, theme)
				fmt.Fprintf(stdout, "%s%s:\n", openDirIcon, d.Name())
			}
			if useGit {
				// Reset git cache per directory when multiple dirs given
			}
			var gitStatus map[string]string
			if useGit {
				gs, err := git.ForDir(d.Name())
				if err != nil {
					log.Printf("modern-ls: git: %v", err)
				}
				gitStatus = gitStatusToStrings(gs)
			}
			entries, subdirs, err := filesystem.ReadDir(d, gitStatus, fsOpts, iconFn)
			d.Close()
			if err != nil {
				log.Printf("modern-ls: partial access to %q: %v", d.Name(), err)
				ec = max(ec, exitMinor)
			}
			_ = subdirs
			sortEntries(entries, sortStrategy)
			r := pickRenderer(longMode, onePerLine)
			if err := r.Render(stdout, entries, terminal.Width(), renderOpts); err != nil {
				log.Printf("modern-ls: render: %v", err)
			}
			if i < len(dirs)-1 {
				fmt.Fprintln(stdout)
			}
		}
	}

	return ec
}

// ── Helper functions ─────────────────────────────────────────────────────────

func pickRenderer(long, oneline bool) renderer.Renderer {
	switch {
	case long:
		return &renderer.Long{}
	case oneline:
		return &renderer.OneLine{}
	default:
		return &renderer.Grid{}
	}
}

func pickSort(bySize, byTime, byExt, byVersion, noSort, reverse bool) sorting.Strategy {
	var s sorting.Strategy
	switch {
	case bySize:
		s = sorting.BySize
	case byTime:
		s = sorting.ByTime
	case byExt:
		s = sorting.ByExt
	case byVersion:
		s = sorting.Natural
	case noSort:
		s = sorting.Unsorted
	default:
		s = sorting.Alpha
	}
	if reverse && !noSort {
		s = sorting.Reversed(s)
	}
	return s
}

func sortEntries(entries []*filesystem.Entry, sortStrat sorting.Strategy) {
	sort.SliceStable(entries, func(i, j int) bool {
		return sortStrat(entries[i], entries[j]) < 0
	})
}

func buildEntries(infos []os.FileInfo, gitStatus map[string]string, opts filesystem.Options, iconFn func(string, string, filesystem.Indicator) (string, string)) []*filesystem.Entry {
	entries := make([]*filesystem.Entry, 0, len(infos))
	for _, fi := range infos {
		full := fi.Name()
		name, ext := filesystem.SplitName(full)
		ind := filesystem.ModeIndicator(fi.Mode())
		gs := ""
		if gitStatus != nil {
			gs = gitStatus[full]
		}
		glyph, color := iconFn(name, ext, ind)
		entries = append(entries, &filesystem.Entry{
			Name:      name,
			Ext:       ext,
			FullName:  full,
			Indicator: ind,
			Size:      fi.Size(),
			ModTime:   fi.ModTime(),
			Mode:      fi.Mode(),
			ModeStr:   fi.Mode().String(),
			GitStatus: gs,
			Icon:      glyph,
			IconColor: color,
		})
	}
	return entries
}

func runRecursive(w io.Writer, d *os.File, fsOpts filesystem.Options, renderOpts renderer.Options, sortStrat sorting.Strategy, iconFn func(string, string, filesystem.Indicator) (string, string), useIcons, useColors bool, theme *themes.Theme) {
	var gitStatus map[string]string
	if fsOpts.Git {
		gs, _ := git.ForDir(d.Name())
		gitStatus = gitStatusToStrings(gs)
	}
	entries, subdirs, err := filesystem.ReadDir(d, gitStatus, fsOpts, iconFn)
	dirName := d.Name()
	d.Close()
	if err != nil {
		log.Printf("modern-ls: partial access to %q: %v", dirName, err)
	}
	sortEntries(entries, sortStrat)
	r := pickRenderer(renderOpts.ShowMode, false)
	r.Render(w, entries, terminal.Width(), renderOpts)

	sort.Strings(subdirs)
	openIcon := dirOpenIcon(useIcons, useColors, theme)
	for _, sub := range subdirs {
		fullSub := filepath.Join(dirName, sub)
		fmt.Fprintf(w, "\n%s%s:\n", openIcon, fullSub)
		f, err := os.Open(fullSub)
		if err != nil {
			log.Printf("modern-ls: cannot access %q: %v", fullSub, err)
			continue
		}
		runRecursive(w, f, fsOpts, renderOpts, sortStrat, iconFn, useIcons, useColors, theme)
	}
}

func dirOpenIcon(useIcons, useColors bool, theme *themes.Theme) string {
	if !useIcons {
		return ""
	}
	glyph := icons.Defaults["diropen"].Glyph
	if !useColors || theme == nil {
		return glyph + " "
	}
	return theme.DirOpen.ANSI() + glyph + themes.Reset + " "
}

func gitStatusToStrings(gs git.RepoStatus) map[string]string {
	if gs == nil {
		return nil
	}
	m := make(map[string]string, len(gs))
	for k, v := range gs {
		m[k] = v.String()
	}
	return m
}

func max(a, b exitCode) exitCode {
	if a > b {
		return a
	}
	return b
}

func printHelp(w io.Writer, fs *pflag.FlagSet) {
	fmt.Fprintln(w, "modern-ls — List information about FILEs with icons and git status.")
	fmt.Fprintln(w, "Sort entries alphabetically if none of -tvSUX is specified.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Usage: modern-ls [OPTIONS] [FILE...]")
	fmt.Fprintln(w)
	fs.PrintDefaults()
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Possible values for --time-style (-T):")
	fmt.Fprintf(w, "  %-11s %q\n", "ANSIC", "Mon Jan _2 15:04:05 2006")
	fmt.Fprintf(w, "  %-11s %q\n", "UnixDate", "Mon Jan _2 15:04:05 MST 2006")
	fmt.Fprintf(w, "  %-11s %q\n", "RFC3339", "2006-01-02T15:04:05Z07:00")
	fmt.Fprintf(w, "  %-11s %q [Default]\n", "Stamp", "Jan _2 15:04:05")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Exit status:")
	fmt.Fprintln(w, "  0  if OK")
	fmt.Fprintln(w, "  1  if minor problems (e.g., cannot access subdirectory)")
	fmt.Fprintln(w, "  2  if serious trouble (e.g., cannot access command-line argument)")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Source: https://github.com/the-mayankjha/modern-ls\n")
}

// Ensure runtime.GOOS doesn't create unused import warnings on Windows.
var _ = runtime.GOOS

func runTree(w io.Writer, d *os.File, currentDepth, maxDepth int, prefix string, fsOpts filesystem.Options, renderOpts renderer.Options, sortStrat sorting.Strategy, iconFn func(string, string, filesystem.Indicator) (string, string), useIcons, useColors bool, theme *themes.Theme) {
	if maxDepth > 0 && currentDepth > maxDepth {
		return
	}

	var gitStatus map[string]string
	if fsOpts.Git {
		gs, _ := git.ForDir(d.Name())
		gitStatus = gitStatusToStrings(gs)
	}

	entries, _, err := filesystem.ReadDir(d, gitStatus, fsOpts, iconFn)
	dirName := d.Name()
	d.Close()
	if err != nil {
		log.Printf("modern-ls: partial access to %q: %v", dirName, err)
	}

	sortEntries(entries, sortStrat)

	for i, e := range entries {
		isLast := i == len(entries)-1

		t := &renderer.Tree{
			Prefix:     prefix,
			IsLast:     isLast,
			TimeFormat: renderOpts.TimeFormat,
		}
		t.RenderEntry(w, e, renderOpts)

		if e.Indicator == filesystem.IndicatorDir {
			fullSub := filepath.Join(dirName, e.Name+e.Ext)
			subD, err := os.Open(fullSub)
			if err != nil {
				log.Printf("modern-ls: cannot access %q: %v", fullSub, err)
				continue
			}
			newPrefix := prefix
			if isLast {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}
			runTree(w, subD, currentDepth+1, maxDepth, newPrefix, fsOpts, renderOpts, sortStrat, iconFn, useIcons, useColors, theme)
		}
	}
}
