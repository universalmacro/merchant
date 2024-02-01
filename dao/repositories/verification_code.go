package repositories

import (
	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/merchant/dao/entities"
)

var verificationCodeRepository = singleton.SingletonFactory(func() *VerificationCodeRepository {
	return &VerificationCodeRepository{}
}, singleton.Lazy)

func GetVerificationCodeRepository() *VerificationCodeRepository {
	return verificationCodeRepository.Get()
}

type VerificationCodeRepository struct {
	*dao.Repository[entities.VerificationCode]
}

func (r *VerificationCodeRepository) FindByPhone(country_code, phone string) *entities.VerificationCode {
	var verificationCode entities.VerificationCode
	if ctx := r.DB.Model(&verificationCode).Where("country_code = ?", country_code).Where("number = ?", phone).Order("craeted_at DESC").First(&verificationCode); ctx.Error != nil {
		return nil
	}
	return &verificationCode
}
