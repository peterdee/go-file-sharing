package cache

import (
	"context"
	"encoding/json"
	"time"
)

const defaultExpiration time.Duration = time.Duration(time.Hour) * 8

var Operations CommonOperations

func (operations *CommonOperations) getFromCache(
	key string,
	target any,
	operationContext context.Context,
) error {
	cachedValue, cacheError := operations.Client.Get(operationContext, key).Result()
	if cacheError != nil {
		return cacheError
	}
	return json.Unmarshal([]byte(cachedValue), &target)
}

func (operations *CommonOperations) removeFromCache(
	key string,
	operationContext context.Context,
) error {
	_, delError := operations.Client.Del(operationContext, key).Result()
	return delError
}

func (operations *CommonOperations) saveToCache(
	key string,
	value any,
	expiration time.Duration,
	operationContext context.Context,
) error {
	bytes, jsonError := json.Marshal(value)
	if jsonError != nil {
		return jsonError
	}
	_, setError := operations.Client.Set(
		operationContext,
		key,
		string(bytes),
		expiration,
	).Result()
	return setError
}

func (operations *CommonOperations) GetFile(
	uid string,
	file any,
	requestContext context.Context,
) error {
	return operations.getFromCache(CreateKey(KeyPrefixes.File, uid), file, requestContext)
}

func (operations *CommonOperations) GetUser(
	uid string,
	user any,
	requestContext context.Context,
) error {
	return operations.getFromCache(CreateKey(KeyPrefixes.User, uid), user, requestContext)
}

func (operations *CommonOperations) RemoveFile(
	uid string,
	requestContext context.Context,
) error {
	return operations.removeFromCache(CreateKey(KeyPrefixes.File, uid), requestContext)
}

func (operations *CommonOperations) RemoveUser(
	uid string,
	requestContext context.Context,
) error {
	return operations.removeFromCache(CreateKey(KeyPrefixes.User, uid), requestContext)
}

func (operations *CommonOperations) SaveFile(
	uid string,
	file any,
	requestContext context.Context,
) error {
	return operations.saveToCache(
		CreateKey(KeyPrefixes.File, uid),
		file,
		defaultExpiration,
		requestContext,
	)
}

func (operations *CommonOperations) SaveUser(
	uid string,
	user any,
	requestContext context.Context,
) error {
	return operations.saveToCache(
		CreateKey(KeyPrefixes.User, uid),
		user,
		defaultExpiration,
		requestContext,
	)
}
