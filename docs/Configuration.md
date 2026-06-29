# Configuration

`modern-ls` allows you to set default preferences globally so you don't have to alias long CLI flags.

## File Location

Create a YAML config file at:
- macOS / Linux: `~/.config/modern-ls/config.yml`
- Windows: `%APPDATA%\modern-ls\config.yml`

You can also pass a custom location at runtime using `--config=/path/to/config.yml`.

## Settings Schema

```yaml
display:
  # The color theme to use (default, catppuccin, dracula, etc)
  theme: "default"
  
  # Set to false to disable all Nerd Font icons
  icons: true
  
  # Set to false to disable all ANSI terminal coloring
  colors: true
  
  # Set to false to disable Git status checking
  git: false
```

Note: Any CLI flag (like `--disable-color` or `-c`) overrides the global configuration.
