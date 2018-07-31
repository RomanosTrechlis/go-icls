// Copyright 2017 The go-icls Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"
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
	if c.flags == nil {
		c.flags = make([]*flag, 0)
	}

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
	n := fmt.Sprintf("%s\n", c.name)
	hasHelpFlag := false
	for _, f := range c.flags {
		n += fmt.Sprint(f)
		if f.name == "h" || f.alias == "help" {
			hasHelpFlag = true
		}
	}
	if !hasHelpFlag {
		n += fmt.Sprintf("\t-%s\t--%s\t\t%s", "h", "help", "prints out information about the command")
	}
	n += fmt.Sprintf("\n%s\n", c.description)
	return n
}
