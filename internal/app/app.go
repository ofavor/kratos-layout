package app

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewGreeterAppService,
	NewMyEventAppService,
	// TODO: add new service here
)
