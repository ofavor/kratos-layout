package infrastructure

import (
	"strings"

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
)

// ProviderSet is infra providers.
var ProviderSet = wire.NewSet(
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
	return dbgorm.NewDatabase(c.Components.Database.Driver, c.Components.Database.Dns, c.Components.Database.EncKey, strings.ToLower(c.Logging.Level) == "debug")
}

func NewCache(c *conf.Bootstrap) cache.Cache {
	return caredis.NewCache(c.Components.Redis.Addr, c.Components.Redis.Password, c.Components.Redis.Db, c.Components.Redis.Prefix)
}

func NewEvent(c *conf.Bootstrap) event.EventBus {
	if c.Components.Event == nil {
		return nil
	}
	switch c.Components.Event.Type {
	case "kafka":
		return evtkafka.NewEventBus(c.Components.Kafka.Brokers, c.Components.Event.BufferSize, c.Components.Event.Group)
	case "redis":
		return evtredis.NewEventBus(c.Components.Redis.Addr, c.Components.Redis.Password, c.Components.Redis.Db, c.Components.Event.BufferSize, c.Components.Event.Group)
	case "memory":
		return evtmem.NewEventBus(c.Components.Event.BufferSize)
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
