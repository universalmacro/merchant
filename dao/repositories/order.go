package repositories

import (
	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/ioc"
)

var GetOrderRepository = singleton.EagerSingleton(func() *OrderRepository {
	return &OrderRepository{
		Repository: dao.NewRepository[entities.Order](ioc.GetDBInstance()),
	}
})

type OrderRepository struct {
	*dao.Repository[entities.Order]
}
