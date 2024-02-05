package services

import (
	"errors"

	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/merchant/dao/repositories"
	"github.com/universalmacro/merchant/ioc"
	"github.com/universalmacro/merchant/services/models"
)

var sessionServiceSingleton = singleton.SingletonFactory(newSessionServices, singleton.Eager)

func GetSessionService() *SessionService {
	return sessionServiceSingleton.Get()
}

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

func (self *SessionService) TokenVerification(token string) models.Account {
	claims, err := ioc.GetJwtSigner().VerifyJwt(token)
	if err != nil {
		return nil
	}
	t, ok := claims["type"].(string)
	if !ok {
		return nil
	}
	if t == "MAIN" {
		merchant := self.merchantService.GetMerchantByAccount(claims["merchantId"].(string))
		return merchant
	} else {
		return nil
	}
}
