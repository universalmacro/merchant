package models

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/common/utils"
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/dao/repositories"
	"github.com/universalmacro/merchant/ioc"
)

type Merchant struct {
	Entity *entities.Merchant
}

func (m *Merchant) MerchantId() uint {
	return m.ID()
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

func (m *Merchant) Granted(account Account) bool {
	return m.ID() == account.MerchantId()
}

func (m *Merchant) CreateSession() string {
	sessionId := sessionIdGenerator.String()
	expired := time.Now().Add(time.Hour * 24 * 7).Unix()
	claims := Claims{
		Type:           "MAIN",
		ID:             sessionId,
		MerchantID:     m.ID(),
		StandardClaims: jwt.StandardClaims{Id: sessionId, ExpiresAt: expired},
	}
	jwt, _ := ioc.GetJwtSigner().SignJwt(claims)
	return jwt
}

func (m *Merchant) CreateSpace(name string) *Space {
	entity, _ := repositories.GetSpaceRepository().Create(&entities.Space{MerchantId: m.ID(), Name: name})
	return &Space{entity}
}

func (m *Merchant) ListSpaces() dao.List[Space] {
	list, _ := repositories.GetSpaceRepository().Pagination(0, 10000)
	spaces := make([]Space, len(list.Items))
	for i := range list.Items {
		spaces[i] = Space{&list.Items[i]}
	}
	return dao.List[Space]{Items: spaces, Pagination: list.Pagination}
}
