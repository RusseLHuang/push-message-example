package pushregistry

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var ctx = context.Background()
var redisClient *redis.Client

func InitClientConnection() {
	if redisClient != nil {
		return
	}

	registryHost := viper.Get("pushRegistryHost")
	registryPort := viper.Get("pushRegistryPort")
	uri := fmt.Sprintf("%s:%v", registryHost, registryPort)

	redisClient = redis.NewClient(&redis.Options{
		Addr:     uri,
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
		fmt.Println("clientID does not exist", clientID)
	} else if err != nil {
		panic(err)
	}

	return val
}
