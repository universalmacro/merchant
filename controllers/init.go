package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/universalmacro/common/server"
	api "github.com/universalmacro/merchant-api-interfaces"
	_ "github.com/universalmacro/merchant/services"
)

var router *gin.Engine
var VERSION = "0.0.1"

func Init(addr ...string) {
	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
	var merchantController = newMerchantController()
	var sessionController = newSessionController()
	var verificationController = newVerificationController()
	var spaceController = newSpaceController()
	var orderController = newOrderController()
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, ApiKey")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})
	server.MetricsMiddleware(router)
	router.GET("/version", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"version": VERSION})
	})
	api.SessionApiBinding(router, sessionController)
	api.MerchantApiBinding(router, merchantController)
	api.VerificationApiBinding(router, verificationController)
	api.SpaceApiBinding(router, spaceController)
	api.OrderApiBinding(router, orderController)
	router.Run(addr...)
}
