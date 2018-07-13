package parse

import (
	"testing"
)

func TestParse(t *testing.T) {
	var test = []struct {
		cmd                 string
		cmdName				string
		numOfFlags          int
		numOfNonEmptyValues int
	}{
		{"get -d dir -f filename", "get", 2, 2},
		{"del -d dir -f filename -e", "del", 3, 2},
		{"rem -d dir -f filename -e -m This is one", "rem",4, 3},
		{"add -d dir -f filename -e -m \"This is one\"", "add",4, 3},
	}

	for _, tt := range test {
		c, _ := Parse(tt.cmd)
		if c.Command != tt.cmdName {
			t.Errorf("expected command name to be '%s', instead got '%s'", tt.cmdName, c.Command)
		}
		if len(c.Flags) != tt.numOfFlags {
			t.Errorf("expected %d number of flags, instead got %d", tt.numOfFlags, len(c.Flags))
		}

		count := 0
		for _, v := range c.Flags {
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
