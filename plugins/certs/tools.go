package certs

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"math/rand"
)

const (
	allowedChars     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	allowedCharsSize = len(allowedChars)
)

// returns a securely generated random string.
func randomString(length int) string {
	b := make([]byte, length)

	for i := range b {
		c := rand.Intn(allowedCharsSize)
		b[i] = allowedChars[c]
	}

	return string(b)
}

func mmd5(str string) string {
	var buf = bytes.NewBufferString(str)
	hash := md5.New()
	hash.Write(buf.Bytes())
	sign := hash.Sum(nil)

	return hex.EncodeToString(sign)
}

func checkHeaderSign(header map[string]interface{}) (slt, sin, v string, ok bool) {

	salt, ok := header[saltKey]
	if !ok {
		return
	}

	sign, ok := header[signKey]
	if !ok {
		return
	}

	ver, ok := header[verKey]
	if !ok {
		return
	}

	slt, sin, v = salt.(string), sign.(string), ver.(string)
	return
}
