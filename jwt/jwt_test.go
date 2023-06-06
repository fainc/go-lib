package jwt

import (
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/fainc/go-lib/crypto/ecdsa_crypto"
)

var (
	pub       = "-----BEGIN ECDSA Public Key-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEYXjlYuQoLpgmpQs4BS9NNaSA35dB\nZAnXseazK4lTo759VtfodjB8ji8DJZt/ixNZP0eHNOo0h9EGyFJc4ycxIQ==\n-----END ECDSA Public Key-----"
	pri       = "-----BEGIN ECDSA Private Key-----\nMHcCAQEEIKNcl4q6P2tJRhes3C6cKAsDc3fn8gg3SNeJREHjGjMRoAoGCCqGSM49\nAwEHoUQDQgAEYXjlYuQoLpgmpQs4BS9NNaSA35dBZAnXseazK4lTo759VtfodjB8\nji8DJZt/ixNZP0eHNOo0h9EGyFJc4ycxIQ==\n-----END ECDSA Private Key-----"
	secret256 = "12345678123456781234567812345678"
)

func TestIssue(_ *testing.T) {
	jwtAlgo := AlgoHS256 // AlgoES256
	private, err := ecdsa_crypto.ParsePrivateKeyFromPEM(pri)
	if err != nil {
		return
	}
	token, _, err := Issuer(IssuerConf{
		JwtSecret: secret256,
		JwtAlgo:   jwtAlgo,
		// CryptoAlgo: AlgoAES,
		JwtPrivate: private,
		// CryptoSecret: secret256,
	}).Publish(&IssueParams{
		UserID:   "10000",
		Subject:  "Auth",
		Audience: []string{"a.com", "b.com"},
		Duration: 1000 * time.Second,
		Ext:      "{'name':'xxx','email':'xxxx'}",
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(token)
	public, err := ecdsa_crypto.ParsePublicKeyFromPEM(pub)
	if err != nil {
		return
	}
	validate, err := Parser(ParserConf{
		JwtAlgo:   jwtAlgo,
		JwtPublic: public,
		JwtSecret: secret256,
		// CryptoSecret: secret256,
	}).Validate(ValidateParams{Token: "Bearer " + token, Audience: "a.com"})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(validate)
}

func TestName(t *testing.T) {
	validate, err := Parser(ParserConf{
		JwtAlgo:   AlgoHS256,
		JwtSecret: "fjycsp2023",
	}).ParseRaw("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1dWlkIjoxLCJzY29wZSI6Im1hbmFnZXIiLCJpc3MiOiJqd3RIZWxwZXIiLCJleHAiOjE2ODYxMTMzMjcsIm5iZiI6MTY4NjAyNjkyNywiaWF0IjoxNjg2MDI2OTI3fQ.vhMSBTgRqmB__cwjg6T9eX2FCPRUx4TKPPPwm2WdJ-k")
	if err != nil {
		g.Dump(err.Error())
		return
	}
	g.Dump(gconv.Time(validate["exp"]).String())
}
