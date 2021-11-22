package mongo

import (
	"context"
	"time"

	mdls "customer-support/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SaveIncident(msg mdls.Incident) (mdls.Incident, error) {
	var result mdls.Incident
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	db := mongoCN.Database(database)
	col := db.Collection(incidentCollection)

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

func CloseIncident(incidentId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := mongoCN.Database(database)
	col := db.Collection(incidentCollection)

	updateString := bson.M{
		"$set": map[string]interface{}{"status": 2},
	}

	objId, _ := primitive.ObjectIDFromHex(incidentId)

	_, err := col.UpdateOne(ctx, bson.M{
		"_id": bson.M{"$eq": objId},
	}, updateString)

	if err != nil {
		return err
	}

	return nil

}
