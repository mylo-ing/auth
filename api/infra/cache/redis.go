package cache

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

// Initializes the Redis client from environment variables
func InitRedis() *redis.Client {
	host := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
	if host == "" {
		log.Fatalf("You must specificy a Redis host environment variable.")
	}

	//Set the cache db number. Keeping session cache separate from entity cache.
	dbStr := os.Getenv("REDIS_SESSION_DB")
	dbNum, err := strconv.Atoi(dbStr)
	if err != nil {
		dbNum = 0
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       dbNum,
	})

	// test connection
	_, err = rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	log.Println("Connected to Redis on", host)
	return rdb
}

// SetValue stores a string value in Redis with an expiration
func SetValue(rdb *redis.Client, key, value string, expiration time.Duration) error {
	return rdb.Set(Ctx, key, value, expiration).Err()
}

// GetValue retrieves a string value from Redis
func GetValue(rdb *redis.Client, key string) (string, error) {
	return rdb.Get(Ctx, key).Result()
}

// DeleteKey removes a key from Redis
func DeleteKey(rdb *redis.Client, key string) error {
	return rdb.Del(Ctx, key).Err()
}
