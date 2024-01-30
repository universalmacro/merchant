package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/universalmacro/common/config"
	"github.com/universalmacro/common/server"
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

// DeleteMerchant implements merchantapiinterfaces.MerchantApi.
func (*MerchantController) DeleteMerchant(ctx *gin.Context) {
	panic("unimplemented")
}

// GetMerchant implements merchantapiinterfaces.MerchantApi.
func (*MerchantController) GetMerchant(ctx *gin.Context) {
	panic("unimplemented")
}

// GetSelfMerchant implements merchantapiinterfaces.MerchantApi.
func (*MerchantController) GetSelfMerchant(ctx *gin.Context) {
	panic("unimplemented")
}

// ListMerchants implements merchantapiinterfaces.MerchantApi.
func (*MerchantController) ListMerchants(ctx *gin.Context) {
	if !ApiKeyAuth(ctx) {
		return
	}
	index, limit := server.IndexAndLimit(ctx)
	list := services.GetMerchantService().ListMerchants(index, limit)
	result := make([]api.Merchant, len(list.Items))
	for i, merchant := range list.Items {
		result[i] = ConvertMerchant(merchant)
	}
	ctx.JSON(http.StatusOK, api.MerchantList{
		Items: result,
		Pagination: api.Pagination{
			Index: list.Pagination.Index,
			Limit: list.Pagination.Limit,
			Total: list.Pagination.Total,
		},
	})
}

// UpdateMerchant implements merchantapiinterfaces.MerchantApi.
func (*MerchantController) UpdateMerchant(ctx *gin.Context) {
	panic("unimplemented")
}

// UpdateMerchantPassword implements merchantapiinterfaces.MerchantApi.
func (c *MerchantController) UpdateMerchantPassword(ctx *gin.Context) {
	id := server.UintID(ctx, "id")
	merchant := c.merchantService.GetMerchant(id)
	if merchant == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "merchant not found"})
		return
	}
	var request api.UpdatePasswordRequest
	ctx.ShouldBindJSON(&request)
	if !ApiKeyAuth(ctx) {
		return
	}

	merchant.UpdatePassword(request.Password)
	ctx.JSON(http.StatusNoContent, nil)
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

var secretKey = config.GetString("node.secretKey")

func ApiKeyAuth(ctx *gin.Context) bool {
	// return true
	var headers Headers
	ctx.ShouldBindHeader(&headers)
	if secretKey != headers.ApiKey {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return false
	}
	return true
}
