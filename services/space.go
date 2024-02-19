package services

import (
	"errors"

	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/common/utils"
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/dao/repositories"
)

var spaceServiceSingleton = singleton.SingletonFactory(newSpaceService, singleton.Eager)

func GetSpaceService() *SpaceService {
	return spaceServiceSingleton.Get()
}

func newSpaceService() *SpaceService {
	return &SpaceService{
		spaceRepository: repositories.GetSpaceRepository(),
	}
}

type SpaceService struct {
	spaceRepository *repositories.SpaceRepository
}

func (self *SpaceService) GetSpace(spaceId uint) *Space {
	s, _ := self.spaceRepository.GetById(spaceId)
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
