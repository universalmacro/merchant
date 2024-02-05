package entities

import (
	"github.com/universalmacro/common/auth"
	"github.com/universalmacro/common/dao/data"
	"github.com/universalmacro/common/snowflake"
	"gorm.io/gorm"
)

type Merchant struct {
	gorm.Model
	Account         string `gorm:"type:varchar(255);uniqueIndex;"`
	ShortMerchantId string `gorm:"type:varchar(255);uniqueIndex;"`
	auth.Password
	*data.PhoneNumber
	Name        string `gorm:"type:varchar(255);"`
	Description string `gorm:"type:varchar(255);"`
}

var merchantIdGenerator = snowflake.NewIdGenertor(0)

func (a *Merchant) BeforeCreate(tx *gorm.DB) (err error) {
	a.Model.ID = merchantIdGenerator.Uint()
	return err
}

type SubAccount struct {
	gorm.Model
	MerchantId uint   `gorm:"index:unique,composite:merchantId_account"`
	Account    string `gorm:"type:varchar(255);index:unique,composite:merchantId_account"`
	auth.Password
	*data.PhoneNumber
}
