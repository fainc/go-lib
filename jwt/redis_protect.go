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

// RedisProtect redis守护 cusRdb: 自定义rdb实例，不传则尝试使用 redis.Singleton().GetClient() 获取
func RedisProtect(cusRdb ...*goRedis.Client) (p *redisProtect, err error) {
	var rdb *goRedis.Client
	if cusRdb != nil && cusRdb[0] != nil {
		rdb = cusRdb[0]
	} else {
		rdb, err = redis.Singleton().GetClient()
		if err != nil {
			return
		}
	}
	return &redisProtect{rdb}, nil
}

// Revoke 吊销jwt
func (rec *redisProtect) Revoke(jti, userID string, exp time.Duration) (err error) {
	_, err = rec.rdb.SetEx(context.Background(), "jwt_block_"+jti, userID, exp).Result()
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
