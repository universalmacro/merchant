package repositories

import (
	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/ioc"
	"gorm.io/gorm"
)

var GetMerchantRepository = singleton.EagerSingleton(func() *MerchantRepository {
	return &MerchantRepository{
		Repository: dao.NewRepository[entities.Merchant](ioc.GetDBInstance()),
	}
})

type MerchantRepository struct {
	*dao.Repository[entities.Merchant]
}

func (m *MerchantRepository) GetByAccount(account string) (*entities.Merchant, *gorm.DB) {
	return m.FindOne("account = ?", account)
}

var GetSubAccountRepository = singleton.EagerSingleton(
	func() *SubAccountRepository {
		return &SubAccountRepository{}
	})

type SubAccountRepository struct {
	*dao.Repository[entities.SubAccount]
}
