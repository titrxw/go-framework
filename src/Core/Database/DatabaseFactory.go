package database

import (
	config "github.com/titrxw/go-framework/src/Core/Config"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type DbLogger struct {
	logger.Writer
	logger *zap.Logger
}

func (this *DbLogger) Printf(info string, vs ...interface{}) {
	var _vs []zap.Field
	for k, v := range vs {
		_vs = append(_vs, zap.Reflect(string(rune(k)), v))
	}
	this.logger.Info(info, _vs...)
}

type DatabaseFactory struct {
	dbMu          sync.Mutex
	dbResolverMap map[string]func() *gorm.DB
	dbMap         map[string]*gorm.DB
	logger        *zap.Logger
}

func NewDatabaseFactory(logger *zap.Logger) *DatabaseFactory {
	return &DatabaseFactory{
		dbMap:         make(map[string]*gorm.DB),
		dbResolverMap: make(map[string]func() *gorm.DB),
		logger:        logger,
	}
}

func (this *DatabaseFactory) Channel(channel string) *gorm.DB {
	db, exists := this.dbMap[channel]
	if exists {
		return db
	}

	this.dbMu.Lock()
	defer this.dbMu.Unlock()

	db, exists = this.dbMap[channel]
	if exists {
		return db
	}

	dbResover, exists := this.dbResolverMap[channel]
	if !exists {
		panic("db channel " + channel + " not exists")
	}

	this.dbMap[channel] = dbResover()

	return this.dbMap[channel]
}

func (this *DatabaseFactory) makeDb(databaseConfig config.Database) *gorm.DB {
	//"user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dns := databaseConfig.User + ":" + databaseConfig.Password + "@tcp(" + databaseConfig.Host + ":" + strconv.Itoa(databaseConfig.Port) + ")/" + databaseConfig.DbName + "?charset=" + databaseConfig.Charset + "&parseTime=True&loc=Local"
	newLogger := logger.New(
		&DbLogger{
			logger: this.logger,
		},
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // Slow SQL threshold
			LogLevel:                  logger.Info,            // Log level
			IgnoreRecordNotFoundError: true,                   // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,                  // Disable color
		},
	)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dns,   // data source name
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   databaseConfig.Prefix,
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	mysqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}

	mysqlDb.SetMaxOpenConns(databaseConfig.PoolSize)

	return db
}

func (this *DatabaseFactory) Register(maps map[string]config.Database) {
	for key, value := range maps {
		this.dbResolverMap[key] = func() *gorm.DB {
			return this.makeDb(value)
		}
	}
}
