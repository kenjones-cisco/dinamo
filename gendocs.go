//go:build ignore
// +build ignore

package main

import (
	"log"

	"github.com/spf13/cobra/doc"

	"github.com/kenjones-cisco/dinamo/commands"
)

func main() {
	cmd := commands.NewCommandCLI()
	cmd.DisableAutoGenTag = true
	if err := doc.GenMarkdownTree(cmd, "./docs"); err != nil {
		log.Fatal(err)
	}
}
