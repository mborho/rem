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

	match, _ := regexp.MatchString("/home/[^/]+/"+file.filename, file.filepath)
	if !match {
		t.Errorf("Filepath not in $home: %s", file.filepath)
	}

	// test with global flag, not non-existant local file
	file.global = false
	err = file.setPath()
	if err != nil {
		t.Error("Error when setting non-exsiting global path.")
	}

	match, _ = regexp.MatchString("/home/[^/]+/"+file.filename, file.filepath)
	if !match {
		t.Errorf("Local non-existant path not pointing to home dir: %s", file.filepath)
	}

	// test with existing local file
	f, err := os.Create(".rem_test")
	defer removeTestFile(f)

	err = file.setPath()
	match, _ = regexp.MatchString("/home/.+/"+file.filename+"$", file.filepath)
	if !match {
		t.Errorf("Local non-existant path not pointing to home dir: %s", file.filepath)
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

	if _, err := os.Stat(file.filepath); os.IsNotExist(err) {
		t.Errorf("Local file was not created: %s", file.filepath)
	}
	defer removeTestFile(file.file)

}

func TestSetFile(t *testing.T) {

	file := &File{
		filename: ".rem_test_set_file",
		global:   false,
	}
	if err := file.createLocalFile(); err != nil {
		t.Errorf("Error creating local file: %s", file.filepath)
	}

	if err := file.setFile(false); err != nil {
		t.Errorf("Error opening file: %s", file.filepath)
	}

	if _, err := os.Stat(file.filepath); os.IsNotExist(err) {
		t.Errorf("Local file was not created: %s", file.filepath)
	}

	fileInfo, _ := file.file.Stat()
	mode := int(0600)
	if fileInfo.Mode() != os.FileMode(mode) {
		t.Errorf("Wrong perms: %s", fileInfo.Mode())
	}
	defer removeTestFile(file.file)

}

func TestSetFileNoAppend(t *testing.T) {

	file := &File{
		filename: ".rem_test_set_file_no_append",
		global:   false,
	}
	if err := file.createLocalFile(); err != nil {
		t.Errorf("Error creating local file: %s", file.filepath)
	}

	if err := file.setFile(true); err != nil {
		t.Errorf("Error opening file: %s", file.filepath)
	}

	if _, err := os.Stat(file.filepath); os.IsNotExist(err) {
		t.Errorf("Local file was not created: %s", file.filepath)
	}

	fileInfo, _ := file.file.Stat()
	mode := int(0600)
	if fileInfo.Mode() != os.FileMode(mode) {
		t.Errorf("Wrong perms: %s", fileInfo.Mode())
	}
	defer removeTestFile(file.file)
}
