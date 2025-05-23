package cache

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"

	"file-sharing/database"
)

type fileService struct {
	client    *redis.Client
	keyPrefix string
}

var FileService fileService

func (instance *fileService) New(client *redis.Client, keyPrefix string) {
	instance.client = client
	instance.keyPrefix = keyPrefix
}

func (instance *fileService) Del(
	operationContext context.Context,
	uid string,
) error {
	_, delError := instance.client.Del(
		operationContext,
		CreateKey(instance.keyPrefix, uid),
	).Result()
	return delError
}

func (instance *fileService) DelMany(
	operationContext context.Context,
	uids ...string,
) error {
	keys := make([]string, len(uids))
	for index, uid := range uids {
		keys[index] = CreateKey(instance.keyPrefix, uid)
	}
	_, delError := instance.client.Del(
		operationContext,
		keys...,
	).Result()
	return delError
}

func (instance *fileService) Get(
	operationContext context.Context,
	uid string,
	destination *database.FileModel,
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

func (instance *fileService) Set(
	operationContext context.Context,
	file database.FileModel,
) error {
	bytes, jsonError := json.Marshal(file)
	if jsonError != nil {
		return jsonError
	}
	_, setError := instance.client.Set(
		operationContext,
		CreateKey(instance.keyPrefix, file.Uid),
		string(bytes),
		defaultExpiration,
	).Result()
	return setError
}
