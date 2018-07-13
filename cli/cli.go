package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/RomanosTrechlis/go-icls/parse"
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
			cmd, err := cli.parse(scanner.Text())
			if err != nil {
				fmt.Fprintf(os.Stderr, "Command '%s' coudn't be parse: %v", err)
				continue
			}
			if cmd.Command == "quit" {
				cli.quit()
				return
			}
			f(*cmd)
			fmt.Fprintf(os.Stdout, "> ")
		}
	}()
	<-cli.closeChan
}

func (cli *CLI) parse(cmd string) (*Command, error) {
	cName, cFlags, err := parse.Parse(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to parse command: %v", err)
	}
	return &Command{cName, cFlags}, nil
}

func (cli *CLI) quit() {
	close(cli.closeChan)
}
