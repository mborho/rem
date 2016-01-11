package main

import (
	"fmt"
	"regexp"
	"text/tabwriter"
)

type Line struct {
	line string
	cmd  string
	tag  string
}

// Read incoming string into Line struct.
func (l *Line) read(line string) {
	re := regexp.MustCompile("^#([^ ]+)?#")
	l.line = line
	if tagMatch := re.FindSubmatch([]byte(line)); tagMatch != nil {
		l.tag = string(tagMatch[1])
		l.cmd = line[len(l.tag)+2:]
	} else {
		// no tag found, simple command
		l.cmd = line
	}
}

// Prints line to tabwriter.
func (l *Line) print(w *tabwriter.Writer, index int, withTag bool) {
	if withTag {
		tag := ""
		if tag = l.tag; tag == "" {
			tag = " - "
		}
		fmt.Fprintf(w, " %d\t%s\t%s\n", index, tag, l.cmd)
	} else {
		fmt.Fprintf(w, " %d\t%s\n", index, l.cmd)
	}
}
