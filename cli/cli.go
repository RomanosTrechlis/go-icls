package cli

import (
	"bufio"
	"fmt"
	"os"
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
	cmd, flags, err := cli.parse(textCmd)
	if err != nil {
		return false, err
	}
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

func (cli *CLI) FlagValue(flag, c string, flags map[string]string) string {
	cmd := cli.Command(c)
	f := cmd.getFlag(flag)
	if s, ok := flags[f.name]; ok {
		return s
	}
	return flags[f.alias]
}

func (cli *CLI) parse(cmd string) (string, map[string]string, error) {
	cmd = strings.Trim(cmd, " ")
	return parse.Parse(cmd)
}

func (cli *CLI) quit() {
	close(cli.closeChan)
}
