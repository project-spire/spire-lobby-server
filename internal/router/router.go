package router

import (
	"github.com/gin-gonic/gin"
	"spire/lobby/internal/core"
)

func NewRouter(ctx *core.Context) *gin.Engine {
	r := gin.Default()

	r.GET("/auth/bot", func(c *gin.Context) { HandleAuthBot(c, ctx) })

	return r
}
