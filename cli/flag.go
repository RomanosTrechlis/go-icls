// Copyright 2017 The go-icls Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"
	"bytes"
	"text/tabwriter"
)

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
	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 1, 8, 4, ' ', tabwriter.TabIndent)

	var name, alias, req = "", "", ""
	if f.name != "" {
		name = fmt.Sprintf( "-%s", f.name)
	}
	if f.alias != "" {
		alias = fmt.Sprintf("--%s", f.alias)
	}
	desc := fmt.Sprintf( "%s", f.description)
	if f.isRequired {
		req = fmt.Sprintf("(required: %v)", f.isRequired)
	}
	fmt.Fprintf(w, "\t%s\t%s\n\t\t\t\t%s %s\n", name, alias, desc, req)
	w.Flush()

	return string(buf.Bytes())
}
