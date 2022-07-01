package console

import "github.com/spf13/cobra"

type Console struct {
	Handler *cobra.Command
}

func NewConsole() *Console {
	return &Console{
		Handler: &cobra.Command{Use: "app"},
	}
}

func (console *Console) RegisterCommand(command CommandInterface) {
	handler := &cobra.Command{
		Use:   command.GetName(),
		Short: command.GetShortCut(),
		Long:  command.GetDescription(),
		Run:   command.Handle,
	}

	command.Configure(handler)

	console.Handler.AddCommand(handler)
}

func (console *Console) Run() {
	console.Handler.Execute()
}
