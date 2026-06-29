// Command modern-ls is a beautiful, fast replacement for the Unix ls command.
// It displays file listings with Nerd Font icons, git status integration,
// and full color support.
//
// Usage: modern-ls [OPTIONS] [FILE...]
//
//go:generate go run ../../generator/...
package main

import (
	"log"
	"os"

	_ "github.com/the-mayankjha/modern-ls/internal/themes" // register all themes

	"github.com/the-mayankjha/modern-ls/internal/cli"
)

func main() {
	log.SetPrefix("modern-ls: ")
	log.SetFlags(0)

	code := cli.Run(os.Args[1:], os.Stdout)
	os.Exit(int(code))
}
