package redis

import (
	mdls "customer-support/models"
)

func GetOnlineUsers() ([]mdls.User, error) {
	result := make([]mdls.User, 0)
	users, err := rdb.HGetAll(ctx, "onlineUsers").Result()
	if err != nil {
		return result, err
	}

	for _, val := range users {
		user := mdls.User{Username: val}
		result = append(result, user)
	}

	return result, nil
}

func SetOfflineUser(username string) error {
	return rdb.HDel(ctx, "onlineUsers", username).Err()
}
