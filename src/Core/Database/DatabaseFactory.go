package database

import (
	"strconv"
	"sync"
	"time"

	config "github.com/titrxw/go-framework/src/Core/Config"

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

func (dbLogger *DbLogger) Printf(info string, vs ...interface{}) {
	var _vs []zap.Field
	for k, v := range vs {
		_vs = append(_vs, zap.Reflect(string(rune(k)), v))
	}
	dbLogger.logger.Info(info, _vs...)
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

func (databaseFactory *DatabaseFactory) Channel(channel string) *gorm.DB {
	db, exists := databaseFactory.dbMap[channel]
	if exists {
		return db
	}

	databaseFactory.dbMu.Lock()
	defer databaseFactory.dbMu.Unlock()

	db, exists = databaseFactory.dbMap[channel]
	if exists {
		return db
	}

	dbResover, exists := databaseFactory.dbResolverMap[channel]
	if !exists {
		panic("db channel " + channel + " not exists")
	}

	databaseFactory.dbMap[channel] = dbResover()

	return databaseFactory.dbMap[channel]
}

func (databaseFactory *DatabaseFactory) makeDb(databaseConfig config.Database) *gorm.DB {
	//"user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dns := databaseConfig.User + ":" + databaseConfig.Password + "@tcp(" + databaseConfig.Host + ":" + strconv.Itoa(databaseConfig.Port) + ")/" + databaseConfig.DbName + "?charset=" + databaseConfig.Charset + "&parseTime=True&loc=Local"
	newLogger := logger.New(
		&DbLogger{
			logger: databaseFactory.logger,
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

func (databaseFactory *DatabaseFactory) Register(maps map[string]config.Database) {
	for key, value := range maps {
		databaseFactory.dbResolverMap[key] = func() *gorm.DB {
			return databaseFactory.makeDb(value)
		}
	}
}
