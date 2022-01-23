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

func GetActiveIncidents(agentId string) ([]*mdls.ActiveIncident, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := mongoCN.Database(database)
	col := db.Collection(incidentCollection)

	var results = make([]*mdls.ActiveIncident, 0)
	opts := make([]bson.M, 0)

	opts = append(opts, bson.M{
		"$match": bson.M{
			"status":  1,
			"agentId": agentId,
		},
	})

	opts = append(opts, bson.M{
		"$addFields": bson.M{
			"customerIdF": bson.M{
				"$toObjectId": "$customerId",
			},
		},
	})

	opts = append(opts, bson.M{
		"$lookup": bson.M{
			"from":         customerCollection,
			"localField":   "customerIdF",
			"foreignField": "_id",
			"as":           "customerInfo",
		},
	})

	opts = append(opts, bson.M{
		"$project": bson.M{
			"_id":                 1,
			"agentId":             1,
			"conversationId":      1,
			"customerId":          1,
			"status":              1,
			"customerInfo._id":    1,
			"customerInfo.name":   1,
			"customerInfo.origin": 1,
		},
	})

	opts = append(opts, bson.M{
		"$unwind": "$customerInfo",
	})

	cursor, err := col.Aggregate(ctx, opts)

	if err != nil {
		return results, err
	}

	for cursor.Next(context.TODO()) {

		var dir mdls.ActiveIncident
		err := cursor.Decode(&dir)
		if err != nil {
			return results, err
		}
		results = append(results, &dir)
	}
	return results, nil
}

func GetActiveConversation(customerId string) ([]*mdls.ActiveConversationInc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := mongoCN.Database(database)
	col := db.Collection(incidentCollection)

	var results = make([]*mdls.ActiveConversationInc, 0)
	opts := make([]bson.M, 0)

	opts = append(opts, bson.M{
		"$match": bson.M{
			"status":     1,
			"customerId": customerId,
		},
	})

	opts = append(opts, bson.M{
		"$addFields": bson.M{
			"incidentIdStr": bson.M{
				"$toString": "$_id",
			},
		},
	})

	opts = append(opts, bson.M{
		"$lookup": bson.M{
			"from":         messageCollection,
			"localField":   "incidentIdStr",
			"foreignField": "incidentId",
			"as":           "messages",
		},
	})

	opts = append(opts, bson.M{
		"$project": bson.M{
			"_id":            1,
			"agentId":        1,
			"conversationId": 1,
			"customerId":     1,
			"status":         1,
			"messages":       1,
		},
	})

	cursor, err := col.Aggregate(ctx, opts)

	if err != nil {
		return results, err
	}

	for cursor.Next(context.TODO()) {

		var msg mdls.ActiveConversationInc
		err := cursor.Decode(&msg)
		if err != nil {
			return results, err
		}
		results = append(results, &msg)
	}
	return results, nil
}
