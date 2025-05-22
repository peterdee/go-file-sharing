package cache

import "github.com/redis/go-redis/v9"

type CommonOperations struct {
	Client *redis.Client
}

type KeyPrefixesStruct struct {
	File string
	User string
}
