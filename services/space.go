package services

import (
	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/merchant/dao/repositories"
	"github.com/universalmacro/merchant/services/models"
)

var spaceServiceSingleton = singleton.SingletonFactory(newSpaceService, singleton.Eager)

func GetSpaceService() *SpaceService {
	return spaceServiceSingleton.Get()
}

func newSpaceService() *SpaceService {
	return &SpaceService{
		spaceRepository: repositories.GetSpaceRepository(),
	}
}

type SpaceService struct {
	spaceRepository *repositories.SpaceRepository
}

func (self *SpaceService) GetSpace(spaceId uint) *models.Space {
	s, _ := self.spaceRepository.GetById(spaceId)
	if s == nil {
		return nil
	}
	return &models.Space{Space: s}
}
