package redis

import (
	"errors"

	mdls "customer-support/models"
	mdb "customer-support/mongo"

	"github.com/go-redis/redis/v8"
)

const setConversationClient = "conversation:client:id:"
const hmapConversation = "conversation:id:"

func GetOrCreateConversation(customerId string) (string, error) {
	val, err := rdb.SMembers(ctx, setConversationClient+customerId).Result()
	if err == redis.Nil || len(val) == 0 {
		newConv := mdls.Conversation{CustomerId: customerId, Status: 1}
		newConv, err = mdb.SaveConversation(newConv)
		if err != nil {
			return "", err
		}

		var conversationID string = newConv.ID.Hex()

		err = rdb.SAdd(ctx, setConversationClient+customerId, conversationID).Err()
		if err != nil {
			return "", err
		}

		conv := map[string]interface{}{"customerId": customerId}

		err = rdb.HSet(ctx, hmapConversation+conversationID, conv).Err()
		if err != nil {
			return "", err
		}

		return conversationID, nil

	} else if err != nil {
		return "", err
	}
	return val[0], nil
}

func CloseConversation(customerId string) error {
	clientKey := setConversationClient + customerId
	clientConv, err := rdb.SMembers(ctx, clientKey).Result()

	if err != nil || len(clientConv) == 0 {
		return errors.New("conversation does not exists")
	}

	var conversationId string = clientConv[0]

	errI := CloseIncident(customerId)
	if err != nil {
		return errI
	}

	count := rdb.Del(ctx, clientKey)
	if count.Val() == 0 {
		return errors.New("error to delete redis conversation")
	}

	count = rdb.Del(ctx, hmapConversation+conversationId)
	if count.Val() == 0 {
		return errors.New("error to delete redis conversation")
	}

	err = mdb.CloseConversation(conversationId)
	if err != nil {
		return err
	}
	return nil
}
