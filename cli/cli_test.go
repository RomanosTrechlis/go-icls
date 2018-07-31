// Copyright 2017 The go-icls Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"testing"
)

func createCLI() *CLI {
	c := New()
	g := c.New("get", "get gets", func(flags map[string]string) error {
		return nil
	})
	g.StringFlag("f", "", "", false)
	g.IntFlag("g", "", "", false)
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
		{"get -f", false, false},
		{"asdf", false, true},
		{"get -h", false, false},
		{"get --help", false, false},
	}
	for _, tt := range test {
		b, err := c.Execute(tt.line)
		if b && !tt.quit {
			t.Errorf("expected '%v', got '%v'", tt.quit, b)
		}
		if err == nil && tt.err {
			t.Errorf("expected error, got no error")
		}
		if err != nil && !tt.err {
			t.Errorf("expected no error, got '%v'", err)
		}
	}
}
