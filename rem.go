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
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"
)

type Rem struct {
	path            string
	lines           []*Line
	hasTags         bool
	printBeforeExec bool
	File
}

func (r *Rem) appendLine(line, tag string) error {
	// Append line to the history file
	r.setFile(true)
	defer r.Close()

	if tag != "" {
		line = fmt.Sprintf("#%s#%s\n", tag, line)
	} else {
		line = fmt.Sprintf("%s\n", line)
	}
	//if _, err := r.WriteString(line); err != nil {
	if err := r.write(line); err != nil {
		panic(err)
	}
	return nil
}

func (r *Rem) executeIndex(index int) error {
	line, err := r.getLine(index)
	if err != nil {
		return err
	}
	return line.execute(r.printBeforeExec)
}

func (r *Rem) executeTag(tag string) error {
	for _, line := range r.lines {
		if line.tag == tag {
			return line.execute(r.printBeforeExec)
		}
	}
	return errors.New("Tag not found.")
}

func (r *Rem) filterLines(filter string) error {
	// Print lines filtered by string (regular expression).
	for x, line := range r.lines {
		matched, err := regexp.MatchString("(?i)"+filter, line.cmd)
		if err != nil {
			return nil
		}
		if matched {
			fmt.Printf(" %d  %s\n", x, line.cmd)
		}
	}
	return nil
}

func (r *Rem) getLine(index int) (*Line, error) {
	// Returns command by index.
	if len(r.lines) <= index {
		return nil, errors.New("Index out of range.")
	}
	return r.lines[index], nil
}

func (r *Rem) getTabWriter() *tabwriter.Writer {
	// Returns new instance of a tabwriter
	return tabwriter.NewWriter(os.Stdout, 1, 0, 2, ' ', tabwriter.DiscardEmptyColumns)
}

func (r *Rem) printAllLines() {
	// Print saved lines enumerated
	w := r.getTabWriter()

	// print out, ignore tags if no tags are present
	for x, line := range r.lines {
		line.print(w, x, r.hasTags)
	}
	w.Flush()
}

func (r *Rem) printLine(index int) error {
	// Print saved cmd by line
	line, err := r.getLine(index)
	if err != nil {
		return err
	}
	fmt.Println(line.cmd)
	return nil
}

func (r *Rem) printTag(tag string) error {
	for _, line := range r.lines {
		if line.tag == tag {
			fmt.Println(line.cmd)
			return nil
		}
	}
	return errors.New("Tag not found.")
}

func (r *Rem) read() error {
	r.setPath()
	// Read lines from the history file.
	lines := []*Line{}

	// read history
	r.setFile(false)
	defer r.Close()

	// read lines
	scanner := bufio.NewScanner(r.file)
	for scanner.Scan() {
		// parse line
		l := &Line{}
		l.read(scanner.Text())
		lines = append(lines, l)

		// tags in files?
		if l.tag != "" {
			r.hasTags = true
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	r.lines = lines
	return nil
}

func (r *Rem) removeLine(index int) error {
	// Removes a line from the rem file at given index.
	lines := []string{}
	// check line exists
	if index >= len(r.lines) {
		return errors.New("Line does not exist!")
	}
	// build new slices
	for _, line := range append(r.lines[:index], r.lines[index+1:]...) {
		lines = append(lines, line.line)
	}
	newLines := []byte{}
	if len(lines) > 0 {
		newLines = append([]byte(strings.Join(lines, "\n")), byte('\n'))
	}
	err := ioutil.WriteFile(r.filepath, newLines, 0644)
	return err
}

func (r *Rem) readFromStdIn() string {
	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		newInput := scanner.Text()
		if len(newInput) == 0 {
			break
		}
		// handle multiline-inputs
		lines = append(lines, strings.TrimSpace(strings.TrimSuffix(newInput, "\\")))
	}
	return strings.Join(lines, " ")
}
