package entities

import (
	"github.com/universalmacro/common/auth"
	"github.com/universalmacro/common/dao/data"
	"gorm.io/gorm"
)

type Merchant struct {
	gorm.Model
	Account         string `gorm:"type:varchar(255);uniqueIndex;"`
	ShortMerchantId string `gorm:"type:varchar(255);uniqueIndex;"`
	Name            string `gorm:"type:varchar(255);"`
	Description     string `gorm:"type:varchar(255);"`
}

type Account struct {
	gorm.Model
	MerchantId uint
	auth.Password
	*data.PhoneNumber
}
