# go-icls

[![Build Status](https://travis-ci.org/RomanosTrechlis/go-icls.svg?branch=master)](https://travis-ci.org/RomanosTrechlis/go-icls)
[![Go Report Card](https://goreportcard.com/badge/github.com/RomanosTrechlis/go-icls)](https://goreportcard.com/report/github.com/RomanosTrechlis/go-icls)
[![codecov](https://codecov.io/gh/RomanosTrechlis/go-icls/branch/master/graph/badge.svg)](https://codecov.io/gh/RomanosTrechlis/go-icls)


go-icls is a small and simple library for creating cli applications.

## Overview

### 1. Execute
go-icls provides a simple parser and command structure for executing single commands.

### 2. Run
go-icls provides an interactive command line interface that takes user input as commands and parses it.

### 3. Reading config files
go-icls provides the functionality to create a configuration structure by reading single *.properties files 
or by reading entire folder containing multiple *.properties files.

## Structure
go-icls uses a simple structure for building the command tree. **CLI**, **Command** and **Flag** provide the necessary
functionality.

### CLI
**CLI** is the parent of everything. Contains a map of command names and a **Command** structure.

Any command contains several key elements.

    app get -d dir -f file

In the previous example the *app* is the application name, the *get* is the command, and *-d*/*-f* are the flag 
keys followed by their respective values.

### Command
**Command** is the pivot point of any cli application. It keeps a name, an alias, a description, an array of **flag**, 
and more importantly a handler.

### Flag
**Flag** is a way of providing specific functionality on a broader command. It can be a key-value pair or a single key.
It must begin with either '-' or '--".

### Handler
**Handler** is a function that must be passed into the command.

```go
func(flags map[string]string) error
```

The command then executes this function when called.

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
	get := c.New("get", "get gets", "get gets", func(flags map[string]string) error {
		fmt.Println("This is the get command")
		return nil
	})
	// add 'get' command flags
	get.StringFlag("d", "", "", "directory name", false)
	get.StringFlag("t", "tetetetetetetet", "", "directory name", true)

	// add 'put' command to cli
	put := c.New("puttertesttest", "putter puts", "putter puts", func(flags map[string]string) error {
		i, err := c.IntValue("g", "puttertesttest", flags)
		if err != nil {
			return err
		}
		fmt.Printf("This is the putter command: %d\n", i)
		return nil
	})
	// add 'put' command flags
	put.StringFlag("f", "file", "", "filename to put", true)
	put.IntFlag("g", "int", 1, "int to put", false)

    // empty command with flags only
	base := c.New("", "empty name", "empty name", func(flags map[string]string) error {
		b, _ := c.BoolValue("v", "", flags)
		if b {
			fmt.Println("this is verbose")
		}
		return nil
	})
	base.BoolFlag("v", "verbose", "add prints", false)
	c.Run()
}
```

When executing the program:

    $ go-icls
    > -h
    Usage:
    
            go-icls.exe <command> [options]
    
    Flags:
            -v    --verbose
                                    add prints
    
    Commands:
            get                   get gets
            puttertesttest        putter puts
    
    Use "go-icls.exe <command> -h" for more information about a command.
    > get -h
    usage: get [get flags]
    
    get gets
    
    Flags:
    
            -d
                                    directory name
            -h    --help
                                    prints out information about the command
            -t    --tetetetetetetet
                                    directory name (required: true)
    > get -d test
    This is the get command
    > quit

    $

In order to include go-icls functionality into an application, import:
```go
import "github.com/RomanosTrechlis/go-icls/cli"
```

Use the **quit** command to exit the interactive interface.

## TODO

- [X] Add help command for printing the command tree.
- [X] Separate short and long command descriptions. Short must appear in general command tree help.
- [X] Change the signature of the handler function. Remove the command param since it is known to the user beforehand.
- [X] Validate for required flags.
- [X] Enable default values.
- [ ] Add more tests.
- [ ] Add var name on printed help after non bool flag.
