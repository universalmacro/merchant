package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/universalmacro/common/server"
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
	ctx.JSON(http.StatusCreated, ConvertSpace(space))
}

// DeleteSpace implements merchantapiinterfaces.SpaceApi.
func (self *SpaceController) DeleteSpace(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id := server.UintID(ctx, "id")
	space := self.spaceService.GetSpace(id)
	if space == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if !space.Granted(account) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	space.Delete()
}

// GetSpace implements merchantapiinterfaces.SpaceApi.
func (self *SpaceController) GetSpace(ctx *gin.Context) {
	id := server.UintID(ctx, "id")
	space := self.spaceService.GetSpace(id)
	if space == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	ctx.JSON(http.StatusOK, ConvertSpace(space))
}

// ListSpaces implements merchantapiinterfaces.SpaceApi.
func (*SpaceController) ListSpaces(ctx *gin.Context) {
	panic("unimplemented")
}

// UpdateSpace implements merchantapiinterfaces.SpaceApi.
func (self *SpaceController) UpdateSpace(ctx *gin.Context) {
	admin := getAccount(ctx)
	if admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id := server.UintID(ctx, "id")
	var updateSpaceRequest api.SaveSpaceRequest
	ctx.ShouldBindJSON(&updateSpaceRequest)
	space := self.spaceService.GetSpace(id)
	if space == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if !space.Granted(admin) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	ctx.JSON(http.StatusOK,
		ConvertSpace(space.SetName(updateSpaceRequest.Name).Submit()))
}
