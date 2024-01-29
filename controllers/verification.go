package controllers

import "github.com/gin-gonic/gin"

func newVerificationController() *VerificationController {
	return &VerificationController{}
}

type VerificationController struct {
}

// SendVerificationCode implements merchantapiinterfaces.VerificationApi.
func (*VerificationController) SendVerificationCode(ctx *gin.Context) {
	panic("unimplemented")
}
