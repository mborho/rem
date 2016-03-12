package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestRunDefault(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := run(remfile)
	if err != nil {
		t.Error("Error when calling without argument!")
	}

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	t.Log(out)
}
