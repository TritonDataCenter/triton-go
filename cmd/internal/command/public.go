package command

import "github.com/spf13/cobra"

type SetupFunc func(parent *Command) error

type Command struct {
	Cobra *cobra.Command
	Setup SetupFunc
}
