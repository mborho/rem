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
	"path"
	"path/filepath"
	"regexp"
	"testing"
)

func removeTestFile(f *os.File) {
	f.Close()
	os.Remove(f.Name())
}

func TestSetPath(t *testing.T) {

	file := &File{
		filename: ".rem_test_file",
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
	f, err := os.Create(".rem_test_file")
	defer removeTestFile(f)

	err = file.setPath()
	match, _ = regexp.MatchString("/home/.+/"+file.filename+"$", file.filepath)
	if !match {
		t.Errorf("Local non-existant path not pointing to home dir: %s", file.filepath)
	}
}

func TestTraverseCheckPath(t *testing.T) {

	file := &File{
		filename: ".rem_test_traversed",
		global:   false,
	}

	f, err := os.Create(file.filename)
	defer removeTestFile(f)

	dir, err := os.Getwd()
	if err != nil {
		t.Errorf("Error when creating test dir: %s", dir)
	}

	testDir := path.Join(dir, "tmp")
	if err := os.Mkdir(testDir, 755); err != nil {
		t.Errorf("Error when creating test dir: %s", testDir)
	}

	if err := os.Chdir(testDir); err != nil {
		t.Errorf("Error when switching test dir: %s", testDir)
	}
	defer os.Chdir(dir)

	err = file.setPath()
	if err != nil {
		t.Error("Error when setting path.")
	}

	if file.filepath != path.Join(dir, file.filename) {
		t.Errorf("Error for traversed filepath: %s", file.filepath)

	}

	if err := os.Remove(testDir); err != nil {
		t.Error("Error when removing test dir.")
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

func TestCheckPath(t *testing.T) {

	file := &File{
		filename: ".rem_test",
		global:   false,
	}
	if err := file.createLocalFile(); err != nil {
		t.Errorf("Error creating local file: %s", file.filepath)
	}

	dirToCheck := filepath.Dir(file.filepath)
	if exists := file.checkPath(dirToCheck); exists != true {
		t.Errorf("Path not successfully checked: %s", file.filepath)
	}

	dirToCheck = filepath.Dir(dirToCheck)
	if exists := file.checkPath(dirToCheck); exists != false {
		t.Errorf("Path was successfully checked: %s", file.filepath)
	}
	defer removeTestFile(file.file)
}
