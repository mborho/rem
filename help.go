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

import ()

var (
	help string
)

func init() {
	help = `NAME:
    rem - small tool for remembering things on the command line.

USAGE:
    rem [flags] [command] [argument]

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
    -p - Print command to stdout before executing index/tag.

EXAMPLES:
    rem add ls -la - Adds "ls -la" to list.
    rem -t list add ls -la - Adds "ls -la" to list with tag "list".
    rem list - Executes line tagged with "list" (ls-la)
    rem 2 - Executes line with index number 2.
    rem rm 4 - Removes line 4.
    rem - Lists all stored commands.
    `
}
