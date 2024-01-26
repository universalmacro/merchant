package services

import (
	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/dao/repositories"
	_ "github.com/universalmacro/merchant/dao/repositories"
	"github.com/universalmacro/merchant/services/models"
)

func GetMerchantService() *MerchantService {
	return merchantSingleton.Get()
}

var merchantSingleton = singleton.NewSingleton(newMerchantService, singleton.Lazy)

func newMerchantService() *MerchantService {
	return &MerchantService{
		merchantRepository: repositories.GetMerchantRepository(),
		accountRepository:  repositories.GetAccountRepository(),
	}
}

type MerchantService struct {
	merchantRepository *repositories.MerchantRepository
	accountRepository  *repositories.AccountRepository
}

func (s *MerchantService) CreateMerchant(shortMerchantId, account, password string) *models.Merchant {
	merchant := &entities.Merchant{
		ShortMerchantId: shortMerchantId,
		Account:         account,
	}
	merchant.SetPassword(password)
	_, ctx := s.merchantRepository.Create(merchant)
	if ctx.RowsAffected == 0 {
		return nil
	}
	return &models.Merchant{Entity: merchant}
}

func (s *MerchantService) ListMerchants(index, limit int64) dao.List[models.Merchant] {
	merchantList, _ := s.merchantRepository.Pagination(index, limit)
	result := make([]models.Merchant, len(merchantList.Items))
	for i, merchant := range merchantList.Items {
		result[i] = models.Merchant{Entity: &merchant}
	}
	var list dao.List[models.Merchant]
	list.Items = result
	list.Pagination = merchantList.Pagination
	return list
}
