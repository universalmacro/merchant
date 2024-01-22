package services

import (
	_ "github.com/universalmacro/merchant/dao/repositories"
	"github.com/universalmacro/merchant/ioc"
	"gorm.io/gorm"
)

func newMerchantService(db *gorm.DB) *MerchantService {
	return &MerchantService{db: ioc.GetDBInstance()}
}

type MerchantService struct {
	db *gorm.DB
}
