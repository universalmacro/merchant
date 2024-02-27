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
		printerService:  services.GetPrinterService(),
	}
}

type SpaceController struct {
	spaceService    *services.SpaceService
	merchantService *services.MerchantService
	printerService  *services.PrinterService
}

// GetPrinter implements merchantapiinterfaces.SpaceApi.
func (*SpaceController) GetPrinter(ctx *gin.Context) {
	account := getAccount(ctx)
	printer := services.GetPrinterService().GetPrinter(server.UintID(ctx, "printerId"))
	if printer == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if !printer.Granted(account) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	ctx.JSON(http.StatusOK, ConvertPrinter(printer))
}

// CreatePrinter implements merchantapiinterfaces.SpaceApi.
func (self *SpaceController) CreatePrinter(ctx *gin.Context) {
	account := getAccount(ctx)
	space := grantedSpace(ctx, server.UintID(ctx, "spaceId"), account)
	if space == nil {
		return
	}
	var createPrinterRequest api.SavePrinter
	ctx.ShouldBindJSON(&createPrinterRequest)
	printer := space.CreatePrinter(createPrinterRequest.Name, createPrinterRequest.Sn, string(createPrinterRequest.Type), string(createPrinterRequest.Model))
	ctx.JSON(http.StatusCreated, ConvertPrinter(printer))
}

// DeletePrinter implements merchantapiinterfaces.SpaceApi.
func (self *SpaceController) DeletePrinter(ctx *gin.Context) {
	account := getAccount(ctx)
	printer := self.printerService.GetPrinter(server.UintID(ctx, "printerId"))
	if printer == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if !printer.Granted(account) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	printer.Delete()
	ctx.JSON(http.StatusNoContent, nil)
}

// ListPrinters implements merchantapiinterfaces.SpaceApi.
func (self *SpaceController) ListPrinters(ctx *gin.Context) {
	account := getAccount(ctx)
	space := grantedSpace(ctx, server.UintID(ctx, "spaceId"), account)
	if space == nil {
		return
	}
	printers := space.ListPrinters()
	apiPrinters := make([]api.Printer, len(printers))
	for i := range printers {
		apiPrinters[i] = ConvertPrinter(&printers[i])
	}
	ctx.JSON(http.StatusOK, apiPrinters)
}

// UpdatePrinter implements merchantapiinterfaces.SpaceApi.
func (*SpaceController) UpdatePrinter(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	printer := services.GetPrinterService().GetPrinter(server.UintID(ctx, "printerId"))
	if printer == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if !printer.Granted(account) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	var updatePrinterRequest api.SavePrinter
	ctx.ShouldBindJSON(&updatePrinterRequest)
	printer.Name = updatePrinterRequest.Name
	printer.Sn = updatePrinterRequest.Sn
	printer.Type = string(updatePrinterRequest.Type)
	printer.Model = string(updatePrinterRequest.Model)
	printer.Submit()
	ctx.JSON(http.StatusOK, ConvertPrinter(printer))
}

// CreateTable implements merchantapiinterfaces.SpaceApi.
func (self *SpaceController) CreateTable(ctx *gin.Context) {
	account := getAccount(ctx)
	space := grantedSpace(ctx, server.UintID(ctx, "spaceId"), account)
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
	space := grantedSpace(ctx, server.UintID(ctx, "spaceId"), account)
	if space == nil {
		return
	}
	space.Delete()
}

// GetSpace implements merchantapiinterfaces.SpaceApi.
func (self *SpaceController) GetSpace(ctx *gin.Context) {
	space := self.spaceService.GetSpace(server.UintID(ctx, "spaceId"))
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
	space := grantedSpace(ctx, server.UintID(ctx, "spaceId"), account)
	if space == nil {
		return
	}
	var updateSpaceRequest api.SaveSpaceRequest
	ctx.ShouldBindJSON(&updateSpaceRequest)
	ctx.JSON(http.StatusOK,
		ConvertSpace(space.SetName(updateSpaceRequest.Name).Submit()))
}

func grantedSpace(ctx *gin.Context, spaceId uint, account services.Account) *services.Space {
	space := services.GetSpaceService().GetSpace(spaceId)
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
