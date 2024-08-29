//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/ofavor/kratos-layout/internal/application"
	"github.com/ofavor/kratos-layout/internal/conf"
	"github.com/ofavor/kratos-layout/internal/infrastructure"
	"github.com/ofavor/kratos-layout/internal/interfaces"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

// wireApp init kratos application.
func wireApp(*conf.Bootstrap, *conf.Registry, log.Logger, *tracesdk.TracerProvider) (*kratos.App, func(), error) {
	panic(wire.Build(infrastructure.ProviderSet, interfaces.ProviderSet, application.ProviderSet, newApp))
}
