package lfc

import (
	"errors"

	redis "gopkg.in/redis.v5"
)

var globalRedisClient *RedisClient

type RedisClient struct {
	*redis.Client
}

// redisInstance...
func redisInstance() (instance *redis.Client, err error) {
	if globalRedisClient == nil {
		err = errors.New("global redis client is nil")
		return
	}
	instance = globalRedisClient.Client
	return
}
