package console

import "github.com/spf13/cobra"

type CommandAbstract struct {
	cobra.Command
	CommandInterface
}

func (this *CommandAbstract) Configure(command *cobra.Command) {

}
