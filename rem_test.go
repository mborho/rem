package main

import (
	"io/ioutil"
	"os"
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
	cmds := []byte("ls\n#foo#ls -la\n")
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

func TestPrintLine(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rem := getTestRem(t)
	defer removeRemFile(rem)
	rem.read()

	rem.printLine(1)

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	if string(out) != "ls -la\n" {
		t.Errorf("Wrong line output, got %s", out)
	}
}

func TestPrintLineError(t *testing.T) {
	rem := getTestRem(t)
	defer removeRemFile(rem)
	rem.read()

	err := rem.printLine(2)
	if err == nil {
		t.Error("No error when calling non existant line.")
	}
}

func TestPrintAllLines(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rem := getTestRem(t)
	defer removeRemFile(rem)
	rem.read()

	rem.printAllLines()

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	if string(out) != " 0   -   ls\n 1  foo  ls -la\n" {
		t.Errorf("Wrong line output, got %s", out)
	}
}
