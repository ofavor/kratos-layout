//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/ofavor/kratos-layout/internal/app"
	"github.com/ofavor/kratos-layout/internal/conf"
	"github.com/ofavor/kratos-layout/internal/iface"
	"github.com/ofavor/kratos-layout/internal/infra"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Components, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(infra.ProviderSet, iface.ProviderSet, app.ProviderSet, newApp))
}
