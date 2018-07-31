# go-icls
go-icls is a small and simple library for creating cli applications.

## Overview

### 1. Execute
go-icls provides a simple parser and command structure for executing single commands.

### 2. Run
go-icls provides an interactive command line interface that takes user input as commands and parses it.

### 3. Reading config files
go-icls provides the functionality to create a configuration structure by reading single *.properties files 
or by reading entire folder containing multiple *.properties files.

## How to use
The following example creates a cli with two commands **get** and **put**. Each command has a single flag.

Systemically there is the **help** command that can be used to print a helpful text on the terminal.

```go
package main

import (
	"fmt"

	"github.com/RomanosTrechlis/go-icls/cli"
)

func main() {
	// create a new cli
	c := cli.New()
	// add 'get' command to cli
	get := c.New("get", "get gets", func(cmd string, flags map[string]string) error {
		fmt.Println("This is the get command")
		return nil
	})
	// add 'get' command flags
	get.StringFlag("d", "dir", "directory name", false)

	// add 'put' command to cli
	put := c.New("put", "put puts", func(cmd string, flags map[string]string) error {
		fmt.Println("This is the put command")
		return nil
	})
	// add 'put' command flags
	put.StringFlag("f", "file", "filename to put", true)

	c.Run()
}
```

When executing the program:
```bash
$ go run main.go
> get -h
get
        -d      --dir       directory name (required: false)
        -h      --help      prints out information about the command
get gets
> get -d test
This is the get command
> quit

$
```

Use the **quit** command to exit the interactive interface.

## TODO

- [ ] Add help command for all commands
- [ ] Separate short and long command descriptions
