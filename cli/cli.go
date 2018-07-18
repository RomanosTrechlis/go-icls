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
	closeChan chan struct{}
}

// New creates a CLI struct.
func New() *CLI {
	return &CLI{
		closeChan: make(chan struct{}, 1),
	}
}

// Run begins reading from the Stdin for ever. Parses the
// command given and apply it to the function passed as parameter.
//
// Example:
//
// func main() {
// 	    c := cli.New()
//	    c.Run(Test)
// }
//
// func Test(cmd cli.Command) {
//     if cmd.Command == "get" {
//	       fmt.Fprintf(os.Stdout, "%v", cmd)
//	   }
// }
//
func (cli *CLI) Run(f func(c Command)) {
	scanner := bufio.NewScanner(os.Stdin)
	go func() {
		fmt.Fprintf(os.Stdout, "> ")
		for scanner.Scan() {
			exit, err := cli.RunCommand(scanner.Text(), f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "command '%s' coudn't be parse: %v", err)
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

// RunCommand parses a string and applies the f function. Returns true for exiting.
func (cli *CLI) RunCommand(textCmd string, f func(c Command)) (bool, error) {
	cmd, err := cli.parse(textCmd)
	if err != nil {
		return false, err
	}
	if cmd.Command == "quit" {
		return true, nil
	}
	f(*cmd)
	return false, nil
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
