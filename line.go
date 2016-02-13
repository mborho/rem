// rem - A tool to remember things on the command line.
// Copyright (C) 2015 Martin Borho (martin@borho.net)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
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

func (l *Line) execute() error {
	// Replace the current process with the cmd.
	cmdParts := strings.Split(l.cmd, " ")

	// absolute path to cmd
	execPath, err := exec.LookPath(cmdParts[0])
	if err != nil {
		return err
	}

	// replace the current process
	env := os.Environ()
	syscall.Exec(execPath, cmdParts, env)
	return nil
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
