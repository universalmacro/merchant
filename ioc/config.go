package ioc

import (
	"github.com/spf13/viper"
	"github.com/universalmacro/common/singleton"
)

var GetConfig = singleton.EagerSingleton(func() *viper.Viper {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigName(".env")
	v.SetConfigType("yml")
	v.ReadInConfig()
	return v
})
