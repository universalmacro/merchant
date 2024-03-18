package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/universalmacro/common/server"
	api "github.com/universalmacro/merchant-api-interfaces"
	"github.com/universalmacro/merchant/services"
)

type MemberController struct {
	memberService *services.MerchantService
}

// SignupMember implements merchantapiinterfaces.MemberApi.
func (mc *MemberController) SignupMember(ctx *gin.Context) {
	var signupMemberRequest api.SignupMemberRequest
	ctx.ShouldBindJSON(&signupMemberRequest)
	token, err := mc.memberService.SignupMember(
		server.UintID(ctx, "merchantId"),
		signupMemberRequest.PhoneNumber.CountryCode,
		signupMemberRequest.PhoneNumber.Number,
		signupMemberRequest.VerificationCode,
	)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(201, api.Session{
		Token: token,
	})
}

// CreateMember implements merchantapiinterfaces.MemberApi.
func (mc *MemberController) CreateMember(ctx *gin.Context) {
	panic("unimplemented")
}

// ListMembers implements merchantapiinterfaces.MemberApi.
func (mc *MemberController) ListMembers(ctx *gin.Context) {
	panic("unimplemented")
}
