package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/RomanosTrechlis/go-icls/parse"
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
	if cli.Command(cmd) == nil {
		return false, fmt.Errorf("failed to find command '%s'", cmd)
	}
	if help(flags) {
		fmt.Fprintf(os.Stdout, "%v", cli.Command(cmd))
		return false, nil
	}
	handler := cli.Command(cmd).handler
	if handler == nil {
		return false, fmt.Errorf("there is no handler for the command '%s'", cmd)
	}
	return false, handler(cmd, flags)
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
	c := cli.Command(name)
	if c == nil {
		c = cli.New(name, "", handler)
		return
	}
	c.handler = handler
}

func (cli *CLI) FlagValue(command, flag string, flags map[string]string) interface{} {
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

func (cli *CLI) BoolValue(flag, c string, flags map[string]string) bool {
	cmd := cli.Command(c)
	f := cmd.getFlag(flag)
	if s, ok := flags[f.name]; ok {
		b, _ := strconv.ParseBool(s)
		return b
	}
	b, _ := strconv.ParseBool(flags[f.alias])
	return b
}

func (cli *CLI) IntValue(flag, c string, flags map[string]string) int64 {
	cmd := cli.Command(c)
	f := cmd.getFlag(flag)
	if s, ok := flags[f.name]; ok {
		i, _ := strconv.Atoi(s)
		return int64(i)
	}
	i, _ := strconv.Atoi(flags[f.alias])
	return int64(i)
}

func (cli *CLI) DoubleValue(flag, c string, flags map[string]string) float64 {
	cmd := cli.Command(c)
	f := cmd.getFlag(flag)
	if s, ok := flags[f.name]; ok {
		d, _ := strconv.ParseFloat(s, 64)
		return d
	}
	d, _ := strconv.ParseFloat(flags[f.alias], 64)
	return d
}

func (cli *CLI) parse(cmd string) (string, map[string]string) {
	cmd = strings.Trim(cmd, " ")
	return parse.Parse(cmd)
}

func (cli *CLI) quit() {
	close(cli.closeChan)
}
