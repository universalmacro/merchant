package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/common/server"
	api "github.com/universalmacro/merchant-api-interfaces"
	"github.com/universalmacro/merchant/services"
)

func newSpaceController() *SpaceController {
	return &SpaceController{
		merchantService: services.GetMerchantService(),
		spaceService:    services.GetSpaceService(),
	}
}

type SpaceController struct {
	spaceService    *services.SpaceService
	merchantService *services.MerchantService
}

// CreatePrinter implements merchantapiinterfaces.SpaceApi.
func (self *SpaceController) CreatePrinter(ctx *gin.Context) {
	account := getAccount(ctx)
	space := self.grantedSpace(ctx, server.UintID(ctx, "id"), account)
	if space == nil {
		return
	}
}

// DeletePrinter implements merchantapiinterfaces.SpaceApi.
func (self *SpaceController) DeletePrinter(ctx *gin.Context) {
	account := getAccount(ctx)
	space := self.grantedSpace(ctx, server.UintID(ctx, "id"), account)
	if space == nil {
		return
	}
}

// ListPrinters implements merchantapiinterfaces.SpaceApi.
func (self *SpaceController) ListPrinters(ctx *gin.Context) {
	account := getAccount(ctx)
	space := self.grantedSpace(ctx, server.UintID(ctx, "id"), account)
	if space == nil {
		return
	}
}

// UpdatePrinter implements merchantapiinterfaces.SpaceApi.
func (*SpaceController) UpdatePrinter(ctx *gin.Context) {
	panic("unimplemented")
}

// CreateTable implements merchantapiinterfaces.SpaceApi.
func (self *SpaceController) CreateTable(ctx *gin.Context) {
	account := getAccount(ctx)
	space := self.grantedSpace(ctx, server.UintID(ctx, "id"), account)
	if space == nil {
		return
	}
	if space == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	var createTableRequest api.SaveTableRequest
	ctx.ShouldBindJSON(&createTableRequest)
	table, err := space.CreateTable(createTableRequest.Label)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, ConvertTable(table))
}

// ListTables implements merchantapiinterfaces.SpaceApi.
func (self *SpaceController) ListTables(ctx *gin.Context) {
	account := getAccount(ctx)
	space := self.grantedSpace(ctx, server.UintID(ctx, "id"), account)
	if space == nil {
		return
	}
	tables := space.ListTables()
	apiTables := make([]api.Table, len(tables))
	for i := range tables {
		apiTables[i] = ConvertTable(&tables[i])
	}
	ctx.JSON(http.StatusOK, apiTables)
}

// CreateSpace implements merchantapiinterfaces.SpaceApi.
func (self *SpaceController) CreateSpace(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
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
	space := self.grantedSpace(ctx, server.UintID(ctx, "id"), account)
	if space == nil {
		return
	}
	space.Delete()
}

// GetSpace implements merchantapiinterfaces.SpaceApi.
func (self *SpaceController) GetSpace(ctx *gin.Context) {
	space := self.spaceService.GetSpace(server.UintID(ctx, "id"))
	if space == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	ctx.JSON(http.StatusOK, ConvertSpace(space))
}

// ListSpaces implements merchantapiinterfaces.SpaceApi.
func (self *SpaceController) ListSpaces(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	merchant := self.merchantService.GetMerchant(account.MerchantId())
	if merchant == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if !merchant.Granted(account) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	spaces := merchant.ListSpaces()
	apiSpaces := make([]api.Space, len(spaces))
	for i := range spaces {
		apiSpaces[i] = ConvertSpace(&spaces[i])
	}
	ctx.JSON(http.StatusOK,
		dao.List[api.Space]{Items: apiSpaces})
}

// UpdateSpace implements merchantapiinterfaces.SpaceApi.
func (self *SpaceController) UpdateSpace(ctx *gin.Context) {
	account := getAccount(ctx)
	space := self.grantedSpace(ctx, server.UintID(ctx, "id"), account)
	if space == nil {
		return
	}
	var updateSpaceRequest api.SaveSpaceRequest
	ctx.ShouldBindJSON(&updateSpaceRequest)
	ctx.JSON(http.StatusOK,
		ConvertSpace(space.SetName(updateSpaceRequest.Name).Submit()))
}

func (self *SpaceController) grantedSpace(ctx *gin.Context, spaceId uint, account services.Account) *services.Space {
	space := self.spaceService.GetSpace(spaceId)
	if account == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return nil
	}
	if space == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return nil
	}
	if !space.Granted(account) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return nil
	}
	return space
}
