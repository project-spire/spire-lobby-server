package account

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"spire/lobby/internal/collection"
	"spire/lobby/internal/core"
)

func HandleBotAccountAuth(c *gin.Context, x *core.Context) {
	type Request struct {
		AccountID     bson.ObjectID `json:"account_id" binding:"required"`
		CharacterName string        `json:"character_name" binding:"required"`
	}

	type Response struct {
		Token string `json:"token"`
	}

	var r Request
	if !core.Check(c.Bind(&r), c, http.StatusBadRequest) {
		return
	}

	account, err := collection.FindAccount(x, r.AccountID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			core.Check(err, c, http.StatusUnauthorized)
			return
		}
		core.Check(err, c, http.StatusInternalServerError)
		return
	}

	var characterFound = false
	for _, character := range account.Characters {
		if character.Name == r.CharacterName {
			characterFound = true
			break
		}
	}
	if !characterFound {
		core.Check(err, c, http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"account_id":     r.AccountID.String(),
		"character_name": r.CharacterName,
	})
	signedString, err := token.SignedString([]byte(x.S.AuthKey))
	if !core.Check(err, c, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, Response{Token: signedString})
}
