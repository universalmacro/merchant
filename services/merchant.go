package services

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/common/utils"
	"github.com/universalmacro/common/utils/random"
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/dao/repositories"
	"github.com/universalmacro/merchant/ioc"
)

var GetMerchantService = singleton.EagerSingleton(func() *MerchantService {
	return &MerchantService{
		merchantRepository:   repositories.GetMerchantRepository(),
		subAccountRepository: repositories.GetSubAccountRepository(),
	}
})

type MerchantService struct {
	merchantRepository   *repositories.MerchantRepository
	subAccountRepository *repositories.SubAccountRepository
}

func (s *MerchantService) GetMerchant(merchantId uint) *Merchant {
	merchant, _ := s.merchantRepository.GetById(merchantId)
	if merchant == nil {
		return nil
	}
	return &Merchant{Entity: merchant}
}

func (s *MerchantService) CreateMerchant(shortMerchantId, account, password string) *Merchant {
	merchant := &entities.Merchant{
		ShortMerchantId: shortMerchantId,
		Account:         account,
	}
	merchant.SetPassword(password)
	_, ctx := s.merchantRepository.Create(merchant)
	if ctx.RowsAffected == 0 {
		return nil
	}
	return &Merchant{Entity: merchant}
}

func (m *MerchantService) CreateVerificationCode(merchantId uint, countryCode, phoneNumber string) {
	db := ioc.GetDBInstance()
	db = dao.ApplyOptions(
		db,
		dao.Where("merchant_id = ?", merchantId),
		dao.Where("country_code = ?", countryCode),
		dao.Where("number = ?", phoneNumber),
		dao.Where("created_at > ?", time.Now().Add(-time.Minute*10)))
	var verificationCode entities.VerificationCode
	ctx := db.Find(&verificationCode)
	if ctx.RowsAffected == 0 {
		verificationCode = entities.VerificationCode{
			MerchantId:  merchantId,
			CountryCode: countryCode,
			Number:      phoneNumber,
			Code:        random.RandomNumberString(6),
		}
		db.Create(&verificationCode)
		return
	}
}

func (s *MerchantService) ListMerchants(index, limit int64) dao.List[Merchant] {
	merchantList, _ := s.merchantRepository.Pagination(index, limit)
	result := make([]Merchant, len(merchantList.Items))
	for i := range merchantList.Items {
		result[i] = Merchant{Entity: &merchantList.Items[i]}
	}
	var list dao.List[Merchant]
	list.Items = result
	list.Pagination = merchantList.Pagination
	return list
}

func (s *MerchantService) GetMerchantByAccount(account string) *Merchant {
	merchant, _ := s.merchantRepository.GetByAccount(account)
	if merchant == nil {
		return nil
	}
	return &Merchant{Entity: merchant}
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

func (m *Merchant) Currency() string {
	return m.Entity.Currency
}

func (m *Merchant) SetCurrency(currency string) *Merchant {
	m.Entity.Currency = currency
	return m
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

func (m *Merchant) ListSpaces() []Space {
	list, _ := repositories.GetSpaceRepository().FindMany("merchant_id = ?", m.ID())
	spaces := make([]Space, len(list))
	for i := range list {
		spaces[i] = Space{&list[i]}
	}
	return spaces
}

func (m *Merchant) Submit() {
	repositories.GetMerchantRepository().Update(m.Entity)
}

type SubAccount struct {
	Entity *entities.SubAccount
}

func (s *SubAccount) ID() uint {
	return s.Entity.ID
}

func (s *SubAccount) PasswordMatching(password string) bool {
	return s.Entity.PasswordMatching(password)
}

func (s *SubAccount) UpdatePassword(password string) {
	s.Entity.SetPassword(password)
	repositories.GetSubAccountRepository().Update(s.Entity)
}
