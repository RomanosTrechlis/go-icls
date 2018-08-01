// Copyright 2017 The go-icls Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"
	"sort"
)

// help checks if there are flags -h or --help and return true
func help(flags map[string]string) bool {
	if checkForKeysInMap(flags, "h", "help") {
		return true
	}
	return false
}

// command is defined by the user and holds the all
// the information necessary to run when it is called
type command struct {
	name        string
	shortDesc   string
	description string
	// flags is an array of flag objects that
	// contain information on its of them
	flags []*flag
	// handler is the function executed when the command is called
	handler func(flags map[string]string) error
}

// Flag add a new flag in the command struct
func (c *command) Flag(name, alias, dataType, description string, isRequired bool) {
	flag := &flag{
		name:        name,
		alias:       alias,
		dataType:    dataType,
		description: description,
		isRequired:  isRequired,
	}

	c.flags = append(c.flags, flag)
}

// IntFlag adds an integer type value flag to command.
func (c *command) IntFlag(name, alias, description string, isRequired bool) {
	c.Flag(name, alias, "int", description, isRequired)
}

// FloatFlag adds a float type value flag to command.
func (c *command) FloatFlag(name, alias, description string, isRequired bool) {
	c.Flag(name, alias, "float", description, isRequired)
}

// BoolFlag adds a bool type value flag to command.
func (c *command) BoolFlag(name, alias, description string, isRequired bool) {
	c.Flag(name, alias, "bool", description, isRequired)
}

// StringFlag adds an String type value flag to command.
func (c *command) StringFlag(name, alias, description string, isRequired bool) {
	c.Flag(name, alias, "string", description, isRequired)
}

func (c *command) getFlag(name string) *flag {
	for _, f := range c.flags {
		if f.name == name || f.alias == name {
			return f
		}
	}
	return nil
}

func (c *command) String() string {
	n := fmt.Sprintf("usage: %s [%s flags]\n\n", c.name, c.name)
	n += fmt.Sprintf("%s\n", c.description)
	n += fmt.Sprintf("\nFlags: \n\n")

	flagMap := createFlagMapFromArray(c.flags)
	if _, ok := flagMap["h"]; !ok {
		flagMap["h"] = &flag{
			name: "h",
			alias: "help",
			dataType: "bool",
			description: "prints out information about the command",
			isRequired: false}
	}

	keys := make([]string, 0, len(flagMap))
	for k := range flagMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		f := flagMap[k]
		n += fmt.Sprint(f)
	}

	return n
}

func createFlagMapFromArray(ff []*flag) map[string]*flag {
	flagMap := make(map[string]*flag, len(ff))
	for _, f := range ff {
		flagMap[f.name] = f
	}
	return flagMap
}
