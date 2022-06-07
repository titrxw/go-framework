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

func (this *Console) RegisterCommand(command CommandInterface) {
	handler := &cobra.Command{
		Use:   command.GetName(),
		Short: command.GetShortCut(),
		Long:  command.GetDescription(),
		Run:   command.Handle,
	}

	command.Configure(handler)

	this.Handler.AddCommand(handler)
}

func (this *Console) Run() {
	this.Handler.Execute()
}
