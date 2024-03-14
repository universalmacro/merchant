package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/universalmacro/common/server"
	api "github.com/universalmacro/merchant-api-interfaces"
	"github.com/universalmacro/merchant/services"
	_ "github.com/universalmacro/merchant/services"
)

const VERSION = "0.0.2"

func Init(addr ...string) {
	var router *gin.Engine
	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
	var merchantController = newMerchantController()
	var sessionController = newSessionController()
	var verificationController = newVerificationController()
	var spaceController = newSpaceController()
	var orderController = newOrderController()
	var memberController = &MemberController{
		memberService: services.GetMerchantService(),
	}
	router.GET("/orders/subscription", orderController.OrderSubscription)
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
	api.MemberApiBinding(router, memberController)
	router.Run(addr...)
}
