package cli

import (
	"testing"
)

func initializeCommand() *CLI {
	c := New()
	g := c.New("get", "get gets", "get gets", nil)
	g.StringFlag("f", "string", "", "", false)
	g.IntFlag("g", "int", 0, "", false)
	g.BoolFlag("b", "bool", false, "", false)
	g.FloatFlag("l", "float", 0.0, "", false)
	g.StringFlag("r", "req", "", "", true)
	g.Flag("u", "", "uint64", uint64(1), "", true)
	g.Flag("ff", "", "float32", float32(1.2), "", true)
	return c
}

func TestDefaultValueToString(t *testing.T) {
	var tests = []struct {
		flagName  string
		expString string
		dvType    string
	}{
		{"f", "", "string"},
		{"g", "0", "int"},
		{"b", "false", "bool"},
		{"l", "0.000000", "float64"},
		{"r", "", "string"},
		{"u", "", "uint64"}, // there is no uint conv for default values
		{"ff", "1.200000", "float32"},
	}
	c := initializeCommand()
	flags := c.Command("get").flags
	for _, tt := range tests {
		f := flags[tt.flagName]
		out := f.defaultValueToString()
		if out != tt.expString {
			t.Errorf("Expected %s, got %s", tt.expString, out)
		}
	}
}
