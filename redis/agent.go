package redis

/**
常用redis操作代理方法，不在常用代理方法内的操作，请自行GetClient()操作
*/
import (
	"context"
	"time"
)

// GetWxSat 获取微信Server Access Token
func GetWxSat(appID string) (sat string, err error) {
	return GetStringValue("glib_wx_sat_" + appID)
}

// SetWxSat 写入微信Server Access Token
func SetWxSat(appID, sat string, second int64) (err error) {
	return SetValue("glib_wx_sat_"+appID, sat, time.Second*time.Duration(second))
}

// GetWxUat 获取微信User Access Token
func GetWxUat(openID string) (uat string, err error) {
	return GetStringValue("glib_wx_uat_" + openID)
}

// SetWxUat 写入微信User Access Token
func SetWxUat(openID, uat string, second int64) (err error) {
	return SetValue("glib_wx_uat_"+openID, uat, time.Second*time.Duration(second))
}

// GetWxJat 获取微信 JS API TICKET
func GetWxJat(appID string) (jat string, err error) {
	return GetStringValue("glib_wx_jat_" + appID)
}

// SetWxJat 写入微信 JS API TICKET
func SetWxJat(appID, jat string, second int64) (err error) {
	return SetValue("glib_wx_jat_"+appID, jat, time.Second*time.Duration(second))
}

// GetJwt 获取JWT信息
func GetJwt(jti string) (userID int64, err error) {
	return GetIntValue("glib_jwt_" + jti)
}

// SetJwt 写入JWT信息
func SetJwt(jti string, userID, second int64) (err error) {
	return SetValue("glib_jwt_"+jti, userID, time.Second*time.Duration(second))
}

// DelJwt 删除JWT
func DelJwt(jti string) (err error) {
	_, err = DelKey("glib_jwt_" + jti)
	return
}

// RenewJwt JWT 有效期延期
func RenewJwt(jti string, expire time.Duration) (s bool, err error) {
	key := "glib_jwt_" + jti
	client, err := GetClient()
	if err != nil {
		return
	}
	e := client.TTL(context.Background(), key).Val()
	if e.Seconds() <= 0 {
		return
	}
	ret := client.Expire(context.Background(), key, e+expire)
	if ret.Err() != nil {
		err = ret.Err()
		return
	}
	s = ret.Val()
	return
}
