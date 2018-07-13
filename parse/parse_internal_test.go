package parse

import "testing"

func TestGetFlags(t *testing.T) {
	var test = []struct {
		cmd                 string
		numOfFlags          int
		numOfNonEmptyValues int
	}{
		{"get -d dir -f filename", 2, 2},
		{"get -d dir -f filename -e", 3, 2},
		{"get -d dir -f filename -e -m This is one", 4, 3},
	}

	for _, tt := range test {
		c := new(Command)
		c.Flags = make(map[string]string)
		c.getFlags(tt.cmd)
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
