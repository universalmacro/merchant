package entities

import (
	"github.com/universalmacro/common/snowflake"
	"gorm.io/gorm"
)

type Member struct {
	gorm.Model
	MerchantId  uint   `json:"merchantId" gorm:"index"`
	CountryCode string `json:"countryCode" gorm:"type:CHAR(6);index,composite:phone_number"`
	PhoneNumber string `json:"number" gorm:"type:CHAR(11);index,composite:phone_number;column:phone_number"`
}

func (a *Member) BeforeCreate(tx *gorm.DB) (err error) {
	a.Model.ID = snowflake.NewIdGenertor(0).Uint()
	return err
}
