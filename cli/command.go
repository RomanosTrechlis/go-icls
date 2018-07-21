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
	description string
	// flags is an array of flag objects that
	// contain information on its of them
	flags []*flag
	// handler is the function executed when the command is called
	handler func(cmd string, flags map[string]string) error
}

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

// New creates a command
func (cli *CLI) New(name, description string, handler func(cmd string, flags map[string]string) error) *command {
	cmd := &command{
		name:        name,
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
func (cli *CLI) HandlerFunc(name string, handler func(cmd string, flags map[string]string) error) {
	cli.Command(name).handler = handler
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

func (f *flag) String() string {
	n := "\t"
	if f.name != "" {
		n += fmt.Sprintf("-%s\t", f.name)
	} else {
		n += fmt.Sprintf("\t")
	}
	if f.alias != "" {
		n += fmt.Sprintf("--%s", f.alias)
	} else {
		n += fmt.Sprintf("\t")
	}
	n += fmt.Sprintf("\t%s (required: %v)\n", f.description, f.isRequired)
	return n
}

func checkForKeysInMap(m map[string]string, keys ...string) bool {
	for k := range m {
		for _, key := range keys {
			if k == key {
				return true
			}
		}
	}
	return false
}
