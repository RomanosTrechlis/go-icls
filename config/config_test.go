package config_test

import (
	"github.com/RomanosTrechlis/go-icls/config"
	"testing"
)

func TestGetConfiguration(t *testing.T) {
	testPath := "testdata/"
	c, err := config.GetConfiguration(testPath)
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
		t.Errorf("expected key 'key' to have value %s, instead it has %s", "value.new", c.C["General"]["key"])
	}
}
