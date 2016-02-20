package main

import (
	"bytes"
	"strings"
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

func TestExecute(t *testing.T) {
	l := &Line{
		cmd: "/not/abdcdef",
	}

	err := l.execute()
	if err == nil {
		t.Errorf("Command was executed without error: %s", l.cmd)
	}

}

func TestPrint(t *testing.T) {
	// test without tag
	l := &Line{
		cmd: "echo foobar",
	}
	var b bytes.Buffer
	l.print(&b, 2, false)

	if strings.Compare(b.String(), " 2\techo foobar\n") != 0 {
		t.Error("Line was printed incorrect.")
	}

	// test with tag but without line tag
	l = &Line{
		cmd: "foo",
	}
	var d bytes.Buffer
	l.print(&d, 3, true)
	if strings.Compare(d.String(), " 3\t - \tfoo\n") != 0 {
		t.Error("Line was printed incorrect.")
	}

	// test with tag
	l = &Line{
		tag: "foo",
		cmd: "bar",
	}
	var c bytes.Buffer
	l.print(&c, 4, true)
	if strings.Compare(c.String(), " 4\tfoo\tbar\n") != 0 {
		t.Error("Line with tag was printed incorrect.")
	}
}
