package account

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"spire/lobby/internal/core"
)

func HandleBotAccountCreate(c *gin.Context, x *core.Context) {
	type Request struct {
		BotID uint64 `json:"bot_id" binding:"required"`
	}

	type Response struct {
		AccountID uint64 `json:"account_id"`
	}

	var r Request
	if !core.Check(c.Bind(&r), c, http.StatusBadRequest) {
		return
	}

	ctx := context.Background()
	tx, err := x.P.Begin(ctx)
	if err != nil {
		core.Check(err, c, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	var accountID uint64
	err = tx.QueryRow(ctx, "INSERT INTO accounts DEFAULT VALUES RETURNING id").Scan(&accountID)
	if err != nil {
		core.Check(err, c, http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(ctx, "INSERT INTO bots (id, account_id) VALUES ($1, $2)", r.BotID, accountID)
	if err != nil {
		core.Check(err, c, http.StatusInternalServerError)
		return
	}

	if tx.Commit(ctx) != nil {
		core.Check(err, c, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, Response{
		AccountID: accountID,
	})
}
