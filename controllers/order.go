package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/universalmacro/common/server"
	api "github.com/universalmacro/merchant-api-interfaces"
	"github.com/universalmacro/merchant/services"
)

func newOrderController() *OrderController {
	return &OrderController{
		merchantService: services.GetMerchantService(),
		spaceService:    services.GetSpaceService(),
		tableService:    services.GetTableService(),
	}
}

type OrderController struct {
	merchantService *services.MerchantService
	spaceService    *services.SpaceService
	tableService    *services.TableService
}

// CancelOrder implements merchantapiinterfaces.OrderApi.
func (*OrderController) CancelOrder(ctx *gin.Context) {
	panic("unimplemented")
}

// CreateFood implements merchantapiinterfaces.OrderApi.
func (self *OrderController) CreateFood(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	space := self.spaceService.GetSpace(server.UintID(ctx, "id"))
	if space == nil {
		ctx.JSON(404, gin.H{"error": "not found"})
		return
	}
	if !space.Granted(account) {
		ctx.JSON(403, gin.H{"error": "forbidden"})
		return
	}
	var createFoodRequest api.SaveFoodRequest
	ctx.ShouldBindJSON(&createFoodRequest)
	// food, err := space.CreateFood(createFoodRequest.Name, createFoodRequest.Description, createFoodRequest.Price, createFoodRequest.FixedOffset)
}

// CreateOrder implements merchantapiinterfaces.OrderApi.
func (*OrderController) CreateOrder(ctx *gin.Context) {
	panic("unimplemented")
}

// CreateTable implements merchantapiinterfaces.OrderApi.
func (self *OrderController) CreateTable(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	id := server.UintID(ctx, "id")
	space := self.spaceService.GetSpace(id)
	if space == nil {
		ctx.JSON(404, gin.H{"error": "not found"})
		return
	}
	if !space.Granted(account) {
		ctx.JSON(403, gin.H{"error": "forbidden"})
		return
	}
	var createTableRequest api.SaveTableRequest
	ctx.ShouldBindJSON(&createTableRequest)
	table, err := space.CreateTable(createTableRequest.Label)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(201, ConvertTable(table))
}

// DeleteFood implements merchantapiinterfaces.OrderApi.
func (*OrderController) DeleteFood(ctx *gin.Context) {
	panic("unimplemented")
}

// DeleteTable implements merchantapiinterfaces.OrderApi.
func (self *OrderController) DeleteTable(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	id := server.UintID(ctx, "id")
	table := self.tableService.GetTable(id)
	if table == nil {
		ctx.JSON(404, gin.H{"error": "not found"})
		return
	}
	if !table.Granted(account) {
		ctx.JSON(403, gin.H{"error": "forbidden"})
		return
	}
	table.Delete()
	ctx.JSON(204, nil)
}

// GetFoodById implements merchantapiinterfaces.OrderApi.
func (*OrderController) GetFoodById(ctx *gin.Context) {
	panic("unimplemented")
}

// ListFoods implements merchantapiinterfaces.OrderApi.
func (*OrderController) ListFoods(ctx *gin.Context) {
	panic("unimplemented")
}

// ListTables implements merchantapiinterfaces.OrderApi.
func (*OrderController) ListTables(ctx *gin.Context) {
	panic("unimplemented")
}

// UpdateFood implements merchantapiinterfaces.OrderApi.
func (*OrderController) UpdateFood(ctx *gin.Context) {
	panic("unimplemented")
}

// UpdateFoodImage implements merchantapiinterfaces.OrderApi.
func (*OrderController) UpdateFoodImage(ctx *gin.Context) {
	panic("unimplemented")
}

// UpdateTable implements merchantapiinterfaces.OrderApi.
func (self *OrderController) UpdateTable(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	id := server.UintID(ctx, "id")
	table := self.tableService.GetTable(id)
	if table == nil {
		ctx.JSON(404, gin.H{"error": "not found"})
		return
	}
	if !table.Granted(account) {
		ctx.JSON(403, gin.H{"error": "forbidden"})
		return
	}
	var updateTableRequest api.SaveTableRequest
	ctx.ShouldBindJSON(&updateTableRequest)
	table.SetLabel(updateTableRequest.Label).Submit()
}
