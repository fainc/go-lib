package wechat_sdk

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/os/gcache"

	"github.com/fainc/go-lib/redis"
)

type cache struct{}

var cacheVar = cache{}

func Cache() *cache {
	return &cacheVar
}

var rdb *redis.Operator
var mdb *gcache.Cache

func init() {
	op, err := redis.Singleton().GetOperator()
	if err == nil && op != nil {
		rdb = op
		return
	}
	mdb = gcache.New()
	fmt.Println("wechat_sdk cache init 警告：因获取redis实例失败，将使用内存读写缓存（应用重启丢失数据），redis错误：" + err.Error())
}
func (rec *cache) GetCache(prefix string, id string) (value string, err error) {
	if rdb == nil {
		value = mdb.MustGet(context.Background(), prefix+"_"+id).String()
		return
	}
	value, err = rdb.GetStringValue(prefix + "_" + id)
	return
}

func (rec *cache) SetCache(prefix string, id string, value string, timeout int) (err error) {
	if rdb == nil {
		err = mdb.Set(context.Background(), prefix+"_"+id, value, time.Duration(timeout)*time.Second)
		return
	}
	err = rdb.SetValue(prefix+"_"+id, value, time.Duration(timeout)*time.Second)
	return
}
