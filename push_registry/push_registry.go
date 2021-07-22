package pushregistry

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

type PushRegistry struct {
	client  *redis.Client
	context context.Context
}

func NewPushRegistry() *PushRegistry {
	registryHost := viper.Get("pushRegistryHost")
	registryPort := viper.Get("pushRegistryPort")
	uri := fmt.Sprintf("%s:%v", registryHost, registryPort)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     uri,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &PushRegistry{
		client:  redisClient,
		context: context.Background(),
	}
}

func (pr *PushRegistry) SetPersistentConnectionID(clientID string, nodeIP string) {
	err := pr.client.Set(pr.context, clientID, nodeIP, 0).Err()
	if err != nil {
		panic(err)
	}
}

func (pr *PushRegistry) GetPersistentConnectionID(clientID string) string {
	val, err := pr.client.Get(pr.context, clientID).Result()
	if err == redis.Nil {
		fmt.Println("clientID does not exist", clientID)
		return ""
	} else if err != nil {
		panic(err)
	}

	return val
}
