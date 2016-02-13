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

func main() {
	globalFlag := flag.Bool("g", false, "use global rem file")
	helpFlag := flag.Bool("h", false, "show this help")
	addFlag := flag.Bool("a", false, "add a command")
	tagFlag := flag.String("t", "", "tag for command")
	filter := flag.String("f", "", "List commands by regexp filter.")
	flag.Parse()

	rem := &Rem{
		global: *globalFlag,
		file: &File{
			filename: ".rem",
		},
	}
	rem.read()

	var err error
	var index int
	remCmd := flag.Arg(0)
	switch {
	case (remCmd == "help" || *helpFlag == true):
		fmt.Println(help)
	case remCmd == "here":
		err = rem.file.createLocalFile()
	case remCmd == "clear":
		err = rem.file.clearFile()
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
