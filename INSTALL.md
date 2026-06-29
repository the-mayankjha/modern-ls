# Installation

There are several ways to install `modern-ls`.

## Prerequisites
`modern-ls` relies heavily on [Nerd Fonts](https://www.nerdfonts.com/). For the icons to render correctly in your terminal, you must use a terminal emulator configured with a Nerd Font v3 or higher.

---

## 1. Pre-built Binaries (Recommended)

Once you have pushed a tag (e.g. `v1.0.0`), GoReleaser will automatically build and publish binaries to the [GitHub Releases](https://github.com/the-mayankjha/modern-ls/releases) page.

### macOS (Apple Silicon & Intel)
```bash
# 1. Download the latest tarball (Replace 1.0.0 with the latest version number)
# For Apple Silicon (M1/M2/M3):
curl -LO https://github.com/the-mayankjha/modern-ls/releases/latest/download/modern-ls_1.0.0_darwin_arm64.tar.gz
# For Intel Macs:
# curl -LO https://github.com/the-mayankjha/modern-ls/releases/latest/download/modern-ls_1.0.0_darwin_amd64.tar.gz

# 2. Extract the archive
tar -xzf modern-ls_*_darwin_*.tar.gz

# 3. Make it executable and move to PATH
chmod +x modern-ls
sudo mv modern-ls /usr/local/bin/
```

### Linux (Ubuntu, Debian, Arch, Fedora, etc.)
```bash
# 1. Download the latest tarball (Replace 1.0.0 with the latest version number)
# For standard 64-bit x86 systems:
curl -LO https://github.com/the-mayankjha/modern-ls/releases/latest/download/modern-ls_1.0.0_linux_amd64.tar.gz
# For ARM64 (e.g. Raspberry Pi):
# curl -LO https://github.com/the-mayankjha/modern-ls/releases/latest/download/modern-ls_1.0.0_linux_arm64.tar.gz

# 2. Extract the archive
tar -xzf modern-ls_*_linux_*.tar.gz

# 3. Make it executable and move to PATH
chmod +x modern-ls
sudo mv modern-ls /usr/local/bin/
```

### Windows
1. Download `modern-ls_1.0.0_windows_amd64.zip` from the [Releases](https://github.com/the-mayankjha/modern-ls/releases) page.
2. Extract the `.zip` file.
3. Move `modern-ls.exe` to a folder of your choice (e.g., `C:\Program Files\modern-ls\`).
4. Add that folder to your System `PATH` environment variable.

---

## 2. Using Go (For Developers)

If you have Go 1.22 or higher installed:

```bash
go install github.com/the-mayankjha/modern-ls/cmd/modern-ls@latest
```

Ensure your `$(go env GOPATH)/bin` directory is in your system `$PATH`.

---

## 3. Build from Source (Manual Export)

If you want to manually build and export `modern-ls` for your system (e.g., to test local changes before releasing):

```bash
git clone https://github.com/the-mayankjha/modern-ls.git
cd modern-ls

# This will compile the binary for your CURRENT operating system and architecture
go build -o modern-ls ./cmd/modern-ls

# (Optional) Cross-compile for other systems by setting GOOS and GOARCH:
# Build for Windows:
# GOOS=windows GOARCH=amd64 go build -o modern-ls.exe ./cmd/modern-ls
# Build for Linux:
# GOOS=linux GOARCH=amd64 go build -o modern-ls-linux ./cmd/modern-ls

# Move to your path (macOS/Linux)
sudo mv modern-ls /usr/local/bin/
```

---

## Alias

For the best experience, you can alias `ls` to `modern-ls` in your shell configuration (`~/.zshrc` or `~/.bashrc`):

```bash
alias ls='modern-ls'
```
