// Copyright 2017 The go-icls Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"
	"reflect"
	"sort"
)

// help checks if there are flags -h or --help and return true
func help(flags Flags) bool {
	return checkForKeysInMap(flags, "h", "help")
}

// command is defined by the user and holds the all
// the information necessary to run when it is called
type command struct {
	name        string
	shortDesc   string
	description string
	// flags is a map of flag objects that
	// contain information on them
	flags map[string]*flag
	// handler is the function executed when the command is called
	handler func(flags Flags) error
}

func (c *command) Handler(h func(flags Flags) error) {
	c.handler = h
}

// Flag add a new flag in the command struct
func (c *command) Flag(name, alias, dataType string, defaultValue interface{}, description string, isRequired bool) error {
	if defaultValue != nil && reflect.TypeOf(defaultValue).String() != dataType {
		return fmt.Errorf("default value %v, is of type %s, expecting type %s", defaultValue,
			reflect.TypeOf(defaultValue).String(), dataType)
	}

	flag := &flag{
		name:         name,
		alias:        alias,
		dataType:     dataType,
		defaultValue: defaultValue,
		description:  description,
		isRequired:   isRequired,
	}

	c.flags[name] = flag
	return nil
}

// IntFlag adds an integer type value flag to command.
func (c *command) IntFlag(name, alias string, defaultValue int, description string, isRequired bool) {
	c.Flag(name, alias, "int", defaultValue, description, isRequired)
}

// FloatFlag adds a float type value flag to command.
func (c *command) FloatFlag(name, alias string, defaultValue float64, description string, isRequired bool) {
	c.Flag(name, alias, "float64", defaultValue, description, isRequired)
}

// BoolFlag adds a bool type value flag to command.
func (c *command) BoolFlag(name, alias string, description string, isRequired bool) {
	c.Flag(name, alias, "bool", false, description, isRequired)
}

// StringFlag adds an String type value flag to command.
func (c *command) StringFlag(name, alias string, defaultValue string, description string, isRequired bool) {
	c.Flag(name, alias, "string", defaultValue, description, isRequired)
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
	n := fmt.Sprintf("usage: %s [%s flags]\n\n%s\n\nFlags: \n\n", c.name, c.name, c.description)

	if _, ok := c.flags["h"]; !ok {
		c.flags["h"] = &flag{
			name:        "h",
			alias:       "help",
			dataType:    "bool",
			description: "prints out information about the command",
			isRequired:  false}
	}

	keys := make([]string, 0, len(c.flags))
	for k := range c.flags {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		f := c.flags[k]
		n += fmt.Sprint(f)
	}

	return n
}
