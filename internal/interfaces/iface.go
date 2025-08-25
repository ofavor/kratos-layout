package interfaces

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/google/wire"
	"github.com/ofavor/kratos-layout/internal/conf"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(
	NewGRPCServer,
	NewHTTPServer,
	NewEventHandler,
	// TODO: add new interface service here
)

func newAuthWhiteListMatcher(bc *conf.Bootstrap) selector.MatchFunc {
	whiteList := make(map[string]struct{})
	for _, v := range bc.Auth.Ignores {
		whiteList[v] = struct{}{}
	}
	// whiteList["/helloworld.v1.Greeter/Create"] = struct{}{}
	// whiteList["/helloworld.v1.Greeter/SayHello"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}
