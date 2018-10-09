// Copyright 2017 The go-icls Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parse

import "testing"

func TestGetFlags(t *testing.T) {
	var test = []struct {
		cmd                 string
		numOfFlags          int
		numOfNonEmptyValues int
	}{
		{"get -d dir -f filename", 2, 2},
		{"get -d dir -f filename -e", 3, 2},
		{"get -d dir -f filename -e -m This is one", 4, 3},
	}

	for _, tt := range test {
		flags := getFlags(tt.cmd)
		if len(flags) != tt.numOfFlags {
			t.Errorf("expected %d number of flags, instead got %d", tt.numOfFlags, len(flags))
		}

		count := 0
		for _, v := range flags {
			if v == "" {
				continue
			}
			count++
		}
		if count != tt.numOfNonEmptyValues {
			t.Errorf("expected %d number of non empty values, instead got %d", tt.numOfNonEmptyValues, count)
		}
	}
}

func TestGetCommand(t *testing.T) {
	var test = []struct {
		line string
		cmd  string
	}{
		{"get -f filename", "get"},
		{"test", "test"},
		{"", ""},
	}

	for _, tt := range test {
		cmd := getCommand(tt.line)
		if cmd != tt.cmd {
			t.Errorf("expected '%s', got '%s'", tt.cmd, cmd)
		}
	}
}

func TestParse2(t *testing.T) {
	var test = []struct {
		line string
		cmd  *Command
	}{
		{"get -f filename",
			&Command{"get", map[string]string{"f": "filename"}, nil,},
		},
		{"test", &Command{"test", map[string]string{}, nil}},
		{"", nil},
		{"get file -f filename",
			&Command{"get", nil, &Command{"file",
			map[string]string{"f": "filename"}, nil}},
		},
		{"-help", &Command{"", map[string]string{"help": ""}, nil},},
	}

	for _, tt := range test {
		cmdTree, _ := ParseLine(tt.line)
		if cmdTree.String() != tt.cmd.String() {
			t.Errorf("expected '%s', got '%s'", tt.cmd.String(), cmdTree.String())
		}
	}
}

func TestCommand_String(t *testing.T) {
	c := &Command{
		Name: "first",
		Flags: nil,
		Child: &Command{
			Name: "second",
			Flags: map[string]string{"1": "one"},
			Child: nil,
		},
	}
	r := "{\"name\":\"first\",\"child\":{\"name\":\"second\",\"flags\":{\"1\":\"one\"}}}"
	if c.String() != r {
		t.Fatalf("failed to stringify command: expexted %s, got %s", r, c.String())
	}
}

func BenchmarkParse(b *testing.B) {
	s := "get -d dir -f filename"
	for n := 0; n < b.N; n++ {
		Parse(s)
	}
	s = "get -d dir -f filename -e"
	for n := 0; n < b.N; n++ {
		Parse(s)
	}
	s = "get -d dir -f filename -e -m This is one"
	for n := 0; n < b.N; n++ {
		Parse(s)
	}
}

func BenchmarkGetKeyValue(b *testing.B) {
	s := "-d dir "
	for n := 0; n < b.N; n++ {
		getKeyValue(s)
	}
	s = "-d dir  "
	for n := 0; n < b.N; n++ {
		getKeyValue(s)
	}
	s = " -d dir"
	for n := 0; n < b.N; n++ {
		getKeyValue(s)
	}
}
