package server

import (
	"github.com/spf13/cobra"
	console "github.com/titrxw/go-framework/src/Core/Console"
	server "github.com/titrxw/go-framework/src/Http/Server"
)

type StartCommand struct {
	console.CommandAbstract

	Server *server.Server
}

func (startCommand *StartCommand) GetName() string {
	return "server:start"
}

func (startCommand *StartCommand) GetShortCut() string {
	return ""
}

func (startCommand *StartCommand) GetDescription() string {
	return ""
}

func (startCommand *StartCommand) Handle(cmd *cobra.Command, args []string) {
	startCommand.Server.Start(args[0] + ":" + args[1])
}
