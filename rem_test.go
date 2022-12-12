package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

var (
	testRemFile string
)

func init() {
	testRemFile = ".rem_test"
}

func removeRemFile(r *Rem) {
	removeTestFile(r.file)
}

func getRem(t *testing.T, remStr string) *Rem {
	cmds := []byte(remStr)
	if err := ioutil.WriteFile(testRemFile, cmds, 0644); err != nil {
		t.Fatalf("Cannot create rem testfile, %s", err)
	}

	rem := &Rem{
		File: File{
			filename: testRemFile,
			global:   false,
		},
	}
	return rem
}

func getTestRem(t *testing.T) *Rem {
	remStr := "ls\n#foo#ls -la\necho test\n"
	return getRem(t, remStr)
}

func getTestEmptyRem(t *testing.T) *Rem {
	return getRem(t, "")
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

func TestPrintTag(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rem := getTestRem(t)
	defer removeRemFile(rem)
	rem.read()

	rem.printTag("foo")

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	if string(out) != "ls -la\n" {
		t.Errorf("Wrong line output, got %s", out)
	}
}

func TestPrintTagError(t *testing.T) {
	rem := getTestRem(t)
	defer removeRemFile(rem)
	rem.read()

	err := rem.printTag("4")
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

func TestFilterNoLines(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rem := getTestRem(t)
	defer removeRemFile(rem)
	rem.read()

	rem.filterLines("nothing")

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	if string(out) != "" {
		t.Errorf("Filtered lines, got %s", out)
	}
}

func TestAddLineWithTag(t *testing.T) {
	// add line
	rem := getTestRem(t)
	defer removeRemFile(rem)
	rem.read()

	if err := rem.appendLine("pwd", "test"); err != nil {
		t.Errorf("Error when appending line, got %s", err)
	}

	// check if line was added
	rem = &Rem{
		File: File{
			filename: testRemFile,
			global:   false,
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
	// add lines
	rem := getTestRem(t)
	defer removeRemFile(rem)
	rem.read()

	if err := rem.appendLine("pwd", ""); err != nil {
		t.Errorf("Error when appending line, got %s", err)
	}

	// check if line was added
	rem = &Rem{
		File: File{
			filename: testRemFile,
			global:   false,
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

func TestRemoveLine(t *testing.T) {
	// remove line
	rem := getTestRem(t)
	defer removeRemFile(rem)
	rem.read()

	if rem.lines[2].cmd != "echo test" {
		t.Error("Line does not exist")
	}

	err := rem.removeLine(3)
	if fmt.Sprintf("%s", err) != "Line does not exist!" {
		t.Error("Non existant line exists.")
	}

	if err = rem.removeLine(2); err != nil {
		t.Errorf("Error when removing line, got %s", err)
	}

	// check if line was removed
	rem = &Rem{
		File: File{
			filename: testRemFile,
			global:   false,
		},
	}
	rem.read()

	if len(rem.lines) != 2 {
		t.Error("Line was not removed")
	}
	if rem.lines[0].cmd != "ls" && rem.lines[1].cmd != "ls -la" {
		t.Error("Wrong line was not removed")
	}

	// remove all lines
	if err = rem.removeLine(0); err != nil {
		t.Errorf("Error when removing line, got %s", err)
	}

	// check if line was removed
	rem = &Rem{
		File: File{
			filename: testRemFile,
			global:   false,
		},
	}
	rem.read()

	if err = rem.removeLine(0); err != nil {
		t.Errorf("Error when removing line, got %s", err)
	}

	// check if line was removed
	rem = &Rem{
		File: File{
			filename: testRemFile,
			global:   false,
		},
	}
	rem.read()

	if len(rem.lines) > 0 {
		t.Errorf("Error when removing lines, got %s", err)
	}
}

func TestReadFromStdIn(t *testing.T) {
	rem := getTestRem(t)
	cases := []struct {
		desc     string
		input    string
		expected string
	}{
		{"TestSinglLine", "date && ls -la", "date && ls -la"},
		{
			desc: "TestMultiline",
			input: `date && \
	ls -la`,
			expected: "date && ls -la",
		},
		{
			desc: "TestMultilineNoSpace",
			input: `date  &&\
	ls -la | \
	wc`,
			expected: "date  && ls -la | wc",
		},
	}
	for _, tc := range cases {
		funcDefer, err := mockStdin(t, tc.input)
		if err != nil {
			t.Fatal(err)
		}
		defer funcDefer()

		input := rem.readFromStdIn()
		if input != tc.expected {
			t.Errorf("Wrong input was read from stdIn, want '%s', got '%s'.\n", tc.expected, input)
		}
	}
}

// thankfully used from https://gist.github.com/KEINOS/76857bc6339515d7144e00f17adb1090
func mockStdin(t *testing.T, dummyInput string) (funcDefer func(), err error) {
	t.Helper()

	oldOsStdin := os.Stdin
	tmpfile, err := ioutil.TempFile(t.TempDir(), t.Name())
	if err != nil {
		return nil, err
	}
	content := []byte(dummyInput)
	if _, err := tmpfile.Write(content); err != nil {
		return nil, err
	}
	if _, err := tmpfile.Seek(0, 0); err != nil {
		return nil, err
	}
	// Set stdin to the temp file
	os.Stdin = tmpfile

	return func() {
		// clean up
		os.Stdin = oldOsStdin
		os.Remove(tmpfile.Name())
	}, nil
}
