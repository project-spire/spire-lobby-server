package router

import (
	"github.com/gin-gonic/gin"
	"spire/lobby/internal/core"
)

func NewRouter(ctx *core.Context) *gin.Engine {
	r := gin.Default()

	r.POST("/auth/bot", func(c *gin.Context) { HandleAuthBot(c, ctx) })
	r.POST("/account/bot", func(c *gin.Context) { HandleAccountBot(c, ctx) })
	r.POST("/register/bot", func(c *gin.Context) { HandleRegisterBot(c, ctx) })

	return r
}
