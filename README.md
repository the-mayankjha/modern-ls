<h1 align="center">modern-ls 🚀</h1>

<p align="center">
  <b>A beautiful, extremely fast, modern replacement for Unix <code>ls</code> written in Go.</b>
</p>

<p align="center">
  <img src="https://img.shields.io/github/license/the-mayankjha/modern-ls?style=for-the-badge&color=blue" alt="License">
  <img src="https://img.shields.io/github/v/release/the-mayankjha/modern-ls?style=for-the-badge&color=green" alt="Release">
  <img src="https://img.shields.io/github/go-mod/go-version/the-mayankjha/modern-ls?style=for-the-badge&color=cyan" alt="Go Version">
</p>

---

`modern-ls` is a cross-platform (Windows, macOS, Linux) CLI utility that brings the classic `ls` command into the 21st century. It features native Git integration, highly customizable Nerd Font icons, vibrant color themes, and an incredibly fast recursive tree-view mode.

## ✨ Features

- 🎨 **Beautiful Color Themes:** Built-in support for popular themes like `catppuccin`, `tokyonight`, `gruvbox`, `dracula`, `nord`, and `rose-pine`.
- 🏷️ **Nerd Font Icons:** Extensive file-type recognition with beautiful glyphs (including custom mappings for `.dmg`, `.exe`, `.zip`, `.tar`, `.rar`, and hundreds more).
- 🌲 **Tree View Engine:** Traverse your directories natively with the `--tree` flag, and control traversal using the `--depth` flag!
- 🐙 **Git Integration:** See tracked, untracked, modified, and deleted Git file statuses natively in the long-listing outputs (`-l`).
- ⚡ **Blazing Fast:** Written in Go with minimal dependencies, compiling into a tiny, zero-dependency static binary.
- 🔄 **Safe Upgrades:** Integrated `--upgrade` command pointing directly to your native package managers.

## 📦 Installation

`modern-ls` is packaged for all major operating systems. You can install it using your preferred package manager:

### macOS / Linux (Homebrew)
```bash
brew tap the-mayankjha/homebrew-tap
brew install modern-ls
```

### Windows (Scoop)
```powershell
scoop bucket add mayankjha https://github.com/the-mayankjha/scoop-bucket.git
scoop install modern-ls
```

### Windows (Winget)
```powershell
winget install the-mayankjha.modern-ls
```

### Debian / Ubuntu (APT)
Download the `.deb` release from the [Releases page](https://github.com/the-mayankjha/modern-ls/releases) and run:
```bash
sudo apt install ./modern-ls_*_linux_amd64.deb
```

### Fedora / RHEL (DNF)
Download the `.rpm` release from the [Releases page](https://github.com/the-mayankjha/modern-ls/releases) and run:
```bash
sudo dnf install ./modern-ls_*_linux_amd64.rpm
```

### Arch Linux (Pacman)
Download the `.pkg.tar.zst` release from the [Releases page](https://github.com/the-mayankjha/modern-ls/releases) and run:
```bash
sudo pacman -U ./modern-ls_*_linux_amd64.pkg.tar.zst
```

## 🚀 Usage

`modern-ls` acts as a drop-in replacement for standard `ls`, accepting many of the same standard POSIX flags.

```bash
# Standard listing
modern-ls

# Long listing with git status, human readable sizes, and all hidden files
modern-ls -laD

# Display as a tree (limit depth to 2)
modern-ls --tree --depth=2

# Change the color theme to tokyonight
modern-ls --theme=tokyonight

# Get author and version info
modern-ls --info
```

### Supported Flags
- `-l`, `--long`: Use long listing format
- `-a`, `--all`: Include hidden files (starting with `.`)
- `-A`, `--almost-all`: Do not list implied `.` and `..`
- `-D`, `--git-status`: Print Git status of files
- `-c`, `--disable-color`: Turn off colors
- `--disable-icon`: Turn off file icons
- `--theme <name>`: Switch the color theme (`default`, `catppuccin`, `tokyonight`, `gruvbox`, `dracula`, `nord`, `rose-pine`)
- `--tree`: Recurse into directories as a tree
- `--depth <N>`: Limit the depth of the tree (0 means unlimited)
- `-v`, `--version`: Print version information
- `-i`, `--info`: Print author and project details
- `--upgrade`: View package-manager upgrade instructions

## 🤝 Contributing

Contributions are always welcome! Feel free to open issues or submit Pull Requests.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 👨‍💻 Author

**Mayank Kumar Jha**
- Portfolio: [https://mayankjha.nfks.co.in/](https://mayankjha.nfks.co.in/)
- GitHub: [@the-mayankjha](https://github.com/the-mayankjha)
- LinkedIn: [the-mayankjha](https://linkedin.com/in/the-mayankjha)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
