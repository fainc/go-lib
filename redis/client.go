package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	goRedis "github.com/redis/go-redis/v9"
)

var rdb *goRedis.Client

type NewRedisConf struct {
	Address  string `json:"address"`
	User     string `json:"user"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// Init redis全局初始化
func Init(conf NewRedisConf) (err error) {
	if rdb != nil {
		err = errors.New("golib redis实例已全局初始化，请勿重复初始化")
		return
	}
	d := goRedis.NewClient(&goRedis.Options{
		Addr:     conf.Address,
		Password: conf.Password,
		Username: conf.User,
		DB:       conf.DB,
	})
	dErr := d.Get(context.Background(), "lib_test_init_connect").Err()
	if dErr != nil && !errors.Is(dErr, goRedis.Nil) {
		err = errors.New("GoLib: redis init failed : " + dErr.Error())
		fmt.Println(err.Error())
		return
	}
	rdb = d
	fmt.Println("GoLib: Redis 初始化成功")
	return
}

// GetClient 获取 redis 全局实例
func GetClient() (client *goRedis.Client, err error) {
	if rdb == nil {
		err = errors.New("golib redis实例未全局初始化，请在全局启动时使用golib包的 redis.Init()进行初始化")
		return
	}
	return rdb, nil
}

// NilError 判断结果是否为空
func NilError(err error) bool {
	return errors.Is(err, goRedis.Nil)
}

// GetStringValue 获取redis key的value string
func GetStringValue(key string) (value string, err error) {
	client, err := GetClient()
	if err != nil {
		return
	}
	value, rdbErr := client.Get(context.Background(), key).Result()
	if rdbErr != nil && !NilError(rdbErr) {
		err = rdbErr
		return
	}
	return
}
func SetValue(key string, value interface{}, expire time.Duration) (err error) {
	client, err := GetClient()
	if err != nil {
		return
	}
	ret := client.Set(context.Background(), key, value, expire)
	if ret.Err() != nil {
		err = ret.Err()
		return
	}
	if ret.Val() != "OK" {
		err = errors.New("glib:set redis value not ok")
		return
	}
	return
}
func DelKey(key string) (n int64, err error) {
	client, err := GetClient()
	if err != nil {
		return
	}
	ret := client.Del(context.Background(), key)
	if ret.Err() != nil {
		err = ret.Err()
		return
	}
	n = ret.Val()
	return
}
func Renew(key string, expire time.Duration) (n bool, err error) {
	client, err := GetClient()
	if err != nil {
		return
	}
	s := client.TTL(context.Background(), key).Val()
	if s.Seconds() <= 0 {
		return
	}
	ret := client.Expire(context.Background(), key, expire+s)
	if ret.Err() != nil {
		err = ret.Err()
		return
	}
	n = ret.Val()
	return
}
func GetIntValue(key string) (value int64, err error) {
	client, err := GetClient()
	if err != nil {
		return
	}
	value, rdbErr := client.Get(context.Background(), key).Int64()
	if rdbErr != nil && !NilError(rdbErr) {
		err = rdbErr
		return
	}
	return
}
