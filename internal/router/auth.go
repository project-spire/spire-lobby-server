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

func HandleBotAuth(c *gin.Context, ctx *core.Context) {
	var r AuthRequest
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	row := ctx.Db.QueryRow("SELECT 1 FROM accounts WHERE account_id = ?", r.AccountId)
	var one int
	if err := row.Scan(&one); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"account_id": r.AccountId,
	})
	s, err := t.SignedString([]byte(ctx.Settings.AuthKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token signing failed"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token: s,
	})
}
