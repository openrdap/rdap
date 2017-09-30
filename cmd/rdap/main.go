package main

import (
	"os"

	"github.com/openrdap/rdap"
)

func main() {
	exitCode := rdap.RunCLI(os.Args[1:], os.Stdout, os.Stderr)

	os.Exit(exitCode)
}
