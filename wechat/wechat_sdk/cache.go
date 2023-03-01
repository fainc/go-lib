package wechat_sdk

import (
	"context"
	"errors"
	"time"

	"github.com/gogf/gf/v2/os/gcache"
	"github.com/redis/go-redis/v9"
)

type cache struct{}

var cacheVar = cache{}

func Cache() *cache {
	return &cacheVar
}

var cacheEngine = "memory"
var redisCacheClient *redis.Client
var memoryCacheClient *gcache.Cache

// InitRedisCacheClient 全局初始化Redis缓存（仅需启动时初始化一次,后续初始化返回首次配置），未初始化redis或初始化失败默认使用memory配置
// redisConf : redis://<user>:<pass>@localhost:6379/<db>
func (rec *cache) InitRedisCacheClient(redisConf string) (*redis.Client, error) {
	if redisCacheClient != nil {
		return redisCacheClient, nil
	}
	opt, err := redis.ParseURL(redisConf)
	if err != nil {
		return nil, err
	}
	rcc := redis.NewClient(opt)
	err = rcc.Get(context.Background(), "test_init_connect").Err()
	if err != nil && err != redis.Nil {
		err = errors.New("redis init failed：" + err.Error())
		return nil, err
	}
	redisCacheClient = rcc
	cacheEngine = "redis"
	return redisCacheClient, nil
}

// GetRedisCacheClient 获取redis缓存client
func (rec *cache) getRedisCacheClient() (*redis.Client, error) {
	if redisCacheClient == nil {
		return nil, errors.New("redis未初始化")
	}
	return redisCacheClient, nil
}

// GetMemoryCacheClient 获取内存缓存Client（无需手动初始化）
func (rec *cache) getMemoryCacheClient() (*gcache.Cache, error) {
	if memoryCacheClient != nil {
		return memoryCacheClient, nil
	}
	memoryCacheClient = gcache.New()
	cacheEngine = "memory"
	return memoryCacheClient, nil
}

func (rec *cache) GetEngine() string {
	return cacheEngine
}

func (rec *cache) GetRedisCache(prefix string, id string) (value string, err error) {
	rdb, err := rec.getRedisCacheClient()
	if err != nil {
		return
	}
	res := rdb.Get(context.Background(), prefix+"_"+id)
	if res.Err() != nil && res.Err() != redis.Nil {
		err = errors.New("redis connect error：" + res.Err().Error())
		return
	}
	if res.Err() == redis.Nil {
		value = ""
		return
	}
	value, err = res.Result()
	return
}

func (rec *cache) SetRedisCache(prefix string, id string, value string, timeout int) (err error) {
	rdb, err := rec.getRedisCacheClient()
	if err != nil {
		return
	}
	cmd := rdb.SetEx(context.Background(), prefix+"_"+id, value, time.Duration(timeout)*time.Second)
	_, err = cmd.Result()
	if err != nil {
		return err
	}
	return
}

func (rec *cache) GetMemoryCache(prefix string, id string) (value string, err error) {
	c, err := rec.getMemoryCacheClient()
	if err != nil {
		return
	}
	value = c.MustGet(context.Background(), prefix+"_"+id).String()
	return
}

func (rec *cache) SetMemoryCache(prefix string, id string, value string, timeout int) (err error) {
	c, err := rec.getMemoryCacheClient()
	if err != nil {
		return
	}
	err = c.Set(context.Background(), prefix+"_"+id, value, time.Duration(timeout)*time.Second)
	if err != nil {
		return err
	}
	return
}
