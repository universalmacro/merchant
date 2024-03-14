package services

import (
	"github.com/Dparty/feieyun"
	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/common/utils"
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/dao/repositories"
	"github.com/universalmacro/merchant/ioc"
)

type Printer struct {
	*entities.Printer
}

func (p *Printer) ID() uint {
	return p.Printer.ID
}

func (p *Printer) StringID() string {
	return utils.UintToString(p.ID())
}

func (p *Printer) Delete() {
	repositories.GetPrinterRepository().Delete(p.Printer)
}

func (p *Printer) SpaceID() uint {
	return p.Printer.SpaceID
}

func (p *Printer) Space() *Space {
	return GetSpaceService().GetSpace(p.SpaceID())
}

func (p *Printer) Granted(account Account) bool {
	return p.Space().Granted(account)
}

func (p *Printer) Submit() *Printer {
	repositories.GetPrinterRepository().Update(p.Printer)
	return p
}

func (p *Printer) Print(content feieyun.PrintContent) {
	printerFactory := ioc.GetPrinterFactory()
	printer, _ := printerFactory.Connect(p.Sn)
	printer.Print(content.String(), "")
}

var GetPrinterService = singleton.EagerSingleton(func() *PrinterService {
	return &PrinterService{
		printerRepository: repositories.GetPrinterRepository(),
	}
})

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

func (s *PrinterService) List(options ...dao.Option) []Printer {
	printers, _ := s.printerRepository.List(options...)
	var result []Printer
	for i := range printers {
		result = append(result, Printer{&printers[i]})
	}
	return result
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
