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
	"os/user"
	"path"
	"strconv"
	"strings"
	"syscall"
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
    0.2.0

COMMANDS:
    -h, help - Shows this help.
    -a, add [string] - Adds a command/text.
    rm [index] - Removes line with given index number.
    echo [index] - Displays line with given index number.
    here - Creates a .rem file in the given directory. Default: ~/.rem
    clear - Clears currently active .rem file, ./.rem or ~/.rem
    [index] - Executes line with given index number.

    Run 'rem' without arguments to list all stored commands/strings.

FLAGS:
    -g - Use global rem file ~/.rem

EXAMPLES:
    rem 2 - Executes line with index number 2.
    rem add ls -la - Adds "ls -la" to list.
    rem rm 4 - Removes line 4.
    rem - Lists all stored commands.
    `
}

type RemFile struct {
	path     string
	lines    []string
	appendTo bool
	filename string
	file     *os.File
	global   bool
}

func (r *RemFile) appendLine(line string) error {
	// Append line to the history file
	r.appendTo = true
	r.setFile()
	defer r.file.Close()

	line = fmt.Sprintf("%s\n", line)
	if _, err := r.file.WriteString(line); err != nil {
		panic(err)
	}
	return nil
}

func (r *RemFile) clearFile() error {
	return os.Remove(r.path)
}

func (r *RemFile) createLocalFile() error {
	// Create history file in the current directory.
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	localFile := path.Join(dir, r.filename)
	file, err := os.OpenFile(localFile, os.O_CREATE, 0600)
	defer file.Close()
	if err != nil {
		return err
	}
	return nil
}

func (r *RemFile) execute(index int) error {
	r.read()
	cmdStr, err := r.getLine(index)
	if err != nil {
		return err
	}
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

func (r *RemFile) getLine(index int) (string, error) {
	// Returns command by index.
	if len(r.lines) <= index {
		return "", errors.New("index out of range!")
	}
	return r.lines[index], nil
}

func (r *RemFile) printAllLines() {
	// Print saved lines enumerated
	r.read()
	for x, line := range r.lines {
		fmt.Printf(" %d  %s\n", x, line)
	}
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
	lines := []string{}

	// read history
	r.setFile()
	defer r.file.Close()

	// read lines
	scanner := bufio.NewScanner(r.file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
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
	lines := append(r.lines[:index], r.lines[index+1:]...)
	newLines := append([]byte(strings.Join(lines, "\n")), byte('\n'))
	err := ioutil.WriteFile(r.path, newLines, 0644)
	return err
}

func (r *RemFile) setFile() error {
	// which mode to use to open file
	var openFlags int
	if r.appendTo {
		openFlags = os.O_CREATE | os.O_APPEND | os.O_WRONLY
	} else {
		openFlags = os.O_CREATE | os.O_RDONLY
	}

	// open history file
	file, err := os.OpenFile(r.path, openFlags, 0600)
	if err == nil {
		r.file = file
	}
	return nil
}

func (r *RemFile) setPath() error {
	// ignore current dir if global .rem file is wanted
	if r.global == false {
		// Set path to history file in current dir if one exists
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		localFile := path.Join(dir, r.filename)
		if _, err := os.Stat(localFile); err == nil {
			r.path = localFile
			return nil
		}
	}

	// Set default path to rem's history file
	usr, err := user.Current()
	if err == nil {
		r.path = path.Join(usr.HomeDir, r.filename)
	}
	return err
}

func toInt(str string) int {
	integer, err := strconv.Atoi(str)
	if err != nil {
		return -1
	}
	return integer
}

func exit(msg error) {
	fmt.Println(msg)
	os.Exit(1)
}

func main() {
	globalFlag := flag.Bool("g", false, "use global rem file")
	helpFlag := flag.Bool("h", false, "show this help")
	addFlag := flag.Bool("a", false, "add a command")
	flag.Parse()

	rem := &RemFile{
		filename: ".rem",
		global:   *globalFlag,
	}
	rem.setPath()

	var err error
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
		err = rem.appendLine(strings.Join(flag.Args()[startIndex:], " "))
	case remCmd == "rm":
		err = rem.removeLine(toInt(flag.Arg(1)))
	case remCmd == "echo":
		err = rem.printLine(toInt(flag.Arg(1)))
	case remCmd != "" && toInt(remCmd) > -1:
		err = rem.execute(toInt(remCmd))
	default:
		rem.printAllLines()
	}
	if err != nil {
		exit(err)
	}
}
