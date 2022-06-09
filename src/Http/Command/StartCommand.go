package command

import (
	"github.com/spf13/cobra"
	console "github.com/titrxw/go-framework/src/Core/Console"
	server "github.com/titrxw/go-framework/src/Http/Server"
)

type StartCommand struct {
	console.CommandAbstract

	Server *server.Server
}

func (this *StartCommand) GetName() string {
	return "server:start"
}

func (this *StartCommand) GetShortCut() string {
	return ""
}

func (this *StartCommand) GetDescription() string {
	return ""
}

func (this *StartCommand) Handle(cmd *cobra.Command, args []string) {
	this.Server.Start(args[0] + ":" + args[1])
}
