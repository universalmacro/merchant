package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	api "github.com/universalmacro/merchant-api-interfaces"
	"github.com/universalmacro/merchant/services"
)

func newSpaceController() *SpaceController {
	return &SpaceController{
		merchantService: services.GetMerchantService(),
	}
}

type SpaceController struct {
	merchantService *services.MerchantService
}

// CreateSpace implements merchantapiinterfaces.SpaceApi.
func (self *SpaceController) CreateSpace(ctx *gin.Context) {
	account := getAccount(ctx)
	merchant := self.merchantService.GetMerchant(account.MerchantId())
	if merchant == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var createSpaceRequest api.CreateSpaceRequest
	ctx.ShouldBindJSON(&createSpaceRequest)
	space := merchant.CreateSpace(createSpaceRequest.Name)
	ctx.JSON(http.StatusCreated, ConvertSpace(*space))
}

// DeleteSpace implements merchantapiinterfaces.SpaceApi.
func (*SpaceController) DeleteSpace(ctx *gin.Context) {
	panic("unimplemented")
}

// GetSpace implements merchantapiinterfaces.SpaceApi.
func (*SpaceController) GetSpace(ctx *gin.Context) {
	panic("unimplemented")
}

// ListSpaces implements merchantapiinterfaces.SpaceApi.
func (*SpaceController) ListSpaces(ctx *gin.Context) {
	panic("unimplemented")
}

// UpdateSpace implements merchantapiinterfaces.SpaceApi.
func (*SpaceController) UpdateSpace(ctx *gin.Context) {
	panic("unimplemented")
}
