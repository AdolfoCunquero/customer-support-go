package redis

import (
	"encoding/json"
	"errors"
	"fmt"

	mdls "customer-support/models"
	mdb "customer-support/mongo"

	"github.com/go-redis/redis/v8"
)

func GetOrCreateIncident(clientId string, agentId string) (string, error) {
	val, err := rdb.HGet(ctx, "incidents", clientId).Result()
	if err == redis.Nil {

		currentConv, err := GetOrCreateConversation(clientId)
		if err != nil {
			return "", err
		}

		newInc := mdls.Incident{ConversationId: currentConv, AgentId: agentId, Status: 1}
		newInc, err = mdb.SaveIncident(newInc)
		if err != nil {
			return "", err
		}

		obj := map[string]interface{}{"incidentId": newInc.ID.Hex(), "status": 1}
		args, _ := json.Marshal(obj)
		value := map[string]interface{}{clientId: string(args)}

		err = rdb.HSet(ctx, "incidents", value).Err()

		if err != nil {
			return "", err
		}
		return newInc.ID.Hex(), nil

	} else if err != nil {
		return "", err
	}
	return val, nil
}

func CloseIncident(clientId string) error {
	incStr, err := rdb.HGet(ctx, "incidents", clientId).Result()
	if err != nil {
		return errors.New("incident does not exists")
	}

	count := rdb.HDel(ctx, "incidents", clientId)
	if count.Val() == 0 {
		return errors.New("error to delete redis incident")
	}

	incData := make(map[string]interface{})
	errU := json.Unmarshal([]byte(incStr), &incData)
	if errU != nil {
		return errU
	}

	err = mdb.CloseIncident(fmt.Sprintf("%v", incData["incidentId"]))
	if err != nil {
		return err
	}
	return nil
}
