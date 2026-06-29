# Migration from logo-ls

If you are coming from `logo-ls`, `modern-ls` is designed to feel right at home. 

## Key Improvements
1. **Performance**: `modern-ls` is entirely rewritten and drastically faster due to concurrent directory scanning and a complete re-architecture of the renderer.
2. **Icons**: Icons are properly synchronized with `nvim-web-devicons`, including support for Nerd Font v3 and accurate semantic folder icons.
3. **Cross Platform**: Because it relies solely on `os.FileInfo` parsing natively in Go rather than deeply-coupled Unix `syscall` intercepts, `modern-ls` natively supports Windows, macOS, and Linux out of the box.

## Missing Features
There are a few flags that `logo-ls` had that were intentionally omitted due to edge cases or performance drains:
- Very specific obscure flags may be missing if they are virtually unused. 
- Deep configuration of the core icon map in `~/.config` is removed in favor of the compile-time generator. If you need a new icon, PR it upstream to `nvim-web-devicons` or here to our custom maps, and everyone benefits!

## Flag Equivalents

Most aliases will perfectly transfer:
```bash
modern-ls -lh
modern-ls -alD
modern-ls -1
```
