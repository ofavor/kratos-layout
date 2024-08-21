package repo

import (
	"github.com/ofavor/ddd-go/pkg/db"
	repogorm "github.com/ofavor/ddd-go/pkg/repo/gorm"
	"github.com/ofavor/kratos-layout/internal/domain/entity"
	"github.com/ofavor/kratos-layout/internal/domain/repository"
	"github.com/ofavor/kratos-layout/internal/infra/repo/dao"
	"gorm.io/gorm"
)

type greeterRepo struct {
	repogorm.GormRepo[entity.Greeter, dao.GreeterDao]
}

func NewGreeterRepo(db db.Database) repository.GreeterRepo {
	return &greeterRepo{
		GormRepo: *repogorm.NewRepo(
			db.GetConn().(*gorm.DB),
			func(d *dao.GreeterDao) entity.Greeter {
				return entity.LoadGreeter(d)
			},
		),
	}
}
