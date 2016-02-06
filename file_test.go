// rem - A tool to remember things on the command line.
// Copyright (C) 2016 Martin Borho (martin@borho.net)
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
	_ "fmt"
	"os"
	"regexp"
	"testing"
)

func removeTestFile(f *os.File) {
	f.Close()
	os.Remove(f.Name())
}

func TestSetPath(t *testing.T) {

	file := &File{
		filename: ".rem_test",
	}
	err := file.setPath()
	if err != nil {
		t.Error("Error when setting path.")
	}

	match, _ := regexp.MatchString("/home/[^/]+/"+file.filename, file.path)
	if !match {
		t.Errorf("Filepath not in $home: %s", file.path)
	}

	// test with global flag, not non-existant local file
	file.global = false
	err = file.setPath()
	if err != nil {
		t.Error("Error when setting non-exsiting global path.")
	}

	match, _ = regexp.MatchString("/home/[^/]+/"+file.filename, file.path)
	if !match {
		t.Errorf("Local non-existant path not pointing to home dir: %s", file.path)
	}

	// test with existing local file
	f, err := os.Create(".rem_test")
	defer removeTestFile(f)

	err = file.setPath()
	match, _ = regexp.MatchString("/home/.+/"+file.filename+"$", file.path)
	if !match {
		t.Errorf("Local non-existant path not pointing to home dir: %s", file.path)
	}
}

func TestCreateLocalFile(t *testing.T) {

	file := &File{
		filename: ".rem_test_create_local_file",
		global:   false,
	}
	file.setPath()
	if err := file.createLocalFile(); err != nil {
		t.Errorf("Error when creating local file.")
	}

	if _, err := os.Stat(file.path); os.IsNotExist(err) {
		t.Errorf("Local file was not created: %s", file.path)
	}
	defer removeTestFile(file.file)
}
