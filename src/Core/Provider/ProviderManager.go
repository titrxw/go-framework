package provider

import (
	"github.com/golobby/container/v3/pkg/container"
	database "github.com/titrxw/go-framework/src/Core/Database"
)

type ProviderManager struct {
	Container container.Container
	DbFactory *database.DatabaseFactory
}

func (providerManager *ProviderManager) MakeProvider(abstract ProviderInterface) ProviderInterface {
	abstract.SetContainer(providerManager.Container)
	abstract.SetDbFactory(providerManager.DbFactory)
	return abstract
}
