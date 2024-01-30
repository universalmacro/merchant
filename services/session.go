package services

import (
	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/merchant/dao/repositories"
)

var sessionServiceSingleton = singleton.NewSingleton(newSessionServices, singleton.Eager)

func GetSessionService() *SessionService {
	return sessionServiceSingleton.Get()
}

func newSessionServices() *SessionService {
	return &SessionService{
		merchantRepository: repositories.GetMerchantRepository(),
	}
}

type SessionService struct {
	merchantRepository *repositories.MerchantRepository
}

func (s *SessionService) CreateSession(account, password string, shortMerchantId *string) (string, error) {
	if shortMerchantId == nil {
		return s.CreateMerchantSession(account, password)
	}
	return s.CreateSubAccountSession(account, password, *shortMerchantId)
}

func (s *SessionService) CreateMerchantSession(account, password string) (string, error) {
	return "", nil
}
func (s *SessionService) CreateSubAccountSession(account, password, shortMerchantId string) (string, error) {
	return "", nil
}
