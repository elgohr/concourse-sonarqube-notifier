package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	os.Exit(run(os.Stdout))
}

func run(outWriter io.Writer) int {
	fmt.Fprintf(outWriter, "[]")
	return 0
}