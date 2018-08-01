// Copyright 2017 The go-icls Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util

func Trim(s string) string {
	const space = ' '
	n := len(s)
	l, h := 0, n
	for l < n && s[l] == space {
		l++
	}
	for h > l && s[h-1] == space {
		h--
	}
	return s[l:h]
}
