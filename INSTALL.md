# Installation

There are several ways to install `modern-ls`.

## Prerequisites
`modern-ls` relies heavily on [Nerd Fonts](https://www.nerdfonts.com/). For the icons to render correctly in your terminal, you must use a terminal emulator configured with a Nerd Font v3 or higher.

## 1. Using Go (Recommended for developers)

If you have Go 1.22 or higher installed:

```bash
go install github.com/the-mayankjha/modern-ls/cmd/modern-ls@latest
```

Ensure your `$(go env GOPATH)/bin` directory is in your system `$PATH`.

## 2. Pre-built Binaries (Coming Soon)

Once releases are published on GitHub, you can download the tarballs directly from the [Releases](https://github.com/the-mayankjha/modern-ls/releases) page for:
- macOS (Intel & Apple Silicon)
- Linux (x86_64, ARM64)
- Windows

Just extract the `modern-ls` binary and move it to your system PATH (e.g. `/usr/local/bin`).

## 3. Build from Source

```bash
git clone https://github.com/the-mayankjha/modern-ls.git
cd modern-ls
go build -o modern-ls ./cmd/modern-ls
sudo mv modern-ls /usr/local/bin/
```

## Alias

For the best experience, you can alias `ls` to `modern-ls` in your shell configuration (`~/.zshrc` or `~/.bashrc`):

```bash
alias ls='modern-ls'
```
