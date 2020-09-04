// Copyright 2017 The go-icls Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parse

import (
	"strings"

	"github.com/RomanosTrechlis/go-icls/internal/util"
)

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
// help command prints the same as command -h
// help add
//
func Parse(cmd string) (string, map[string]string) {
	cmdName := getCommand(cmd)
	if cmdName == "help" {
		cmdName = getNextCommand(cmd, cmdName)
		cmd = cmdName + " -h"
	}
	flags := getFlags(cmd)
	return cmdName, flags
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

func getNextCommand(cmd, current string) string {
	ss := strings.Split(cmd, " ")
	l := len(ss)
	for i := 0; i < l; i++ {
		if ss[i] != current {
			continue;
		}
		if i+1 >= l {
			continue
		}

		if strings.HasPrefix(ss[i+1], "-") {
			return ""
		}
		return ss[i+1]
	}
	return ""
}
