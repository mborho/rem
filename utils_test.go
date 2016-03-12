package main

import (
	_ "errors"
	_ "os"
	_ "os/exec"
	"testing"
)

func TestToInt(t *testing.T) {
	// test with number string
	integer, err := toInt("5")
	if err != nil {
		t.Error("Error converting number string in integer")
	}

	if integer != 5 {
		t.Error("Number string convereted to false integer")
	}

	// test with non-number string
	integer2, err := toInt("foobar")
	if err == nil {
		t.Error("No error when converting string which isn't also an integer")
	}

	if integer2 != 0 {
		t.Error("String falsely converted to integer.")
	}

}

/*func TestExit(t *testing.T) {
	exitError := errors.New("Test-Error")
	if os.Getenv("BE_CRASHER") == "1" {
		exit(exitError)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestExit")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("Process ran with err %v, want exit status 1", err)
}*/
