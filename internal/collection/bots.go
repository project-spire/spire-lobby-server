package collection

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"spire/lobby/internal/core"
)

type Bot struct {
	BotID     uint64        `bson:"bot_id"`
	AccountID bson.ObjectID `bson:"account_id"`
}

func FindBot(x *core.Context, botId uint64) (Bot, error) {
	collection := x.Collection("bots")
	filter := bson.D{{"bot_id", botId}}

	var bot Bot
	err := collection.FindOne(context.Background(), filter).Decode(&bot)

	return bot, err
}

func InsertBot(x *core.Context, bot *Bot) (*mongo.InsertOneResult, error) {
	collection := x.Collection("bots")
	return collection.InsertOne(context.Background(), bot)
}
