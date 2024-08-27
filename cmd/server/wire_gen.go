// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/ofavor/kratos-layout/internal/app"
	"github.com/ofavor/kratos-layout/internal/conf"
	"github.com/ofavor/kratos-layout/internal/iface"
	"github.com/ofavor/kratos-layout/internal/infra"
	"github.com/ofavor/kratos-layout/internal/infra/repo"
	"go.opentelemetry.io/otel/sdk/trace"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(server *conf.Server, registry *conf.Registry, components *conf.Components, auth *conf.Auth, logger log.Logger, tracerProvider *trace.TracerProvider) (*kratos.App, func(), error) {
	eventBus := infra.NewEvent(components)
	database := infra.NewDatabase(components)
	greeterRepo := repo.NewGreeterRepo(database)
	greeterAppService := app.NewGreeterAppService(logger, eventBus, greeterRepo)
	grpcServer := iface.NewGRPCServer(logger, tracerProvider, server, auth, greeterAppService)
	httpServer := iface.NewHTTPServer(logger, tracerProvider, server, auth, greeterAppService)
	registrar := iface.NewRegistrar(registry)
	myEventAppService := app.NewMyEventAppService(logger)
	eventHandler := iface.NewEventHandler(logger, eventBus, myEventAppService)
	cache := infra.NewCache(components)
	infraInfra := infra.NewInfra(database, cache, eventBus)
	kratosApp := newApp(logger, grpcServer, httpServer, registrar, eventHandler, infraInfra)
	return kratosApp, func() {
	}, nil
}
