package handlers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"

	"file-sharing/cache"
	"file-sharing/database"
)

func deleteCachedRecord(key string) error {
	return cache.Client.Del(
		context.Background(),
		key,
	).Err()
}

func getCachedRecord(key string) (*database.Files, error) {
	var record database.Files
	cachedValue, cacheError := cache.Client.Get(context.Background(), key).Result()
	if cacheError != nil {
		if cacheError == redis.Nil {
			return nil, nil
		}
		return nil, cacheError
	}
	decodeError := json.Unmarshal([]byte(cachedValue), &record)
	if decodeError != nil {
		return nil, decodeError
	}
	return &record, nil
}

func setCacheValue(key string, value *database.Files) error {
	stringValue, encodeError := json.Marshal(&value)
	if encodeError != nil {
		return encodeError
	}
	cache.Client.Set(
		context.Background(),
		key,
		stringValue,
		time.Duration(time.Hour)*8,
	)
	return nil
}
