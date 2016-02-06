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
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"text/tabwriter"
)

var (
	help string
)

func init() {
	help = `NAME:
    rem - small tool for remembering things on the command line.

USAGE:
    rem [flags] [command] [argument]


VERSION:
    0.7.0

COMMANDS:
    -h, help - Shows this help.
    -a, add [string] - Adds a command/text.
    rm [index] - Removes line with given index number.
    echo [index] - Displays line with given index number.
    -f, filter [regexp] - Filters stored commands by given regular expression.
    here - Creates a .rem file in the given directory. Default: ~/.rem
    clear - Clears currently active .rem file, ./.rem or ~/.rem
    [index|tag] - Executes line with given index number / tag name.

    Run 'rem' without arguments to list all stored commands/strings.

FLAGS:
    -g - Use global rem file ~/.rem
    -t - Tag for command when adding with -a/add.

EXAMPLES:
    rem add ls -la - Adds "ls -la" to list.
    rem -t list add ls -la - Adds "ls -la" to list with tag "list".
    rem list - Executes line tagged with "list" (ls-la)
    rem 2 - Executes line with index number 2.
    rem rm 4 - Removes line 4.
    rem - Lists all stored commands.
    `
}

type RemFile struct {
	path    string
	lines   []*Line
	hasTags bool
	file    *File
	global  bool
}

func (r *RemFile) appendLine(line, tag string) error {
	// Append line to the history file
	r.file.setFile(true)
	defer r.close()

	if tag != "" {
		line = fmt.Sprintf("#%s#%s\n", tag, line)
	} else {
		line = fmt.Sprintf("%s\n", line)
	}
	//if _, err := r.file.WriteString(line); err != nil {
	if err := r.file.write(line); err != nil {
		panic(err)
	}
	return nil
}

func (r *RemFile) clearFile() error {
	return r.file.clearFile()
}

func (r *RemFile) createLocalFile() error {
	return r.file.createLocalFile()
}

func (r *RemFile) execute(cmdStr string) error {
	// Replace the current process with the cmd.
	cmdParts := strings.Split(cmdStr, " ")

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

func (r *RemFile) executeIndex(index int) error {
	r.read()
	cmdStr, err := r.getLine(index)
	if err != nil {
		return err
	}
	return r.execute(cmdStr)
}

func (r *RemFile) executeTag(tag string) error {
	r.read()
	for _, line := range r.lines {
		if line.tag == tag {
			return r.execute(line.cmd)
		}
	}
	return errors.New("Tag not found.")
}

func (r *RemFile) filterLines(filter string) error {
	// Print lines filtered by string (regular expression).
	r.read()
	for x, line := range r.lines {
		matched, err := regexp.MatchString("(?i)"+filter, line.cmd)
		if err != nil {
			return nil
		}
		if matched {
			fmt.Printf(" %d  %s\n", x, line)
		}
	}
	return nil
}

func (r *RemFile) getLine(index int) (string, error) {
	// Returns command by index.
	if len(r.lines) <= index {
		return "", errors.New("Index out of range.")
	}
	return r.lines[index].cmd, nil
}

func (r *RemFile) getTabWriter() *tabwriter.Writer {
	// Returns new instance of a tabwriter
	return tabwriter.NewWriter(os.Stdout, 1, 0, 2, ' ', tabwriter.DiscardEmptyColumns)
}

func (r *RemFile) printAllLines() {
	// Print saved lines enumerated
	w := r.getTabWriter()
	r.read()

	// print out, ignore tags if no tags are present
	for x, line := range r.lines {
		line.print(w, x, r.hasTags)
	}
	w.Flush()
}

func (r *RemFile) printLine(index int) error {
	// Print saved cmd by line
	r.read()
	cmd, err := r.getLine(index)
	if err != nil {
		return err
	}
	fmt.Println(cmd)
	return nil
}

func (r *RemFile) read() error {
	// Read lines from the history file.
	lines := []*Line{}

	// read history
	r.file.setFile(false)
	defer r.file.Close()
	fmt.Println(r.file)
	// read lines
	scanner := bufio.NewScanner(r.file.file)
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

func (r *RemFile) removeLine(index int) error {
	// Removes a line from the rem file at given index.
	r.read()
	lines := []string{}
	for _, line := range append(r.lines[:index], r.lines[index+1:]...) {
		lines = append(lines, line.line)
	}
	newLines := append([]byte(strings.Join(lines, "\n")), byte('\n'))
	err := ioutil.WriteFile(r.path, newLines, 0644)
	return err
}

func (r *RemFile) close() {
	r.file.Close()
}

func (r *RemFile) setPath() error {
	return r.file.setPath()
}

func toInt(str string) (int, error) {
	integer, err := strconv.Atoi(str)
	if err != nil {
		return 0, errors.New("Need index number.")
	}
	return integer, err
}

func exit(msg error) {
	fmt.Println(msg)
	os.Exit(1)
}

func main() {
	globalFlag := flag.Bool("g", false, "use global rem file")
	helpFlag := flag.Bool("h", false, "show this help")
	addFlag := flag.Bool("a", false, "add a command")
	tagFlag := flag.String("t", "", "tag for command")
	filter := flag.String("f", "", "List commands by regexp filter.")
	flag.Parse()

	rem := &RemFile{
		global: *globalFlag,
		file: &File{
			filename: ".rem",
		},
	}
	rem.setPath()
	rem.read()

	var err error
	var index int
	remCmd := flag.Arg(0)
	switch {
	case (remCmd == "help" || *helpFlag == true):
		fmt.Println(help)
	case remCmd == "here":
		err = rem.createLocalFile()
	case remCmd == "clear":
		err = rem.clearFile()
	case (remCmd == "add" || *addFlag == true):
		startIndex := 1
		if *addFlag == true {
			startIndex = 0
		}
		err = rem.appendLine(strings.Join(flag.Args()[startIndex:], " "), *tagFlag)
	case (remCmd == "filter"):
		err = rem.filterLines(strings.Join(flag.Args()[1:], " "))
	case *filter != "":
		err = rem.filterLines(*filter)
	case remCmd == "rm":
		if index, err = toInt(flag.Arg(1)); err == nil {
			err = rem.removeLine(index)
		}
	case remCmd == "echo":
		if index, err = toInt(flag.Arg(1)); err == nil {
			err = rem.printLine(index)
		}
	case remCmd != "":
		if index, err = toInt(remCmd); err == nil {
			err = rem.executeIndex(index)
		} else {
			err = rem.executeTag(remCmd)
		}
	default:
		// if there is a not known rem-cmd it can be assumed it is a tag
		rem.printAllLines()
	}
	if err != nil {
		exit(err)
	}
}
