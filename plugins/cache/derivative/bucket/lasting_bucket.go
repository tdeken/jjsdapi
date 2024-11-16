package bucket

import (
	"context"
	"fmt"
	grds "github.com/go-redis/redis/v8"
	"jjsdapi/plugins/cache/derivative"
	rl "jjsdapi/plugins/cache/derivative/lock"
	"math"
	"time"
)

const (
	defaultBucketLen = 100            //默认桶长度
	defaultLife      = 24 * time.Hour //默认静止生命时长
)

// CheckValue 验证这个值是否插入桶
type CheckValue func(ctx context.Context, val string) bool

// Value 产生的值
type Value func(ctx context.Context) string

// LastingBucket 持久桶
// 一直放在缓存里面，不够就补到够，拿的时候就直接从缓存拿
type LastingBucket struct {
	derivative.Derivative
	val       Value        //生成值得函数
	check     []CheckValue //验证值得函数，只有验证通过才会存到桶中，没设置的时候就验证
	bl        int64        //桶的总长度
	set       bool         //是否已经设置过
	threshold int64        //桶的长度阈值，当剩余长度小于或等于这个值得时候才填充桶
	life      time.Duration
}

// NewLastingBucket 实例一个持久桶
func NewLastingBucket(ctx context.Context, cli *grds.Client, key string, val Value, opts ...LastingOption) *LastingBucket {
	lb := &LastingBucket{
		val:       val,
		check:     nil,
		bl:        defaultBucketLen,
		set:       false,
		threshold: int64(math.Ceil(defaultBucketLen / 3)), //默认是 在3分之一的时候填满
		life:      defaultLife,
	}

	//设置一些参数
	for _, opt := range opts {
		opt(lb)
	}

	//初始化桶
	lb.Init(ctx, cli, key)

	return lb
}

// BucketOnly 通内唯一
func (lb *LastingBucket) BucketOnly() *LastingBucket {
	lb.check = append(lb.check, func(ctx context.Context, val string) bool {
		list := lb.Cli.LRange(ctx, lb.Key, 0, -1).Val()
		for _, v := range list {
			if v == val {
				return false
			}
		}
		return true
	})
	return lb
}

// GetValue 获取持久桶的值
func (lb *LastingBucket) GetValue() (val string) {
	defer lb.ExtendExp(lb.life)
	defer lb.SetValue(true)
	for {
		//拿一个值出来
		val = lb.Cli.LPop(lb.Ctx, lb.Key).Val()
		if val != "" {
			return
		}

		//list里面没有值，那就去设置一下
		lb.SetValue(true)

		//没拿成功等20ms再拿一次
		time.Sleep(20 * time.Millisecond)
	}
}

// RemainLen 桶剩余长度
func (lb *LastingBucket) RemainLen() int64 {
	return lb.Cli.LLen(lb.Ctx, lb.Key).Val()
}

// SetValue 设置持久桶的值
func (lb *LastingBucket) SetValue(sync bool) {
	//一个实例不能设置两次
	if lb.set {
		return
	}
	lb.set = true

	if sync {
		//异步的时候要重新处理ctx值，防止初始的时候ctx被cancel导致的异常
		lb.Ctx = context.Background()
		go lb.pushBucket()
	} else {
		lb.pushBucket()
	}
}

// 是否需要填充
func (lb *LastingBucket) needFilling() (int64, bool) {
	//当前桶的长度
	nowBucketLen := lb.RemainLen()

	//如果当前桶长度还是大于阈值，就不需要填充
	if nowBucketLen > lb.threshold {
		return 0, false
	}

	//当需要创建
	return lb.bl - nowBucketLen, true
}

// 设置桶的值
func (lb *LastingBucket) pushBucket() {
	rl.NewRedisLock(lb.Ctx, lb.Cli, fmt.Sprintf("%s:LOCK", lb.Key)).LockDone(func() {
		createLen, ok := lb.needFilling()
		//不需要填充
		if !ok {
			return
		}

		//当前创建了多少个
		var cn, retry int64
		for {
			//创建一个值
			str := lb.val(lb.Ctx)

			// 验证
			var can = true
			for _, check := range lb.check {
				if can = check(lb.Ctx, str); !can {
					break
				}
			}

			// 验证通过
			if can {
				lb.Cli.RPush(lb.Ctx, lb.Key, str)
				cn++
			}

			//尝试次数
			retry++

			//创建够了直接返回，容错次数不能超过2倍的最大桶长度
			if createLen <= cn || retry >= lb.bl*2 {
				return
			}
		}
	})
}
