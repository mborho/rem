package main

import (
	"regexp"
)

type Line struct {
	line string
	cmd  string
	tag  string
}

func (l *Line) read(line string) {
	re := regexp.MustCompile("^#([^ ]+)?#")
	l.line = line
	if tagMatch := re.FindSubmatch([]byte(line)); tagMatch != nil {
		l.tag = string(tagMatch[1])
		l.cmd = line[len(l.tag)+2:]
	} else {
		// not tag found, simple command
		l.cmd = line
	}
}
