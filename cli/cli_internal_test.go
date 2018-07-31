package cli

import "testing"

func createCLI() *CLI {
	c := New()
	g := c.New("get", "get gets", "get gets", func(flags map[string]string) error {
		return nil
	})
	g.StringFlag("f", "", "", true)
	g.IntFlag("g", "", "", false)
	return c
}

func TestValidateFlags(t *testing.T) {
	c := createCLI()
	var test = []struct {
		name        string
		description string
		flags       map[string]string
		ok          bool
	}{
		{"get", "f is empty string and must be not ok", map[string]string{"f": ""}, false},
		{"get", "f is not empty and must be ok",map[string]string{"f": "file"}, true},
		{"get", "there is no f and must be not ok",map[string]string{"g": "ggg"}, false},
		{"get", "there is no f and must be not ok",map[string]string{"g": ""}, false},
		{"get", "there is f and must be ok",map[string]string{"g": "", "f": "file"}, true},
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
