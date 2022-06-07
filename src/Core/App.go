package core

import (
	appconfig "github.com/titrxw/go-framework/src/Core/Config"
	console "github.com/titrxw/go-framework/src/Core/Console"
	database "github.com/titrxw/go-framework/src/Core/Database"
	exception "github.com/titrxw/go-framework/src/Core/Exception"
	logger "github.com/titrxw/go-framework/src/Core/Logger"
	provider "github.com/titrxw/go-framework/src/Core/Provider"
	redis "github.com/titrxw/go-framework/src/Core/Redis"
	session "github.com/titrxw/go-framework/src/Core/Session"
	"net/http"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/golobby/container/v3/pkg/container"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

type App struct {
	Name              string
	HandlerExceptions *exception.HandlerExceptions
	Config            *appconfig.Config
	Container         container.Container
	RedisFactory      *redis.RedisFactory
	DbFactory         *database.DatabaseFactory
	LoggerFactory     *logger.LoggerFactory
	Session           *session.Session
	ProviderManager   *provider.ProviderManager
	Translator        ut.Translator
	Console           *console.Console
}

func NewApp() *App {
	return &App{}
}

func (this *App) registerExceptionHandler() {
	this.HandlerExceptions = &exception.HandlerExceptions{
		Logger: this.LoggerFactory.Channel("default"),
	}
	this.HandlerExceptions.RegisterExceptionHandle(nil)
}

func (this *App) InitConfig(obj interface{}) {
	config.ClearAll()
	config.WithOptions(config.ParseEnv)
	config.AddDriver(yaml.Driver)

	err := config.LoadFiles("./app_config.yaml")
	if err != nil {
		panic(err)
	}
	err = config.BindStruct("", &obj)
	if err != nil {
		panic(err)
	}
}

func (this *App) registerConfig() {
	this.InitConfig(&this.Config)
	this.Name = this.Config.App.Name
}

func (this *App) registerContainer() {
	this.Container = container.New()
}

func (this *App) registerLogger() {
	this.LoggerFactory = logger.NewLoggerFactory()
	this.LoggerFactory.Register(this.Config.LogMap)
}

func (this *App) registerRedis() {
	this.RedisFactory = redis.NewRedisFactory()
	this.RedisFactory.Register(this.Config.RedisMap)
}

func (this *App) registerDb() {
	logger := this.LoggerFactory.RegisterLogger(appconfig.Log{
		Path:    "db.log",
		Level:   "info",
		MaxDays: 7,
	})

	this.DbFactory = database.NewDatabaseFactory(logger)
	this.DbFactory.Register(this.Config.DatabaseMap)
}

func (this *App) registerSession() {
	this.Session = &session.Session{
		SessionManager: *scs.New(),
	}
	this.Session.ErrorFunc = func(writer http.ResponseWriter, request *http.Request, err error) {

	}
	if this.Config.Session.Lifetime > 0 {
		this.Session.Lifetime = this.Config.Session.Lifetime
	}
	if this.Config.Cookie.Name != "" {
		this.Session.Cookie.Name = this.Config.Cookie.Name
	}
	if this.Config.Cookie.Domain != "" {
		this.Session.Cookie.Domain = this.Config.Cookie.Domain
	}
	if this.Config.Cookie.Path != "" {
		this.Session.Cookie.Path = this.Config.Cookie.Path
	}
	this.Session.Cookie.HttpOnly = this.Config.Cookie.HttpOnly
	this.Session.Cookie.Persist = this.Config.Cookie.Persist
	this.Session.Cookie.Secure = this.Config.Cookie.Secure
	if this.Config.Cookie.SameSite == "Lax" {
		this.Session.Cookie.SameSite = http.SameSiteLaxMode
	} else if this.Config.Cookie.SameSite == "Strict" {
		this.Session.Cookie.SameSite = http.SameSiteStrictMode
	} else if this.Config.Cookie.SameSite == "None" {
		this.Session.Cookie.SameSite = http.SameSiteNoneMode
	} else {
		this.Session.Cookie.SameSite = http.SameSiteDefaultMode
	}

	this.Session.SetStorageResolver(func() scs.Store {
		sessionDb := this.Config.Session.DbConnection
		if sessionDb == "" {
			sessionDb = "default"
		}
		db, err := this.DbFactory.Channel(sessionDb).DB()
		if err != nil {
			panic(err)
		}

		return mysqlstore.New(db)
	})
}

func (this *App) RegisterValidation() {
	uni := ut.New(zh.New())
	lang := this.Config.App.Lang
	if lang == "" {
		lang = "zh"
	}

	this.Translator, _ = uni.GetTranslator(lang)
	_ = zh_translations.RegisterDefaultTranslations(binding.Validator.Engine().(*validator.Validate), this.Translator)
}

func (this *App) registerProvider() {
	err := this.Container.NamedSingleton("provider_manager", func() *provider.ProviderManager {
		return &provider.ProviderManager{
			Container: this.Container,
			DbFactory: this.DbFactory,
		}
	})
	if err != nil {
		panic(err)
	}

	var providerManager *provider.ProviderManager
	err = this.Container.NamedResolve(&providerManager, "provider_manager")
	if err != nil {
		panic(err)
	}

	this.ProviderManager = providerManager
}

func (this *App) RegisterConsole() {
	this.Console = console.NewConsole()
}

func (this *App) Bootstrap() {
	this.RegisterConsole()
	this.registerConfig()
	this.registerContainer()
	this.registerLogger()
	this.registerExceptionHandler()
	this.registerRedis()
	this.registerDb()
	this.registerProvider()
	this.registerSession()
	this.RegisterValidation()
}
