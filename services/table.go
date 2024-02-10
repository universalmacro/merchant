package services

import (
	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/merchant/dao/repositories"
	"github.com/universalmacro/merchant/services/models"
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

func (self *TableService) GetTable(tableId uint) *models.Table {
	t, _ := self.tableService.GetById(tableId)
	if t == nil {
		return nil
	}
	return &models.Table{Table: t}
}
