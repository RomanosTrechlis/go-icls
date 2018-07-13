package cli

// Command holds the actual name of command
// the several flags that might exist.
type Command struct {
	Command string
	Flags   map[string]string
}

func (c Command) String() string {
	s := c.Command + "\n"
	for k, v := range c.Flags {
		s += k + "=\"" + v + "\"\n"
	}
	return s
}

func (c Command) Validate(f func(c Command) error) error {
	return f(c)
}
