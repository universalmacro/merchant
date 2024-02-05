package services

import (
	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/merchant/dao/repositories"
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
