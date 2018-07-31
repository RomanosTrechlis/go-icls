// Copyright 2017 The go-icls Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parse_test

import (
	"github.com/RomanosTrechlis/go-icls/parse"
	"testing"
)

func TestParse(t *testing.T) {
	var test = []struct {
		cmd                 string
		cmdName             string
		numOfFlags          int
		numOfNonEmptyValues int
	}{
		{"get -d dir -f filename", "get", 2, 2},
		{"del -d dir -f filename -e", "del", 3, 2},
		{"rem -d dir -f filename -e --m This is one", "rem", 4, 3},
		{"add -d dir -f filename -e -m \"This is one\"", "add", 4, 3},
	}

	for _, tt := range test {
		cmdName, flags := parse.Parse(tt.cmd)
		if cmdName != tt.cmdName {
			t.Errorf("expected command name to be '%s', instead got '%s'", tt.cmdName, cmdName)
		}
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
