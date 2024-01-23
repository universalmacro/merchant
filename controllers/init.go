package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/universalmacro/common/server"
	api "github.com/universalmacro/merchant-api-interfaces"
	_ "github.com/universalmacro/merchant/services"
)

var router = gin.Default()
var VERSION = "0.0.1"

func Init(addr ...string) {
	var merchantController = newMerchantController()
	var sessionController = newSessionController()
	router.Use(server.CorsMiddleware())
	server.MetricsMiddleware(router)
	api.SessionApiBinding(router, sessionController)
	api.MerchantApiBinding(router, merchantController)
	router.Run(addr...)
}
