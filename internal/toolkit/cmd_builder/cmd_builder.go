package cmdbuilder

import (
	"errors"

	"github.com/google/shlex"
)

type CmdBuilder struct {
	bin  string
	args []string
	err  error
}

func NewCmdBuilder(bin string) *CmdBuilder {
	return &CmdBuilder{
		bin: bin,
	}
}

func (c *CmdBuilder) Add(args ...string) *CmdBuilder {
	c.args = append(c.args, args...)
	return c
}

func (c *CmdBuilder) AddIf(cond bool, args ...string) *CmdBuilder {
	if cond {
		c.Add(args...)
	}
	return c
}

func (c *CmdBuilder) AddIfElse(cond bool, args1 []string, args2 []string) *CmdBuilder {
	if cond {
		c.Add(args1...)
	} else {
		c.Add(args2...)
	}
	return c
}

func (c *CmdBuilder) AddSplit(input string) *CmdBuilder {
	return c.AddSplitWithError(input, errors.New("failed to split arguments"))
}

func (c *CmdBuilder) AddSplitWithError(input string, e error) *CmdBuilder {
	parts, err := shlex.Split(input)
	if err != nil {
		c.err = errors.Join(e, err)
		return c
	}
	return c.Add(parts...)
}

func (c *CmdBuilder) Build() (string, []string, error) {
	return c.bin, c.args, c.err
}
