package services

import (
	"fmt"

	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/common/utils"
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/dao/repositories"
)

func GetOrderService() *OrderService {
	return orderServiceSingleton.Get()
}

var orderServiceSingleton = singleton.SingletonFactory(NewOrderService, singleton.Eager)

func NewOrderService() *OrderService {
	return &OrderService{
		orderRepo: repositories.GetOrderRepository(),
	}
}

type OrderService struct {
	orderRepo *repositories.OrderRepository
}

func (os *OrderService) GetById(id uint) *Order {
	order, _ := os.orderRepo.GetById(id)
	if order == nil {
		return nil
	}
	return &Order{order}
}

func (os *OrderService) List(options ...dao.Option) []Order {
	orders, _ := os.orderRepo.List(options...)
	var result []Order
	for i := range orders {
		result = append(result, Order{&orders[i]})
	}
	return result
}

type Order struct {
	*entities.Order
}

func (o *Order) SetTableLabel(label string) *Order {
	o.Order.TableLabel = &label
	return o
}

func (o *Order) StringID() string {
	return utils.UintToString(o.ID)
}

func (o *Order) Code() string {
	code := o.PickUpCode % 1000
	return fmt.Sprintf("%d", code)
}

func (o *Order) FoodSpecs() []FoodSpec {
	var foods []FoodSpec
	for i := range o.Order.Foods {
		foods = append(foods, NewFoodSpec(o.Order.Foods[i]))
	}
	return foods
}

func (o *Order) PrintKitchen() {

}

func (o *Order) PrintCashier() {

}

func (o *Order) CancelItem(food FoodSpec) *Order {
	foodSpecs := o.FoodSpecs()
	for i := range foodSpecs {
		if foodSpecs[i].Equals(food) {
			o.Order.Foods = append(o.Order.Foods[:i], o.Order.Foods[i+1:]...)
			break
		}
	}
	return o
}

func (o *Order) CancelItems(foods ...FoodSpec) *Order {
	for i := range foods {
		o.CancelItem(foods[i])
	}
	return o
}

func (o *Order) AddItems(foods ...FoodSpec) {

}

func (o *Order) Submit() *Order {
	repositories.GetOrderRepository().Save(o.Order)
	return o
}

func (o *Order) Space() *Space {
	return GetSpaceService().GetSpace(o.Order.SpaceID)
}
