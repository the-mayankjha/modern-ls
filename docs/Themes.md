# Themes

`modern-ls` comes with stunning, carefully-curated color palettes natively built-in.

## Available Themes

- `default`
- `catppuccin` (Macchiato variant)
- `dracula`
- `gruvbox`
- `nord`
- `rose-pine`
- `tokyonight`

## How it works

When a theme is active, it colors the textual representations of the files (the filename itself) depending on its type:
- **Folders**
- **Executables**
- **Symlinks**
- **Regular Files**
- **Hidden Files**
- **Open Directories (Recursive Mode)**

> Note: The icons themselves keep their native `nvim-web-devicons` colors by design. The theme only affects the terminal UI elements and text.

## Setting the theme

You can test themes dynamically:
```bash
modern-ls --theme=nord
modern-ls --theme=catppuccin
```

Or you can make it permanent by putting it in `~/.config/modern-ls/config.yml`:
```yaml
display:
  theme: tokyonight
```
