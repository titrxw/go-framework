package server

import (
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gin-gonic/gin"
	app "github.com/titrxw/go-framework/src/App"
	middleware "github.com/titrxw/go-framework/src/Http/Middleware"
	session "github.com/titrxw/go-framework/src/Http/Session"
	"net/http"
)

type Server struct {
	App *app.App

	ginEngine *gin.Engine
	Session   *session.Session
}

func NewHttpSerer(app *app.App) *Server {
	server := &Server{
		App: app,
	}
	server.RegisterSession()

	return server
}

func (this *Server) RegisterSession() {
	this.Session = &session.Session{
		SessionManager: *scs.New(),
	}
	this.Session.ErrorFunc = func(writer http.ResponseWriter, request *http.Request, err error) {

	}
	if this.App.Config.Session.Lifetime > 0 {
		this.Session.Lifetime = this.App.Config.Session.Lifetime
	}
	if this.App.Config.Cookie.Name != "" {
		this.Session.Cookie.Name = this.App.Config.Cookie.Name
	}
	if this.App.Config.Cookie.Domain != "" {
		this.Session.Cookie.Domain = this.App.Config.Cookie.Domain
	}
	if this.App.Config.Cookie.Path != "" {
		this.Session.Cookie.Path = this.App.Config.Cookie.Path
	}
	this.Session.Cookie.HttpOnly = this.App.Config.Cookie.HttpOnly
	this.Session.Cookie.Persist = this.App.Config.Cookie.Persist
	this.Session.Cookie.Secure = this.App.Config.Cookie.Secure
	if this.App.Config.Cookie.SameSite == "Lax" {
		this.Session.Cookie.SameSite = http.SameSiteLaxMode
	} else if this.App.Config.Cookie.SameSite == "Strict" {
		this.Session.Cookie.SameSite = http.SameSiteStrictMode
	} else if this.App.Config.Cookie.SameSite == "None" {
		this.Session.Cookie.SameSite = http.SameSiteNoneMode
	} else {
		this.Session.Cookie.SameSite = http.SameSiteDefaultMode
	}

	this.Session.SetStorageResolver(func() scs.Store {
		sessionDb := this.App.Config.Session.DbConnection
		if sessionDb == "" {
			sessionDb = "default"
		}
		db, err := this.App.DbFactory.Channel(sessionDb).DB()
		if err != nil {
			panic(err)
		}

		return mysqlstore.New(db)
	})
}

func (this *Server) initGinEngine() {
	gin.SetMode(this.App.Config.App.Env)
	this.ginEngine = gin.Default()
	this.ginEngine.Use(middleware.ExceptionMiddleware{HandlerExceptions: this.App.HandlerExceptions}.Process)
	this.ginEngine.Use(middleware.NewSessionMiddleware(this.Session).Process)
}

func (this *Server) RegisterRouters(register func(engine *gin.Engine)) *Server {
	register(this.ginEngine)
	return this
}

func (this *Server) Start(addr ...string) {
	this.initGinEngine()
	err := this.ginEngine.Run(addr...)

	if err != nil {
		panic(err)
	}
}
