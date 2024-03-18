package entities

import (
	"github.com/universalmacro/common/snowflake"
	"gorm.io/gorm"
)

type Member struct {
	gorm.Model
	MerchantId  uint   `json:"merchantId" gorm:"index:,unique,composite:member_key"`
	CountryCode string `json:"countryCode" gorm:"type:CHAR(6);index:,unique,composite:member_key"`
	PhoneNumber string `json:"number" gorm:"type:CHAR(11);index:,unique,composite:member_key"`
}

func (a *Member) BeforeCreate(tx *gorm.DB) (err error) {
	a.Model.ID = snowflake.NewIdGenertor(0).Uint()
	return err
}
