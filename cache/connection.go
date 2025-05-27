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

const defaultExpiration time.Duration = time.Duration(time.Hour) * 8

func Connect() {
	redisHost := utilities.GetEnv(
		constants.ENV_NAMES.RedisHost,
		constants.DEFAULT_REDIS_HOST,
	)
	redisPort := utilities.GetEnv(
		constants.ENV_NAMES.RedisPort,
		constants.DEFAULT_REDIS_PORT,
	)
	redisPassword := utilities.GetEnv(constants.ENV_NAMES.RedisPassword)

	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
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

	FileService.New(Client, KeyPrefixes.File)
	UserService.New(Client, KeyPrefixes.User)

	log.Println("Redis connection is ready")
}

func CreateKey(prefix, value string) string {
	return fmt.Sprintf("%s-%s", prefix, value)
}
