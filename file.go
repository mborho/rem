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
    "os"
)



type File struct {
	path     string
	filename string
	file     *os.File
	global   bool
}

func (f *File) clearFile() error {
	return os.Remove(f.path)
}

func (f *File) setFile(appendTo booö) error {
	// which mode to use to open file
	var openFlags int
	if appendTo {
		openFlags = os.O_CREATE | os.O_APPEND | os.O_WRONLY
	} else {
		openFlags = os.O_CREATE | os.O_RDONLY
	}

	// open history file
	file, err := os.OpenFile(f.path, openFlags, 0600)
	if err == nil {
		f.file = file
	}
	return nil
}

func (f *File) createLocalFile() error {
	// Create history file in the current directory.
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	localFile := path.Join(dir, f.filename)
	f.file, err := os.OpenFile(localFile, os.O_CREATE, 0600)
	defer f.file.Close()
	if err != nil {
		return err
	}
	return nil
}

func (f *File) setPath() error {
	// ignore current dir if global .rem file is wanted
	if f.global == false {
		// Set path to history file in current dir if one exists
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		localFile := path.Join(dir, f.filename)
		if _, err := os.Stat(localFile); err == nil {
			f.path = localFile
			return nil
		}
	}

	// Set default path to rem's history file
	usr, err := user.Current()
	if err == nil {
		f.path = path.Join(usr.HomeDir, f.filename)
	}
	return err
}

