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

func (this *ProviderAbstract) SetContainer(container container.Container) {
	this.Container = container
}

func (this *ProviderAbstract) SetDbFactory(DbFactory *database.DatabaseFactory) {
	this.DbFactory = DbFactory
}

func (this *ProviderAbstract) RegisterAutoPanic(name string, resolver interface{}) {
	err := this.Container.NamedSingleton(name, resolver)
	if err != nil {
		panic(err)
	}
}
