// Copyright 2017 The go-icls Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"reflect"
	"testing"
)

func TestGetDataTypeFunction(t *testing.T) {
	var test = []struct {
		dt    string
		value string
		out   string
		err   bool
	}{
		{"string", "filename", "string", false},
		{"int", "25", "int64", false},
		{"int", "asd", "int64", true},
		{"float", "55.5", "float64", false},
		{"bool", "true", "bool", false},
	}

	for _, tt := range test {
		f := getDataTypeFunction(tt.dt)
		v, err := f(tt.value)
		if tt.err && err == nil {
			t.Errorf("expected error, got %v", v)
			continue
		}

		if err != nil {
			continue
		}

		if reflect.TypeOf(v).String() != tt.out {
			t.Errorf("expected '%s', got '%s'", tt.out, reflect.TypeOf(v).String())
		}
	}
}

func TestConv(t *testing.T) {
	var test = []struct {
		f     func(s string) (interface{}, error)
		value string
		out   string
		val   interface{}
	}{
		{getDataTypeFunction("string"), "filename", "string", "filename"},
		{getDataTypeFunction("int"), "25", "int64", int64(25)},
		{getDataTypeFunction("int"), "asd", "int64", int64(0)},
		{getDataTypeFunction("float"), "55.5", "float64", 55.5},
		{getDataTypeFunction("bool"), "true", "bool", true},
	}

	for _, tt := range test {
		i, _ := conv(tt.value, tt.f)
		if reflect.TypeOf(i).String() != tt.out {
			t.Errorf("expected '%s', got '%s'", tt.out, reflect.TypeOf(i).String())
		}

		if reflect.ValueOf(i).String() != reflect.ValueOf(tt.val).String() {
			t.Errorf("expected '%v', got '%v'", tt.val, i)
		}
	}
}
