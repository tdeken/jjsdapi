package vercontrol

import (
	"context"
	"fmt"
	grds "github.com/go-redis/redis/v8"
	"jjsdapi/plugins/cache/derivative"
	rl "jjsdapi/plugins/cache/derivative/lock"
	"time"
)

type Done func() error

const (
	defaultExp    = 2 * time.Hour  //默认的版本过期时间
	lockKeyPreFix = "VER_CONTROL:" //锁的key前缀
)

type VerControl struct {
	derivative.Derivative
	exp         time.Duration //版本过期时间
	comparedMsg string        //对比后的内容
	hasNewVer   bool          //是否有新的版本
}

// NewRedisVerControl 实例化一个Redis版本控制
func NewRedisVerControl(ctx context.Context, cli *grds.Client, key string, opts ...Option) *VerControl {
	vc := &VerControl{
		exp: defaultExp,
	}

	vc.Init(ctx, cli, key)

	//设置配置
	for _, opt := range opts {
		opt(vc)
	}

	return vc
}

// Cancel 取消版本，但是有新版本的时候，不会取消
func (c *VerControl) Cancel() error {
	if c.hasNewVer { //有新版本的时候，不能取消版本
		return nil
	}
	return c.Cli.Del(c.Ctx, c.Key).Err()
}

// HasNewVersion 是否有新版本
func (c *VerControl) HasNewVersion(t time.Time) bool {
	last := c.GetLastUpdateTimeVal()
	c.comparedMsg = fmt.Sprintf("当前版本: %d, 最新版本:%d", t.UnixNano(), last.UnixNano())

	c.hasNewVer = last.After(t)

	return c.hasNewVer
}

// GetLastUpdateTime 获取当前最新的版本
func (c *VerControl) GetLastUpdateTime() (time.Time, error) {
	return c.Cli.Get(c.Ctx, c.Key).Time()
}

// GetComparedMsg 获取比较后的信息
func (c *VerControl) GetComparedMsg() string {
	return c.comparedMsg
}

// GetLastUpdateTimeVal 获取当前最新的版本的值
func (c *VerControl) GetLastUpdateTimeVal() time.Time {
	t, _ := c.GetLastUpdateTime()

	return t
}

// Create 创建一个版本
func (c *VerControl) Create(updateTime time.Time, fc Done) error {
	lock := rl.NewRedisLock(c.Ctx, c.Cli, fmt.Sprintf("%s%s", lockKeyPreFix, c.Key)).Lock()
	defer lock.Release()
	if lock.Error() != nil {
		return lock.Error()
	}

	//获取当前的版本修改时间
	oldTime := c.GetLastUpdateTimeVal()

	//旧的版本修改时间小于新的版本时间
	if oldTime.Before(updateTime) {
		c.hasNewVer = true
		c.Cli.Set(c.Ctx, c.Key, updateTime, c.exp)
		return fc()
	}

	return nil
}
