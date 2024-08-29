package application

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewGreeterAppService,
	// TODO: add new service here
)
