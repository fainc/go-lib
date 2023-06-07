package redis

import (
	"errors"

	goRedis "github.com/redis/go-redis/v9"
)

// NilError 判断结果是否为空
func NilError(err error) bool {
	return errors.Is(err, goRedis.Nil)
}
