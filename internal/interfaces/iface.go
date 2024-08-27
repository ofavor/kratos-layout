package interfaces

import (
	"strings"

	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/google/wire"
	"github.com/ofavor/kratos-layout/internal/conf"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(
	NewRegistrar,
	NewGRPCServer,
	NewHTTPServer,
	NewEventHandler,
	// TODO: add new interface service here
)

func NewRegistrar(conf *conf.Registry) registry.Registrar {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: strings.Split(conf.Etcd.Endpoints, ","),
	})
	if err != nil {
		panic(err)
	}
	r := etcd.New(client)
	return r
}
