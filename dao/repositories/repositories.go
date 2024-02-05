package repositories

import (
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/ioc"
)

func init() {
	ioc.GetDBInstance().AutoMigrate(&entities.Merchant{},
		&entities.SubAccount{},
		&entities.VerificationCode{},
		&entities.Space{})
}
