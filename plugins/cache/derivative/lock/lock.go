package lock

import (
	"context"
	"errors"
	grds "github.com/go-redis/redis/v8"
	"github.com/rs/xid"
	"jjsdapi/plugins/cache/derivative"
	"time"
)

const (
	DefaultRetryTime = 20 * time.Millisecond //20ms重新重试一次
	DefaultRetryNum  = 100                   //重试次数
)

var (
	FailTimeoutErr = errors.New("lock timeout")
	CtxCancelErr   = errors.New("context cancel")
)

type ILock struct {
	derivative.Derivative
	value      string        // 锁的值
	exp        time.Duration //锁过期时间
	lockStatus bool          // 锁状态
	err        error         // 错误
}

// SoftLock 软锁
func (l *ILock) SoftLock() bool {
	l.lockStatus = l.Cli.SetNX(l.Ctx, l.Key, l.value, l.exp).Val()
	return l.lockStatus
}

func (l *ILock) init(ctx context.Context, cli *grds.Client, key string) {
	l.Init(ctx, cli, key)
	l.value = xid.New().String()
	l.exp = DefaultRetryNum * DefaultRetryTime * 2

	return
}

// Release 释放锁
func (l *ILock) release() {
	//释放redis锁
	l.lockStatus = false
	//拿缓存里面的值
	val, err := l.Cli.Get(l.Ctx, l.Key).Result()
	if err != nil && !errors.Is(err, grds.Nil) {
		l.err = err
		return
	}

	if val == l.value {
		l.err = l.Cli.Del(l.Ctx, l.Key).Err()
		return
	}

	return
}

// LockDone 上锁成功就执行方法，没上锁就直接过
func (l *ILock) LockDone(fc func()) {
	if l.SoftLock() {
		defer l.release()
		fc()
	}
}

// 锁错误
func (l *ILock) Error() error {
	return l.err
}

func forLock(key string) string {
	return key + ":__for_lock__"
}
