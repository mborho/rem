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
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
	"golang.org/x/sys/unix"
)

type Line struct {
	line     string
	cmd      string
	tag      string
	execFlag string
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

// Edit opens the line in a text editor and returns the edited string.
func (l *Line) edit() (string, error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano" // Fallback to 'nano' as the default editor if $EDITOR is not set
	}

	file, err := ioutil.TempFile("/tmp", "tmp-*.txt")
	if err != nil {
		return "", err
	}
	defer os.Remove(file.Name())

	_, err = file.WriteString(l.cmd)
	if err != nil {
		return "", err
	}
	file.Close()

	cmd := exec.Command(editor, file.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return "", err
	}

	modifiedText, err := ioutil.ReadFile(file.Name())
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(string(modifiedText)), "\r\n", " "), "\n", " "), nil
}

func (l *Line) execute(printCmd bool) error {
	// get the pid of the calling shell
	p, err := process.NewProcess(int32(os.Getppid()))
	if err != nil {
		return err
	}

	// path of calling shell
	callerPath, err := p.Exe()
	if err != nil {
		return err
	}

	// define 'execute' flag if not set
	if l.execFlag == "" {
		l.execFlag = "-c"
	}

	// print cmd before executing
	if printCmd == true {
		fmt.Println(l.cmd)
	}

	// /bin/bash -c "ls -la"
	execParts := []string{callerPath, l.execFlag, l.cmd}

	// replace the current process
	err = unix.Exec(callerPath, execParts, os.Environ())
	if err != nil {
		return err
	}
	return nil
}

// Prints line to tabwriter.
func (l *Line) print(w io.Writer, index int, withTag bool) {
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
