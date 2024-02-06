package models

import (
	"github.com/universalmacro/common/utils"
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/dao/repositories"
)

type Space struct {
	*entities.Space
}

func (s *Space) ID() uint {
	return s.Space.ID
}

func (s *Space) StringID() string {
	return utils.UintToString(s.ID())
}

func (s *Space) Granted(account Account) bool {
	return account.MerchantId() == s.MerchantId
}

func (s *Space) Name() string {
	return s.Space.Name
}

func (s *Space) SetName(name string) *Space {
	s.Space.Name = name
	return s
}

func (s *Space) Submit() *Space {
	repo := repositories.GetSpaceRepository()
	repo.Save(s.Space)
	return s
}

func (s *Space) Delete() {
	repo := repositories.GetSpaceRepository()
	repo.Delete(s.Space)
}
