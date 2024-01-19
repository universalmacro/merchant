package entities

import "gorm.io/gorm"

type Merchant struct {
	gorm.Model
	Name string `gorm:"type:varchar(255);"`
}
