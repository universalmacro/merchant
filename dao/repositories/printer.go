package repositories

import (
	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/ioc"
)

var printerRepository = singleton.SingletonFactory(func() *PrinterRepository {
	return &PrinterRepository{
		Repository: dao.NewRepository[entities.Printer](ioc.GetDBInstance()),
	}
}, singleton.Lazy)

func GetPrinterRepository() *PrinterRepository {
	return printerRepository.Get()
}

type PrinterRepository struct {
	*dao.Repository[entities.Printer]
}
