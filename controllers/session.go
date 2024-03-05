package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	api "github.com/universalmacro/merchant-api-interfaces"
	"github.com/universalmacro/merchant/services"
)

func newSessionController() *SessionController {
	return &SessionController{
		sessionService: services.GetSessionService(),
	}
}

type SessionController struct {
	sessionService *services.SessionService
}

// CreateSession implements merchantapiinterfaces.SessionApi.
func (sc *SessionController) CreateSession(ctx *gin.Context) {
	var request api.CreateSessionRequest
	ctx.ShouldBindJSON(&request)
	token, err := sc.sessionService.CreateSession(*request.Account, *request.Password, request.ShortMerchantId)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    "Unauthorized",
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, api.Session{
		Token: token,
	})
}
