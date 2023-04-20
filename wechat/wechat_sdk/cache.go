package wechat_sdk

import (
	"context"
	"errors"
	"fmt"
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
		return nil, errors.New("wechat_sdk: redis缓存已初始化，请勿重复初始化")
	}
	if memoryCacheClient != nil {
		return nil, errors.New("wechat_sdk: redis缓存初始化失败，已初始化memory缓存，与redis缓存冲突")
	}
	opt, err := redis.ParseURL(redisConf)
	if err != nil {
		return nil, err
	}
	rcc := redis.NewClient(opt)
	err = rcc.Get(context.Background(), "test_init_connect").Err()
	if err != nil && err != redis.Nil {
		err = errors.New("wechat_sdk: redis init failed : " + err.Error())
		return nil, err
	}
	redisCacheClient = rcc
	cacheEngine = "redis"
	fmt.Println("wechat_sdk 提示：redis缓存已初始化成功，缓存数据通过redis存储")
	return redisCacheClient, nil
}

// GetRedisCacheClient 获取redis缓存client
func (rec *cache) getRedisCacheClient() (*redis.Client, error) {
	if redisCacheClient == nil {
		return nil, errors.New("wechat_sdk:redis未初始化")
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
	fmt.Println("wechat_sdk 提示：您正在使用MemoryCache缓存(accessToken等)，应用重启缓存会丢失，如需持久化缓存请使用 redis或远程缓存中心 ，确定使用内存缓存可忽略本提示")
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
