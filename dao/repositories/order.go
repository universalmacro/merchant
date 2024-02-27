package repositories

import (
	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/ioc"
)

var orderRepositorySingleton = singleton.SingletonFactory(newOrderRepository, singleton.Lazy)

func GetOrderRepository() *OrderRepository {
	return orderRepositorySingleton.Get()
}

func newOrderRepository() *OrderRepository {
	return &OrderRepository{
		Repository: dao.NewRepository[entities.Order](ioc.GetDBInstance()),
	}
}

type OrderRepository struct {
	*dao.Repository[entities.Order]
}
