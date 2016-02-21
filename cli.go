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
	"flag"
	"fmt"
	"strings"
)

// Reads command line arguments and runs rem.
func run() error {
	// read flags
	globalFlag := flag.Bool("g", false, "use global rem file")
	helpFlag := flag.Bool("h", false, "show this help")
	addFlag := flag.Bool("a", false, "add a command")
	tagFlag := flag.String("t", "", "tag for command")
	filter := flag.String("f", "", "List commands by regexp filter.")
	flag.Parse()

	// build rem type
	rem := &Rem{
		global: *globalFlag,
		File: File{
			filename: ".rem",
		},
	}
	rem.read()

	// check flags and run specific method.
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
	return err
}
