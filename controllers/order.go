package controllers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/common/server"
	"github.com/universalmacro/common/utils"
	api "github.com/universalmacro/merchant-api-interfaces"
	"github.com/universalmacro/merchant/controllers/factories"
	"github.com/universalmacro/merchant/dao/entities"
	"github.com/universalmacro/merchant/services"
)

func newOrderController() *OrderController {
	return &OrderController{
		merchantService: services.GetMerchantService(),
		spaceService:    services.GetSpaceService(),
		tableService:    services.GetTableService(),
		foodService:     services.GetFoodService(),
		orderService:    services.GetOrderService(),
	}
}

type OrderController struct {
	merchantService *services.MerchantService
	spaceService    *services.SpaceService
	tableService    *services.TableService
	foodService     *services.FoodService
	orderService    *services.OrderService
}

// CreateBill implements merchantapiinterfaces.OrderApi.
func (oc *OrderController) CreateBill(ctx *gin.Context) {
	// account := getAccount(ctx)
	// var createBillRequest api.CreateBillRequest
	// ctx.ShouldBindJSON(&createBillRequest)

}

// GetBill implements merchantapiinterfaces.OrderApi.
func (oc *OrderController) GetBill(ctx *gin.Context) {
	panic("unimplemented")
}

// ListBills implements merchantapiinterfaces.OrderApi.
func (oc *OrderController) ListBills(ctx *gin.Context) {
	panic("unimplemented")
}

// PrintBill implements merchantapiinterfaces.OrderApi.
func (oc *OrderController) PrintBill(ctx *gin.Context) {
	panic("unimplemented")
}

// UpdateOrderTableLabel implements merchantapiinterfaces.OrderApi.
func (oc *OrderController) UpdateOrderTableLabel(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	order := oc.orderService.GetById(server.UintID(ctx, "orderId"))
	if order == nil {
		ctx.JSON(404, gin.H{"error": "not found"})
		return
	}
	if !order.Space().Granted(account) {
		ctx.JSON(403, gin.H{"error": "forbidden"})
		return
	}
	var updateOrderTableLabelRequest api.UpdateOrderTableLabelRequest
	ctx.ShouldBindJSON(&updateOrderTableLabelRequest)
	order.SetTableLabel(updateOrderTableLabelRequest.TableLabel).Submit()
	ctx.JSON(200, ConvertOrder(order))
}

// AddOrder implements merchantapiinterfaces.OrderApi.
func (oc *OrderController) AddOrder(ctx *gin.Context) {
	order := oc.orderService.GetById(server.UintID(ctx, "orderId"))
	if order == nil {
		ctx.JSON(404, gin.H{"error": "not found"})
		return
	}
	var addOrderRequest api.AddOrderRequest
	ctx.ShouldBindJSON(&addOrderRequest)
	// order.AddItems()
}

// ListFoodPrinters implements merchantapiinterfaces.OrderApi.
func (oc *OrderController) ListFoodPrinters(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	id := server.UintID(ctx, "id")
	food := oc.foodService.GetById(id)
	if food == nil {
		ctx.JSON(404, gin.H{"error": "not found"})
		return
	}
	if !food.Granted(account) {
		ctx.JSON(403, gin.H{"error": "forbidden"})
		return
	}
	printers := food.Printers()
	result := make([]api.Printer, len(printers))
	for i := range printers {
		result[i] = ConvertPrinter(&printers[i])
	}
	ctx.JSON(200, result)
}

// UpdateFoodPrinters implements merchantapiinterfaces.OrderApi.
func (*OrderController) UpdateFoodPrinters(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	id := server.UintID(ctx, "id")
	food := services.GetFoodService().GetById(id)
	if food == nil {
		ctx.JSON(404, gin.H{"error": "not found"})
		return
	}
	if !food.Granted(account) {
		ctx.JSON(403, gin.H{"error": "forbidden"})
		return
	}
	var updateFoodPrintersRequest api.UpdateFoodPrintersRequest
	ctx.ShouldBindJSON(&updateFoodPrintersRequest)
	var printerIds []uint
	for i := range updateFoodPrintersRequest.PrinterIds {
		printerIds = append(printerIds, utils.StringToUint(updateFoodPrintersRequest.PrinterIds[i]))
	}
	food.SetPrinters(printerIds...).Submit()
	printers := food.Printers()
	result := make([]api.Printer, len(printers))
	for i := range printers {
		result[i] = ConvertPrinter(&printers[i])
	}
	ctx.JSON(200, result)
}

// ListFoodCategories implements merchantapiinterfaces.OrderApi.
func (oc *OrderController) ListFoodCategories(ctx *gin.Context) {
	space := oc.spaceService.GetSpace(server.UintID(ctx, "spaceId"))
	if space == nil {
		ctx.JSON(404, gin.H{"error": "not found"})
		return
	}
	ctx.JSON(200, space.FoodCategories())
}

// UpdateFoodCategories implements merchantapiinterfaces.OrderApi.
func (*OrderController) UpdateFoodCategories(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	space := grantedSpace(ctx, server.UintID(ctx, "spaceId"), account)
	if space == nil {
		return
	}
	var updateFoodCategoriesRequest []string
	ctx.ShouldBindJSON(&updateFoodCategoriesRequest)
	space.SetFoodCategories(updateFoodCategoriesRequest...)
	space.Submit()
	ctx.JSON(200, space.FoodCategories())
}

// CancelOrder implements merchantapiinterfaces.OrderApi.
func (*OrderController) CancelOrder(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	order := services.GetOrderService().GetById(server.UintID(ctx, "orderId"))
	if order == nil {
		ctx.JSON(404, gin.H{"error": "order not found"})
		return
	}
	if !order.Space().Granted(account) {
		ctx.JSON(403, gin.H{"error": "forbidden"})
		return
	}
	var cancelOrderRequest api.CancelOrderRequest
	ctx.ShouldBindJSON(&cancelOrderRequest)
	foodSpecs := factories.NewFoodSpecs(cancelOrderRequest.Foods)
	order.CancelItems(foodSpecs...)
	order.Submit()
	ctx.JSON(200, ConvertOrder(order))
}

// CreateFood implements merchantapiinterfaces.OrderApi.
func (oc *OrderController) CreateFood(ctx *gin.Context) {
	account := getAccount(ctx)
	space := grantedSpace(ctx, server.UintID(ctx, "spaceId"), account)
	if space == nil {
		return
	}
	var createFoodRequest api.SaveFoodRequest
	ctx.ShouldBindJSON(&createFoodRequest)
	food, err := updateFood(createFoodRequest, services.NewEmptyFood())
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	space.CreateFood(food)
	ctx.JSON(201, ConvertFood(food))
}

// CreateOrder implements merchantapiinterfaces.OrderApi.
func (oc *OrderController) CreateOrder(ctx *gin.Context) {
	account := getAccount(ctx)
	space := oc.spaceService.GetSpace(server.UintID(ctx, "spaceId"))
	if space == nil {
		ctx.JSON(404, gin.H{"error": "space not found"})
		return
	}
	var createOrderRequest api.CreateOrderRequest
	ctx.ShouldBindJSON(&createOrderRequest)
	if len(createOrderRequest.Foods) == 0 {
		ctx.JSON(400, gin.H{"error": "no foods"})
		return
	}
	order := space.CreateOrder(
		account,
		createOrderRequest.TableLabel,
		factories.NewFoodSpecs(createOrderRequest.Foods),
	)
	order.PrintKitchen()
	order.PrintCashier()
	ctx.JSON(201, ConvertOrder(&order))
}

// DeleteFood implements merchantapiinterfaces.OrderApi.
func (oc *OrderController) DeleteFood(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	id := server.UintID(ctx, "id")
	food := oc.foodService.GetById(id)
	if food == nil {
		ctx.JSON(404, gin.H{"error": "not found"})
		return
	}
	if !food.Granted(account) {
		ctx.JSON(403, gin.H{"error": "forbidden"})
		return
	}
	food.Delete()
	ctx.JSON(204, nil)
}

// DeleteTable implements merchantapiinterfaces.OrderApi.
func (oc *OrderController) DeleteTable(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	id := server.UintID(ctx, "id")
	table := oc.tableService.GetTable(id)
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
func (oc *OrderController) GetFoodById(ctx *gin.Context) {
	food := oc.foodService.GetById(server.UintID(ctx, "id"))
	if food == nil {
		ctx.JSON(404, gin.H{"error": "not found"})
		return
	}
	ctx.JSON(200, ConvertFood(food))
}

// ListFoods implements merchantapiinterfaces.OrderApi.
func (oc *OrderController) ListFoods(ctx *gin.Context) {
	space := oc.spaceService.GetSpace(server.UintID(ctx, "spaceId"))
	if space == nil {
		ctx.JSON(404, gin.H{"error": "not found"})
		return
	}
	foods := space.Foods()
	result := make([]api.Food, len(foods))
	for i := range foods {
		result[i] = ConvertFood(&foods[i])
	}
	ctx.JSON(200, result)
}

// ListTables implements merchantapiinterfaces.OrderApi.
func (oc *OrderController) ListTables(ctx *gin.Context) {
	space := oc.spaceService.GetSpace(server.UintID(ctx, "spaceId"))
	if space == nil {
		ctx.JSON(404, gin.H{"error": "not found"})
		return
	}
	tables := space.ListTables()
	result := make([]api.Table, len(tables))
	for i := range tables {
		result[i] = ConvertTable(&tables[i])
	}
	ctx.JSON(200, result)
}

// UpdateFood implements merchantapiinterfaces.OrderApi.
func (oc *OrderController) UpdateFood(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	food := oc.foodService.GetById(server.UintID(ctx, "id"))
	if food == nil {
		ctx.JSON(404, gin.H{"error": "not found"})
		return
	}
	if !food.Granted(account) {
		ctx.JSON(403, gin.H{"error": "forbidden"})
		return
	}
	var updateFoodRequest api.SaveFoodRequest
	ctx.ShouldBindJSON(&updateFoodRequest)
	_, err := updateFood(updateFoodRequest, food)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	food.Submit()
	ctx.JSON(200, ConvertFood(food))
}

// UpdateFoodImage implements merchantapiinterfaces.OrderApi.
func (oc *OrderController) UpdateFoodImage(ctx *gin.Context) {
	// TODO: implement
	account := getAccount(ctx)
	if account == nil {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	id := server.UintID(ctx, "id")
	food := oc.foodService.GetById(id)
	if food == nil {
		ctx.JSON(404, gin.H{"error": "not found"})
		return
	}
	if !food.Granted(account) {
		ctx.JSON(403, gin.H{"error": "forbidden"})
		return
	}
	file, _ := ctx.FormFile("file")
	if file == nil {
		ctx.JSON(400, gin.H{"error": "no file"})
		return
	}
	food.UpdateImage(file)
	ctx.JSON(200, ConvertFood(food))
}

// CreateTable implements merchantapiinterfaces.OrderApi.
func (oc *OrderController) CreateTable(ctx *gin.Context) {
	account := getAccount(ctx)
	space := grantedSpace(ctx, server.UintID(ctx, "spaceId"), account)
	if space == nil {
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

// UpdateTable implements merchantapiinterfaces.OrderApi.
func (oc *OrderController) UpdateTable(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	id := server.UintID(ctx, "id")
	table := oc.tableService.GetTable(id)
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

func (oc *OrderController) ListOrders(ctx *gin.Context) {
	account := getAccount(ctx)
	space := grantedSpace(ctx, server.UintID(ctx, "spaceId"), account)
	if space == nil {
		return
	}
	var options []dao.Option
	if startAt, err := strconv.Atoi(ctx.Query("startAt")); err == nil {
		options = append(options, dao.Where("created_at >= ?", time.Unix(int64(startAt), 0)))
	}
	if endAt, err := strconv.Atoi(ctx.Query("endAt")); err == nil {
		options = append(options, dao.Where("created_at <= ?", time.Unix(int64(endAt), 0)))
	}
	options = append(options, dao.Where("space_id = ?", space.ID()))
	if statuses := ctx.QueryArray("statuses"); len(statuses) > 0 {
		options = append(options, dao.Where("status IN (?)", statuses))
	}
	if tableLabels := ctx.QueryArray("tableLabels"); len(tableLabels) > 0 {
		options = append(options, dao.Where("table_label IN (?)", tableLabels))
	}
	orders := oc.orderService.List(options...)
	result := make([]api.Order, len(orders))
	for i := range orders {
		result[i] = ConvertOrder(&orders[i])
	}
	ctx.JSON(200, result)
}

func updateFood(saveFoodRequest api.SaveFoodRequest, food *services.Food) (*services.Food, error) {
	food.SetName(
		saveFoodRequest.Name).
		SetDescription(saveFoodRequest.Description).
		SetPrice(saveFoodRequest.Price).
		SetFixedOffset(saveFoodRequest.FixedOffset).
		SetImage(saveFoodRequest.Image).
		SetCategories(saveFoodRequest.Categories...)
	if saveFoodRequest.Status != nil {
		food.SetStatus(string(*saveFoodRequest.Status))
	}
	attributes := saveFoodRequest.Attributes
	food.SetAttributes([]entities.Attribute{})
	for i := range attributes {
		if len(attributes[i].Options) == 0 {
			continue
		}
		options := make([]entities.Option, len(attributes[i].Options))
		for j := range attributes[i].Options {
			var extra int64
			if attributes[i].Options[j].Extra != nil {
				extra = *attributes[i].Options[j].Extra
			}
			options[j] = entities.Option{
				Label: attributes[i].Options[j].Label,
				Extra: extra,
			}
		}
		_, err := food.AddAttribute(attributes[i].Label, options...)
		if err != nil {
			return nil, err
		}
	}
	return food, nil
}
