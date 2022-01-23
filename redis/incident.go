package redis

import (
	"errors"

	mdls "customer-support/models"
	mdb "customer-support/mongo"

	"github.com/go-redis/redis/v8"
)

const setIncidentClient = "incident:client:id:"
const setIncidentAgent = "incident:agent:id:"
const hmapIncident = "incident:id:"

func GetOrCreateIncident(clientId string, agentId string) (string, error) {
	val, err := rdb.SMembers(ctx, setIncidentClient+clientId).Result()
	if err == redis.Nil || len(val) == 0 {

		currentConv, err := GetOrCreateConversation(clientId)
		if err != nil {
			return "", err
		}

		newInc := mdls.Incident{ConversationId: currentConv, AgentId: agentId, Status: 1, CustomerId: clientId}
		newInc, err = mdb.SaveIncident(newInc)
		if err != nil {
			return "", err
		}

		var incidentId string = newInc.ID.Hex()

		err = rdb.SAdd(ctx, setIncidentClient+clientId, incidentId).Err()
		if err != nil {
			return "", err
		}

		inc := map[string]interface{}{"clientId": clientId, "agentId": agentId}

		err = rdb.HSet(ctx, hmapIncident+incidentId, inc).Err()
		if err != nil {
			return "", err
		}

		err = rdb.SAdd(ctx, setIncidentAgent+agentId, incidentId).Err()
		if err != nil {
			return "", err
		}

		return incidentId, nil

	} else if err != nil {
		return "", err
	}
	return val[0], nil
}

func CloseIncident(clientId string) error {
	clientKey := setIncidentClient + clientId
	clientInc, err := rdb.SMembers(ctx, clientKey).Result()

	if err != nil || len(clientInc) == 0 {
		return errors.New("incident does not exists")
	}

	incidentId := clientInc[0]
	agentId := rdb.HGet(ctx, hmapIncident+incidentId, "agentId").Val()

	count := rdb.Del(ctx, clientKey)
	if count.Val() == 0 {
		return errors.New("error to delete redis incident")
	}

	count = rdb.Del(ctx, hmapIncident+incidentId)
	if count.Val() == 0 {
		return errors.New("error to delete redis incident")
	}

	count = rdb.SRem(ctx, setIncidentAgent+agentId, incidentId)
	if count.Val() == 0 {
		return errors.New("error to delete redis incident")
	}

	err = mdb.CloseIncident(incidentId)
	if err != nil {
		return err
	}
	return nil
}
