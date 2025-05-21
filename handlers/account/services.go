package account

import (
	"context"
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"file-sharing/cache"
	"file-sharing/database"
)

func getUserFromCache(uid string, requestContext context.Context) (string, error) {
	return cache.Client.Get(
		requestContext,
		cache.CreateKey(cache.KeyPrefixes.Account, uid),
	).Result()
}

func getUserFromDatabase(
	uid string,
	user *database.Users,
	requestContext context.Context,
) error {
	return database.UsersCollection.FindOne(
		requestContext,
		bson.M{"uid": uid},
	).Decode(&user)
}

func removeUserFromCache(uid string, requestContext context.Context) error {
	_, delError := cache.Client.Del(
		requestContext,
		cache.CreateKey(cache.KeyPrefixes.Account, uid),
	).Result()
	return delError
}

func saveUserToCache(uid string, user any, requestContext context.Context) error {
	userBytes, jsonError := json.Marshal(user)
	if jsonError != nil {
		return jsonError
	}
	_, setError := cache.Client.Set(
		requestContext,
		cache.CreateKey(cache.KeyPrefixes.Account, uid),
		string(userBytes),
		time.Duration(time.Hour)*8,
	).Result()
	return setError
}
