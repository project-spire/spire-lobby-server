package route

import (
	"github.com/gin-gonic/gin"
	"spire/lobby/internal/core"
	"spire/lobby/internal/route/bot/account"
	"spire/lobby/internal/route/bot/character"
)

func NewRouter(ctx *core.Context) *gin.Engine {
	r := gin.Default()

	r.GET("/ping")

	r.POST("/bot/account/auth", func(c *gin.Context) { account.HandleBotAccountAuth(c, ctx) })
	r.POST("/bot/account/create", func(c *gin.Context) { account.HandleBotAccountMe(c, ctx) })
	r.POST("/bot/account/me", func(c *gin.Context) { account.HandleBotAccountCreate(c, ctx) })
	r.POST("/bot/character/create", func(c *gin.Context) { character.HandleBotCharacterCreate(c, ctx) })
	r.POST("/bot/character/list", func(c *gin.Context) { character.HandleBotCharacterList(c, ctx) })

	return r
}
