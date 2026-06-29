# Icon Generator

`modern-ls` does **not** rely on static hardcoded icon maps written by hand, nor does it require downloading and parsing JSON/Lua at runtime.

Instead, we use a Go code-generator pattern inside the `generator/` package.

## How it works

1. We vendor `nvim-web-devicons` (which is written in Lua) into `generator/vendor`.
2. Running `go generate ./...` in the root will invoke the `generator/generate.go` script.
3. This script parses the Lua tables using a custom AST parser.
4. It outputs an intermediate `assets/icons.json` which is useful for debugging.
5. It immediately writes `internal/icons/icons_generated.go`, formatting the icons into native Go maps.

## Updating Icons

To pull the absolute latest icons from the upstream repository:
```bash
# 1. Run the script to fetch and vendor the latest Lua files
./scripts/update-icons.sh

# 2. Run the generator to rebuild the internal go code
go generate ./...
```

You can then commit the changes to `icons_generated.go`.
