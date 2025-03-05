package character

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"spire/lobby/internal/core"
)

func HandleBotCharacterCreate(c *gin.Context, x *core.Context) {
	type Request struct {
		AccountID     int64  `json:"account_id" binding:"required"`
		CharacterName string `json:"character_name" binding:"required"`
	}

	type Response struct{}

	var r Request
	if !core.Check(c.Bind(&r), c, http.StatusBadRequest) {
		return
	}

	err := x.P.QueryRow(context.Background(), "INSERT INTO characters (name) VALUES ($1)", r.CharacterName).Scan()
	if err != nil {
		core.Check(err, c, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, Response{})
}
