<h1 align="center">modern-ls 🚀</h1>

<p align="center">
  <b>A beautiful, extremely fast, modern replacement for Unix <code>ls</code> written in Go.</b>
</p>

<p align="center">
  <a href="https://github.com/the-mayankjha/modern-ls/releases"><img src="https://img.shields.io/github/v/tag/the-mayankjha/modern-ls?sort=semver&style=for-the-badge&color=green&label=Version" alt="Version"></a>
  <a href="https://github.com/the-mayankjha/modern-ls/actions"><img src="https://img.shields.io/github/actions/workflow/status/the-mayankjha/modern-ls/release.yml?style=for-the-badge" alt="Build Status"></a>
  <img src="https://img.shields.io/github/go-mod/go-version/the-mayankjha/modern-ls?style=for-the-badge&color=cyan" alt="Go Version">
  <a href="https://github.com/the-mayankjha/modern-ls/blob/main/LICENSE"><img src="https://img.shields.io/github/license/the-mayankjha/modern-ls?style=for-the-badge&color=blue" alt="License"></a>
  <a href="https://github.com/the-mayankjha/modern-ls/stargazers"><img src="https://img.shields.io/github/stars/the-mayankjha/modern-ls?style=for-the-badge&color=yellow" alt="Stars"></a>
</p>

---

`modern-ls` is a cross-platform (Windows, macOS, Linux) CLI utility that brings the classic `ls` command into the 21st century. It features native Git integration, highly customizable Nerd Font and Emoji icons, vibrant color themes, and an incredibly fast recursive tree-view mode.

<img width="1200" height="800" alt="demo" src="https://github.com/user-attachments/assets/7a10a465-4173-4863-be3f-1ae7bbfcea8e" />


## ✨ Features

- 🎨 **Beautiful Color Themes:** Built-in support for popular themes like `catppuccin` (default), `tokyonight`, `gruvbox`, `dracula`, `nord`, and `rose-pine`.
- 🏷️ **Nerd Font & Emoji Icons:** Extensive file-type recognition with beautiful glyphs (including custom native Emoji mappings for `.bin`, `.mp4`, `.jpg`, `.env` and hundreds more).
- 🌲 **Tree View Engine:** Traverse your directories natively with the `--tree` flag, and control traversal using the `--depth` flag!
- 🐙 **Git Integration:** See tracked, untracked, modified, and deleted Git file statuses natively in the long-listing outputs (`-l`).
- ⚡ **Blazing Fast:** Written in Go with minimal dependencies, compiling into a tiny, zero-dependency static binary.
- 🔄 **Safe Upgrades:** Integrated `--upgrade` command pointing directly to your native package managers.

---

## 🎨 Themes

`modern-ls` ships with a vibrant set of color themes out of the box, with `catppuccin` set as the default. Every single file, directory, and Git status is color-coordinated for maximum readability!

<div align="center">
  <img alt="Catppuccin Theme (Default)" src="https://github.com/user-attachments/assets/600a67d6-0f41-4e17-a4de-d8ca5a5ac42f" />
  <p><em>Catppuccin Theme (Default)</em></p>
</div>

<div align="center">
  <img alt="Rosé Pine Theme" src="https://github.com/user-attachments/assets/aa45b0fa-522a-4e23-807a-5f880c6388ce" />
  <p><em>Rosé Pine Theme</em></p>
</div>

<div align="center">
  <img alt="Tokyo Night Theme" src="https://github.com/user-attachments/assets/b52a7702-9e13-49f4-93ed-df47cf542834" />
  <p><em>Tokyo Night Theme</em></p>
</div>

<div align="center">
  <img alt="Gruvbox Theme" src="https://github.com/user-attachments/assets/f667b8dd-0abb-4193-8dc2-05d8a84fb447" />
  <p><em>Gruvbox Theme</em></p>
</div>

<img width="1200" height="800" alt="themes" src="https://github.com/user-attachments/assets/35eeaea0-80de-4adc-8c40-79a238529e7c" />


To switch themes on the fly, just use the `--theme` flag:
```bash
modern-ls --theme=tokyonight
```

---

## 💡 Nerd Fonts & Emojis

The application features deep integration with both standard Unicode emojis and Nerd Font icons. Whether it's a `📦` for `package.json`, a `🎥` for video files, or a `🔑` for your `.env` secrets—`modern-ls` renders them natively right inside your terminal.

<div align="center">
  <img alt="Nerd Fonts and Emojis" src="https://github.com/user-attachments/assets/7e0657c6-fb4d-48ab-a0c7-0ff0a1078b37" />
  <p><em>Native Unicode Emojis & Nerd Fonts in action</em></p>
</div>


---

## 📦 Installation

`modern-ls` is packaged for all major operating systems. You can install it using your preferred package manager:

### macOS / Linux (Homebrew)
```bash
brew tap the-mayankjha/homebrew-tap
brew trust the-mayankjha/tap
brew install --cask modern-ls
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

### Linux (APT / DNF / Pacman)
Download the `.deb`, `.rpm`, or `.pkg.tar.zst` release from the [Releases page](https://github.com/the-mayankjha/modern-ls/releases) and run your native package manager installer (e.g., `sudo apt install ./modern-ls_linux_amd64.deb`).

---

## ⬆️ How to Upgrade

If you installed `modern-ls` using **Homebrew**, remember that Homebrew locally caches external taps. To ensure you grab the absolute latest version with all the newest themes and icons, always run `update` before `upgrade`:

```bash
brew update && brew upgrade modern-ls
```

---

## 🔒 macOS Security Guidelines

Since `modern-ls` is currently distributed via a custom Homebrew tap rather than the official App Store, macOS Gatekeeper might block it from running the first time with an "unidentified developer" warning.
<div align="center">
  <img alt="Unidentified Developer Warning" src="https://github.com/user-attachments/assets/9575886a-50d8-45c0-958e-0b493130dc9d" width="600" />
  <p><em>macOS Unidentified Developer Warning</em></p>
</div>

**To allow `modern-ls` to run:**
1. Open **System Settings** on your Mac.
2. Navigate to **Privacy & Security**.
3. Scroll down to the **Security** section.
4. You will see a message saying `modern-ls` was blocked. Click **"Allow Anyway"**.

<div align="center">
  <img alt="macOS Privacy & Security Settings" src="https://github.com/user-attachments/assets/5f930cd9-3248-401d-855b-1409434a249f" width="600" />
  <p><em>Allowing the application in Privacy & Security</em></p>
</div>

Once allowed, `modern-ls` will run smoothly without any further prompts!

---

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
- `--theme <name>`: Switch the color theme (`catppuccin`, `tokyonight`, `gruvbox`, `dracula`, `nord`, `rose-pine`)
- `--tree`: Recurse into directories as a tree
- `--depth <N>`: Limit the depth of the tree (0 means unlimited)
- `-v`, `--version`: Print version information
- `-i`, `--info`: Print author and project details
- `--upgrade`: View package-manager upgrade instructions
