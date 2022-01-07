// Copyright 2017 The go-icls Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"bytes"
	"fmt"
	"reflect"
	"text/tabwriter"
)

// flag holds information on specific flags
type flag struct {
	name         string
	alias        string
	dataType     string
	defaultValue interface{}
	description  string
	isRequired   bool
}

func (f *flag) defaultValueToString() string {
	value := f.defaultValue
	valueType := reflect.TypeOf(value).String()
	switch valueType {
	case "string":
		return fmt.Sprintf("%s", value)
	case "int":
		return fmt.Sprintf("%d", value)
	case "float32":
		return fmt.Sprintf("%f", value)
	case "float64":
		return fmt.Sprintf("%f", value)
	case "bool":
		return fmt.Sprintf("%t", value)
	default:
		return ""
	}
}

func (f *flag) String() string {
	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 1, 8, 4, ' ', tabwriter.TabIndent)

	var name, alias, req = "", "", ""
	if f.name != "" {
		name = fmt.Sprintf("-%s", f.name)
	}
	if f.alias != "" {
		alias = fmt.Sprintf("--%s", f.alias)
	}
	if f.isRequired {
		req = fmt.Sprintf("(required: %v)", f.isRequired)
	}
	fmt.Fprintf(w, "\t%s\t%s\n\t\t\t\t%s %s\n", name, alias, f.description, req)
	w.Flush()

	return buf.String()
}
