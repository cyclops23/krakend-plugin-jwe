package main

import (
	"errors"

	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
)

func decryptAndDecode(jweToken *string, decryptionKey string, sharedSecret string) (map[string]interface{}, error) {
	token, err := decrypt(jweToken, decryptionKey)
	if err != nil {
		return nil, err
	}
	return decode(token, sharedSecret)
}

func decode(token *jwt.JSONWebToken, sharedSecret string) (map[string]interface{}, error) {
	padded := make([]byte, 32)
	copy(padded, sharedSecret)
	out := map[string]interface{}{}
	if err := token.Claims(padded, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func decrypt(jweToken *string, decryptionKey string) (*jwt.JSONWebToken, error) {
	if jweToken == nil {
		return nil, errors.New("missing token")
	}

	jwe, err := jose.ParseEncrypted(*jweToken)
	if err != nil {
		return nil, err
	}

	decryptedKey, err := jwe.Decrypt(decryptionKey)
	if err != nil {
		return nil, err
	}

	jwtParsed, err := jwt.ParseSigned(string(decryptedKey))
	if err != nil {
		return nil, err
	}

	return jwtParsed, nil
}
