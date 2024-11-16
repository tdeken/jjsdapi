package lock

import (
	"context"
	grds "github.com/go-redis/redis/v8"
	"time"
)

type RedisLock struct {
	ILock
	retryNum  int
	retryTime time.Duration
}

// NewRedisLock 实例化一个redis锁, 有特定参数，用 SetConfig 方法实现
func NewRedisLock(ctx context.Context, cli *grds.Client, key string, opts ...Option) *RedisLock {
	l := &RedisLock{
		retryNum:  DefaultRetryNum,
		retryTime: DefaultRetryTime,
	}

	// 初始化锁
	l.init(ctx, cli, key)

	//设置配置
	for _, opt := range opts {
		opt(&l.ILock)
	}

	return l
}

// Lock redis加锁
func (l *RedisLock) Lock() *RedisLock {
	retry := 1
	for {
		if l.SoftLock() {
			break
		}

		if retry > l.retryNum {
			l.err = FailTimeoutErr
			break
		}
		retry++
		time.Sleep(l.retryTime)
	}
	return l
}

// Release 释放锁
func (l *RedisLock) Release() {
	if !l.lockStatus {
		return
	}

	defer l.release()

	return
}
