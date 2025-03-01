package character

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"spire/lobby/internal/collection"
	"spire/lobby/internal/core"
)

func HandleBotCharacterCreate(c *gin.Context, x *core.Context) {
	type Request struct {
		AccountID     bson.ObjectID `json:"account_id" binding:"required"`
		CharacterName string        `json:"character_name" binding:"required"`
	}

	type Response struct{}

	var r Request
	if !core.Check(c.Bind(&r), c, http.StatusBadRequest) {
		return
	}

	account, err := collection.FindAccount(x, r.AccountID)
	if err != nil {
		core.Check(err, c, http.StatusUnauthorized)
		return
	}

	//TODO: Check duplicated character name
	err = collection.InsertCharacter(x, &account, &collection.Character{
		Name: r.CharacterName,
	})
	if err != nil {
		core.Check(err, c, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, Response{})
}
