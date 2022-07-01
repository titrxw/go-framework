package app

import (
	"github.com/asaskevich/EventBus"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/golobby/container/v3/pkg/container"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	appconfig "github.com/titrxw/go-framework/src/Core/Config"
	console "github.com/titrxw/go-framework/src/Core/Console"
	database "github.com/titrxw/go-framework/src/Core/Database"
	exception "github.com/titrxw/go-framework/src/Core/Exception"
	logger "github.com/titrxw/go-framework/src/Core/Logger"
	provider "github.com/titrxw/go-framework/src/Core/Provider"
	redis "github.com/titrxw/go-framework/src/Core/Redis"
)

type App struct {
	Name              string
	HandlerExceptions *exception.HandlerExceptions
	Config            *appconfig.Config
	Container         container.Container
	Event             EventBus.Bus
	RedisFactory      *redis.RedisFactory
	DbFactory         *database.DatabaseFactory
	LoggerFactory     *logger.LoggerFactory
	ProviderManager   *provider.ProviderManager
	Translator        ut.Translator
	Console           *console.Console
}

func NewApp() *App {
	return &App{}
}

func (app *App) registerExceptionHandler() {
	app.HandlerExceptions = &exception.HandlerExceptions{
		Logger: app.LoggerFactory.Channel("default"),
	}
}

func (app *App) InitConfig(obj interface{}) {
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

func (app *App) registerConfig() {
	app.InitConfig(&app.Config)
	app.Name = app.Config.App.Name
}

func (app *App) registerContainer() {
	app.Container = container.New()
}

func (app *App) registerEvent() {
	app.Event = EventBus.New()
}

func (app *App) registerLogger() {
	app.LoggerFactory = logger.NewLoggerFactory()
	app.LoggerFactory.Register(app.Config.LogMap)
}

func (app *App) registerRedis() {
	app.RedisFactory = redis.NewRedisFactory()
	app.RedisFactory.Register(app.Config.RedisMap)
}

func (app *App) registerDb() {
	logger := app.LoggerFactory.RegisterLogger(appconfig.Log{
		Path:    "db.log",
		Level:   "info",
		MaxDays: 7,
	})

	app.DbFactory = database.NewDatabaseFactory(logger)
	app.DbFactory.Register(app.Config.DatabaseMap)
}

func (app *App) registerValidation() {
	uni := ut.New(zh.New())
	lang := app.Config.App.Lang
	if lang == "" {
		lang = "zh"
	}

	app.Translator, _ = uni.GetTranslator(lang)
	_ = zh_translations.RegisterDefaultTranslations(binding.Validator.Engine().(*validator.Validate), app.Translator)
}

func (app *App) registerProvider() {
	err := app.Container.NamedSingleton("provider_manager", func() *provider.ProviderManager {
		return &provider.ProviderManager{
			Container: app.Container,
			DbFactory: app.DbFactory,
		}
	})
	if err != nil {
		panic(err)
	}

	var providerManager *provider.ProviderManager
	err = app.Container.NamedResolve(&providerManager, "provider_manager")
	if err != nil {
		panic(err)
	}

	app.ProviderManager = providerManager
}

func (app *App) registerConsole() {
	app.Console = console.NewConsole()
}

func (app *App) Bootstrap() {
	app.registerConsole()
	app.registerConfig()
	app.registerContainer()
	app.registerEvent()
	app.registerLogger()
	app.registerExceptionHandler()
	app.registerRedis()
	app.registerDb()
	app.registerProvider()
	app.registerValidation()
}
