package parse

import (
	"strings"
	"fmt"
)

// Command holds the actual name of command
// the several flags that might exist.
type Command struct {
	Command string
	Flags   map[string]string
}

// Parse parses the command given and return a Command struct
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
func Parse(cmd string) (*Command, error) {
	cmdName, err := getCommand(cmd)
	if err != nil {
		return nil, err
	}

	command := &Command{
		Command: cmdName,
		Flags: make(map[string]string),
	}

	err = command.getFlags(cmd)
	return command, nil
}

func (c *Command) getFlags(cmd string) error {
	f := strings.Split(cmd, " -")
	for i := 1; i < len(f); i++ {
		flag := f[i]
		k, v := getKeyValue(flag)
		c.Flags[k] = v
	}
	return nil
}

func getKeyValue(s string) (string, string) {
	ff := strings.Split(s, " ")
	if len(ff) == 1 {
		return strings.Trim(ff[0], " "), ""
	}

	if len(ff) > 2 {
		key := strings.Trim(ff[0], " ")
		val := ""
		for i := 1; i < len(ff); i++ {
			val += ff[i] + " "
		}
		val = strings.Trim(val, " ")
		val = strings.Replace(val, "\"", "", -1)
		return key, val
	}
	return strings.Trim(ff[0], " "), strings.Trim(ff[1], " ")
}

func getCommand(cmd string) (string, error) {
	ss := strings.Split(cmd, " ")
	if len(ss) == 0 {
		return "", fmt.Errorf("failed to read command")
	}
	return ss[0], nil
}
