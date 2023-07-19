package driver

import "github.com/redis/go-redis/v9"

type RedisDB struct {
	Redis *redis.Client
}

func ConnectRedis(address, password string, db int) *RedisDB {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	return &RedisDB{
		Redis: rdb,
	}
}
