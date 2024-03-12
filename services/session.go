package services

import (
	"errors"
	"strings"

	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/common/utils"
	"github.com/universalmacro/merchant/dao/repositories"
	"github.com/universalmacro/merchant/ioc"
)

var GetSessionService = singleton.EagerSingleton(newSessionServices)

func newSessionServices() *SessionService {
	return &SessionService{
		merchantService:    GetMerchantService(),
		merchantRepository: repositories.GetMerchantRepository(),
	}
}

type SessionService struct {
	merchantService    *MerchantService
	merchantRepository *repositories.MerchantRepository
}

func (s *SessionService) CreateSession(account, password string, shortMerchantId *string) (string, error) {
	if shortMerchantId == nil {
		return s.CreateMerchantSession(account, password)
	}
	return s.CreateSubAccountSession(account, password, *shortMerchantId)
}

var ErrAccountNotFound = errors.New("account not found")
var ErrPasswordNotMatch = errors.New("password not match")

func (s *SessionService) CreateMerchantSession(account, password string) (string, error) {
	merchant := s.merchantService.GetMerchantByAccount(account)
	if merchant == nil {
		return "", ErrAccountNotFound
	}
	if !merchant.PasswordMatching(password) {
		return "", ErrPasswordNotMatch
	}
	return merchant.CreateSession(), nil
}
func (s *SessionService) CreateSubAccountSession(account, password, shortMerchantId string) (string, error) {
	return "", nil
}

func (s *SessionService) TokenVerification(token string) (Account, error) {
	sp := strings.Split(token, " ")
	if len(sp) != 2 {
		return nil, errors.New("invalid token")
	}
	jwt := sp[1]
	claims, err := ioc.GetJwtSigner().VerifyJwt(jwt)
	if err != nil {
		return nil, err
	}
	t, ok := claims["type"].(string)
	if !ok {
		return nil, errors.New("invalid token")
	}
	if t == "MAIN" {
		id, _ := claims["merchantId"].(string)
		merchant := s.merchantService.GetMerchant(utils.StringToUint(id))
		return merchant, nil
	} else {
		return nil, errors.New("invalid token")
	}
}
