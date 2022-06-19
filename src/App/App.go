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

func (this *App) registerExceptionHandler() {
	this.HandlerExceptions = &exception.HandlerExceptions{
		Logger: this.LoggerFactory.Channel("default"),
	}
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

func (this *App) registerEvent() {
	this.Event = EventBus.New()
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

func (this *App) registerValidation() {
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

func (this *App) registerConsole() {
	this.Console = console.NewConsole()
}

func (this *App) Bootstrap() {
	this.registerConsole()
	this.registerConfig()
	this.registerContainer()
	this.registerEvent()
	this.registerLogger()
	this.registerExceptionHandler()
	this.registerRedis()
	this.registerDb()
	this.registerProvider()
	this.registerValidation()
}
