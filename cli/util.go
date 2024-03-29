// Copyright 2017 The go-icls Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"strconv"
)

func getDataTypeFunction(dt string) func(string) (interface{}, error) {
	switch dt {
	case "int":
		return func(s string) (interface{}, error) {
			i, err := strconv.Atoi(s)
			return int64(i), err
		}
	case "string":
		return func(s string) (interface{}, error) {
			return s, nil
		}
	case "bool":
		return func(s string) (interface{}, error) {
			b, err := strconv.ParseBool(s)
			return bool(b), err
		}
	case "float":
		return func(s string) (interface{}, error) {
			f, err := strconv.ParseFloat(s, 64)
			return float64(f), err
		}
	default:
		return nil
	}
}

func conv(t string, f func(s string) (interface{}, error)) (interface{}, error) {
	return f(t)
}

func checkForKeysInMap(m Flags, keys ...string) bool {
	for k := range m {
		for _, key := range keys {
			if k == key {
				return true
			}
		}
	}
	return false
}
