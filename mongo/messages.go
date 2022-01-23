package mongo

import (
	"context"
	"time"

	mdls "customer-support/models"

	"go.mongodb.org/mongo-driver/bson"
)

func SaveMessage(msg mdls.Message) (mdls.Message, error) {
	var result mdls.Message
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	db := mongoCN.Database(database)
	col := db.Collection(messageCollection)

	inserted, errInsrt := col.InsertOne(ctx, msg)

	if errInsrt != nil {
		return result, err
	}

	err = col.FindOne(ctx, bson.M{"_id": inserted.InsertedID}).Decode(&result)

	if err != nil {
		return result, err
	}

	return result, nil
}
