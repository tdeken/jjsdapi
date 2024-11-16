package lock

import (
	"context"
	"fmt"
	grds "github.com/go-redis/redis/v8"
	cmap "github.com/orcaman/concurrent-map/v2"
	"time"
)

var subMap = cmap.New[*suber]()

const (
	subLockTopic = "---sub_lock_topic---"
	okStatus     = "ok"
)

type suber struct {
	subs []chan subRes
}

type subRes struct {
	status int
	err    error
}

const (
	SomeErr  = iota // 通知解锁
	DoLockOk        // 上锁成功
	NoticeOk        // 通知解锁
	ReLockOk        // 重锁成功
)

type SubLock struct {
	ILock
	status int
}

// NewSubLock 实例化订阅锁
// 当锁被释放之后，发布消息通知该订阅锁的客户端
func NewSubLock(ctx context.Context, cli *grds.Client, key string, opts ...Option) *SubLock {
	l := &SubLock{
		status: SomeErr,
	}

	// 初始化锁
	l.init(ctx, cli, key)

	//设置配置
	for _, opt := range opts {
		opt(&l.ILock)
	}

	return l
}

// Lock 锁
func (l *SubLock) Lock() (lc *SubLock) {
	lc = l
	exp := l.exp + 2*time.Second
	if l.SoftLock() {
		l.status = DoLockOk
		return
	}

	ch := make(chan subRes, 1)
	var redisSub bool

	var newSub = &suber{
		subs: make([]chan subRes, 0, 50),
	}
	subMap.Upsert(l.topic(), newSub, func(exist bool, valueInMap *suber, newValue *suber) *suber {
		if !exist {
			redisSub = true
			return newValue
		}

		valueInMap.subs = append(valueInMap.subs, ch)

		return valueInMap
	})

	if redisSub {
		l.redisSub(exp)
		return
	}

	l.memorySub(ch, exp)

	return
}

// 内存订阅
func (l *SubLock) memorySub(ch chan subRes, exp time.Duration) {
	tc, cancel := context.WithTimeout(context.Background(), exp)
	defer cancel()

	for {
		select {
		case <-l.Ctx.Done():
			l.err = CtxCancelErr
			return
		case <-tc.Done():
			l.err = FailTimeoutErr
			return
		case msg := <-ch:
			l.status = msg.status
			l.err = msg.err
			return
		default:
			if l.SoftLock() {
				l.status = ReLockOk
				return
			}
			time.Sleep(DefaultRetryTime)
		}
	}

}

// redis订阅
func (l *SubLock) redisSub(exp time.Duration) {
	tc, cancel := context.WithTimeout(context.Background(), exp)
	defer cancel()

	// 订阅主题
	rsb := l.Cli.Subscribe(tc, l.topic())
	defer rsb.Close()
	_, err := rsb.ReceiveTimeout(tc, exp)
	if err != nil {
		l.err = err
		return
	}

	for {
		select {
		case <-tc.Done(): // timeout
			l.err = FailTimeoutErr
			return
		case msg := <-rsb.Channel(): // 接收一个消息
			// 拿到状态和数据
			if msg.Payload == okStatus {
				l.pubMemory()
				return
			}

		default:
			if l.SoftLock() {
				l.pubMemory()
				return
			}
			time.Sleep(DefaultRetryTime)
		}
	}
}

// 通知内存锁
func (l *SubLock) pubMemory() {
	group, ok := subMap.Pop(l.topic())
	if !ok {
		return
	}

	for _, v := range group.subs {
		v <- subRes{status: NoticeOk}
	}
}

// Release 释放锁
func (l *SubLock) Release() {
	if !l.lockStatus {
		return
	}

	l.release()

	// 订阅通知
	_, err := l.Cli.Publish(l.Ctx, l.topic(), okStatus).Result()
	if err != nil {
		l.err = err
		return
	}

	return
}

// Status 解锁状态
func (l *SubLock) Status() int {
	return l.status
}

// 订阅key
func (l *SubLock) topic() string {
	return fmt.Sprintf("%s:%s", l.Key, subLockTopic)
}
