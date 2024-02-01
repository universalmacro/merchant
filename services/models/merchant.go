package models

import (
	"time"

	"github.com/universalmacro/common/utils"
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/dao/repositories"
)

type Merchant struct {
	Entity *entities.Merchant
}

func (m *Merchant) ID() uint {
	return m.Entity.ID
}

func (m *Merchant) StringID() string {
	return utils.UintToString(m.Entity.ID)
}

func (m *Merchant) Account() string {
	return m.Entity.Account
}

func (m *Merchant) ShortMerchantId() string {
	return m.Entity.ShortMerchantId
}

func (m *Merchant) Name() string {
	return m.Entity.Name
}

func (m *Merchant) Description() string {
	return m.Entity.Description
}

func (m *Merchant) CreatedAt() time.Time {
	return m.Entity.CreatedAt
}

func (m *Merchant) UpdatedAt() time.Time {
	return m.Entity.UpdatedAt
}

func (m *Merchant) PasswordMatching(password string) bool {
	return m.Entity.PasswordMatching(password)
}

func (m *Merchant) UpdatePassword(password string) {
	m.Entity.SetPassword(password)
	repo := repositories.GetMerchantRepository()
	repo.Update(m.Entity)
}

func (m *Merchant) Verification(countryCode, phoneNumber, code string) bool {
	repo := repositories.GetVerificationCodeRepository()
	verificationCode := repo.FindByPhone(countryCode, phoneNumber)
	if verificationCode == nil {
		return false
	}
	if verificationCode.CreatedAt.Add(time.Minute * 10).After(time.Now()) {
		return false
	}
	verificationCode.Tries++
	repo.Update(verificationCode)
	if verificationCode.Tries >= 10 {
		return false
	}
	return verificationCode.Code == code
}

func (m *Merchant) CreateMember() {

}

func (m *Merchant) CreateSession() string {
	return ""
}
