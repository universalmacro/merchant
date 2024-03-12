package repositories

import (
	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/ioc"
)

var GetSpaceRepository = singleton.EagerSingleton(func() *SpaceRepository {
	return &SpaceRepository{
		Repository: dao.NewRepository[entities.Space](ioc.GetDBInstance()),
	}
})

type SpaceRepository struct {
	*dao.Repository[entities.Space]
}
