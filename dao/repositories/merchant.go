package repositories

import (
	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/ioc"
)

var merchantRepository = singleton.NewSingleton(func() *MerchantRepository {
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

var accountRepository = singleton.NewSingleton(func() *AccountRepository {
	return &AccountRepository{}
}, singleton.Lazy)

func GetAccountRepository() *AccountRepository {
	return accountRepository.Get()
}

type AccountRepository struct {
	*dao.Repository[entities.SubAccount]
}
