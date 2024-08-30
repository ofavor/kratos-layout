package interfaces

import (
	"context"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/google/wire"
	nacsdkcl "github.com/nacos-group/nacos-sdk-go/clients"
	nacsdkco "github.com/nacos-group/nacos-sdk-go/common/constant"
	nacsdkvo "github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/ofavor/kratos-layout/internal/conf"
	etcdsdk "go.etcd.io/etcd/client/v3"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(
	// NewRegistrar,
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

func parseNacosEndpoints(conf string) []nacsdkco.ServerConfig {
	addrs := strings.Split(conf, ",")
	ret := make([]nacsdkco.ServerConfig, 0, len(addrs))
	for _, addr := range addrs {
		vv := strings.Split(addr, ":")
		port, _ := strconv.ParseUint(vv[1], 10, 64)
		ret = append(ret, nacsdkco.ServerConfig{
			IpAddr: vv[0],
			Port:   port,
		})
	}
	return ret
}

func NewRegistrar(conf *conf.Registry) registry.Registrar {
	switch conf.Type {
	case "etcd":
		client, err := etcdsdk.New(etcdsdk.Config{
			Endpoints: strings.Split(conf.Etcd.Endpoints, ","),
		})
		if err != nil {
			panic(err)
		}
		return etcd.New(client)
	case "nacos":
		addrs := parseNacosEndpoints(conf.Nacos.Endpoints)
		client, err := nacsdkcl.NewNamingClient(
			nacsdkvo.NacosClientParam{
				ServerConfigs: addrs,
			},
		)
		if err != nil {
			panic(err)
		}
		return nacos.New(client)
	default:
		panic("unknown registry type")
	}
}
