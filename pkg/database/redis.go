package database

import "github.com/go-redis/redis/v9"

type RedisDB struct {
	DB *redis.Client
}

var Redis = &RedisDB{}

func NewRedisClient() *RedisDB {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})

	Redis.DB = rdb
	return Redis
}