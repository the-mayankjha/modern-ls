// Command generator parses vendored nvim-web-devicons Lua tables, writes
// assets/icons.json, and then generates internal/icons/icons_generated.go so
// that the runtime never needs to do I/O to look up an icon.
//
// Run from the project root:
//
//	go run ./generator/...
package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/the-mayankjha/modern-ls/generator/exporter"
	"github.com/the-mayankjha/modern-ls/generator/parser"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("generator: ")

	root, err := projectRoot()
	if err != nil {
		log.Fatalf("cannot determine project root: %v", err)
	}
	log.Printf("project root: %s", root)

	// ── 1. Parse vendor Lua files ──────────────────────────────────────────
	vendorDir := filepath.Join(root, "generator", "vendor", "nvim-web-devicons")

	byFilename, err := parseLuaFile(filepath.Join(vendorDir, "icons_by_filename.lua"))
	if err != nil {
		log.Fatalf("parse icons_by_filename.lua: %v", err)
	}
	log.Printf("parsed %d filename entries", len(byFilename))

	byExtension, err := parseLuaFile(filepath.Join(vendorDir, "icons_by_file_extension.lua"))
	if err != nil {
		log.Fatalf("parse icons_by_file_extension.lua: %v", err)
	}

	// Add/override custom icons for specific file types
	byExtension["dmg"] = parser.Icon{Name: "dmg", Glyph: "💿", Color: "#E5E5E5"}
	byExtension["exe"] = parser.Icon{Name: "exe", Glyph: "🪟", Color: "#0078D7"}
	byExtension["tar"] = parser.Icon{Name: "tar", Glyph: "📦", Color: "#ECA517"}
	byExtension["zip"] = parser.Icon{Name: "zip", Glyph: "📦", Color: "#ECA517"}
	byExtension["rar"] = parser.Icon{Name: "rar", Glyph: "📦", Color: "#ECA517"}

	log.Printf("parsed %d extension entries", len(byExtension))

	// ── 2. Build defaults and folders ──────────────────────────────────────
	defaults := buildDefaults()
	byFolder := buildFolderIcons()

	// ── 3. Write assets/icons.json ────────────────────────────────────────
	db := exporter.IconDB{
		SchemaVersion: 1,
		ByFilename:    byFilename,
		ByExtension:   byExtension,
		ByFolder:      byFolder,
		Defaults:      defaults,
	}

	assetsDir := filepath.Join(root, "assets")
	if err := os.MkdirAll(assetsDir, 0o755); err != nil {
		log.Fatalf("mkdir assets: %v", err)
	}

	jsonPath := filepath.Join(assetsDir, "icons.json")
	if err := writeJSON(jsonPath, db); err != nil {
		log.Fatalf("write icons.json: %v", err)
	}
	log.Printf("wrote %s", jsonPath)

	// ── 4. Read back and generate Go source ───────────────────────────────
	dbRead, err := readJSON(jsonPath)
	if err != nil {
		log.Fatalf("read icons.json: %v", err)
	}

	goSrc, err := generateGo(dbRead)
	if err != nil {
		log.Fatalf("generate Go source: %v", err)
	}

	iconsDir := filepath.Join(root, "internal", "icons")
	if err := os.MkdirAll(iconsDir, 0o755); err != nil {
		log.Fatalf("mkdir internal/icons: %v", err)
	}

	goPath := filepath.Join(iconsDir, "icons_generated.go")
	if err := os.WriteFile(goPath, goSrc, 0o644); err != nil {
		log.Fatalf("write icons_generated.go: %v", err)
	}
	log.Printf("wrote %s", goPath)
}

// projectRoot walks up from the working directory until it finds go.mod.
// This ensures the generator works correctly regardless of where it is invoked.
func projectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found in any ancestor directory of %s", dir)
		}
		dir = parent
	}
}

// parseLuaFile opens path and delegates to parser.ParseLua.
func parseLuaFile(path string) (map[string]parser.Icon, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parser.ParseLua(f)
}

// buildDefaults returns the hardcoded fallback icons for generic file kinds.
// These are not present in nvim-web-devicons and must be supplied manually.
func buildDefaults() map[string]parser.Icon {
	return map[string]parser.Icon{
		// Generic regular file
		"file": {
			Name:  "file",
			Glyph: "\uf016", // nf-fa-file (U+F016)
			Color: "#C5C8C6",
		},
		// Generic directory (closed)
		"dir": {
			Name:  "dir",
			Glyph: "\uf07b", // nf-fa-folder (U+F07B)
			Color: "#81A2BE",
		},
		// Generic directory (open / expanded)
		"diropen": {
			Name:  "diropen",
			Glyph: "\uf07c", // nf-fa-folder_open (U+F07C)
			Color: "#81A2BE",
		},
		// Dot-file (hidden regular file)
		"hiddenfile": {
			Name:  "hiddenfile",
			Glyph: "\uf016", // nf-fa-file (U+F016)
			Color: "#8C9440",
		},
		// Hidden directory
		"hiddendir": {
			Name:  "hiddendir",
			Glyph: "\uf07b", // nf-fa-folder (U+F07B)
			Color: "#5F819D",
		},
		// Executable file
		"exe": {
			Name:  "exe",
			Glyph: "\uf013", // nf-fa-cog (U+F013)
			Color: "#B5BD68",
		},
		".config": {Name: ".config", Glyph: "\uf013", Color: "#c5c8c6"},
		"admin":   {Name: "admin", Glyph: "\uf023", Color: "#81a2be"}, // nf-fa-lock
		"roles":   {Name: "roles", Glyph: "\uf007", Color: "#b294bb"}, // nf-fa-user
		"users":   {Name: "users", Glyph: "\uf0c0", Color: "#b5bd68"}, // nf-fa-users
	}
}

// buildFolderIcons returns a hardcoded map of specific folder icons.
func buildFolderIcons() map[string]parser.Icon {
	return map[string]parser.Icon{
		"web":          {Name: "web", Glyph: "\uf0ac", Color: "#519aba"},
		"server":       {Name: "server", Glyph: "\uf233", Color: "#81a2be"},
		"client":       {Name: "client", Glyph: "\uf108", Color: "#5f819d"},
		"models":       {Name: "models", Glyph: "\uf1b2", Color: "#f2a272"},
		"scripts":      {Name: "scripts", Glyph: "\uf121", Color: "#8c9440"},
		"public":       {Name: "public", Glyph: "\uf0ac", Color: "#519aba"},
		"src":          {Name: "src", Glyph: "\uf121", Color: "#5f819d"},
		"test":         {Name: "test", Glyph: "\uf0c3", Color: "#de935f"},
		"tests":        {Name: "tests", Glyph: "\uf0c3", Color: "#de935f"},
		"docs":         {Name: "docs", Glyph: "\uf02d", Color: "#81a2be"},
		"assets":       {Name: "assets", Glyph: "\uf1fc", Color: "#c5c8c6"},
		"images":       {Name: "images", Glyph: "\uf03e", Color: "#b5bd68"},
		"img":          {Name: "img", Glyph: "\uf03e", Color: "#b5bd68"},
		"views":        {Name: "views", Glyph: "\uf06e", Color: "#8abeb7"},
		"styles":       {Name: "styles", Glyph: "\ue749", Color: "#519aba"},
		"css":          {Name: "css", Glyph: "\ue749", Color: "#519aba"},
		"js":           {Name: "js", Glyph: "\ue781", Color: "#f5c06f"},
		"node_modules": {Name: "node_modules", Glyph: "\ue5fa", Color: "#cc3e44"},
		".github":      {Name: ".github", Glyph: "\uf09b", Color: "#c5c8c6"},
		".git":         {Name: ".git", Glyph: "\uf1d3", Color: "#f54d27"},
		"bin":          {Name: "bin", Glyph: "\uf013", Color: "#b5bd68"},
		"config":       {Name: "config", Glyph: "\uf013", Color: "#c5c8c6"},
		".config":      {Name: ".config", Glyph: "\uf013", Color: "#c5c8c6"},
		"admin":        {Name: "admin", Glyph: "\uf023", Color: "#81a2be"}, // nf-fa-lock
		"roles":        {Name: "roles", Glyph: "\uf007", Color: "#b294bb"}, // nf-fa-user
		"users":        {Name: "users", Glyph: "\uf0c0", Color: "#b5bd68"}, // nf-fa-users
	}
}

// writeJSON writes db as indented JSON to path (creating/truncating the file).
func writeJSON(path string, db exporter.IconDB) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return exporter.WriteJSON(f, db)
}

// readJSON opens path and delegates to exporter.ReadJSON.
func readJSON(path string) (exporter.IconDB, error) {
	f, err := os.Open(path)
	if err != nil {
		return exporter.IconDB{}, err
	}
	defer f.Close()
	return exporter.ReadJSON(f)
}

// ── Code generation ──────────────────────────────────────────────────────────

// iconEntry is the template model for one map initialiser entry.
type iconEntry struct {
	Key   string
	Glyph string
	Color string
}

// goTemplate is the text/template source for icons_generated.go.
// The outer backtick-delimited string ensures actual glyph characters are
// embedded verbatim (no escape sequences needed in the generated source).
const goTemplate = `// Code generated by go generate; DO NOT EDIT.
// Source: go run ./generator/...

package icons

// IconEntry holds a single Nerd Font icon and its associated foreground colour.
type IconEntry struct {
	// Glyph is the Nerd Font v3 character for this icon.
	Glyph string
	// Color is the recommended hex foreground colour, e.g. "#F05340".
	Color string
}

var (
	// ByFilename maps exact filenames (e.g. "go.mod", "Dockerfile") to icons.
	ByFilename map[string]IconEntry
	// ByExtension maps file extensions without the leading dot (e.g. "go", "ts") to icons.
	ByExtension map[string]IconEntry
	// ByFolder maps specific folder names to icons.
	ByFolder map[string]IconEntry
	// Defaults holds fallback icons keyed by generic kind:
	//   "file", "dir", "diropen", "hiddenfile", "hiddendir", "exe".
	Defaults map[string]IconEntry
)

func init() {
	ByFilename = map[string]IconEntry{
{{- range .ByFilename}}
		{{printf "%q" .Key}}: {Glyph: "{{.Glyph}}", Color: "{{.Color}}"},
{{- end}}
	}

	ByExtension = map[string]IconEntry{
{{- range .ByExtension}}
		{{printf "%q" .Key}}: {Glyph: "{{.Glyph}}", Color: "{{.Color}}"},
{{- end}}
	}

	ByFolder = map[string]IconEntry{
{{- range .ByFolder}}
		{{printf "%q" .Key}}: {Glyph: "{{.Glyph}}", Color: "{{.Color}}"},
{{- end}}
	}

	Defaults = map[string]IconEntry{
{{- range .Defaults}}
		{{printf "%q" .Key}}: {Glyph: "{{.Glyph}}", Color: "{{.Color}}"},
{{- end}}
	}
}
`

// templateData holds sorted slices for deterministic template execution.
type templateData struct {
	ByFilename  []iconEntry
	ByExtension []iconEntry
	ByFolder    []iconEntry
	Defaults    []iconEntry
}

// generateGo renders goTemplate with db and gofmt-formats the result.
func generateGo(db exporter.IconDB) ([]byte, error) {
	data := templateData{
		ByFilename:  toEntries(db.ByFilename),
		ByExtension: toEntries(db.ByExtension),
		ByFolder:    toEntries(db.ByFolder),
		Defaults:    toEntries(db.Defaults),
	}

	tmpl, err := template.New("icons").Parse(goTemplate)
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("execute template: %w", err)
	}

	// Run gofmt so the file is always canonical.
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		// Include the raw source in the error to aid debugging.
		return nil, fmt.Errorf("gofmt: %w\n--- raw source ---\n%s", err, buf.String())
	}
	return formatted, nil
}

// toEntries converts a map to a sorted slice of iconEntry values.
func toEntries(m map[string]parser.Icon) []iconEntry {
	keys := exporter.SortedKeys(m)
	entries := make([]iconEntry, 0, len(keys))
	for _, k := range keys {
		ic := m[k]
		// Sanitise the glyph: strip surrounding double quotes if the Lua parser
		// captured them, and replace any raw tab/newline that would break the
		// generated string literal.
		glyph := strings.TrimSpace(ic.Glyph)
		entries = append(entries, iconEntry{
			Key:   k,
			Glyph: glyph,
			Color: ic.Color,
		})
	}
	return entries
}
