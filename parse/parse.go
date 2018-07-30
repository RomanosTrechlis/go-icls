package parse

import (
	"fmt"
	"strings"

	"github.com/RomanosTrechlis/go-icls/util"
)

// Parse parses the command given and return the command and the flags
//
// Example commands:
//
// command with key-value flags
// get -d dir -f filename
//
// command with key-value and key only flags
// del -d dir -f filename -e
//
// command with flag value of several words
// rem -d dir -f filename -e -m This is one
//
// command with flag value in quotation marks
// add -d dir -f filename -e -m "This is one"
//
func Parse(cmd string) (string, map[string]string, error) {
	cmdName, err := getCommand(cmd)
	if err != nil {
		return "", nil, err
	}

	flags := getFlags(cmd)
	return cmdName, flags, nil
}

func getFlags(cmd string) map[string]string {
	f := strings.Split(cmd, " -")
	flags := make(map[string]string)
	for i := 1; i < len(f); i++ {
		flag := f[i]
		k, v := getKeyValue(flag)
		flags[k] = v
	}
	return flags
}

func getKeyValue(s string) (string, string) {
	ff := strings.Split(s, " ")
	if len(ff) == 1 {
		return util.Trim(ff[0]), ""
	}

	if len(ff) > 2 {
		key := util.Trim(ff[0])
		val := ""
		for i := 1; i < len(ff); i++ {
			val += ff[i] + " "
		}
		val = util.Trim(val)
		val = strings.Replace(val, "\"", "", -1)
		return key, val
	}
	return util.Trim(ff[0]), util.Trim(ff[1])
}

func getCommand(cmd string) (string, error) {
	ss := strings.Split(cmd, " ")
	if len(ss) == 0 {
		return "", fmt.Errorf("failed to read command")
	}
	return ss[0], nil
}
