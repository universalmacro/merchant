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
	var createMemberRequest api.CreateMemberRequest
	ctx.ShouldBindJSON(&createMemberRequest)
	mc.memberService.SignupMember(
		server.UintID(ctx, "merchantId"),
		createMemberRequest.PhoneNumber.CountryCode,
		createMemberRequest.PhoneNumber.Number,
		*createMemberRequest.VerificationCode,
	)
}

// CreateMember implements merchantapiinterfaces.MemberApi.
func (mc *MemberController) CreateMember(ctx *gin.Context) {
	panic("unimplemented")
}

// ListMembers implements merchantapiinterfaces.MemberApi.
func (mc *MemberController) ListMembers(ctx *gin.Context) {
	panic("unimplemented")
}
