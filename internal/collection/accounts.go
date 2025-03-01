package collection

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"spire/lobby/internal/core"
)

type Account struct {
	ID         bson.ObjectID `bson:"_id"`
	Characters []Character   `bson:"characters,omitempty"`
}

type Character struct {
	Name string `bson:"name"`
}

func FindAccount(x *core.Context, id bson.ObjectID) (Account, error) {
	coll := x.Collection("accounts")
	filter := bson.D{{"_id", id}}

	var account Account
	err := coll.FindOne(context.Background(), filter).Decode(&account)

	return account, err
}

func InsertAccount(x *core.Context) (*mongo.InsertOneResult, error) {
	coll := x.Collection("accounts")
	return coll.InsertOne(context.Background(), bson.D{})
}

func InsertCharacter(x *core.Context, account *Account, character *Character) error {
	coll := x.Collection("accounts")
	filter := bson.D{{"_id", account.ID}}
	update := bson.D{{"$push", bson.D{{"characters", character}}}}

	_, err := coll.UpdateOne(context.Background(), filter, update)
	return err
}
