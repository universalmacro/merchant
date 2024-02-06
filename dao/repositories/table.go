package repositories

import (
	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/ioc"
)

var tableRepository = singleton.SingletonFactory(func() *TableRepository {
	return &TableRepository{
		Repository: dao.NewRepository[entities.Table](ioc.GetDBInstance()),
	}
}, singleton.Lazy)

func GetTableRepository() *TableRepository {
	return tableRepository.Get()
}

type TableRepository struct {
	*dao.Repository[entities.Table]
}
