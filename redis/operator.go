package redis

import (
	"context"
	"errors"
	"time"

	goRedis "github.com/redis/go-redis/v9"
)

type operator struct {
	rdb *goRedis.Client
}

// Operator redis操作封装，通过封装操作不需要判断err is nil,空数据返回默认值
func Operator(rdb *goRedis.Client) *operator {
	if rdb == nil {
		panic("redis operator: rdb is nil")
	}
	return &operator{rdb}
}

func (rec *operator) GetStringValue(key string) (value string, err error) {
	value, rdbErr := rec.rdb.Get(context.Background(), key).Result()
	if rdbErr != nil && !NilError(rdbErr) {
		return "", rdbErr
	}
	return
}
func (rec *operator) SetValue(key string, value interface{}, expire time.Duration) (err error) {
	ret := rec.rdb.Set(context.Background(), key, value, expire)
	if ret.Err() != nil {
		return ret.Err()
	}
	if ret.Val() != "OK" {
		return errors.New("glib:set redis value not ok")
	}
	return
}
func (rec *operator) DelKey(key string) (n int64, err error) {
	ret := rec.rdb.Del(context.Background(), key)
	if ret.Err() != nil {
		err = ret.Err()
		return
	}
	n = ret.Val()
	return
}
func (rec *operator) Renew(key string, expire time.Duration) (n bool, err error) {
	s := rec.rdb.TTL(context.Background(), key).Val()
	if s.Seconds() <= 0 {
		return
	}
	ret := rec.rdb.Expire(context.Background(), key, expire+s)
	if ret.Err() != nil {
		err = ret.Err()
		return
	}
	n = ret.Val()
	return
}
func (rec *operator) GetIntValue(key string) (value int64, err error) {
	value, rdbErr := rec.rdb.Get(context.Background(), key).Int64()
	if rdbErr != nil && !NilError(rdbErr) {
		err = rdbErr
		return
	}
	return
}
