package account

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"spire/lobby/internal/collection"
	"spire/lobby/internal/core"
)

func HandleBotAccountCreate(c *gin.Context, x *core.Context) {
	type Request struct {
		BotID uint64 `json:"bot_id" binding:"required"`
	}

	type Response struct {
		AccountID bson.ObjectID `json:"account_id"`
	}

	var r Request
	if !core.Check(c.Bind(&r), c, http.StatusBadRequest) {
		return
	}

	session, err := x.StartSession()
	if err != nil {
		core.Check(err, c, http.StatusInternalServerError)
		return
	}
	defer session.EndSession(context.Background())

	var accountID bson.ObjectID
	_, err = session.WithTransaction(context.Background(), func(ctx context.Context) (interface{}, error) {
		res, err := collection.InsertAccount(x)
		if err != nil {
			return nil, err
		}
		accountID = res.InsertedID.(bson.ObjectID)

		_, err = collection.InsertBot(x, &collection.Bot{
			BotID:     r.BotID,
			AccountID: accountID})
		if err != nil {
			return nil, err
		}

		return nil, session.CommitTransaction(ctx)
	})
	if err != nil {
		_ = session.AbortTransaction(context.Background())
		core.Check(err, c, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, Response{
		AccountID: accountID,
	})
}
