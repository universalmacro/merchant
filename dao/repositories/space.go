package repositories

import (
	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/ioc"
)

var spaceRepository = singleton.SingletonFactory(func() *SpaceRepository {
	return &SpaceRepository{
		Repository: dao.NewRepository[entities.Space](ioc.GetDBInstance()),
	}
}, singleton.Lazy)

func GetSpaceRepository() *SpaceRepository {
	return spaceRepository.Get()
}

type SpaceRepository struct {
	*dao.Repository[entities.Space]
}
