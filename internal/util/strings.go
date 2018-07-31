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
