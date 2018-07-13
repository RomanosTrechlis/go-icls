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
		flags := getFlags(tt.cmd)
		if len(flags) != tt.numOfFlags {
			t.Errorf("expected %d number of flags, instead got %d", tt.numOfFlags, len(flags))
		}

		count := 0
		for _, v := range flags {
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
