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
	spaceService    *services.SpaceService
	merchantService *services.MerchantService
}

// CreateTable implements merchantapiinterfaces.SpaceApi.
func (*SpaceController) CreateTable(ctx *gin.Context) {
	panic("unimplemented")
}

// ListTables implements merchantapiinterfaces.SpaceApi.
func (*SpaceController) ListTables(ctx *gin.Context) {
	panic("unimplemented")
}

// CreateSpace implements merchantapiinterfaces.SpaceApi.
func (self *SpaceController) CreateSpace(ctx *gin.Context) {
	account := getAccount(ctx)
	merchant := self.merchantService.GetMerchant(account.MerchantId())
	if merchant == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var createSpaceRequest api.SaveSpaceRequest
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
func (self *SpaceController) UpdateSpace(ctx *gin.Context) {
	// admin := getAccount(ctx)
	// if admin == nil {
	// 	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	// 	return
	// }
	// id := server.UintID(ctx, "id")
	// var updateSpaceRequest api.SaveSpaceRequest
	// ctx.ShouldBindJSON(&updateSpaceRequest)
	// space := self.spaceService.GetSpace(id)
	// if space == nil {
	// 	ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	// 	return
	// }
	// space.Name = updateSpaceRequest.Name
}
