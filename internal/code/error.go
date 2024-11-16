package code

import (
	"encoding/json"
	"errors"
)

type Error struct {
	Form   string `json:"form"`   //来源
	Code   int32  `json:"code"`   //错误码
	Detail string `json:"detail"` //错误信息
}

var Ok = &Error{
	Code:   OK,
	Detail: OKMsg,
}

// GetCode 状态码
func (e *Error) GetCode() int32 {
	if e == nil {
		return OK
	}
	return e.Code
}

// GetDetail 状态码说明
func (e *Error) GetDetail() string {
	if e == nil {
		return OKMsg
	}
	return e.Detail
}

func (e *Error) Error() string {
	if e == nil {
		return "nil"
	}
	b, _ := json.Marshal(e)
	return string(b)
}

// NewError 实例化一个错误
func NewError(code int32, detail string) *Error {
	return &Error{Code: code, Detail: detail}
}

// As 切换成错误
func As(err error) (e *Error, ok bool) {
	ok = errors.As(err, &e)
	if !ok {
		return
	}
	return err.(*Error), true
}
