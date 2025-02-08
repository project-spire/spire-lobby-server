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
	AccountId   uint64 `json:"account_id" binding:"required"`
	CharacterId uint64 `json:"character_id" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func HandleBotAuth(c *gin.Context, x *core.Context) {
	var r AuthRequest
	if !check(c.Bind(&r), c, http.StatusBadRequest) {
		return
	}

	row := x.D.QueryRow("SELECT 1 FROM accounts WHERE account_id = ?", r.AccountId)
	var one int
	if err := row.Scan(&one); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			check(err, c, http.StatusUnauthorized)
			return
		}
		check(err, c, http.StatusInternalServerError)
		return
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"account_id":   r.AccountId,
		"character_id": r.CharacterId,
	})
	s, err := t.SignedString([]byte(x.S.AuthKey))
	if !check(err, c, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token: s,
	})
}
