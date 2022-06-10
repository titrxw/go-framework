package route

import (
	"github.com/gin-gonic/gin"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	console "github.com/titrxw/go-framework/src/Core/Console"
	"os"
)

type ListCommand struct {
	console.CommandAbstract

	GinEngine *gin.Engine
}

func (this *ListCommand) GetName() string {
	return "route:list"
}

func (this *ListCommand) GetShortCut() string {
	return ""
}

func (this *ListCommand) GetDescription() string {
	return ""
}

func (this *ListCommand) Handle(cmd *cobra.Command, args []string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Path", "Method", "Handler"})
	t.AppendSeparator()
	for index, route := range this.GinEngine.Routes() {
		t.AppendRow([]interface{}{index, route.Path, route.Method, route.Handler})
	}
	t.Render()
}
