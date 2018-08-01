// Copyright 2017 The go-icls Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli_test

import (
	"testing"

	"github.com/RomanosTrechlis/go-icls/cli"
)

func createCLI() *cli.CLI {
	c := cli.New()
	g := c.New("get", "get gets", "get gets", func(flags map[string]string) error {
		if _, ok := flags["g"]; ok {
			_, err := c.IntValue("g", "get", flags)
			if err != nil {
				return err
			}
		}
		if _, ok := flags["int"]; ok {
			_, err := c.IntValue("int", "get", flags)
			if err != nil {
				return err
			}
		}
		_ = c.StringValue("f", "get", flags)
		if _, ok := flags["f"]; ok {
			_, err := c.FlagValue("get", "f", flags)
			if err != nil {
				return err
			}
		}
		if _, ok := flags["string"]; ok {
			_, err := c.FlagValue("get", "string", flags)
			if err != nil {
				return err
			}
		}
		if _, ok := flags["b"]; ok {
			_, err := c.BoolValue("b", "get", flags)
			if err != nil {
				return err
			}
		}
		if _, ok := flags["bool"]; ok {
			_, err := c.BoolValue("bool", "get", flags)
			if err != nil {
				return err
			}
		}
		if _, ok := flags["l"]; ok {
			_, err := c.DoubleValue("l", "get", flags)
			if err != nil {
				return err
			}
		}
		if _, ok := flags["float"]; ok {
			_, err := c.DoubleValue("float", "get", flags)
			if err != nil {
				return err
			}
		}
		return nil
	})
	g.StringFlag("f", "string", "", "", false)
	g.IntFlag("g", "int", 0, "", false)
	g.BoolFlag("b", "bool", false, "", false)
	g.FloatFlag("l", "float", 0.0, "", false)
	g.StringFlag("r", "req", "", "", true)

	c.HandlerFunc("put", func(flags map[string]string) error {
		return nil
	})

	c.New("test", "test tests", "test tests", nil)
	c.HandlerFunc("test", func(flags map[string]string) error {
		return nil
	})

	c.New("nohandler", "", "", nil)
	return c
}

func TestCLI_Execute(t *testing.T) {
	c := createCLI()
	var test = []struct {
		line string
		quit bool
		err  bool
	}{
		{"quit", true, false},
		{"get -f -r t", false, false},
		{"asdf", false, true},
		{"get -h", false, false},
		{"get --help", false, false},
		{"--help", false, false},
		{"get -g fail  -r t", false, true},
		{"get -g 1  -r t", false, false},
		{"get -b -r t", false, false},
		{"get -l 1.0 -r t", false, false},
		{"get -l fail -r t", false, true},

		{"get --int 1 -r t", false, false},
		{"get --bool true -r t", false, false},
		{"get --float 1.0 -r t", false, false},
		{"get --string success -r t", false, false},
		{"get --string success", false, true},
		{"nohandler", false, true},
	}
	for _, tt := range test {
		b, err := c.Execute(tt.line)
		if b && !tt.quit {
			t.Errorf("expected '%v', got '%v'", tt.quit, b)
		}
		if err == nil && tt.err {
			t.Errorf("expected error, got no error: %s", tt.line)
		}
		if err != nil && !tt.err {
			t.Errorf("expected no error, got '%v'", err)
		}
	}
}
