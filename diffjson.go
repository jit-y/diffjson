package main

import (
	"errors"
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
	if err := run(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	var (
		labels  labelNames
		unified int
	)

	flags := flag.NewFlagSet("diffjson", flag.ContinueOnError)
	flags.Var(&labels, "L", "label")
	flags.IntVar(&unified, "u", 3, "unified")
	flags.Parse(os.Args[1:])

	args := flags.Args()
	if len(args) < 2 {
		return errors.New(usage)
	}

	j1, err := prettyPrint(args[0])
	if err != nil {
		return err
	}
	j2, err := prettyPrint(args[1])
	if err != nil {
		return err
	}

	diffs := lineDiffs(j1, j2)
	formatter := newDefaultFormatter(diffs)
	out := colorable.NewColorable(os.Stdout)

	fmt.Fprint(out, formatter.diffString())

	return nil
}

func prettyPrint(pathToFile string) (string, error) {
	abs, err := filepath.Abs(pathToFile)
	if err != nil {
		return "", err
	}

	file, err := os.Open(abs)
	if err != nil {
		return "", err
	}

	return ppjson.NewPrinter(file).Pretty()
}

func lineDiffs(a, b string) []diffmatchpatch.Diff {
	dmp := diffmatchpatch.New()
	c1, c2, la := dmp.DiffLinesToChars(a, b)
	diffs := dmp.DiffMain(c1, c2, false)

	return dmp.DiffCharsToLines(diffs, la)
}
