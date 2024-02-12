package repositories

import (
	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/ioc"
)

var foodRepositorySingleton = singleton.SingletonFactory(newFoodRepository, singleton.Lazy)

func GetFoodRepository() *FoodRepository {
	return foodRepositorySingleton.Get()
}

func newFoodRepository() *FoodRepository {
	return &FoodRepository{
		Repository: dao.NewRepository[entities.Food](ioc.GetDBInstance()),
	}
}

type FoodRepository struct {
	*dao.Repository[entities.Food]
}
