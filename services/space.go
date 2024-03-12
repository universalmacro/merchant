package services

import (
	"errors"

	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/common/utils"
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/dao/repositories"
)

var GetSpaceService = singleton.EagerSingleton(newSpaceService)

func newSpaceService() *SpaceService {
	return &SpaceService{
		spaceRepository: repositories.GetSpaceRepository(),
	}
}

type SpaceService struct {
	spaceRepository *repositories.SpaceRepository
}

func (sp *SpaceService) GetSpace(spaceId uint) *Space {
	s, _ := sp.spaceRepository.GetById(spaceId)
	if s == nil {
		return nil
	}
	return &Space{Space: s}
}

type Space struct {
	*entities.Space
}

func (s *Space) ID() uint {
	return s.Space.ID
}

func (s *Space) StringID() string {
	return utils.UintToString(s.ID())
}

func (s *Space) Children() []Space {
	children, _ := repositories.GetSpaceRepository().FindMany("parent_id = ?", s.ID())
	result := make([]Space, len(children))
	for i := range children {
		result[i] = Space{&children[i]}
	}
	return result
}

func (s *Space) Granted(account Account) bool {
	return account.MerchantId() == s.MerchantId
}

func (s *Space) Name() string {
	return s.Space.Name
}

func (s *Space) SetName(name string) *Space {
	s.Space.Name = name
	return s
}

func (s *Space) Submit() *Space {
	repo := repositories.GetSpaceRepository()
	repo.Save(s.Space)
	return s
}

func (s *Space) Delete() {
	repo := repositories.GetSpaceRepository()
	repo.Delete(s.Space)
}

func (s *Space) CreateTable(label string) (*Table, error) {
	tableRepo := repositories.GetTableRepository()
	table, _ := tableRepo.FindMany("space_id = ? AND label = ?", s.ID(), label)
	if len(table) > 0 {
		return nil, errors.New("table label already exists")
	}
	t := &entities.Table{
		SpaceAsset: entities.SpaceAsset{
			SpaceID: s.ID(),
		},
		Label: label,
	}
	tableRepo.Save(t)
	return &Table{t}, nil
}

func (s *Space) ListTables() []Table {
	tableRepo := repositories.GetTableRepository()
	tables, _ := tableRepo.FindMany("space_id = ?", s.ID())
	result := make([]Table, len(tables))
	for i := range tables {
		result[i] = Table{&tables[i]}
	}
	return result
}

func (s *Space) Foods() []Food {
	foodRepo := repositories.GetFoodRepository()
	foods, _ := foodRepo.FindMany("space_id = ?", s.ID())
	result := make([]Food, len(foods))
	for i := range foods {
		result[i] = Food{&foods[i]}
	}
	return result
}

func (s *Space) CreateFood(food *Food) (*Food, error) {
	food.SpaceID = s.ID()
	food.Create()
	return food, nil
}

func (s *Space) CreatePrinter(name, sn, printerType, model string) *Printer {
	printer := &entities.Printer{
		SpaceAsset: entities.SpaceAsset{
			SpaceID: s.ID(),
		},
		Name:  name,
		Sn:    sn,
		Type:  printerType,
		Model: model,
	}
	repositories.GetPrinterRepository().Save(printer)
	return &Printer{printer}
}

func (s *Space) ListPrinters() []Printer {
	printers, _ := repositories.GetPrinterRepository().FindMany("space_id = ?", s.ID())
	result := make([]Printer, len(printers))
	for i := range printers {
		result[i] = Printer{&printers[i]}
	}
	return result
}

func (s *Space) FoodCategories() []string {
	return s.Space.FoodCategories
}

func (s *Space) SetFoodCategories(categories ...string) *Space {
	mapCategories := make(map[string]struct{})
	var foodCategories dao.StringArray
	for _, category := range categories {
		if _, ok := mapCategories[category]; ok {
			continue
		}
		mapCategories[category] = struct{}{}
		foodCategories = append(foodCategories, category)
	}
	s.Space.FoodCategories = foodCategories
	return s
}

func (s *Space) CreateOrder(account Account, tableLabel *string, foods []FoodSpec) Order {
	var foodSpace entities.FoodSpces
	for i := range foods {
		foodSpace = append(foodSpace, entities.FoodSpec{
			Food: *foods[i].Food.Food,
			Spec: foods[i].Spec.Spec,
		})
	}
	order := &entities.Order{
		SpaceAsset: entities.SpaceAsset{
			SpaceID: s.ID(),
		},
		TableLabel: tableLabel,
		Status:     "SUBMITTED",
		Foods:      foodSpace,
	}
	orderRepo := repositories.GetOrderRepository()
	orders := GetOrderService().List(
		dao.Where("space_id = ?", s.ID()),
		dao.OrderBy("created_at", true),
		dao.Limit(1),
	)
	var pickUpCode int64
	if len(orders) > 0 {
		pickUpCode = orders[0].PickUpCode + 1
	}
	order.PickUpCode = pickUpCode
	orderRepo.Save(order)
	o := Order{order}
	o.PrintCashier()
	o.PrintKitchen()
	return o
}

func (space *Space) ListOrders() []Order {
	orderRepo := repositories.GetOrderRepository()
	orders, _ := orderRepo.FindMany("space_id = ?", space.ID())
	result := make([]Order, len(orders))
	for i := range orders {
		result[i] = Order{&orders[i]}
	}
	return result
}
