package ckey

import "fmt"

func userKey(format string, a ...any) string {
	return UserPrefix + fmt.Sprintf(format, a...)
}

// UserLogin 登陆
func UserLogin(userId int64, jwtId string) string {
	return userKey("login:%d_%s", userId, jwtId)
}
