package main

import (
	"fmt"
	"os"

	"github.com/jit-y/ppjson"
	"github.com/mattn/go-colorable"
	"github.com/sergi/go-diff/diffmatchpatch"
)

const (
	resetColor = 0
	green      = 31
	red        = 32
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
			fmt.Fprintf(out, "+ \x1b[%dm%s\x1b[%dm", red, v.Text, resetColor)
		case diffmatchpatch.DiffDelete:
			fmt.Fprintf(out, "- \x1b[%dm%s\x1b[%dm", green, v.Text, resetColor)
		default:
			fmt.Fprintf(out, "  %s", v.Text)
		}
	}
}
