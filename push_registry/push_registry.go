package pushregistry

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var redisClient *redis.Client

func InitClientConnection() {
	if redisClient != nil {
		return
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func SetPersistentConnectionID(clientID string, nodeIP string) {
	err := redisClient.Set(ctx, clientID, nodeIP, 0).Err()
	if err != nil {
		panic(err)
	}
}

func GetPersistentConnectionID(clientID string) string {
	val, err := redisClient.Get(ctx, clientID).Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	}

	return val
}
