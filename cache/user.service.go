package cache

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"

	"file-sharing/database"
)

type userService struct {
	client    *redis.Client
	keyPrefix string
}

var UserService userService

func (instance *userService) New(client *redis.Client, keyPrefix string) {
	instance.client = client
	instance.keyPrefix = keyPrefix
}

func (instance *userService) Del(
	operationContext context.Context,
	uid string,
) error {
	_, delError := instance.client.Del(
		operationContext,
		CreateKey(instance.keyPrefix, uid),
	).Result()
	return delError
}

func (instance *userService) Get(
	operationContext context.Context,
	uid string,
	destination *database.UserModel,
) error {
	cachedValue, cacheError := instance.client.Get(
		operationContext,
		CreateKey(instance.keyPrefix, uid),
	).Result()
	if cacheError != nil {
		return cacheError
	}
	return json.Unmarshal([]byte(cachedValue), &destination)
}

func (instance *userService) Set(
	operationContext context.Context,
	user database.UserModel,
) error {
	bytes, jsonError := json.Marshal(user)
	if jsonError != nil {
		return jsonError
	}
	_, setError := instance.client.Set(
		operationContext,
		CreateKey(instance.keyPrefix, user.Uid),
		string(bytes),
		defaultExpiration,
	).Result()
	return setError
}
