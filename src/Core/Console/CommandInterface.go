package console

import "github.com/spf13/cobra"

type CommandInterface interface {
	GetName() string
	GetShortCut() string
	GetDescription() string
	Configure(command *cobra.Command)
	Handle(cmd *cobra.Command, args []string)
}
