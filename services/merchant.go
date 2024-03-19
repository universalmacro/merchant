package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/common/sms/models"
	"github.com/universalmacro/common/sms/tencent"
	"github.com/universalmacro/common/utils"
	"github.com/universalmacro/common/utils/random"
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/dao/repositories"
	"github.com/universalmacro/merchant/ioc"
)

const MAX_TRIES = 10

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

func (ms *MerchantService) SignupMember(merchantId uint, countryCode, phoneNumber, code string) (string, error) {
	merchant := ms.GetMerchant(merchantId)
	if merchant == nil {
		return "", errors.New("merchant not found")
	}
	verificationCode := ms.GetVerificationCode(merchantId, countryCode, phoneNumber)
	if verificationCode == nil {
		return "", errors.New("verification code not found")
	}
	if verificationCode.Tries > MAX_TRIES {
		return "", errors.New("verification code has been tried too many times")
	}
	verificationCode.Tries++
	db := ioc.GetDBInstance()
	db.Save(verificationCode)
	if verificationCode.Code != code {
		return "", errors.New("verification code not matching")
	}
	if verificationCode.Tries >= 10 {
		return "", errors.New("verification code has been tried too many times")
	}
	db.Delete(verificationCode)
	var member entities.Member
	ctx := db.Find(&member, "merchant_id = ? AND country_code = ? AND phone_number = ?", merchantId, countryCode, phoneNumber)
	if ctx.RowsAffected > 0 {
		m := Member{member}
		return m.GenerateToken(), nil
	}
	member = entities.Member{
		MerchantId:  merchantId,
		CountryCode: countryCode,
		PhoneNumber: phoneNumber,
	}
	db.Create(&member)
	m := Member{member}
	token := m.GenerateToken()
	return token, nil
}

var ErrVerificationCodeHasBeenSent = errors.New("verification code has been sent")

func (m *MerchantService) GetVerificationCode(merchantId uint, countryCode, phoneNumber string) *entities.VerificationCode {
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
		return nil
	}
	return &verificationCode
}

func (m *MerchantService) CreateVerificationCode(merchantId uint, countryCode, phoneNumber string) error {
	merchant := m.GetMerchant(merchantId)
	verificationCode := m.GetVerificationCode(merchantId, countryCode, phoneNumber)
	code := random.RandomNumberString(6)
	if verificationCode != nil {
		return ErrVerificationCodeHasBeenSent
	}
	verificationCode = &entities.VerificationCode{
		MerchantId:  merchantId,
		CountryCode: countryCode,
		Number:      phoneNumber,
		Code:        code,
	}
	db := ioc.GetDBInstance()
	db.Create(verificationCode)
	smsSender := ioc.GetSmsSender()
	config := ioc.GetConfig()
	err := smsSender.SendWithConfig(models.PhoneNumber{
		AreaCode: countryCode,
		Number:   phoneNumber,
	}, tencent.Config{
		TemplateId: config.GetString("tencent.sms.templateId"),
		AppId:      config.GetString("tencent.sms.appId"),
	}, []string{code, "創建" + merchant.Name() + "會員帳號"})
	if err != nil {
		return errors.New("sms send failed")
	}
	return nil
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
