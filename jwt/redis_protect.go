package jwt

import (
	"context"
	"time"

	"github.com/fainc/go-lib/redis"
)

type redisProtect struct{}

var redisProtectVar = redisProtect{}

func RedisProtect() *redisProtect {
	return &redisProtectVar
}

// Publish 发布redis jwt白名单
func (rec *redisProtect) Publish(jti, userID string, exp time.Duration) (err error) {
	rdb, err := redis.GetClient()
	if err != nil {
		return
	}
	rdb.SetEx(context.Background(), "jwt_"+jti, userID, exp)
	return
}
