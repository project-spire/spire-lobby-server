package router

import (
	"github.com/gin-gonic/gin"
	"spire/lobby/internal/core"
)

func NewRouter(ctx *core.Context) *gin.Engine {
	r := gin.Default()

	r.POST("/bot/auth", func(c *gin.Context) { HandleBotAuth(c, ctx) })
	r.POST("/bot/account", func(c *gin.Context) { HandleBotAccount(c, ctx) })
	r.POST("/bot/register", func(c *gin.Context) { HandleBotRegister(c, ctx) })
	r.POST("/bot/character/list", func(c *gin.Context) { HandleBotCharacterList(c, ctx) })
	r.POST("/bot/character/create", func(c *gin.Context) { HandleBotCharacterCreate(c, ctx) })

	return r
}
