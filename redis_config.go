package lfc

import (
	"time"

	redis "gopkg.in/redis.v5"
)

type redisConfig struct {
	Address      string `json:"address" toml:"address" yaml:"address"`
	Password     string `json:"password" toml:"password" yaml:"password"`
	DB           int    `json:"db" toml:"db" yaml:"db"`
	PoolSize     int    `json:"poolsize" toml:"poolsize" yaml:"poolsize"`
	DialTimeout  int    `json:"DialTimeout" toml:"dial_timeout" yaml:"dial_timeout"`    // 毫秒
	ReadTimeout  int    `json:"ReadTimeout" toml:"read_timeout" yaml:"read_timeout"`    // 毫秒
	WriteTimeout int    `json:"WriteTimeout" toml:"write_timeout" yaml:"write_timeout"` // 毫秒
}

func (rc *redisConfig) newRedis() *RedisClient {
	options := &redis.Options{
		Addr:         rc.Address,
		Password:     rc.Address,
		DB:           rc.DB,
		PoolSize:     rc.PoolSize,
		DialTimeout:  time.Duration(rc.DialTimeout) * time.Microsecond,
		ReadTimeout:  time.Duration(rc.ReadTimeout) * time.Microsecond,
		WriteTimeout: time.Duration(rc.WriteTimeout) * time.Microsecond,
	}
	redisClient := &RedisClient{redis.NewClient(options)}
	return redisClient
}
