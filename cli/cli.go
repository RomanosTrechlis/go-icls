package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/RomanosTrechlis/go-icls/parse"
	"strings"
)

// CLI holds the closing channel.
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
	cmd, err := cli.parse(textCmd)
	if err != nil {
		return false, err
	}
	if cmd.Command == "quit" {
		return true, nil
	}
	if cli.Command(cmd.Command) == nil {
		return false, fmt.Errorf("failed to find command '%s'", cmd.Command)
	}
	if cmd.help() {
		fmt.Fprintf(os.Stdout, "%v", cli.Command(cmd.Command))
		return false, nil
	}
	handler := cli.Command(cmd.Command).handler
	if handler == nil {
		return false, fmt.Errorf("there is no handler for the command '%s'", cmd.Command)
	}
	return false, handler(*cmd)
}

func (cli *CLI) FlagValue(flag string, c Command) string {
	cmd := cli.Command(c.Command)
	f := cmd.getFlag(flag)
	if s, ok := c.Flags[f.name]; ok {
		return s
	}
	return c.Flags[f.alias]
}

func (cli *CLI) parse(cmd string) (*Command, error) {
	cmd = strings.Trim(cmd, " ")
	cName, cFlags, err := parse.Parse(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to parse command: %v", err)
	}
	return &Command{cName, cFlags}, nil
}

func (cli *CLI) quit() {
	close(cli.closeChan)
}
