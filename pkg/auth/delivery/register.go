package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/auth"
)

func RegisterHTTPEndpoints(router *gin.RouterGroup, usecase auth.UseCase) {
	h := newHandler(usecase)

	router.POST("/sign-up", h.signUp)
	router.POST("/sign-in", h.signIn)
}
