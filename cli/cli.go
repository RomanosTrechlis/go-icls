// Copyright 2017 The go-icls Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/RomanosTrechlis/go-icls/parse"
)

// CLI holds the closing channel and the defined commands.
type CLI struct {
	commands  map[string]*command
	closeChan chan struct{}
}

type Flags map[string]string

// BoolValue returns true if flag name or alias is passed as a flag in the executed command
func BoolValue(flag, alias string, flags Flags) (bool, error) {
	_, ok := flags[flag]
	if ok {
		return true, nil
	}
	_, ok = flags[alias]
	if ok {
		return true, nil
	}
	return false, fmt.Errorf("non existing flag: %s/%s", flag, alias)
}

// IntValue returns the integer value of a flag passed in the executed command
func IntValue(flag, alias string, flags Flags) (int, error) {
	v, ok := flags[flag]
	if ok {
		return strconv.Atoi(v)
	}
	v, ok = flags[alias]
	if ok {
		return strconv.Atoi(v)
	}
	return 0, fmt.Errorf("non existing flag: %s/%s", flag, alias)
}

// FloatValue returns the float value of a flag passed in the executed command
func FloatValue(flag, alias string, flags Flags) (float64, error) {
	v, ok := flags[flag]
	if ok {
		return strconv.ParseFloat(v, 64)
	}
	v, ok = flags[alias]
	if ok {
		return strconv.ParseFloat(v, 64)
	}
	return 0.0, fmt.Errorf("non existing flag: %s/%s", flag, alias)
}

// StringValue returns the string value of a flag passed in the executed command
func StringValue(flag, alias string, flags Flags) (string, error) {
	v, ok := flags[flag]
	if ok {
		return v, nil
	}
	v, ok = flags[alias]
	if ok {
		return v, nil
	}
	return "", fmt.Errorf("non existing flag: %s/%s", flag, alias)
}

// New creates a CLI struct.
func New() *CLI {
	return &CLI{
		commands:  make(map[string]*command),
		closeChan: make(chan struct{}, 1),
	}
}

// Run begins reading from the Stdin for ever. Parses the
// command given and apply it to the command handler.
func (cli *CLI) Run() {
	scanner := bufio.NewScanner(os.Stdin)
	go func() {
		fmt.Fprintf(os.Stdout, "> ")
		for scanner.Scan() {
			exit, err := cli.Execute(scanner.Text())
			if exit {
				cli.quit()
				return
			}

			if err != nil && err.Error() != "" {
				fmt.Fprintf(os.Stderr, "command failed: %v\n", err)
			}

			fmt.Fprintf(os.Stdout, "> ")
		}
	}()
	<-cli.closeChan
}

// Execute parses a string and applies the f function. Returns true for exiting.
func (cli *CLI) Execute(textCmd string) (bool, error) {
	trimedCmd := strings.Trim(textCmd, " ")
	if trimedCmd == "" {
		return false, nil
	}
	cmd, flags := cli.parse(trimedCmd)
	if cmd == "quit" || cmd == "q" {
		return true, nil
	}
	if cli.Command(cmd) == nil && !help(flags) {
		return false, fmt.Errorf("failed to find command '%s'", cmd)
	}
	if help(flags) {
		cli.printHelp(cmd, flags)
		return false, nil
	}
	if _, ok := cli.validateFlags(cmd, flags); !ok {
		cli.printHelp(cmd, flags)
		return false, fmt.Errorf("")
	}
	handler := cli.Command(cmd).handler
	if handler == nil {
		return false, fmt.Errorf("there is no handler for the command '%s'", cmd)
	}
	return false, handler(flags)
}

// New creates a command
func (cli *CLI) New(name, shortDesc, description string, handler func(flags Flags) error) *command {
	cmd := &command{
		name:        name,
		shortDesc:   shortDesc,
		description: description,
		flags:       make(map[string]*flag),
		handler:     handler,
	}
	cli.commands[name] = cmd
	return cmd
}

// New creates a command
func (cli *CLI) Simple(name, shortDesc, description string) *command {
	cmd := &command{
		name:        name,
		shortDesc:   shortDesc,
		description: description,
		handler:     emptyHandler(),
		flags:       make(map[string]*flag),
	}
	cli.commands[name] = cmd
	return cmd
}

// Command returns a command reference
func (cli *CLI) Command(name string) *command {
	cmd, ok := cli.commands[name]
	if !ok {
		return nil
	}
	return cmd
}

// HandlerFunc adds a handler to the specific command
func (cli *CLI) HandlerFunc(commandName string, handler func(flags Flags) error) {
	c := cli.Command(commandName)
	if c == nil {
		cli.New(commandName, "", "", handler)
		return
	}
	c.handler = handler
}

// FlagValue returns the value from the flag list.
func (cli *CLI) FlagValue(command, flag string, flags Flags) (interface{}, error) {
	cmd := cli.Command(command)
	f := cmd.getFlag(flag)
	s := cli.getValueFromFlag(f, flags)
	return conv(s, getDataTypeFunction(f.dataType))
}

// StringValue returns the string value from the flag list.
func (cli *CLI) StringValue(flag, c string, flags Flags) string {
	cmd := cli.Command(c)
	f := cmd.getFlag(flag)
	s := cli.getValueFromFlag(f, flags)
	return s
}

// BoolValue returns the bool value from the flag list.
func (cli *CLI) BoolValue(flag, c string, flags Flags) (bool, error) {
	cmd := cli.Command(c)
	f := cmd.getFlag(flag)
	s := cli.getValueFromFlag(f, flags)
	return strconv.ParseBool(s)
}

// IntValue returns the int value from the flag list.
func (cli *CLI) IntValue(flag, c string, flags Flags) (int, error) {
	cmd := cli.Command(c)
	f := cmd.getFlag(flag)
	s := cli.getValueFromFlag(f, flags)
	i, err := strconv.Atoi(s)
	return i, err
}

// DoubleValue returns the float64 value from the flag list.
func (cli *CLI) DoubleValue(flag, c string, flags Flags) (float64, error) {
	cmd := cli.Command(c)
	f := cmd.getFlag(flag)
	if f == nil {
		return 0.0, fmt.Errorf("couldn't find flag '%s' in command tree", flag)
	}
	s := cli.getValueFromFlag(f, flags)
	return strconv.ParseFloat(s, 64)
}

func (cli *CLI) getValueFromFlag(flag *flag, flags Flags) string {
	if s, ok := flags[flag.name]; ok {
		if flag.dataType == "bool" {
			return "true"
		}
		return s
	}
	if s, ok := flags[flag.alias]; ok {
		if flag.dataType == "bool" {
			return "true"
		}
		return s
	}

	return flag.defaultValueToString()
}

func (cli *CLI) parse(cmd string) (string, Flags) {
	cmd = strings.Trim(cmd, " ")
	return parse.Parse(cmd)
}

func (cli *CLI) quit() {
	close(cli.closeChan)
}

func (cli *CLI) printHelp(cmd string, flags Flags) {
	if cmd == "" {
		fmt.Fprintf(os.Stdout, "%v\n", cli)
		return
	}
	fmt.Fprintf(os.Stdout, "%v", cli.Command(cmd))
}

func (cli *CLI) String() string {
	app := os.Args[0]
	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 1, 8, 8, ' ', tabwriter.TabIndent)
	fmt.Fprintf(w, "Usage:\n\n\t%s <command> [options]\n\n", app)

	// when there is an empty command prints the flags
	empty := cli.Command("")
	if empty != nil {
		fmt.Fprintf(w, "Flags:\n")
		for _, flag := range empty.flags {
			fmt.Fprintf(w, "%s\n", flag)
		}
	}

	fmt.Fprintf(w, "Commands:\n")
	for k, v := range cli.commands {
		// when it's an empty command skip the printing as a command
		if k == "" {
			continue
		}
		fmt.Fprintf(w, "\t%s\t%s\n", k, v.shortDesc)
	}
	fmt.Fprintf(w, "\nUse \"%s <command> -h\" for more information about a command.", app)
	w.Flush()

	return buf.String()
}

func (cli *CLI) validateFlags(cmd string, flags Flags) (string, bool) {
	c := cli.Command(cmd)
	for _, f := range c.flags {
		if !f.isRequired {
			continue
		}
		// additionally the flag must have a value
		if fl, ok := flags[f.name]; !ok || fl == "" {
			return f.name, false
		}
	}
	return "", true
}

func emptyHandler() func(flags Flags) error {
	return func(flags Flags) error {
		return nil
	}
}
