package jwt

import (
	"fmt"
	"testing"
	"time"

	"github.com/fainc/go-lib/crypto/ecdsa_crypto"
)

func TestGenKey(_ *testing.T) {
	var pub, pri string
	var err error
	if pub, pri, err = ecdsa_crypto.GenKey(); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(pub)
	fmt.Println(pri)
}
func TestIssue(_ *testing.T) {
	pub, pri, _ := ecdsa_crypto.GenKey()
	secret := "12345678123456781234567812345678"
	private, err := ecdsa_crypto.ParsePrivateKeyFromPEM(pri)
	if err != nil {
		return
	}
	public, err := ecdsa_crypto.ParsePublicKeyFromPEM(pub)
	if err != nil {
		return
	}
	token, _, err := Issuer(IssuerConf{
		JwtSecret:    secret,
		JwtAlgo:      AlgoES256,
		CryptoAlgo:   AlgoAES,
		JwtPrivate:   private,
		CryptoSecret: secret,
	}).Publish(&IssueParams{
		UserID:    "1",
		Subject:   "Auth",
		Audience:  []string{"a.com", "b.com"},
		Duration:  1 * time.Second,
		NotBefore: time.Now().Add(2 * time.Second),
		Ext:       "hello",
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Token=", token)
	c, err := Parser(ParserConf{JwtSecret: secret,
		JwtAlgo:      AlgoES256,
		JwtPublic:    public,
		CryptoSecret: secret}).Validate(ValidateParams{
		Token:     "Bearer " + token,
		Audience:  "b.com",
		Leeway:    10 * time.Second,
		LeewayNbf: true,
		LifeCycle: 10 * time.Second,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("UserID=", c.UserID)
	fmt.Println("Ext=", c.Ext)
}
