// Copyright 2017 The go-icls Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parse

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/RomanosTrechlis/go-icls/internal/util"
)

type Command struct {
	Name string `json:"name,omitempty"`
	Flags map[string]string `json:"flags,omitempty"`
	Child *Command `json:"child,omitempty"`
}

func (c *Command) String() string {
	b, _ := json.Marshal(&c)
	return fmt.Sprintf("%s", string(b))
}

// Parse parses the command given and return the command and the flags
//
// Example commands:
//
// command with key-value flags
// get -d dir -f filename
//
// command with key-value and key only flags
// del -d dir -f filename -e
//
// command with flag value of several words
// rem -d dir -f filename -e -m This is one
//
// command with flag value in quotation marks
// add -d dir -f filename -e -m "This is one"
//
// command with inner command
// add file -f filename
// add repo -r reponame
//
// Main assumption is that the only the last command has flags
//
// add file -d directory
//
// It has two commands. Command add that can't have any flags and command file that
// has a flag "d". Command "file" is a child of command "add"
func Parse(cmd string) (string, map[string]string) {
	cmdName := getCommand(cmd)
	flags := getFlags(cmd)
	return cmdName, flags
}

func ParseLine(line string) (*Command, error) {
	if line == "" {
		return nil, errors.New("line is empty, can't create command tree")
	}

	ss := strings.Split(line, " ")

	startCmd := strings.Trim(ss[0], " ")
	if strings.HasPrefix(startCmd, "-") {
		return &Command{"", getFlags(strings.Join(ss, " ")), nil}, nil
	}

	c := &Command{Name: startCmd}
	current := c
	for i := 1; i < len(ss); i++ {
		cmd := current
		if strings.HasPrefix(ss[i], "-") {
			cmd.Flags = getFlags(strings.Join(ss[i:], " "))
			break;
		}
		cmd.Child = &Command{Name: ss[i]}
		current = cmd.Child
	}
	return c, nil
}

func getFlags(cmd string) map[string]string {
	if strings.HasPrefix(cmd, "-") {
		// add a space as prefix in order to use the split function
		cmd = " " + cmd
	}

	f := strings.Split(cmd, " -")
	flags := make(map[string]string)
	for i := 1; i < len(f); i++ {
		// this is the case that a command is passed with '--' as prefix
		if strings.HasPrefix(f[i], "-") {
			f[i] = strings.Replace(f[i], "-", "", 1)
		}
		flag := f[i]
		k, v := getKeyValue(flag)
		flags[k] = v
	}
	return flags
}

func getKeyValue(s string) (string, string) {
	ff := strings.Split(s, " ")
	if len(ff) == 1 {
		return util.Trim(ff[0]), ""
	}

	if len(ff) > 2 {
		key := util.Trim(ff[0])
		val := ""
		for i := 1; i < len(ff); i++ {
			val += ff[i] + " "
		}
		val = util.Trim(val)
		val = strings.Replace(val, "\"", "", -1)
		return key, val
	}
	return util.Trim(ff[0]), util.Trim(ff[1])
}

func getCommand(cmd string) string {
	ss := strings.Split(cmd, " ")
	if strings.HasPrefix(ss[0], "-") {
		return ""
	}
	return ss[0]
}
