package main

import (
	"testing"
)

var (
	remfile string
)

func init() {
	remfile = ".rem_test"
}

func getTestRem(t *testing.T) *Rem {
	rem := &Rem{
		global: false,
		File: File{
			filename: remfile,
		},
	}
	if err := rem.createLocalFile(); err != nil {
		t.Fatalf("Cannot create local rem testfile, %s", err)
	}
	return rem
}

func TestRead(t *testing.T) {
	rem := getTestRem(t)
	defer removeTestFile(rem.file)

	err := rem.read()
	if err != nil {
		t.Errorf("Error reading: %s", err)
	}

	if rem.filepath == "" {
		t.Error("Filepath not set")
	}

	if rem.file == nil {
		t.Error("File not opened")
	}
}
