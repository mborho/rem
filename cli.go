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

var (
	globalFlag *bool
	helpFlag   *bool
	addFlag    *bool
	tagFlag    *string
	printFlag  *bool
	filter     *string
)

func init() {
	// read flags
	globalFlag = flag.Bool("g", false, "use global rem file")
	helpFlag = flag.Bool("h", false, "show this help")
	addFlag = flag.Bool("a", false, "add a command")
	tagFlag = flag.String("t", "", "tag for command")
	printFlag = flag.Bool("p", false, "print command before executing")
	filter = flag.String("f", "", "List commands by regexp filter.")
}

// Reads command line arguments and runs rem.
func run(remfile string) error {
	flag.Parse()

	// build rem type
	rem := &Rem{
		File: File{
			filename: remfile,
			global:   *globalFlag,
		},
		printBeforeExec: *printFlag,
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
		toAdd := strings.Join(flag.Args()[startIndex:], " ")
		if toAdd == "" {
			// read line from stdIn
			toAdd = rem.readFromStdIn()
		}
		err = rem.appendLine(toAdd, *tagFlag)
	case (remCmd == "filter"):
		err = rem.filterLines(strings.Join(flag.Args()[1:], " "))
	case *filter != "":
		err = rem.filterLines(*filter)
	case remCmd == "rm":
		if index, err = toInt(flag.Arg(1)); err == nil {
			err = rem.removeLine(index)
		}
	case remCmd == "echo":
		if flag.Arg(1) != "" {
			if index, err = toInt(flag.Arg(1)); err == nil {
				err = rem.printLine(index)
			} else {
				err = rem.printTag(flag.Arg(1))
			}
		}
	case remCmd != "":
		if index, err = toInt(remCmd); err == nil {
			err = rem.executeIndex(index)
		} else {
			err = rem.executeTag(remCmd)
		}
	default:
		rem.printAllLines()
		if len(rem.lines) == 0 {
			// show help if nothing was found
			fmt.Println(help)
		}
	}
	return err
}
