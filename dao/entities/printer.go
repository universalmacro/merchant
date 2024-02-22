package entities

import (
	"github.com/universalmacro/common/snowflake"
	"gorm.io/gorm"
)

type Printer struct {
	SpaceAsset
	Name  string `gorm:"type:varchar(64)"`
	Sn    string `gorm:"type:varchar(64)"`
	Type  string `gorm:"type:varchar(64)"`
	Model string `gorm:"type:varchar(64)"`
}

func (a *Printer) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = snowflake.NewIdGenertor(0).Uint()
	return err
}
