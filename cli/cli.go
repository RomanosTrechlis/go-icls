// Copyright 2017 The go-icls Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/RomanosTrechlis/go-icls/parse"
	"github.com/RomanosTrechlis/go-icls/internal/util"
	"text/tabwriter"
	"bytes"
)

// CLI holds the closing channel and the defined commands.
type CLI struct {
	commands  map[string]*command
	closeChan chan struct{}
}

// New creates a CLI struct.
func New() *CLI {
	return &CLI{
		commands:  make(map[string]*command, 0),
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
			if err != nil {
				fmt.Fprintf(os.Stderr, "command failed: %v\n", err)
				fmt.Fprintf(os.Stdout, "> ")
				continue
			}
			if exit {
				cli.quit()
				return
			}
			fmt.Fprintf(os.Stdout, "> ")
		}
	}()
	<-cli.closeChan
}

// Execute parses a string and applies the f function. Returns true for exiting.
func (cli *CLI) Execute(textCmd string) (bool, error) {
	cmd, flags := cli.parse(textCmd)
	if cmd == "quit" {
		return true, nil
	}
	if cli.Command(cmd) == nil && !help(flags) {
		return false, fmt.Errorf("failed to find command '%s'", cmd)
	}
	if help(flags) {
		cli.printHelp(cmd, flags)
		return false, nil
	}
	flag, ok := cli.validateFlags(cmd, flags)
	if !ok {
		return false, fmt.Errorf("flag '%s' is required", flag)
	}
	handler := cli.Command(cmd).handler
	if handler == nil {
		return false, fmt.Errorf("there is no handler for the command '%s'", cmd)
	}
	return false, handler(flags)
}

// New creates a command
func (cli *CLI) New(name, shortDesc, description string, handler func(flags map[string]string) error) *command {
	cmd := &command{
		name:        name,
		shortDesc:   shortDesc,
		description: description,
		flags:       make([]*flag, 0),
		handler:     handler,
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
func (cli *CLI) HandlerFunc(commandName string, handler func(flags map[string]string) error) {
	c := cli.Command(commandName)
	if c == nil {
		c = cli.New(commandName, "", "", handler)
		return
	}
	c.handler = handler
}

func (cli *CLI) FlagValue(command, flag string, flags map[string]string) (interface{}, error) {
	cmd := cli.Command(command)
	f := cmd.getFlag(flag)
	if s, ok := flags[f.name]; ok {
		return conv(s, getDataTypeFunction(f.dataType))
	}
	return conv(flags[f.alias], getDataTypeFunction(f.dataType))
}

func (cli *CLI) StringValue(flag, c string, flags map[string]string) string {
	cmd := cli.Command(c)
	f := cmd.getFlag(flag)
	if s, ok := flags[f.name]; ok {
		return s
	}
	return flags[f.alias]
}

func (cli *CLI) BoolValue(flag, c string, flags map[string]string) (bool, error) {
	cmd := cli.Command(c)
	f := cmd.getFlag(flag)
	if s, ok := flags[f.name]; ok {
		return strconv.ParseBool(s)
	}
	return strconv.ParseBool(flags[f.alias])
}

func (cli *CLI) IntValue(flag, c string, flags map[string]string) (int64, error) {
	cmd := cli.Command(c)
	f := cmd.getFlag(flag)
	if s, ok := flags[f.name]; ok {
		i, err := strconv.Atoi(s)
		return int64(i), err
	}
	i, err := strconv.Atoi(flags[f.alias])
	return int64(i), err
}

func (cli *CLI) DoubleValue(flag, c string, flags map[string]string) (float64, error) {
	cmd := cli.Command(c)
	f := cmd.getFlag(flag)
	if s, ok := flags[f.name]; ok {
		return strconv.ParseFloat(s, 64)
	}
	return strconv.ParseFloat(flags[f.alias], 64)
}

func (cli *CLI) parse(cmd string) (string, map[string]string) {
	cmd = util.Trim(cmd)
	return parse.Parse(cmd)
}

func (cli *CLI) quit() {
	close(cli.closeChan)
}

func (cli *CLI) printHelp(cmd string, flags map[string]string) {
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
	fmt.Fprintf(w, "Commands:\n")
	for k, v := range cli.commands {
		fmt.Fprintf(w, "\t%s\t%s\n", k, v.shortDesc)
	}
	fmt.Fprintf(w, "\nUse \"%s <command> -h\" for more information about a command.", app)
	w.Flush()

	return string(buf.Bytes())
}

func (cli *CLI) validateFlags(cmd string, flags map[string]string) (string, bool) {
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
