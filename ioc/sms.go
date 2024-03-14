package ioc

import (
	"github.com/universalmacro/common/singleton"
	"github.com/universalmacro/common/sms/tencent"
)

var GetSmsSender = singleton.EagerSingleton(func() *tencent.SmsSender {
	config := GetConfig()
	return tencent.NewSmsSender(
		config.GetString("tencent.region"),
		config.GetString("tencent.secretId"),
		config.GetString("tencent.secretKey"),
	)
})
