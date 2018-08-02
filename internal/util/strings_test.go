// Copyright 2017 The go-icls Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util_test

import (
	"testing"

	"github.com/RomanosTrechlis/go-icls/internal/util"
)

func TestTrim(t *testing.T) {
	var test = []struct {
		input  string
		output string
	}{
		{" test", "test"},
		{" test ", "test"},
		{"test ", "test"},
		{"test 1 ", "test 1"},
	}

	for _, tt := range test {
		o := util.Trim(tt.input)
		if o != tt.output {
			t.Errorf("expected '%s', got '%s'", tt.output, o)
		}
	}
}
