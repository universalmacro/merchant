package services

import (
	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/dao/repositories"
)

type Bill struct {
	*entities.Bill
}

func (b *Bill) Granted(account Account) bool {
	return b.MerchantId == account.MerchantId()
}

func (b *Bill) Orders() []Order {
	var orders []Order
	orderRepo := repositories.GetOrderRepository()
	os, _ := orderRepo.List(dao.Where("bill_id = ?", b.ID))
	for i := range os {
		orders = append(orders, Order{&os[i]})
	}
	return orders
}

func (b *Bill) Print() {

}
