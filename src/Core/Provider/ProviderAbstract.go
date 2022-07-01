package provider

import (
	"github.com/golobby/container/v3/pkg/container"
	database "github.com/titrxw/go-framework/src/Core/Database"
)

type ProviderAbstract struct {
	ProviderInterface
	Container container.Container
	Config    interface{}
	DbFactory *database.DatabaseFactory
}

func (providerAbstract *ProviderAbstract) SetContainer(container container.Container) {
	providerAbstract.Container = container
}

func (providerAbstract *ProviderAbstract) SetDbFactory(DbFactory *database.DatabaseFactory) {
	providerAbstract.DbFactory = DbFactory
}

func (providerAbstract *ProviderAbstract) RegisterAutoPanic(name string, resolver interface{}) {
	err := providerAbstract.Container.NamedSingleton(name, resolver)
	if err != nil {
		panic(err)
	}
}
