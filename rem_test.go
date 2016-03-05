package main

import (
	"io/ioutil"
	"testing"
)

var (
	remfile string
)

func init() {
	remfile = ".rem_test"
}

func removeRemFile(r *Rem) {
	removeTestFile(r.file)
}

func getTestRem(t *testing.T) *Rem {
	cmds := []byte("ls\n#foo#bla\n")
	if err := ioutil.WriteFile(remfile, cmds, 0644); err != nil {
		t.Fatalf("Cannot create rem testfile, %s", err)
	}

	rem := &Rem{
		global: false,
		File: File{
			filename: remfile,
		},
	}
	/*if err := rem.createLocalFile(); err != nil {
		t.Fatalf("Cannot create local rem testfile, %s", err)
	}*/
	return rem
}

func TestRead(t *testing.T) {
	rem := getTestRem(t)
	defer removeRemFile(rem)

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
