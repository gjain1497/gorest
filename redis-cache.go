package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisCache struct {
	host    string
	db      int
	expires time.Duration
}

func NewRedisCache(host string, db int, exp time.Duration) UserCache {
	return &redisCache{
		host:    host,
		db:      db,
		expires: exp,
	}
}

func (cache *redisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.host,
		Password: "",
		DB:       cache.db,
	})
}

func (cache *redisCache) Set(key string, value User) {
	client := cache.getClient()
	ctx := context.Background()

	json, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}

	client.Set(ctx, key, json, cache.expires*time.Second)
}

func (cache *redisCache) Get(key string) *User {
	client := cache.getClient()
	ctx := context.Background()

	val, err := client.Get(ctx, key).Result()
	if err != nil {
		return nil
	}

	user := User{}
	err = json.Unmarshal([]byte(val), &user)

	if err != nil {
		panic(err)
	}
	return &user
}
