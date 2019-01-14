package main

import (
	"fmt"
	"os"

	"github.com/jit-y/ppjson"
	"github.com/mattn/go-colorable"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <file> <file>\n", os.Args[0])
		os.Exit(1)
	}

	dmp := diffmatchpatch.New()
	out := colorable.NewColorable(os.Stdout)

	file1, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	file2, err := os.Open(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	j1, _ := ppjson.NewPrinter(file1).Pretty()
	j2, _ := ppjson.NewPrinter(file2).Pretty()

	c1, c2, la := dmp.DiffLinesToChars(j1, j2)
	diffs := dmp.DiffMain(c1, c2, false)

	result := dmp.DiffCharsToLines(diffs, la)

	for _, v := range result {
		switch v.Type {
		case diffmatchpatch.DiffInsert:
			fmt.Fprint(out, "+ "+"\x1B[32m"+v.Text+"\x1B[0m")
		case diffmatchpatch.DiffDelete:
			fmt.Fprint(out, "- "+"\x1B[31m"+v.Text+"\x1B[0m")
		default:
			fmt.Fprint(out, "  "+v.Text)
		}
	}
}
