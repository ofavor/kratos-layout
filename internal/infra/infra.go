package infra

import (
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
	"github.com/ofavor/kratos-layout/internal/domain/repository"
	"github.com/ofavor/kratos-layout/internal/infra/repo"
	"github.com/ofavor/kratos-layout/internal/infra/repo/dao"
)

// ProviderSet is infra providers.
var ProviderSet = wire.NewSet(NewDatabase, NewCache, NewEvent, repo.NewGreeterRepo)

type Infra struct {
	db    db.Database
	cache cache.Cache
	event event.EventBus

	greeterRepo repository.GreeterRepo
}

func NewDatabase(c *conf.Components) db.Database {
	return dbgorm.NewDatabase(c.Database.Driver, c.Database.Dns, c.Database.EncKey, c.Database.Debug)
}

func NewCache(c *conf.Components) cache.Cache {
	return caredis.NewCache(c.Redis.Addr, c.Redis.Password, c.Redis.Db, c.Redis.Prefix)
}

func NewEvent(c *conf.Components) event.EventBus {
	if c.Event == nil {
		return nil
	}
	switch c.Event.Type {
	case "kafka":
		return evtkafka.NewEventBus(c.Kafka.Brokers, c.Event.BufferSize, c.Event.Group)
	case "redis":
		return evtredis.NewEventBus(c.Redis.Addr, c.Redis.Password, c.Redis.Db, c.Event.BufferSize, c.Event.Group)
	case "memory":
		return evtmem.NewEventBus(c.Event.BufferSize)
	}
	return nil
}

func NewInfra(db db.Database, cache cache.Cache, event event.EventBus, greeterRepo repository.GreeterRepo) *Infra {
	return &Infra{
		db:    db,
		cache: cache,
		event: event,

		greeterRepo: greeterRepo,
	}
}

func (i *Infra) Initialize() error {
	// DB
	i.db.RegisterModels([]interface{}{
		&dao.GreeterDao{},
	})

	// TODO
	return nil
}

func (i *Infra) GetDatabase() db.Database {
	return i.db
}

func (i *Infra) GetCache() cache.Cache {
	return i.cache
}

func (i *Infra) GetEvent() event.EventBus {
	return i.event
}

func (i *Infra) GetGreeterRepo() repository.GreeterRepo {
	return i.greeterRepo
}
