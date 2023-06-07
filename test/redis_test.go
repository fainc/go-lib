package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/fainc/go-lib/redis"
)

// redis 初始化
func TestSingleton_Init(_ *testing.T) {
	err := redis.Singleton().Init(redis.NewRedisConf{
		Address:  "",
		User:     "",
		Password: "",
		DB:       0,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("test init redis success")
}

// 通过单例直接调用操作封装
func TestSingleton_GetOperator(t *testing.T) {
	op, err := redis.Singleton().GetOperator()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	val, err := op.GetIntValue("abc")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(val)
}

// 独立调用操作封装
func TestSingleton_Operator2(t *testing.T) {
	cl, err := redis.Singleton().GetClient()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	val, err := redis.Operator(cl).GetIntValue("key")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(val)
}

// 获取单例底层redis实例
func TestSingleton_GetClient(t *testing.T) {
	cl, err := redis.Singleton().GetClient()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	val, err := cl.Get(context.Background(), "key").Result()
	if err != nil {
		fmt.Println("is result nil:", redis.NilError(err))
		fmt.Println(err.Error())
		return
	}
	fmt.Println(val)
}
