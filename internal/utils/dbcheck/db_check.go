package dbcheck

import (
	"errors"
	grds "github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// DbError 是否是要直接返回的db错误 (gorm.ErrRecordNotFound is not a error)
func DbError(err error) error {
	//没找到不算一个错误，判断返回结构体未nil来做处理
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return err
}

// RdError 是否是要直接返回的redis错误 (grds.Nil is not a error)
func RdError(err error) error {
	//没找到不算一个错误，判断返回结构体未nil来做处理
	if errors.Is(err, grds.Nil) {
		return nil
	}

	return err
}
