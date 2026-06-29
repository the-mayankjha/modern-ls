# Architecture

`modern-ls` is designed with a strict separation of concerns, ensuring high maintainability, performance, and cross-platform consistency.

## Core Packages

- **`cmd/modern-ls`**: The main binary entrypoint. Delegates immediately to `internal/cli`.
- **`internal/cli`**: Parses CLI arguments (using `pflag`), orchestrates configuration loading, constructs the core `filesystem.Options`, runs the directory traversal, and connects it to the `renderer`.
- **`internal/filesystem`**: The heart of the program. It abstracts `os.FileInfo` into `filesystem.Entry` structs. It contains optimized parallel directory scanning.
- **`internal/icons`**: Exposes the `Lookup(name, ext, ind)` function. This package holds no I/O logic at runtime. All icon mappings are embedded in `icons_generated.go` via `go generate`.
- **`internal/renderer`**: Implements the `Renderer` interface (`Grid`, `OneLine`, `Long`). Uses ANSI escape sequences to format the display strings.
- **`internal/themes`**: Houses statically defined native palettes (Catppuccin, Nord, etc.) for syntax highlighting and terminal coloring.
- **`generator`**: A standalone tool that builds the icon dictionaries at compile-time by parsing upstream `nvim-web-devicons` Lua files and outputting standard Go maps.

## Data Flow
1. **Config & CLI parsing**: Load YAML settings, apply CLI flags on top.
2. **Path Traversal**: Read target directories via `filesystem.ReadDir()`.
3. **Stat + Icon Resolution**: Fetch metadata via `os.Stat` and resolve icons synchronously as we create `filesystem.Entry` instances.
4. **Git Resolution**: Spawn `go-git` synchronously to fetch untracked/modified mappings.
5. **Sorting**: Sort the array of `*filesystem.Entry` pointers.
6. **Rendering**: The renderer iterates through the slice and flushes the formatted ANSI strings to `os.Stdout`.
