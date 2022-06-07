package provider

import (
	"github.com/golobby/container/v3/pkg/container"
	database "github.com/titrxw/go-framework/src/Core/Database"
)

type ProviderInterface interface {
	Register(options interface{})
	SetContainer(container container.Container)
	SetDbFactory(DbFactory *database.DatabaseFactory)
}
