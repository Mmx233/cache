package cache

import (
	"context"
	"encoding/json"
	"time"
)

type RedisHelper struct {
	Redis
	Timeout time.Duration
	Key     string
}

func (a *RedisHelper) INC() (int64, error) {
	num, e := a.Redis.INC(a.Key)
	if e != nil {
		return 0, e
	}
	return num, a.SetExpr()
}

func (a *RedisHelper) Cache(v string) error {
	return a.Redis.Cache(a.Key, v, a.Timeout)
}

func (a *RedisHelper) CacheStruct(v interface{}) error {
	j, e := json.Marshal(v)
	if e != nil {
		return e
	}
	return a.Cache(string(j))
}

func (a *RedisHelper) Read() (string, error) {
	v, e := a.Redis.Read(a.Key)
	if e != nil {
		return "", e
	}
	_ = a.SetExpr()
	return v, nil
}

func (a *RedisHelper) ReadStruct(v interface{}) error {
	t, e := a.Read()
	if e != nil {
		return e
	}
	_ = a.SetExpr()
	return json.Unmarshal([]byte(t), v)
}

func (a *RedisHelper) Del() error {
	return a.Redis.Del(a.Key)
}

func (a *RedisHelper) SetExpr() error {
	return a.DB.Expire(context.Background(), a.Key, a.Timeout).Err()
}

func (a *RedisHelper) Flush() {
	a.DB.FlushDB(context.Background())
}
