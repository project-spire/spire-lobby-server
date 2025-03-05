package account

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"spire/lobby/internal/core"
)

func HandleBotAccountMe(c *gin.Context, x *core.Context) {
	type Request struct {
		BotID int64 `json:"bot_id" binding:"required"`
	}

	type Response struct {
		Found     bool  `json:"found"`
		AccountID int64 `json:"account_id"`
	}

	var r Request
	if !core.Check(c.Bind(&r), c, http.StatusBadRequest) {
		return
	}

	found := true
	var accountID int64 = 0
	err := x.P.QueryRow(context.Background(), "SELECT account_id FROM bots WHERE id=$1", r.BotID).Scan(&accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			found = false
		} else {
			core.Check(err, c, http.StatusInternalServerError)
			return
		}
	}

	c.JSON(http.StatusOK, Response{
		Found:     found,
		AccountID: accountID})
}
