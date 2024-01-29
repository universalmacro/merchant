package entities

import "gorm.io/gorm"

type VerificationCode struct {
	gorm.Model
	MerchantId  uint   `gorm:"index,composite:phone_number"`
	CountryCode string `json:"countryCode" gorm:"type:CHAR(6);index,composite:phone_number"`
	Number      string `json:"number" gorm:"type:CHAR(11);index,composite:phone_number"`
	Code        string `json:"code" gorm:"type:CHAR(6)"`
	Tries       int64  `json:"tries"`
}
