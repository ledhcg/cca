// Command cca runs multiple Claude Code accounts on one machine. See
// internal/cli for the command implementations.
package main

import (
	"fmt"
	"os"

	"github.com/ledhcg/cca/internal/app"
	"github.com/ledhcg/cca/internal/cli"
)

func main() {
	a, err := app.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, "✗", err)
		os.Exit(1)
	}
	os.Exit(cli.Run(a, os.Args[1:]))
}
