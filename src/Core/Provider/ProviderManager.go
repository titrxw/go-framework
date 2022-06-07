package provider

import (
	"github.com/golobby/container/v3/pkg/container"
	database "github.com/titrxw/go-framework/src/Core/Database"
)

type ProviderManager struct {
	Container container.Container
	DbFactory *database.DatabaseFactory
}

func (this *ProviderManager) MakeProvider(abstract ProviderInterface) ProviderInterface {
	abstract.SetContainer(this.Container)
	abstract.SetDbFactory(this.DbFactory)
	return abstract
}
