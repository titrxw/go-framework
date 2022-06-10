package server

import (
	"github.com/gin-gonic/gin"
	app "github.com/titrxw/go-framework/src/App"
	session "github.com/titrxw/go-framework/src/Http/Session"
)

type Server struct {
	App *app.App

	GinEngine *gin.Engine
	Session   *session.Session
}

func NewHttpSerer(app *app.App) *Server {
	server := &Server{
		App: app,
	}
	server.initGinEngine()

	return server
}

func (this *Server) initGinEngine() {
	gin.SetMode(this.App.Config.App.Env)
	this.GinEngine = gin.Default()
}

func (this *Server) RegisterRouters(register func(engine *gin.Engine)) *Server {
	register(this.GinEngine)
	return this
}

func (this *Server) Start(addr ...string) {
	err := this.GinEngine.Run(addr...)

	if err != nil {
		panic(err)
	}
}
