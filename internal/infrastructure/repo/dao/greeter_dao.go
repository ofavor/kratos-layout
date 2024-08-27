package dao

import (
	"github.com/ofavor/kratos-layout/internal/domain/vo"
	"gorm.io/gorm"
)

type GreeterDao struct {
	gorm.Model
	Name      string     `gorm:"type:varchar(50);not null;default:'';comment:Name"`
	Greeting  string     `gorm:"type:varchar(50);not null;default:'';comment:Greeting"`
	Address   string     `gorm:"type:varchar(1024);comment:Address"`
	AddressVo vo.Address `gorm:"-"`
}

func (d *GreeterDao) TableName() string {
	return "hw_gretter"
}

func (d *GreeterDao) BeforeSave(tx *gorm.DB) error {
	d.Address = d.AddressVo.ToJson()
	return nil
}

func (d *GreeterDao) AfterFind(tx *gorm.DB) error {
	d.AddressVo = vo.NewAddressFromJson(d.Address)
	return nil
}
