package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/universalmacro/common/fault"
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

// SendMerchantVerificationCode implements merchantapiinterfaces.MerchantApi.
func (mc *MerchantController) SendMerchantVerificationCode(ctx *gin.Context) {
	var createVerificationCodeRequest api.CreateVerificationCodeRequest
	ctx.ShouldBindJSON(&createVerificationCodeRequest)
	err := mc.merchantService.CreateVerificationCode(
		server.UintID(ctx, "merchantId"),
		createVerificationCodeRequest.PhoneNumber.CountryCode,
		createVerificationCodeRequest.PhoneNumber.Number)
	if err != nil {
		switch err {
		case services.ErrVerificationCodeHasBeenSent:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": services.ErrVerificationCodeHasBeenSent})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusCreated, nil)
}

// GetSelfMerchantConfig implements merchantapiinterfaces.MerchantApi.
func (*MerchantController) GetSelfMerchantConfig(ctx *gin.Context) {
	panic("unimplemented")
}

// UpdateSelfMerchantConfig implements merchantapiinterfaces.MerchantApi.
func (*MerchantController) UpdateSelfMerchantConfig(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		fault.GinHandler(ctx, fault.ErrUnauthorized)
		return
	}
	var request api.MerchantConfig
	ctx.ShouldBindJSON(&request)
	merchant := services.GetMerchantService().GetMerchant(account.MerchantId())
	if !merchant.Granted(account) {
		fault.GinHandler(ctx, fault.ErrPermissionDenied)
		return
	}
	if request.Currency != nil {
		merchant.SetCurrency(string(*request.Currency))
	}
	merchant.Submit()
	ctx.JSON(http.StatusNoContent, request)
}

// DeleteSelfContactForm implements merchantapiinterfaces.MerchantApi.
func (*MerchantController) DeleteSelfContactForm(ctx *gin.Context) {
	panic("unimplemented")
}

// GetSelfContactForm implements merchantapiinterfaces.MerchantApi.
func (*MerchantController) GetSelfContactForm(ctx *gin.Context) {
	panic("unimplemented")
}

// ListSelfContactForm implements merchantapiinterfaces.MerchantApi.
func (*MerchantController) ListSelfContactForm(ctx *gin.Context) {
	panic("unimplemented")
}

// SendContactForm implements merchantapiinterfaces.MerchantApi.
func (*MerchantController) SendContactForm(ctx *gin.Context) {
	panic("unimplemented")
}

// ListSelfMembers implements merchantapiinterfaces.MerchantApi.
func (*MerchantController) ListSelfMembers(ctx *gin.Context) {
	panic("unimplemented")
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
	for i := range list.Items {
		result[i] = ConvertMerchant(list.Items[i])
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
	id := server.UintID(ctx, "merchantId")
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
