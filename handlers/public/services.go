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

func removeFromCache(key string) error {
	_, delError := cache.Client.Del(context.Background(), key).Result()
	return delError
}

func saveToCache(key string, value interface{}) error {
	encoded, encodeError := json.Marshal(&value)
	if encodeError != nil {
		return encodeError
	}
	_, setError := cache.Client.Set(
		context.Background(),
		key,
		encoded,
		time.Duration(time.Hour)*8,
	).Result()
	return setError
}
