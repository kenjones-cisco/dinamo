package main

import (
	"log"

	"github.com/spf13/cobra/doc"

	"github.com/kenjones-cisco/dinamo/cmd/dinamo/commands"
)

func main() {
	err := doc.GenMarkdownTree(commands.NewCommandCLI(), "./docs/usage")
	if err != nil {
		log.Fatal(err)
	}
}
