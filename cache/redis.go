package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"

	"file-sharing/constants"
	"file-sharing/utilities"
)

var Client *redis.Client

var NilError error = redis.Nil

var KeyPrefixes = KeyPrefixesStruct{
	File: "file",
	User: "user",
}

func Connect() {
	redisHost := utilities.GetEnv(
		constants.ENV_NAMES.RedisHost,
		constants.DEFAULT_REDIS_HOST,
	)
	redisPassword := utilities.GetEnv(constants.ENV_NAMES.RedisPassword)

	Client = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		DB:       0,
	})

	context := context.Background()
	for i := 1; i <= 6; i += 1 {
		pingError := Client.Ping(context).Err()
		if pingError == nil {
			break
		}
		if i == 6 {
			log.Fatal(pingError)
		}

		log.Printf("Redis connection failed, retry in %d seconds", i)
		time.Sleep(time.Duration(i) * time.Second)
	}

	Operations.Client = Client

	log.Println("Redis connection is ready")
}

func CreateKey(prefix, value string) string {
	return fmt.Sprintf("%s-%s", prefix, value)
}
