package repository

import (
	"github.com/ofavor/ddd-go/pkg/repo"
	"github.com/ofavor/kratos-layout/internal/domain/entity"
	"github.com/ofavor/kratos-layout/internal/infrastructure/repo/dao"
)

type GreeterFilter struct {
	Id   []int64
	Name string
}

func (f *GreeterFilter) Conditions() map[string]interface{} {
	conds := make(map[string]interface{})
	if len(f.Id) > 0 {
		conds["id"] = f.Id
	}
	if f.Name != "" {
		conds["name"] = f.Name
	}
	return conds
}

type GreeterRepo interface {
	repo.Repository[entity.Greeter, dao.GreeterDao]
}
