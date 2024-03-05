package main

import (
	"fmt" //nolint:typecheck
	"io" //nolint:typecheck
	"os"
	"path/filepath" //nolint:typecheck
	"strings" //nolint:typecheck
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles) //nolint:typecheck
	if err != nil {
		panic(err.Error())
	}
}
