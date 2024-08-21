package repository

import (
	"github.com/ofavor/ddd-go/pkg/repo"
	"github.com/ofavor/kratos-layout/internal/domain/entity"
	"github.com/ofavor/kratos-layout/internal/infra/repo/dao"
)

type GreeterRepo interface {
	repo.Repository[entity.Greeter, dao.GreeterDao]
}
