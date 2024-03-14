package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/Dparty/feieyun"
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
	bill, err := createBillHelper(os.db, true, ac, amount, orderIds...)
	return bill, err
}

func (os *OrderService) PrintBill(ac Account, amount uint, orderIds ...uint) (*Bill, error) {
	bill, err := createBillHelper(os.db, false, ac, amount, orderIds...)
	if err != nil {
		return nil, err
	}
	bill.Print()
	return bill, err
}

func (os *OrderService) ListBills(options ...dao.Option) []Bill {
	db := ioc.GetDBInstance()
	var billEntities []entities.Bill
	db = dao.ApplyOptions(db, options...)
	db.Find(&billEntities)
	var result []Bill
	for i := range billEntities {
		result = append(result, Bill{&billEntities[i]})
	}
	return result
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
	groups := foodSpecsGroupHelper(o.FoodSpecs()...)
	timestring := time.Now().Add(time.Hour * 8).Format("2006-01-02 15:04")
	var tableLabel string
	if o.TableLabel != nil {
		tableLabel = *o.TableLabel
	} else {
		tableLabel = "無桌號"
	}
	for _, group := range groups {
		var pc feieyun.PrintContent
		pc.AddLines(&feieyun.CenterBold{Content: &feieyun.Text{Content: fmt.Sprintf("餐號: %s", o.Code())}},
			&feieyun.CenterBold{Content: &feieyun.Text{Content: fmt.Sprintf("桌號: %s", tableLabel)}},
			&feieyun.Bold{Content: &feieyun.Text{Content: fmt.Sprintf("%s X%d", group.Name(), group.Amount)}})
		for _, option := range group.FoodSpec.Spec.Spec {
			pc.AddLines(&feieyun.Bold{Content: &feieyun.Text{Content: fmt.Sprintf("-  %sX%d", option.Optioned, group.Amount)}})
		}
		pc.AddLines(&feieyun.Text{Content: timestring})
		for _, printer := range group.Printers() {
			printer.Print(pc)
		}
	}
}

func (o *Order) PrintCashier() {
	// printers := o.Space().ListPrinters(dao.Where("type = ?", "CASHIER"))
	// for _, printer := range printers {
	// groups := foodSpecsGroupHelper(o.FoodSpecs()...)
	// timestring := time.Now().Add(time.Hour * 8).Format("2006-01-02 15:04")
	// var printContent feieyun.PrintContent
	// printContent.AddLines(&feieyun.CenterBold{Content: &feieyun.Text{Content: o.Space().Name()}})
	// printContent.AddLines(&feieyun.CenterBold{Content: &feieyun.Text{Content: fmt.Sprintf("餐號: %d", bill.PickUpCode)}})
	// printContent.AddLines(&feieyun.CenterBold{Content: &feieyun.Text{Content: fmt.Sprintf("桌號: %s", table.Label)}})
	// printContent.AddDiv(p.Width())
	// }
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

type FoodSpecGroup struct {
	FoodSpec
	Amount uint
}

func foodSpecsGroupHelper(foods ...FoodSpec) []FoodSpecGroup {
	var groups []FoodSpecGroup
	for _, food := range foods {
		if len(groups) == 0 {
			groups = append(groups, FoodSpecGroup{food, 1})
		} else {
			for i, group := range groups {
				if group.FoodSpec.Equals(food) {
					groups[i].Amount++
					break
				}
				if i == len(groups)-1 {
					groups = append(groups, FoodSpecGroup{food, 1})
				}
			}
		}
	}
	return groups
}

func createBillHelper(db *gorm.DB, completed bool, ac Account, amount uint, orderIds ...uint) (*Bill, error) {
	var orderEntities []entities.Order
	db.Find(&orderEntities, orderIds)
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
	if completed {
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
		if completed {
			orderEntities[i].Status = "COMPLETED"
		}
		err := db.Save(&orderEntities[i]).Error
		if err != nil {
			db.Rollback()
			return nil, err
		}
	}
	if completed {
		db.Commit()
	}
	return bill, nil
}
