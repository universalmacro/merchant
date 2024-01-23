package controllers

import "github.com/gin-gonic/gin"

func newMerchantController() *MerchantController {
	return &MerchantController{}
}

type MerchantController struct {
}

// createMerchant implements merchantapiinterfaces.MerchantApi.
func (*MerchantController) CreateMerchant(ctx *gin.Context) {
	panic("unimplemented")
}
