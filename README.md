# modern-ls 🚀

A beautiful, extremely fast, modern replacement for Unix `ls`, written in Go.

![License](https://img.shields.io/github/license/the-mayankjha/modern-ls)
![Release](https://img.shields.io/github/v/release/the-mayankjha/modern-ls)
![Build](https://img.shields.io/github/actions/workflow/status/the-mayankjha/modern-ls/ci.yml)

`modern-ls` is a highly-optimized directory lister that features:
- 🎨 **Icons**: Real `nvim-web-devicons` support baked natively, rendering Nerd Font v3 icons for every file, folder, and extension.
- 🌳 **Semantic Folders**: Folders like `src`, `web`, `admin`, `scripts`, etc., get custom icons for easy spotting.
- 🐙 **Git Status**: Instantly highlights untracked, modified, added, or deleted files.
- ⚡ **Blazing Fast**: Written in pure Go, taking heavy advantage of concurrent directory traversal and optimized grid rendering to beat standard Python or Ruby tools.
- 🛠️ **Cross Platform**: Works flawlessly on macOS (Apple Silicon natively supported), Linux, and Windows.

## Installation

See [INSTALL.md](INSTALL.md) for detailed instructions across all platforms.

### Quick Start (macOS & Linux)
If you have Go installed, you can instantly install `modern-ls`:
```bash
go install github.com/the-mayankjha/modern-ls/cmd/modern-ls@latest
```

Alternatively, download the pre-compiled binaries from the [GitHub Releases](https://github.com/the-mayankjha/modern-ls/releases) page for Windows, macOS, and Linux, extract them, and move them to your system `PATH`.

## Features & Usage

Run `modern-ls` in any directory.

```bash
# Basic listing
modern-ls

# Long listing with human-readable sizes
modern-ls -lh

# Show all hidden files, sorted by time
modern-ls -at

# Display Git statuses
modern-ls -D
```

By default, the grid view is used. You can force one-per-line using `-1`.

### Colors and Themes

`modern-ls` includes beautiful native themes:
- `default`
- `catppuccin`
- `dracula`
- `gruvbox`
- `nord`
- `tokyonight`
- `rose-pine`

You can use the `--theme` flag to test them, or set them globally in `~/.config/modern-ls/config.yml`.

## Configuration

You can fully customize modern-ls. Create `~/.config/modern-ls/config.yml`:

```yaml
display:
  theme: catppuccin
  icons: true
  colors: true
  git: true
```

## License

MIT License. See [LICENSE](LICENSE) for details.
Written by Mayank Jha.
