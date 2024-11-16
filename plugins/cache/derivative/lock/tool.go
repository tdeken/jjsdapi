package lock

import (
	"context"
	"errors"
	grds "github.com/go-redis/redis/v8"
	"math/rand"
	"time"
)

type ToolData struct {
	Data   []byte
	Err    error
	Expire time.Duration
}

// GetSet 带锁的getset
func GetSet(ctx context.Context, cli *grds.Client, key string, find func() ToolData, opts ...Option) (data []byte, err error) {
	data, err = cli.Get(ctx, key).Bytes()
	if err == nil || !errors.Is(err, grds.Nil) {
		return
	}

	lk := forLock(key)
	lc := NewSubLock(ctx, cli, lk, opts...).Lock()
	defer lc.Release()
	if lc.Error() != nil {
		return
	}

	data, err = cli.Get(ctx, key).Bytes()
	if err == nil || !errors.Is(err, grds.Nil) {
		return
	}

	td := find()
	data, err = td.Data, td.Err
	if err != nil {
		return
	}

	expire := td.Expire
	if expire == 0 {
		expire = time.Duration(rand.Intn(60)+60) * time.Minute
	}

	err = cli.Set(ctx, key, data, expire).Err()

	return
}
