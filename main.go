package main

import (
	"os"

	"github.com/ewan-valentine/reqstat/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

