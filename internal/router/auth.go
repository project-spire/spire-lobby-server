package router

import (
	"database/sql"
	"errors"
	"net/http"
	"spire/lobby/internal/core"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthRequest struct {
	AccountId uint64 `json:"account_id" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func HandleAuthBot(c *gin.Context, ctx *core.Context) {
	var request AuthRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	queryChan := make(chan error, 1)
	go func() {
		row := ctx.Db.QueryRow("SELECT 1 FROM accounts WHERE account_id = ?", request.AccountId)
		queryChan <- row.Scan()
		close(queryChan)
	}()
	if err := <-queryChan; err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"account_id": request.AccountId,
	})
	s, err := t.SignedString(ctx.Settings.AuthKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token signing failed"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token: s,
	})
}
