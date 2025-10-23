package main

import (
	"os"

	"cloudamqp-cli/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
