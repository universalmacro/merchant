package entities

import (
	"github.com/universalmacro/common/auth"
	"gorm.io/gorm"
)

type Merchant struct {
	gorm.Model
	Name string `gorm:"type:varchar(255);"`
}

type Account struct {
	gorm.Model
	auth.Password
}
