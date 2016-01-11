package main

import (
	"testing"
)

func TestParseLine(t *testing.T) {
	line := "#foo#blafasel blubb"
	l := &Line{}
	l.read(line)

	if l.tag != "foo" {
		t.Errorf("Can't detect tag for %s", line)
	}

	if l.cmd != "blafasel blubb" {
		t.Errorf("Can't detect tag for %s", line)
	}

	l = &Line{}
	line = "ls -la"
	l.read(line)
	if l.tag != "" {
		t.Errorf("Tag was found in %s", line)
	}

	if l.cmd != "ls -la" {
		t.Errorf("Command was not found in %s", line)
	}

}
