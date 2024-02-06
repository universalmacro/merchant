package entities

import (
	"github.com/universalmacro/common/snowflake"
	"gorm.io/gorm"
)

type SpaceAsset struct {
	gorm.Model
	SpaceID uint `gorm:"index"`
}

type Space struct {
	gorm.Model
	MerchantId uint   `gorm:"index"`
	Name       string `gorm:"type:varchar(255);"`
}

var spaceIdGenerator = snowflake.NewIdGenertor(0)

func (a *Space) BeforeCreate(tx *gorm.DB) (err error) {
	a.Model.ID = spaceIdGenerator.Uint()
	return err
}
