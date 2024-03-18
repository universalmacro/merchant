package repositories

import (
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/ioc"
)

func init() {
	db := ioc.GetDBInstance()
	db.AutoMigrate(&entities.Merchant{})
	db.AutoMigrate(&entities.SubAccount{})
	db.AutoMigrate(&entities.VerificationCode{})
	db.AutoMigrate(&entities.Space{})
	db.AutoMigrate(&entities.Table{})
	db.AutoMigrate(&entities.Food{})
	db.AutoMigrate(&entities.Printer{})
	db.AutoMigrate(&entities.Order{})
	db.AutoMigrate(&entities.Bill{})
	db.AutoMigrate(&entities.Member{})
}
