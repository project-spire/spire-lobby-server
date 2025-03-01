package account

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"spire/lobby/internal/collection"
	"spire/lobby/internal/core"
)

func HandleBotAccountMe(c *gin.Context, x *core.Context) {
	type Request struct {
		BotID uint64 `json:"bot_id" binding:"required"`
	}

	type Response struct {
		Found     bool          `json:"found"`
		AccountID bson.ObjectID `json:"account_id"`
	}

	var r Request
	if !core.Check(c.Bind(&r), c, http.StatusBadRequest) {
		return
	}

	found := true
	bot, err := collection.FindBot(x, r.BotID)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			core.Check(err, c, http.StatusInternalServerError)
			return
		}

		found = false
	}

	c.JSON(http.StatusOK, Response{
		Found:     found,
		AccountID: bot.AccountID})
}
