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
	cmds := []byte("ls\n#foo#ls -la\necho test\n")
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

	err := rem.printLine(4)
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

	if string(out) != " 0   -   ls\n 1  foo  ls -la\n 2   -   echo test\n" {
		t.Errorf("Wrong line output, got %s", out)
	}
}

func TestFilterLines(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rem := getTestRem(t)
	defer removeRemFile(rem)
	rem.read()

	rem.filterLines("ls")

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	if string(out) != " 0  ls\n 1  ls -la\n" {
		t.Errorf("Wrong line output, got %s", out)
	}
}

func TestAddLineWithTag(t *testing.T) {
	rem := getTestRem(t)
	defer removeRemFile(rem)
	rem.read()

	if err := rem.appendLine("pwd", "test"); err != nil {
		t.Errorf("Error when appending line, got %s", err)
	}

	rem = &Rem{
		global: false,
		File: File{
			filename: remfile,
		},
	}
	rem.read()
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rem.printLine(3)

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	if string(out) != "pwd\n" {
		t.Errorf("Wrong line output, got %s", out)
	}

	if rem.lines[3].tag != "test" {
		t.Errorf("Wrong tag saved, got %s", rem.lines[3].tag)
	}
}

func TestAddLineWithoutTag(t *testing.T) {
	rem := getTestRem(t)
	defer removeRemFile(rem)
	rem.read()

	if err := rem.appendLine("pwd", ""); err != nil {
		t.Errorf("Error when appending line, got %s", err)
	}

	rem = &Rem{
		global: false,
		File: File{
			filename: remfile,
		},
	}
	rem.read()
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rem.printLine(3)

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	if string(out) != "pwd\n" {
		t.Errorf("Wrong line output, got %s", out)
	}

	if rem.lines[3].tag != "" {
		t.Errorf("Tag was saved, got %s", rem.lines[3].tag)
	}
}
