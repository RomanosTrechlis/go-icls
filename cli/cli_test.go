package cli

import (
	"testing"
)

func createCLI() *CLI {
	c := New()
	g := c.New("get", "get gets", func(cmd string, flags map[string]string) error {
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
