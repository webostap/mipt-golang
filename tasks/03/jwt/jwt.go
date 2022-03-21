package jwt

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"hash"
	"time"
)

type SignMethod string

const (
	HS256 SignMethod = "HS256"
	HS512 SignMethod = "HS512"
)

var (
	ErrInvalidSignMethod      = errors.New("invalid sign method")
	ErrSignatureInvalid       = errors.New("signature invalid")
	ErrTokenExpired           = errors.New("token expired")
	ErrSignMethodMismatched   = errors.New("sign method mismatched")
	ErrConfigurationMalformed = errors.New("configuration malformed")
	ErrInvalidToken           = errors.New("invalid token")
)

func Encode(data interface{}, opts ...Option) ([]byte, error) {
	conf := readOptions(opts)
	now := timeFunc()
	if conf.Expires != nil && (conf.TTL != nil || conf.Expires.Before(now)) {
		return nil, ErrConfigurationMalformed
	}

	headerI := map[string]interface{}{
		"alg": "HS256",
		"typ": "JWT",
	}
	hf, err := getSignMethod(conf.SignMethod)
	if err != nil {
		return nil, ErrInvalidSignMethod
	}
	headerI["alg"] = conf.SignMethod

	payloadI := map[string]interface{}{
		"d":   data,
		"exp": float64(0),
	}
	if conf.Expires != nil {
		payloadI["exp"] = conf.Expires.Unix()
	} else if conf.TTL != nil {
		payloadI["exp"] = now.Add(*conf.TTL).Unix()
	} else {
		delete(payloadI, "exp")
	}

	header, _ := json.Marshal(headerI)
	payload, _ := json.Marshal(payloadI)

	var token bytes.Buffer
	token.WriteString(base64.RawURLEncoding.EncodeToString(header))
	token.WriteString(".")
	token.WriteString(base64.RawURLEncoding.EncodeToString(payload))

	sig := getHmac(hf, token.Bytes(), conf.Key)

	token.WriteString(".")
	token.Write(sig)

	return token.Bytes(), nil
}

func Decode(token []byte, data interface{}, opts ...Option) error {
	parts := bytes.Split(token, []byte("."))

	if len(parts) != 3 {
		return ErrInvalidToken
	}
	header, err := base64toMap(parts[0])
	if err != nil {
		return ErrInvalidToken
	}

	payload, err := base64toMap(parts[1])
	if err != nil {
		return ErrInvalidToken
	}

	dti := data.(*map[string]interface{})
	*dti = payload["d"].(map[string]interface{})

	if header["typ"] != "JWT" {
		return ErrInvalidToken
	}

	conf := readOptions(opts)

	if header["alg"] != string(conf.SignMethod) {
		return ErrSignMethodMismatched
	}
	hf, err := getSignMethod(conf.SignMethod)
	if err != nil {
		return ErrInvalidSignMethod
	}

	var hap bytes.Buffer
	hap.Write(parts[0])
	hap.WriteString(".")
	hap.Write(parts[1])

	sig := getHmac(hf, hap.Bytes(), conf.Key)

	if bytes.Compare(sig, parts[2]) != 0 {
		return ErrSignatureInvalid
	}

	if payload["exp"] != nil {
		exp := time.Time{}
		if timeFunc().After(exp) {
			return ErrTokenExpired
		}
	}
	return nil
}
func readOptions(opts []Option) config {
	var conf config
	for _, opt := range opts {
		opt(&conf)
	}
	return conf
}

func base64toMap(tok []byte) (map[string]interface{}, error) {
	dst := make([]byte, base64.RawURLEncoding.DecodedLen(len(tok)))
	_, err := base64.RawURLEncoding.Decode(dst, tok)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(dst, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func getSignMethod(SignMethod SignMethod) (func() hash.Hash, error) {
	switch SignMethod {
	case HS256:
		return sha256.New, nil
	case HS512:
		return sha512.New, nil
	default:
		return nil, ErrInvalidSignMethod
	}
}
func getHmac(hf func() hash.Hash, data []byte, secret []byte) []byte {
	h := hmac.New(hf, secret)
	h.Write(data)

	hs := h.Sum(nil)
	sig := make([]byte, base64.RawURLEncoding.EncodedLen(len(hs)))
	base64.RawURLEncoding.Encode(sig, hs)

	return sig
}

// To mock time in tests
var timeFunc = time.Now
