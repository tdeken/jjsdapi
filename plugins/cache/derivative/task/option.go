package task

import (
	"time"
)

// SetConsumeTaskOption 设置一些可配参数
type SetConsumeTaskOption func(o *Option)

type Option struct {
	taskIdKey    string                   //任务链路ID的key名字
	traceIdKey   string                   //链路ID的key名字
	ackMode      ackMode                  //ack模式
	nilWaitTime  time.Duration            //当数据拿完之后等待多少时间才消亡
	afterDone    func(c *Consume)         //消费完成后处理
	beforeAction func(c *Data)            //数据处理之前处理
	afterAction  func(c *Data, err error) //数据处理之后处理
}

// SetAckMode 设置ack模式
func SetAckMode(mode ackMode) SetConsumeTaskOption {
	return func(o *Option) {
		o.ackMode = mode
	}
}

// SetNilWaiteTime 设置当数据拿完之后等待多少时间才消亡
func SetNilWaiteTime(t time.Duration) SetConsumeTaskOption {
	return func(o *Option) {
		o.nilWaitTime = t
	}
}

// SetBeforeAction 设置消费之前处理
func SetBeforeAction(fn func(c *Data)) SetConsumeTaskOption {
	return func(o *Option) {
		o.beforeAction = fn
	}
}

// SetAfterAction 设置消费之前处理
func SetAfterAction(fn func(c *Data, err error)) SetConsumeTaskOption {
	return func(o *Option) {
		o.afterAction = fn
	}
}

// SetAfterDone 设置消费完成后处理
func SetAfterDone(fn func(c *Consume)) SetConsumeTaskOption {
	return func(o *Option) {
		o.afterDone = fn
	}
}

// SetTraceIdKey 设置消费完成后处理
func SetTraceIdKey(key string) SetConsumeTaskOption {
	return func(o *Option) {
		o.traceIdKey = key
	}
}

// SetTaskIdKey 设置消费完成后处理
func SetTaskIdKey(key string) SetConsumeTaskOption {
	return func(o *Option) {
		o.taskIdKey = key
	}
}

// 任务名称前缀，数据都放到这个前缀的key里面
func taskPrefix(taskName string) string {
	return "redis_task:" + taskName
}

// 任务锁前缀，用于保证每一个任务都都只被系统运行一次的锁
func taskLockSuffix(taskName string) string {
	return taskPrefix(taskName) + ":lock"
}

// 任务ack前缀，数据ack的临时数据表
func ackTaskSuffix(taskName string) string {
	return taskPrefix(taskName) + ":ack"
}
