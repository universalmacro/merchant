package services

import (
	"errors"

	"github.com/universalmacro/common/dao"
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
	table, _ := tableRepo.List(
		dao.Where("space_id = ?", s.ID()),
		dao.Where("label = ?", label),
	)
	if len(table) > 0 {
		return nil, errors.New("tableLabel already exists")
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

func (s *Space) GetTable(label string) *Table {
	tableRepo := repositories.GetTableRepository()
	table, _ := tableRepo.List(
		dao.Where("space_id = ?", s.ID()),
		dao.Where("label = ?", label),
	)
	if len(table) == 0 {
		return nil
	}
	return &Table{&table[0]}
}

func (s *Space) ListTables() []Table {
	tableRepo := repositories.GetTableRepository()
	tables, _ := tableRepo.List(
		dao.Where("space_id = ?", s.ID()),
	)
	result := make([]Table, len(tables))
	for i := range tables {
		result[i] = Table{&tables[i]}
	}
	return result
}

func (s *Space) CreateFood(name string, description string, price, fixedOffset *int64) {

}
