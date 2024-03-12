package repositories

import (
	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/ioc"
)

var GetPrinterRepository = singleton.EagerSingleton(func() *PrinterRepository {
	return &PrinterRepository{
		Repository: dao.NewRepository[entities.Printer](ioc.GetDBInstance()),
	}
})

type PrinterRepository struct {
	*dao.Repository[entities.Printer]
}
