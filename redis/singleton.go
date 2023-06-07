package redis

import (
	"context"
	"errors"

	goRedis "github.com/redis/go-redis/v9"
)

type singleton struct {
}

var singletonVar = singleton{}

func Singleton() *singleton {
	return &singletonVar
}

var rdb *goRedis.Client

type NewRedisConf struct {
	Address  string `json:"address"`
	User     string `json:"user"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// Init 初始化redis实例到内存
func (rec *singleton) Init(conf NewRedisConf) (err error) {
	if rdb != nil { // 忽略重新初始化
		return errors.New("redis.singleton:全局实例已初始化，请勿重新初始化")
	}
	d := goRedis.NewClient(&goRedis.Options{
		Addr:     conf.Address,
		Password: conf.Password,
		Username: conf.User,
		DB:       conf.DB,
	})
	dErr := d.Get(context.Background(), "redis_test_init_connect").Err()
	if dErr != nil && !errors.Is(dErr, goRedis.Nil) {
		return errors.New("redis.singleton:全局实例初始化测试连接失败 : " + dErr.Error())
	}
	rdb = d
	return
}

// GetClient 获取 redis 底层实例，可根据官方文档操作redis
func (rec *singleton) GetClient() (client *goRedis.Client, err error) {
	if rdb == nil {
		err = errors.New("redis.singleton:全局实例未初始化，请先进行初始化")
		return
	}
	return rdb, nil
}

// GetOperator 获取 redis 封装操作实例
func (rec *singleton) GetOperator() (client *Operator, err error) {
	d, err := rec.GetClient()
	if err != nil {
		return
	}
	return NewOperator(d), nil
}
