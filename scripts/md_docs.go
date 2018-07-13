package main

import (
	"log"

	"github.com/spf13/cobra/doc"

	"github.com/kenjones-cisco/dinamo/commands"
)

func main() {
	if err := doc.GenMarkdownTree(commands.NewCommandCLI(), "./docs/usage"); err != nil {
		log.Fatal(err)
	}
}
