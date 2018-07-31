// Copyright 2017 The go-icls Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import "fmt"

// flag holds information on specific flags
type flag struct {
	name string
	// for now allies does nothing
	alias string
	// for now dataType does nothing
	dataType    string
	description string
	isRequired  bool
}

func (f *flag) String() string {
	n := "\t"
	if f.name != "" {
		n += fmt.Sprintf("-%s\t", f.name)
	} else {
		n += fmt.Sprintf("\t")
	}
	if f.alias != "" {
		n += fmt.Sprintf("--%s", f.alias)
	} else {
		n += fmt.Sprintf("\t")
	}
	n += fmt.Sprintf("\t%s (required: %v)\n", f.description, f.isRequired)
	return n
}
