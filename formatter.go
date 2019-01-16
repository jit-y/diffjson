package main

import (
	"bytes"
	"fmt"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type formatter interface {
	diffString()
}

type defaultFormatter struct {
	diffs []diffmatchpatch.Diff
	buf   *bytes.Buffer
}

func newDefaultFormatter(diffs []diffmatchpatch.Diff) *defaultFormatter {
	var buf bytes.Buffer
	return &defaultFormatter{
		diffs: diffs,
		buf:   &buf,
	}
}

func (d *defaultFormatter) diffString() string {
	for _, v := range d.diffs {
		switch v.Type {
		case diffmatchpatch.DiffInsert:
			fmt.Fprintf(d.buf, "+ \x1b[%dm%s\x1b[%dm", red, v.Text, resetColor)
		case diffmatchpatch.DiffDelete:
			fmt.Fprintf(d.buf, "- \x1b[%dm%s\x1b[%dm", green, v.Text, resetColor)
		}
	}

	return d.buf.String()
}
