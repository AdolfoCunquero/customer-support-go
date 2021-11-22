package redis

import (
	"errors"
	"fmt"

	mdls "customer-support/models"
	mdb "customer-support/mongo"

	"github.com/go-redis/redis/v8"
)

func GetOrCreateConversation(clientId string) (string, error) {
	val, err := rdb.HGet(ctx, "conversations", clientId).Result()
	if err == redis.Nil {
		newConv := mdls.Conversation{ClientId: clientId, Status: 1}
		newConv, err = mdb.SaveConversation(newConv)
		if err != nil {
			return "", err
		}

		value := map[string]interface{}{clientId: newConv.ID.Hex()}

		err = rdb.HSet(ctx, "conversations", value).Err()
		if err != nil {
			return "", err
		}
		return newConv.ID.Hex(), nil

	} else if err != nil {
		return "", err
	}
	return val, nil
}

func CloseConversation(clientId string) error {
	conversationId, err := rdb.HGet(ctx, "conversations", clientId).Result()
	if err != nil {
		return errors.New("conversation does not exists")
	}

	errI := CloseIncident(clientId)
	if err != nil {
		return errI
	}

	count := rdb.HDel(ctx, "conversations", clientId)
	if count.Val() == 0 {
		return errors.New("error to delete redis conversation")
	}

	err = mdb.CloseConversation(conversationId)
	if err != nil {
		return err
	}
	return nil
}

func SetExample() {
	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}

	//Output: key value
	//key2 does not exist

}
