package account

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/jackc/pgx/v5"
	"spire/lobby/internal/core"
)

func HandleBotAccountAuth(c *gin.Context, x *core.Context) {
	type Request struct {
		AccountID   uint64 `json:"account_id" binding:"required"`
		CharacterID uint64 `json:"character_id" binding:"required"`
	}

	type Response struct {
		Token string `json:"token"`
	}

	var r Request
	if !core.Check(c.Bind(&r), c, http.StatusBadRequest) {
		return
	}

	var characterId uint64
	err := x.P.QueryRow(context.Background(), "SELECT id FROM characters WHERE id=$1 AND account_id=$2", r.CharacterID, r.AccountID).Scan(&characterId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			core.Check(err, c, http.StatusUnauthorized)
			return
		}
		core.Check(err, c, http.StatusInternalServerError)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aid": strconv.FormatUint(r.AccountID, 10),
		"cid": strconv.FormatUint(r.CharacterID, 10),
		"prv": "None",

		"exp": jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	})
	signedString, err := token.SignedString([]byte(x.S.AuthKey))
	if !core.Check(err, c, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, Response{Token: signedString})
}
