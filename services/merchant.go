package services

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/common/utils"
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/dao/repositories"
	_ "github.com/universalmacro/merchant/dao/repositories"
	"github.com/universalmacro/merchant/ioc"
	"github.com/universalmacro/merchant/services/models"
)

func GetMerchantService() *MerchantService {
	return merchantSingleton.Get()
}

var merchantSingleton = singleton.SingletonFactory(newMerchantService, singleton.Lazy)

func newMerchantService() *MerchantService {
	return &MerchantService{
		merchantRepository:   repositories.GetMerchantRepository(),
		subAccountRepository: repositories.GetSubAccountRepository(),
	}
}

type MerchantService struct {
	merchantRepository   *repositories.MerchantRepository
	subAccountRepository *repositories.SubAccountRepository
}

func (s *MerchantService) GetMerchant(merchantId uint) *models.Merchant {
	merchant, _ := s.merchantRepository.GetById(merchantId)
	if merchant == nil {
		return nil
	}
	return &models.Merchant{Entity: merchant}
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
	for i := range merchantList.Items {
		result[i] = models.Merchant{Entity: &merchantList.Items[i]}
	}
	var list dao.List[models.Merchant]
	list.Items = result
	list.Pagination = merchantList.Pagination
	return list
}

func (s *MerchantService) GetMerchantByAccount(account string) *models.Merchant {
	merchant, _ := s.merchantRepository.GetByAccount(account)
	if merchant == nil {
		return nil
	}
	return &models.Merchant{Entity: merchant}
}

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
		MerchantID:     m.StringID(),
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

type SubAccount struct {
	Entity *entities.SubAccount
}

func (s *SubAccount) ID() uint {
	return s.Entity.ID
}

func (s *SubAccount) PasswordMatching(password string) bool {
	return s.PasswordMatching(password)
}

func (s *SubAccount) UpdatePassword(password string) {
	s.Entity.SetPassword(password)
	repositories.GetSubAccountRepository().Update(s.Entity)
}
