//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/ofavor/kratos-layout/internal/biz"
	"github.com/ofavor/kratos-layout/internal/conf"
	"github.com/ofavor/kratos-layout/internal/data"
	"github.com/ofavor/kratos-layout/internal/server"
	"github.com/ofavor/kratos-layout/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
