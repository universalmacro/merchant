package models

import (
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/dao/repositories"
)

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
