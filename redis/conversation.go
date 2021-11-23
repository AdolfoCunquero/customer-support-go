package redis

import (
	"errors"

	mdls "customer-support/models"
	mdb "customer-support/mongo"

	"github.com/go-redis/redis/v8"
)

const setConversationClient = "conversation:client:id:"
const hmapConversation = "conversation:id:"

func GetOrCreateConversation(clientId string) (string, error) {
	val, err := rdb.SMembers(ctx, setConversationClient+clientId).Result()
	if err == redis.Nil || len(val) == 0 {
		newConv := mdls.Conversation{ClientId: clientId, Status: 1}
		newConv, err = mdb.SaveConversation(newConv)
		if err != nil {
			return "", err
		}

		var conversationID string = newConv.ID.Hex()

		err = rdb.SAdd(ctx, setConversationClient+clientId, conversationID).Err()
		if err != nil {
			return "", err
		}

		conv := map[string]interface{}{"clientId": clientId}

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

func CloseConversation(clientId string) error {
	clientKey := setConversationClient + clientId
	clientConv, err := rdb.SMembers(ctx, clientKey).Result()

	if err != nil || len(clientConv) == 0 {
		return errors.New("conversation does not exists")
	}

	var conversationId string = clientConv[0]

	errI := CloseIncident(clientId)
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
