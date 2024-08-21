package app

import (
	"github.com/google/wire"
	"github.com/ofavor/kratos-layout/internal/infra"
)

var ProviderSet = wire.NewSet(infra.NewInfra, NewGreeterAppService, NewMyEventAppService)
