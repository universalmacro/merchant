package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	api "github.com/universalmacro/merchant-api-interfaces"
)

func newSessionController() *SessionController {
	return &SessionController{}
}

type SessionController struct {
}

// CreateSession implements merchantapiinterfaces.SessionApi.
func (*SessionController) CreateSession(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, api.Session{
		Token: "test-token",
	})
}
