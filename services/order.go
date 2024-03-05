package services

import (
	"fmt"

	"github.com/universalmacro/common/utils"
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/dao/repositories"
)

type Order struct {
	*entities.Order
}

func (o *Order) StringID() string {
	return utils.UintToString(o.ID)
}

func (o *Order) Code() string {
	code := o.PickUpCode % 1000
	return fmt.Sprintf("%d", code)
}

func (o *Order) FoodSpec() []FoodSpec {
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

func (o *Order) CancelItems() {

}

func (o *Order) AddItems(foods []FoodSpec) {

}

func (o *Order) Submit() *Order {
	repositories.GetOrderRepository().Save(o.Order)
	return o
}
