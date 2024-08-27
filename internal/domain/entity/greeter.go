package entity

import (
	"github.com/ofavor/ddd-go/pkg/entity"
	"github.com/ofavor/kratos-layout/internal/domain/vo"
	"github.com/ofavor/kratos-layout/internal/infrastructure/repo/dao"
)

type HelloworldId uint

type Greeter interface {
	entity.Entity[dao.GreeterDao]
	GetId() HelloworldId
	GetName() string
	GetGreeting() string
	GetAddress() vo.Address

	SayHello() string
	ChangeGreeting(string) error
}

type greeter struct {
	dao *dao.GreeterDao
}

func NewGreeter(name, greeting string, address vo.Address) (Greeter, error) {
	g := &greeter{
		dao: &dao.GreeterDao{
			Name:      name,
			Greeting:  greeting,
			AddressVo: address,
		},
	}
	return g, nil
}

func LoadGreeter(data *dao.GreeterDao) Greeter {
	return &greeter{
		dao: data,
	}
}

func (h *greeter) IsNew() bool {
	return h.dao.ID == 0
}

func (h *greeter) DAO() *dao.GreeterDao {
	return h.dao
}

func (h *greeter) GetId() HelloworldId {
	return HelloworldId(h.dao.ID)
}

func (h *greeter) GetGreeting() string {
	return h.dao.Greeting
}

func (h *greeter) GetName() string {
	return h.dao.Name
}

func (h *greeter) GetAddress() vo.Address {
	return h.dao.AddressVo
}

func (h *greeter) SayHello() string {
	return h.dao.Name + " says " + h.dao.Greeting
}

func (h *greeter) ChangeGreeting(greeting string) error {
	h.dao.Greeting = greeting
	return nil
}
