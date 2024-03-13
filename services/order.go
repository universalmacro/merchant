package services

import (
	"fmt"

	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/common/utils"
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/dao/repositories"
	"github.com/universalmacro/merchant/ioc"
	"gorm.io/gorm"
)

var GetOrderService = singleton.EagerSingleton(func() *OrderService {
	return &OrderService{
		db:        ioc.GetDBInstance(),
		orderRepo: repositories.GetOrderRepository(),
	}
})

type OrderService struct {
	db        *gorm.DB
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

func (os *OrderService) CreateOrder(ac Account, spcesId uint, amount uint, orderId ...uint) *Bill {
	space := GetSpaceService().GetSpace(spcesId)
	if space == nil {
		return nil
	}
	if !space.Granted(ac) {
		return nil
	}
	orderEntities, _ := os.orderRepo.List(dao.Where("space_id = ?", spcesId), dao.Where("id IN (?)", orderId))
	billEntity := entities.Bill{
		CashierID: ac.ID(),
		Amount:    amount,
	}
	db := os.db.Begin()
	db.Create(&billEntity)
	bill := &Bill{&billEntity}
	for i := range orderEntities {
		orderEntities[i].BillId = billEntity.ID
		err := db.Save(&orderEntities[i]).Error
		if err != nil {
			db.Rollback()
			return nil
		}
	}
	db.Commit()
	return bill
}

type Order struct {
	*entities.Order
}

func (o *Order) Granted(account Account) bool {
	return account.MerchantId() == o.Space().MerchantId
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
	if len(o.Order.Foods) == 0 {
		o.Order.Status = "CANCELLED"
	}
	return o
}

func (o *Order) AddItem(food FoodSpec) *Order {
	o.Order.Foods = append(o.Order.Foods, entities.FoodSpec{
		Food: *food.Food.Food,
		Spec: food.Spec.Spec,
	})
	return o
}

func (o *Order) AddItems(foods ...FoodSpec) *Order {
	for i := range foods {
		o.AddItem(foods[i])
	}
	return o
}

func (o *Order) Submit() *Order {
	repositories.GetOrderRepository().Save(o.Order)
	return o
}

func (o *Order) Space() *Space {
	return GetSpaceService().GetSpace(o.Order.SpaceID)
}

type Bill struct {
	*entities.Bill
}
