package certs

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/xid"
	"time"
)

const (
	userV1 = "v1"
)

var userSecret = map[string]string{
	userV1: "SXtoJkrMT7UtDR2rYdID",
}

type User struct {
	UserId int64 `json:"user_id"`
}

type UserClaims struct {
	*User
	jwt.RegisteredClaims
}

func NewUser() *UserClaims {
	return &UserClaims{}
}

// Token 获取token
func (u *UserClaims) Token(user *User) (string, error) {
	u.User = user
	u.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:    issuer,
		Subject:   fmt.Sprint(user.UserId),
		ExpiresAt: jwt.NewNumericDate(u.Expired()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        xid.New().String(),
	}

	var key = randomString(6)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, u)

	//加密校验字符串
	token.Header[saltKey] = key
	token.Header[signKey] = mmd5(key + userSecret[userV1])
	token.Header[verKey] = userV1

	return token.SignedString([]byte(key))
}

// Parse 解析token
func (u *UserClaims) Parse(tokenString string) (err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		salt, sign, ver, ok := checkHeaderSign(token.Header)
		if !ok || mmd5(salt+userSecret[ver]) != sign {
			return nil, errors.New("key invalid")
		}

		return []byte(salt), nil
	})

	if err == nil && token.Valid {
		user, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return InvalidToken
		}

		exp, _ := user.GetExpirationTime()
		if exp.Before(time.Now()) {
			return ExpiredErr
		}

		b, _ := json.Marshal(user)
		_ = json.Unmarshal(b, u)
		return
	}

	return InvalidToken
}

func (u *UserClaims) Expired() time.Time {
	return time.Now().AddDate(0, 0, 7)
}
