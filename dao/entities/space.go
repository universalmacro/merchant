package entities

import "gorm.io/gorm"

type Space struct {
	gorm.Model
	MerchantId uint   `gorm:"index"`
	Name       string `gorm:"type:varchar(255);"`
}
