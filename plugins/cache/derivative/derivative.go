package derivative

import (
	"context"
	grds "github.com/go-redis/redis/v8"
	"time"
)

// Derivative redis衍生组件的基础结构体
type Derivative struct {
	Ctx context.Context //上下文
	Cli *grds.Client    //客户端

	Key string //redis的key
}

// Init 初始化衍生组件
func (d *Derivative) Init(ctx context.Context, cli *grds.Client, key string) {
	d.Ctx, d.Key, d.Cli = ctx, key, cli
}

// ExtendExp 延长缓存生命
func (d *Derivative) ExtendExp(exp time.Duration) {
	d.Cli.Expire(d.Ctx, d.Key, exp)
}
