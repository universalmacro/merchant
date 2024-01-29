package controllers

import (
	"net/http"

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
	var verificationController = newVerificationController()
	router.Use(server.CorsMiddleware())
	server.MetricsMiddleware(router)
	router.GET("/version", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"version": VERSION})
	})
	api.SessionApiBinding(router, sessionController)
	api.MerchantApiBinding(router, merchantController)
	api.VerificationApiBinding(router, verificationController)
	router.Run(addr...)
}
