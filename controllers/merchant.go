package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/universalmacro/common/config"
	api "github.com/universalmacro/merchant-api-interfaces"
	"github.com/universalmacro/merchant/services"
)

func newMerchantController() *MerchantController {
	return &MerchantController{
		merchantService: services.GetMerchantService(),
	}
}

type MerchantController struct {
	merchantService *services.MerchantService
}

type Headers struct {
	ApiKey string
}

// createMerchant implements merchantapiinterfaces.MerchantApi.
func (c *MerchantController) CreateMerchant(ctx *gin.Context) {
	if !ApiKeyAuth(ctx) {
		return
	}
	var request api.CreateMerchantRequest
	ctx.ShouldBindJSON(&request)
	merchant := c.merchantService.CreateMerchant(request.ShortMerchantId, request.Account, request.Password)
	if merchant == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	ctx.JSON(http.StatusCreated, api.Merchant{
		ShortMerchantId: merchant.ShortMerchantId(),
		Account:         merchant.Account(),
		CreatedAt:       merchant.CreatedAt().Unix(),
		UpdatedAt:       merchant.UpdatedAt().Unix(),
	})
}

func ApiKeyAuth(ctx *gin.Context) bool {
	secretKey := config.GetString("node.secretKey")
	var headers Headers
	ctx.ShouldBindHeader(&headers)
	if secretKey != headers.ApiKey {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return false
	}
	return true
}
