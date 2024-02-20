package services

import (
	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/dao/repositories"
)

type Printer struct {
	*entities.Printer
}

func (self *Printer) Delete() {
	repositories.GetPrinterRepository().Delete(self.Printer)
}

func (self *Printer) SpaceID() uint {
	return self.Printer.SpaceID
}

func (self *Printer) Space() *Space {
	return GetSpaceService().GetSpace(self.SpaceID())
}

func (self *Printer) Granted(account Account) bool {
	return self.Space().Granted(account)
}

func GetPrinterService() *MerchantService {
	return merchantSingleton.Get()
}

var printerSingleton = singleton.SingletonFactory(newMerchantService, singleton.Lazy)

func newPrinterService() *MerchantService {
	return &MerchantService{
		merchantRepository:   repositories.GetMerchantRepository(),
		subAccountRepository: repositories.GetSubAccountRepository(),
	}
}

type PrinterService struct {
	printerRepository *repositories.PrinterRepository
}

func (s *PrinterService) GetPrinter(printerId uint) *Printer {
	printer, _ := s.printerRepository.GetById(printerId)
	if printer == nil {
		return nil
	}
	return &Printer{printer}
}

func (s *PrinterService) CreatePrinter(name, sn string) *Printer {
	printer := &entities.Printer{
		Name: name,
	}
	_, ctx := s.printerRepository.Create(printer)
	if ctx.RowsAffected == 0 {
		return nil
	}
	return &Printer{printer}
}
