package services

import (
	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/dao/repositories"
)

var tableServiceSingleton = singleton.SingletonFactory(newTableService, singleton.Eager)

func GetTableService() *TableService {
	return tableServiceSingleton.Get()
}

func newTableService() *TableService {
	return &TableService{
		tableService: repositories.GetTableRepository(),
	}
}

type TableService struct {
	tableService *repositories.TableRepository
}

func (self *TableService) GetTable(tableId uint) *Table {
	t, _ := self.tableService.GetById(tableId)
	if t == nil {
		return nil
	}
	return &Table{Table: t}
}

type Table struct {
	*entities.Table
}

func (t *Table) ID() uint {
	return t.Table.ID
}

func (t *Table) Label() string {
	return t.Table.Label
}

func (t *Table) SetLabel(label string) *Table {
	t.Table.Label = label
	return t
}

func (t *Table) Submit() *Table {
	repo := repositories.GetTableRepository()
	repo.Save(t.Table)
	return t
}

func (t *Table) Delete() {
	repo := repositories.GetTableRepository()
	repo.Delete(t.Table)
}

func (t *Table) SpaceID() uint {
	return t.Table.SpaceID
}

func (t *Table) GetSpace() *Space {
	space := GetSpaceService().GetSpace(t.SpaceID())
	return space
}

func (t *Table) Granted(account Account) bool {
	return t.GetSpace().Granted(account)
}
