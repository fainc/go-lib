package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/fainc/go-lib/gf2/token"
	"github.com/fainc/go-lib/jwt"
)

func TestJwtIssue(_ *testing.T) {
	token, id, err := token.Helper().Publish(token.PublishParams{
		UserID:   100,
		Duration: 7 * time.Hour,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(token)
	fmt.Println(id)
}

func TestJwtRevoke(_ *testing.T) {
	err := jwt.RedisProtect().Revoke("360738DB-2094-4041-BA29-29469D50A717", time.Second*10)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	revoked, err := jwt.RedisProtect().IsRevoked("D925BBD9-53B9-4A65-A2E1-3578A35ACE62")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(revoked)
}

func TestJwtRevoked(_ *testing.T) {
	revoked, err := jwt.RedisProtect().IsRevoked("D925BBD9-53B9-4A65-A2E1-3578A35ACE62")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(revoked)
}
