package repositories

import (
	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/ioc"
	"gorm.io/gorm"
)

var merchantRepository = singleton.SingletonFactory(func() *MerchantRepository {
	return &MerchantRepository{
		Repository: dao.NewRepository[entities.Merchant](ioc.GetDBInstance()),
	}
}, singleton.Lazy)

func GetMerchantRepository() *MerchantRepository {
	return merchantRepository.Get()
}

type MerchantRepository struct {
	*dao.Repository[entities.Merchant]
}

func (m *MerchantRepository) GetByAccount(account string) (*entities.Merchant, *gorm.DB) {
	return m.FindOne("account = ?", account)
}

var subAccountRepository = singleton.SingletonFactory(func() *SubAccountRepository {
	return &SubAccountRepository{}
}, singleton.Lazy)

func GetSubAccountRepository() *SubAccountRepository {
	return subAccountRepository.Get()
}

type SubAccountRepository struct {
	*dao.Repository[entities.SubAccount]
}
