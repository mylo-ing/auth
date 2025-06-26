// api/infra/cache/redis.go
package cache

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisCodeCache struct {
	rdb *redis.Client
	ctx context.Context
}

func OpenRedis() CodeCache {
	addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	pwd := os.Getenv("REDIS_PASSWORD")
	dbNum, _ := strconv.Atoi(os.Getenv("REDIS_SESSION_DB"))

	var cc CodeCache
	for i := 1; i <= 10; i++ {
		cc = NewRedis(addr, pwd, dbNum)
		if err := cc.SetValue("healthz", "ok", time.Second); err == nil {
			cc.DeleteKey("healthz")
			log.Printf("connected to Redis %s (DB=%d)", addr, dbNum)
			return cc
		}
		log.Printf("Redis not ready (%d/10)", i)
		time.Sleep(2 * time.Second)
	}
	log.Fatalf("Redis unreachable at %s", addr)
	return nil
}

func NewRedis(addr, pwd string, db int) CodeCache {
	rdb := redis.NewClient(&redis.Options{Addr: addr, Password: pwd, DB: db})
	return &redisCodeCache{rdb: rdb, ctx: context.Background()}
}

func (c *redisCodeCache) SetValue(k, v string, ttl time.Duration) error {
	return c.rdb.Set(c.ctx, k, v, ttl).Err()
}
func (c *redisCodeCache) GetValue(k string) (string, error) {
	return c.rdb.Get(c.ctx, k).Result()
}
func (c *redisCodeCache) DeleteKey(k string) error {
	return c.rdb.Del(c.ctx, k).Err()
}
