package infrastructure

import (
	"strings"

	"strconv"

	"github.com/google/wire"
	"github.com/ofavor/ddd-go/pkg/cache"
	caredis "github.com/ofavor/ddd-go/pkg/cache/redis"
	"github.com/ofavor/ddd-go/pkg/db"
	dbgorm "github.com/ofavor/ddd-go/pkg/db/gorm"
	"github.com/ofavor/ddd-go/pkg/event"
	evtkafka "github.com/ofavor/ddd-go/pkg/event/kafka"
	evtmem "github.com/ofavor/ddd-go/pkg/event/memory"
	evtredis "github.com/ofavor/ddd-go/pkg/event/redis"
	"github.com/ofavor/kratos-layout/internal/conf"
	"github.com/ofavor/kratos-layout/internal/infrastructure/repo"
	"github.com/ofavor/kratos-layout/internal/infrastructure/repo/dao"

	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2/registry"
	nacsdkcl "github.com/nacos-group/nacos-sdk-go/clients"
	nacsdkco "github.com/nacos-group/nacos-sdk-go/common/constant"
	nacsdkvo "github.com/nacos-group/nacos-sdk-go/vo"
	etcdsdk "go.etcd.io/etcd/client/v3"
)

// ProviderSet is infra providers.
var ProviderSet = wire.NewSet(
	// NewRegistrar,
	NewDatabase,
	NewCache,
	NewEvent,
	NewInfra,
	repo.NewGreeterRepo,
	// ddd-go AUTO GENERATE SLOT, DO NOT UPDATE/DELETE new repo
	// TODO: add new infrastructure component here
)

type Infra struct {
	db    db.Database
	cache cache.Cache
	event event.EventBus
}

func NewDatabase(c *conf.Bootstrap) db.Database {
	dc := c.Components.Database
	return dbgorm.NewDatabase(dc.Driver, dc.Dsn, dc.EncKey, strings.ToLower(c.Logging.Level) == "debug")
}

func NewCache(c *conf.Bootstrap) cache.Cache {
	rc := c.Components.Redis
	return caredis.NewCache(rc.Addr, rc.Password, rc.Db, rc.Prefix)
}

func NewEvent(c *conf.Bootstrap) event.EventBus {
	ec := c.Components.Event
	if ec == nil {
		return nil
	}
	switch ec.Type {
	case "kafka":
		kc := c.Components.Kafka
		return evtkafka.NewEventBus(kc.Brokers, ec.BufferSize, ec.Group)
	case "redis":
		rc := c.Components.Redis
		return evtredis.NewEventBus(rc.Addr, rc.Password, rc.Db, ec.BufferSize, ec.Group)
	case "memory":
		return evtmem.NewEventBus(ec.BufferSize)
	}
	return nil
}

func NewInfra(db db.Database, cache cache.Cache, event event.EventBus) *Infra {
	return &Infra{
		db:    db,
		cache: cache,
		event: event,
	}
}

func (i *Infra) Initialize() error {
	// DB
	i.db.RegisterModels([]interface{}{
		&dao.GreeterDao{},
		// ddd-go AUTO GENERATE SLOT, DO NOT UPDATE/DELETE new dao
		// TODO add new dao model here
	})

	// TODO
	return nil
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
