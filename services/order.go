package services

import (
	"fmt"

	"github.com/universalmacro/common/utils"
	"github.com/universalmacro/merchant/dao/entities"
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

func (o *Order) FoodSpec() {
	var foods []Food
	for i := range o.Order.Foods {
		foods = append(foods, Food{&o.Order.Foods[i].Food})
	}
}
