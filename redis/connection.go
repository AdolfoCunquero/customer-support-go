package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var rdb = ConnectionDB()
var ctx = context.Background()

func ConnectionDB() *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	fmt.Println("Connected to Redis")
	return rdb
}
