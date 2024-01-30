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
func (*SessionController) CreateSession(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, api.Session{
		Token: "test-token",
	})
}
