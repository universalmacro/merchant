package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/universalmacro/common/server"
)

var router = gin.Default()
var VERSION = "0.0.1"

func Init(addr ...string) {
	router.Use(server.CorsMiddleware())
	server.MetricsMiddleware(router)
	router.Run(addr...)
}
