// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/ofavor/kratos-layout/internal/application"
	"github.com/ofavor/kratos-layout/internal/conf"
	"github.com/ofavor/kratos-layout/internal/infrastructure"
	"github.com/ofavor/kratos-layout/internal/infrastructure/repo"
	"github.com/ofavor/kratos-layout/internal/interfaces"
	"go.opentelemetry.io/otel/sdk/trace"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(server *conf.Server, registry *conf.Registry, components *conf.Components, auth *conf.Auth, logger log.Logger, tracerProvider *trace.TracerProvider) (*kratos.App, func(), error) {
	eventBus := infrastructure.NewEvent(components)
	database := infrastructure.NewDatabase(components)
	greeterRepo := repo.NewGreeterRepo(database)
	greeterAppService := application.NewGreeterAppService(logger, eventBus, greeterRepo)
	grpcServer := interfaces.NewGRPCServer(logger, tracerProvider, server, auth, greeterAppService)
	httpServer := interfaces.NewHTTPServer(logger, tracerProvider, server, auth, greeterAppService)
	registrar := interfaces.NewRegistrar(registry)
	myEventAppService := application.NewMyEventAppService(logger)
	eventHandler := interfaces.NewEventHandler(logger, eventBus, myEventAppService)
	cache := infrastructure.NewCache(components)
	infra := infrastructure.NewInfra(database, cache, eventBus)
	app := newApp(logger, grpcServer, httpServer, registrar, eventHandler, infra)
	return app, func() {
	}, nil
}
