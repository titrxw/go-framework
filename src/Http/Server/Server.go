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

func (server *Server) initGinEngine() {
	gin.SetMode(server.App.Config.App.Env)
	server.GinEngine = gin.Default()
}

func (server *Server) RegisterRouters(register func(engine *gin.Engine)) *Server {
	register(server.GinEngine)
	return server
}

func (server *Server) Start(addr ...string) {
	err := server.GinEngine.Run(addr...)

	if err != nil {
		panic(err)
	}
}
