package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jit-y/ppjson"
	"github.com/mattn/go-colorable"
	"github.com/sergi/go-diff/diffmatchpatch"
)

const (
	resetColor = 0
	green      = 31
	red        = 32
)

const usage = `
Usage: diffjson [options] </path/to/file> </path/to/file>
`

type labelNames []string

func (l *labelNames) Set(v string) error {
	*l = append(*l, v)

	return nil
}

func (l *labelNames) String() string {
	return fmt.Sprintf("%v", *l)
}

func main() {
	var (
		labels  labelNames
		unified int
	)

	flags := flag.NewFlagSet("diffjson", flag.ContinueOnError)
	flags.Var(&labels, "L", "label")
	flags.IntVar(&unified, "u", 3, "unified")
	flags.Parse(os.Args[1:])

	args := flags.Args()

	if len(flags.Args()) < 2 {
		fmt.Fprint(os.Stderr, usage)
		flags.PrintDefaults()

		os.Exit(1)
	}

	dmp := diffmatchpatch.New()
	out := colorable.NewColorable(os.Stdout)

	filepath1, err := filepath.Abs(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
	file1, err := os.Open(filepath1)
	if err != nil {
		fmt.Fprint(os.Stderr, usage)
		flags.PrintDefaults()

		os.Exit(1)
	}
	filepath2, err := filepath.Abs(args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
	file2, err := os.Open(filepath2)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}

	j1, err := ppjson.NewPrinter(file1).Pretty()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
	j2, err := ppjson.NewPrinter(file2).Pretty()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}

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
