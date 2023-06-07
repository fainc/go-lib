package jwt

import (
	"context"
	"time"

	goRedis "github.com/redis/go-redis/v9"

	"github.com/fainc/go-lib/redis"
)

type redisProtect struct {
	rdb *goRedis.Client // 通过 goRedis Client 操作 redis
}

// RedisProtect redis守护 cusRdb: 自定义rdb实例，不传则尝试通过单例 redis.Singleton().GetClient() 获取，单例也未初始化会panic
func RedisProtect(cusRdb ...*goRedis.Client) (p *redisProtect) {
	var rdb *goRedis.Client
	if cusRdb != nil && cusRdb[0] != nil {
		rdb = cusRdb[0]
	} else {
		var err error
		if rdb, err = redis.Singleton().GetClient(); err != nil {
			panic(err)
		}
	}
	return &redisProtect{rdb}
}

// Publish 发布jwt（白名单）
func (rec *redisProtect) Publish(jti, userID string, exp time.Duration) (err error) {
	_, err = rec.rdb.SetEx(context.Background(), "jwt_pub_"+jti, userID, exp).Result()
	if err != nil {
		return err
	}
	return
}

// IsPublished 判断是否已发布jwt
func (rec *redisProtect) IsPublished(jti string) (is bool, err error) {
	n, err := rec.rdb.Exists(context.Background(), "jwt_pub_"+jti).Result()
	if err != nil {
		return
	}
	return n > 0, nil
}

// Revoke 吊销jwt(黑名单)
func (rec *redisProtect) Revoke(jti string, exp time.Duration) (err error) {
	detail := time.Now().Format(time.DateTime) + "|" + time.Now().Add(exp).Format(time.DateTime)
	_, err = rec.rdb.SetEx(context.Background(), "jwt_block_"+jti, detail, exp).Result()
	if err != nil {
		return err
	}
	return
}

// IsRevoked 判断是否吊销的jwt
func (rec *redisProtect) IsRevoked(jti string) (is bool, err error) {
	n, err := rec.rdb.Exists(context.Background(), "jwt_block_"+jti).Result()
	if err != nil {
		return
	}
	return n > 0, nil
}
