package certs

import "errors"

const issuer = "certified"

const saltKey = "salt"
const signKey = "sign"
const verKey = "ver"

var ExpiredErr = errors.New("token expired")
var InvalidToken = errors.New("invalid token")
