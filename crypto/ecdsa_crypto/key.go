package ecdsa_crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

var (
	errKeyMustBePEMEncoded = errors.New("invalid key: Key must be a PEM encoded PKCS1 or PKCS8 key")
	errNotECPublicKey      = errors.New("key is not a valid ECDSA public key")
	errNotECPrivateKey     = errors.New("key is not a valid ECDSA private key")
)

func GenKey() (pubStr, priStr string, err error) {
	generateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return
	}
	derText, err := x509.MarshalECPrivateKey(generateKey)
	if err != nil {
		return
	}
	pubDerText, err := x509.MarshalPKIXPublicKey(&generateKey.PublicKey)
	if err != nil {
		return
	}
	priBlock := pem.Block{
		Type:  "ECDSA Private Key",
		Bytes: derText,
	}
	pri := pem.EncodeToMemory(&priBlock)
	pubBlock := pem.Block{
		Type:  "ECDSA Public Key",
		Bytes: pubDerText,
	}
	pub := pem.EncodeToMemory(&pubBlock)
	if pub == nil || pri == nil {
		err = errors.New("gen key failed")
		return
	}
	pubStr = string(pub)
	priStr = string(pri)
	return
}

func ParsePrivateKeyFromPEM(pemStr string) (*ecdsa.PrivateKey, error) {
	key := []byte(pemStr)
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, errKeyMustBePEMEncoded
	}

	// Parse the key
	var parsedKey interface{}
	if parsedKey, err = x509.ParseECPrivateKey(block.Bytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
			return nil, err
		}
	}

	var pkey *ecdsa.PrivateKey
	var ok bool
	if pkey, ok = parsedKey.(*ecdsa.PrivateKey); !ok {
		return nil, errNotECPrivateKey
	}

	return pkey, nil
}

// ParsePublicKeyFromPEM parses a PEM encoded PKCS1 or PKCS8 public key
func ParsePublicKeyFromPEM(pemStr string) (*ecdsa.PublicKey, error) {
	key := []byte(pemStr)
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, errKeyMustBePEMEncoded
	}

	// Parse the key
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
			parsedKey = cert.PublicKey
		} else {
			return nil, err
		}
	}

	var pkey *ecdsa.PublicKey
	var ok bool
	if pkey, ok = parsedKey.(*ecdsa.PublicKey); !ok {
		return nil, errNotECPublicKey
	}

	return pkey, nil
}
