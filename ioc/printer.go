package ioc

import (
	"github.com/Dparty/feieyun"
	"github.com/universalmacro/common/singleton"
)

var GetPrinterFactory = singleton.EagerSingleton(func() *feieyun.PrinterFactory {
	config := GetConfig()
	return &feieyun.PrinterFactory{
		Config: feieyun.FeieyunConfig{
			User: config.GetString("feieyun.user"),
			Ukey: config.GetString("feieyun.ukey"),
			Url:  config.GetString("feieyun.url"),
		},
	}
})
