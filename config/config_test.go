// Copyright 2017 The go-icls Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config_test

import (
	"os"
	"testing"

	"github.com/RomanosTrechlis/go-icls/config"
)

func TestGetConfigurationFromSingleFile(t *testing.T) {
	testPath := "testdata/"
	c, err := config.GetConfigurationFromSingleFile(testPath + "test.properties")
	if err != nil {
		t.Fatalf("failed to read properties file correctly: %v", err)
	}

	if len(c.C) != 2 {
		t.Errorf("expected to have 2 sections, instead it has %d", len(c.C))
	}
	if len(c.C["General"]) != 3 {
		t.Errorf("expected 'General' section to have 3 key-value pairs, instead it has %d", len(c.C["General"]))
	}
	if len(c.C["Work"]) != 3 {
		t.Errorf("expected 'Work' section to have 3 key-value pairs, instead it has %d", len(c.C["General"]))
	}
	if c.C["General"]["key"] == "value" {
		t.Errorf("expected key 'key' to have value '%s', instead it has '%s'", "value.new", c.C["General"]["key"])
	}

	_, err = config.GetConfigurationFromSingleFile(testPath + "nonexistentfile.properties")
	if err == nil {
		t.Fatalf("expected to fail with error, but didn't: %s", testPath+"nonexistentfile.properties")
	}

	_, err = config.GetConfigurationFromSingleFile(testPath + "test_no_file.properties")
	if err == nil {
		t.Error("expecting error, got no error")
	}
}

func TestGetConfigurationFromDir(t *testing.T) {
	testPath := "testdata/"
	c, err := config.GetConfigurationFromDir(testPath)
	if err != nil {
		t.Fatalf("failed to read properties file correctly: %v", err)
	}

	if len(c.C) != 3 {
		t.Errorf("expected to have 3 sections, instead it has %d", len(c.C))
	}
	if len(c.C["General"]) != 4 {
		t.Errorf("expected 'General' section to have 4 key-value pairs, instead it has %d", len(c.C["General"]))
	}
	if len(c.C["Work"]) != 3 {
		t.Errorf("expected 'Work' section to have 3 key-value pairs, instead it has %d", len(c.C["General"]))
	}
	if c.C["General"]["key"] == "value" {
		t.Errorf("expected key 'key' to have value '%s', instead it has '%s'", "value.new", c.C["General"]["key"])
	}
	if c.C["Specs"]["v"] != "100m/s" {
		t.Errorf("expected key 'v' to have value '%s', instead it has '%s'", "100m/s", c.C["Specs"]["v"])
	}

	testPath = "testdata/nopath/"
	_, err = config.GetConfigurationFromDir(testPath)
	if err == nil {
		t.Error("expecting error, got no error")
	}

	// clean-up 'cause GetConfigurationFromDir creates given folder
	os.Remove(testPath)
}

func BenchmarkGetConfigurationFromSingleFile(b *testing.B) {
	testPath := "testdata/"
	for n := 0; n < b.N; n++ {
		config.GetConfigurationFromSingleFile(testPath + "test.properties")
	}
}

func BenchmarkGetConfigurationFromDir(b *testing.B) {
	testPath := "testdata/"
	for n := 0; n < b.N; n++ {
		config.GetConfigurationFromDir(testPath)
	}
}
