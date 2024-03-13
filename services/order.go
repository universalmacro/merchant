package services

import (
	"errors"
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

func (os *OrderService) CreateBill(ac Account, amount uint, orderIds ...uint) (*Bill, error) {
	return createBillHelper(os.db, true, ac, amount, orderIds...)
}

func (os *OrderService) PrintBill(ac Account, amount uint, orderIds ...uint) (*Bill, error) {
	bill, err := createBillHelper(os.db, false, ac, amount, orderIds...)
	if err != nil {
		return nil, err
	}
	bill.Print()
	return bill, err
}

func (os *OrderService) GetBill(id uint) *Bill {
	db := ioc.GetDBInstance()
	var billEntity entities.Bill
	db.Find(&billEntity, id)
	return &Bill{&billEntity}
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

func (o *Order) CancelItem(food FoodSpec) (*Order, error) {
	if o.Status != "SUBMITTED" {
		return nil, fmt.Errorf("order status is not submitted")
	}
	foodSpecs := o.FoodSpecs()
	for i := range foodSpecs {
		if foodSpecs[i].Equals(food) {
			o.Order.Foods = append(o.Order.Foods[:i], o.Order.Foods[i+1:]...)
			break
		}
	}
	return o, nil
}

func (o *Order) CancelItems(foods ...FoodSpec) (*Order, error) {
	for i := range foods {
		_, err := o.CancelItem(foods[i])
		if err != nil {
			return o, err
		}
	}
	if len(o.Order.Foods) == 0 {
		o.Order.Status = "CANCELLED"
	}
	return o, nil
}

func (o *Order) AddItem(food FoodSpec) (*Order, error) {
	if o.Status != "SUBMITTED" {
		return o, fmt.Errorf("order status is not submitted")
	}
	o.Order.Foods = append(o.Order.Foods, entities.FoodSpec{
		Food: *food.Food.Food,
		Spec: food.Spec.Spec,
	})
	return o, nil
}

func (o *Order) AddItems(foods ...FoodSpec) (*Order, error) {
	for i := range foods {
		_, err := o.AddItem(foods[i])
		if err != nil {
			return o, err
		}
	}
	return o, nil
}

func (o *Order) Submit() *Order {
	repositories.GetOrderRepository().Save(o.Order)
	return o
}

func (o *Order) Space() *Space {
	return GetSpaceService().GetSpace(o.Order.SpaceID)
}

func createBillHelper(db *gorm.DB, submit bool, ac Account, amount uint, orderIds ...uint) (*Bill, error) {
	var orderEntities []entities.Order
	db.Find(&orderEntities, dao.Where("id IN (?)", orderIds))
	if len(orderEntities) == 0 {
		return nil, errors.New("order ids is empty")
	}
	space := GetSpaceService().GetSpace(orderEntities[0].SpaceID)
	if !space.Granted(ac) {
		return nil, errors.New("permission denied")
	}
	billEntity := entities.Bill{
		MerchantId: ac.MerchantId(),
		CashierID:  ac.ID(),
		Amount:     amount,
		SpaceID:    space.ID(),
	}
	if submit {
		db = db.Begin()
	}
	db.Create(&billEntity)
	bill := &Bill{&billEntity}
	for i, orderEntity := range orderEntities {
		if orderEntity.Status != "SUBMITTED" {
			db.Rollback()
			return nil, errors.New("order status is not submitted")
		}
		if space.ID() != orderEntity.SpaceID {
			db.Rollback()
			return nil, errors.New("order is not in the same space as the bill")
		}
		orderEntities[i].BillId = &billEntity.ID
		orderEntities[i].Status = "COMPLETED"
		err := db.Save(&orderEntities[i]).Error
		if err != nil {
			db.Rollback()
			return nil, errors.New("create order failed")
		}
	}
	if submit {
		db.Commit()
	}
	return bill, nil
}
