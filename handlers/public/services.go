package public

import (
	"context"
	"encoding/json"
	"path/filepath"
	"time"

	"file-sharing/cache"
	"file-sharing/constants"
	"file-sharing/utilities"
)

func createFilePath(uid string) string {
	return filepath.Join(
		utilities.GetEnv(
			constants.ENV_NAMES.UplaodsDirectoryName,
			constants.DEFAULT_UPLOADS_DIRECTORY_NAME,
		),
		uid,
	)
}

func getFromCache(uid string, requestContext context.Context) (string, error) {
	return cache.Client.Get(
		requestContext,
		cache.CreateKey(cache.KeyPrefixes.File, uid),
	).Result()
}

func removeFromCache(uid string, requestContext context.Context) error {
	_, delError := cache.Client.Del(
		context.Background(),
		cache.CreateKey(cache.KeyPrefixes.File, uid),
	).Result()
	return delError
}

func saveToCache(uid string, value any, requestContext context.Context) error {
	encoded, encodeError := json.Marshal(&value)
	if encodeError != nil {
		return encodeError
	}
	_, setError := cache.Client.Set(
		requestContext,
		cache.CreateKey(cache.KeyPrefixes.File, uid),
		string(encoded),
		time.Duration(time.Hour)*8,
	).Result()
	return setError
}
