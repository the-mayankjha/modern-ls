# Contributing to modern-ls

Thank you for your interest in contributing to `modern-ls`! We aim to build a robust, beautiful CLI tool and welcome your help.

## Development Setup

1. Make sure you have Go 1.22+ installed.
2. Clone the repository and run `go mod download`.
3. To build the CLI locally:
```bash
go build -o modern-ls ./cmd/modern-ls
```

## Adding Icons

We pull most file extensions and filename icons automatically from `nvim-web-devicons`.
If an icon is missing:
1. Contribute it to `nvim-web-devicons` repository directly.
2. We periodically run `./scripts/update-icons.sh` and `go generate ./...` to pull their latest dictionaries into our codebase.

If you want to add a *semantic folder* (like `src/`, `admin/`, etc.), you can edit `generator/generate.go` in the `buildFolderIcons` function and send a PR.

## Submitting Pull Requests

1. Fork the repo and create your branch from `main`.
2. Ensure you run `go fmt ./...` and `go test ./...` before submitting.
3. Write a clear, concise PR description.

## Bug Reports & Feature Requests
Use GitHub Issues to report bugs or request features. Please include your OS, Terminal Emulator, and standard `ls` equivalent command if reporting a bug.
