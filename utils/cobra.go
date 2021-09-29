package utils

import (
	"github.com/spf13/cobra"
)

type FlagParser func() error

type CommandRegister struct {
	Command   *cobra.Command
	ParseFlag FlagParser
}

func (c *CommandRegister) Register(cmd *cobra.Command) error {
	cmd.AddCommand(c.Command)
	return c.ParseFlag()
}
