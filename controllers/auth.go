package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/universalmacro/common/config"
	"github.com/universalmacro/merchant/services"
)

var secretKey = config.GetString("node.secretKey")

type Headers struct {
	ApiKey        string
	Authorization string
}

func ApiKeyAuth(ctx *gin.Context) bool {
	var headers Headers
	ctx.ShouldBindHeader(&headers)
	if secretKey != headers.ApiKey {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return false
	}
	return true
}

func getAccount(ctx *gin.Context) services.Account {
	var headers Headers
	ctx.ShouldBindHeader(&headers)
	if headers.Authorization == "" {
		return nil
	}
	account, _ := services.GetSessionService().TokenVerification(headers.Authorization)
	return account
}
