package cache

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"time"
)

type Redis struct {
	DB *redis.Client
}

func (a *Redis) NewHelper(key string, timeout time.Duration) *RedisHelper {
	return &RedisHelper{
		Redis:   *a,
		Timeout: timeout,
		Key:     key,
	}
}

func (a *Redis) INC(key string) (int64, error) {
	return a.DB.Incr(context.Background(), key).Result()
}

func (a *Redis) Cache(key string, v string, t time.Duration) error {
	return a.DB.Set(context.Background(), key, v, t).Err()
}

func (a *Redis) CacheStruct(key string, v interface{}, t time.Duration) error {
	j, e := json.Marshal(v)
	if e != nil {
		return e
	}
	return a.Cache(key, string(j), t)
}

func (a *Redis) Read(key string) (string, error) {
	s, e := a.DB.Get(context.Background(), key).Result()
	return s, e
}

func (a *Redis) ReadStruct(key string, v interface{}) error {
	t, e := a.Read(key)
	if e != nil {
		return e
	}
	return json.Unmarshal([]byte(t), v)
}

func (a *Redis) Del(key string) error {
	return a.DB.Del(context.Background(), key).Err()
}

func (a *Redis) SetExpr(key string, t time.Duration) error {
	return a.DB.Expire(context.Background(), key, t).Err()
}

func (a *Redis) Flush() {
	a.DB.FlushDB(context.Background())
}
