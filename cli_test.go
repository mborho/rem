package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestRunDefault(t *testing.T) {
	// create test file
	rem := getTestRem(t)
	defer removeRemFile(rem)
	rem.read()

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := run(testRemFile)
	if err != nil {
		t.Error("Error when calling without argument!")
	}
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	if string(out) != " 0   -   ls\n 1  foo  ls -la\n 2   -   echo test\n" {
		t.Errorf("Wrong line output, got %s", out)
	}
}

func TestRunRmAdd(t *testing.T) {
	// create test file
	rem := getTestRem(t)
	defer removeRemFile(rem)
	rem.read()

	// remove line with index 1
	os.Args = []string{"", "rm", "1"}
	err := run(testRemFile)
	if err != nil {
		t.Error("Error when calling without argument!")
	}

	// add line
	os.Args = []string{"", "add", "pwd"}
	err = run(testRemFile)
	if err != nil {
		t.Error("Error when calling without argument!")
	}

	// list all commands
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	os.Args = []string{""}
	err = run(testRemFile)

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	if string(out) != " 0  ls\n 1  echo test\n 2  pwd\n" {
		t.Errorf("Wrong line output, got %s", out)
	}
}

/*func TestRunAddFlag(t *testing.T) {
	// create test file
	rem := getTestRem(t)
	defer removeRemFile(rem)
	rem.read()

	// add line
    tagStr := ""
	tagFlag = &tagStr
	os.Args = []string{"", "-a", "pwd"}
	err := run(testRemFile+"1")
	if err != nil {
		t.Error("Error when adding line!")
	}

	// list all commands
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	os.Args = []string{""}
	err = run(testRemFile+"1")

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
    t.Log(string(out))
	if string(out) != " 0   -   ls\n 1  foo  ls -la\n 2   -   echo test\n 3  -  pwd\n" {
		t.Errorf("Wrong line output, got %s", out)
	}
}*/

func TestRunPrintLine(t *testing.T) {
	// create test file
	rem := getTestRem(t)
	defer removeRemFile(rem)
	rem.read()

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	os.Args = []string{"", "echo", "1"}

	err := run(testRemFile)
	if err != nil {
		t.Error("Error when echoing line!")
	}
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	if string(out) != "ls -la\n" {
		t.Errorf("Wrong line output, got %s", out)
	}
}

func TestRunFilter(t *testing.T) {
	// create test file
	rem := getTestRem(t)
	defer removeRemFile(rem)
	rem.read()

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	os.Args = []string{"", "filter", "echo"}

	err := run(testRemFile)
	if err != nil {
		t.Error("Error when filter lines!")
	}
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	if string(out) != " 2  echo test\n" {
		t.Errorf("Wrong filter output, got %s", out)
	}
}

func TestRunFilterFlag(t *testing.T) {
	// create test file
	rem := getTestRem(t)
	defer removeRemFile(rem)
	rem.read()

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	os.Args = []string{"", "-f", "echo"}

	err := run(testRemFile)
	if err != nil {
		t.Error("Error when filter lines!")
	}
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	if string(out) != " 2  echo test\n" {
		t.Errorf("Wrong filter output, got %s", out)
	}
}

func TestRunHereClear(t *testing.T) {
	if _, err := os.Stat(testRemFile); os.IsNotExist(err) == false {
		t.Log("Local testRemFile already created!")
	}

	// create local rem file
	globalBool := false
	globalFlag = &globalBool

	os.Args = []string{"", "here"}
	err := run(testRemFile)
	if err != nil {
		t.Error("Error when creating local testRemFile!")
	}
	if _, err := os.Stat(testRemFile); os.IsNotExist(err) {
		t.Log("Local testRemFile was not created!")
	}

	// clear local rem file
	os.Args = []string{"", "clear"}
	err = run(testRemFile)
	if err != nil {
		t.Error("Error when removing local testRemFile!")
	}
	if _, err := os.Stat(testRemFile); os.IsNotExist(err) == false {
		t.Log("Local testRemFile was not cleared!")
	}

}

func TestRunHelp(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	helpBool := true
	helpFlag = &helpBool

	err := run(testRemFile)
	if err != nil {
		t.Error("Error when calling without argument!")
	}

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	if strings.Contains(string(out), help) == false {
		t.Error("Help not  correct!")
	}
}
