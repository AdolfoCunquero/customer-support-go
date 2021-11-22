package mongo

import (
	"context"
	"time"

	mdls "customer-support/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SaveConversation(msg mdls.Conversation) (mdls.Conversation, error) {
	var result mdls.Conversation
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	db := mongoCN.Database(database)
	col := db.Collection(conversationCollection)

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

func CloseConversation(conversationId string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := mongoCN.Database(database)
	col := db.Collection(conversationCollection)

	updateString := bson.M{
		"$set": map[string]interface{}{"status": 2},
	}

	objId, _ := primitive.ObjectIDFromHex(conversationId)

	_, err := col.UpdateOne(ctx, bson.M{
		"_id": bson.M{"$eq": objId},
	}, updateString)

	if err != nil {
		return err
	}

	return nil

}
