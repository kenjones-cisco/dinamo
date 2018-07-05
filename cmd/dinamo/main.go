package main

import (
	"os"
	"runtime"

	"github.com/kenjones-cisco/dinamo/cmd/dinamo/commands"
)

// entrypoint for the CLI
func main() {
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	commands.Execute()
}
