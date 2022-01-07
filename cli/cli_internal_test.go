package cli

import "testing"

func createCLI() *CLI {
	c := New()
	g := c.New("get", "get gets", "get gets", func(flags Flags) error {
		return nil
	})
	g.StringFlag("f", "", "", "", true)
	g.IntFlag("g", "", 1, "", false)
	return c
}

func TestValidateFlags(t *testing.T) {
	c := createCLI()
	var test = []struct {
		name        string
		description string
		flags       Flags
		ok          bool
	}{
		{"get", "f is empty string and must be not ok", Flags{"f": ""}, false},
		{"get", "f is not empty and must be ok", Flags{"f": "file"}, true},
		{"get", "there is no f and must be not ok", Flags{"g": "ggg"}, false},
		{"get", "there is no f and must be not ok", Flags{"g": ""}, false},
		{"get", "there is f and must be ok", Flags{"g": "", "f": "file"}, true},
	}
	for _, tt := range test {
		_, ok := c.validateFlags(tt.name, tt.flags)
		if ok && !tt.ok {
			t.Errorf("expected not ok, got ok for: %s", tt.description)
		}
		if !ok && tt.ok {
			t.Errorf("expected ok, got not ok for: %s", tt.description)
		}
	}
}
