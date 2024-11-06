package cache

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"

	"file-exchange/constants"
	"file-exchange/utilities"
)

var Client *redis.Client

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

		log.Printf("Redis connection failed, repeating in %d seconds", i)
		time.Sleep(time.Duration(i) * time.Second)
	}

	log.Println("Redis connection is ready")
}
