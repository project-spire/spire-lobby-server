package character

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"spire/lobby/internal/core"
)

func HandleBotCharacterCreate(c *gin.Context, x *core.Context) {
	type Request struct {
		AccountID     uint64 `json:"account_id" binding:"required"`
		CharacterName string `json:"character_name" binding:"required"`
	}

	type Response struct {
		CharacterID uint64 `json:"character_id"`
	}

	var r Request
	if !core.Check(c.Bind(&r), c, http.StatusBadRequest) {
		return
	}

	var characterID uint64
	err := x.P.QueryRow(context.Background(),
		"INSERT INTO characters (account_id, name) VALUES ($1, $2) RETURNING id",
		r.AccountID, r.CharacterName).Scan(&characterID)
	if err != nil {
		core.Check(err, c, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, Response{})
}
