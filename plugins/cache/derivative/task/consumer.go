package task

import (
	"context"
	"errors"
	"fmt"
	grds "github.com/go-redis/redis/v8"
	"github.com/rs/xid"
	lock3 "jjsdapi/plugins/cache/derivative/lock"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ack模式
type ackMode int8

var (
	AbnormalErr  = errors.New("任务获取异常")
	ReenterLimit = errors.New("重试已超过上限")
	HasTaskDoing = errors.New("任务已经在执行")
	waitDeath    = errors.New("等待消亡")
)

const (
	AutoAck        ackMode = iota //自动ACK
	ManualAck                     //手动ACK
	DoneSuccessAck                //操作正常ACK即（TaskDone的err == nil才ack）
)

const (
	defaultNilWaitTime = 200 * time.Millisecond      //300毫秒秒后消亡
	heartbeatTime      = 10 * time.Minute            //心跳检测时间
	lockExp            = heartbeatTime + time.Minute //锁的过期时间，延长一分钟，防止临界值问题
)

// Consume 消费者结构体
type Consume struct {
	*RedisTask
	AckName string           //响应hash表
	tid     string           //运行的任务ID
	lock    *lock3.RedisLock //任务锁
	exit    chan error       //退出通道，用于通知任务退出
	history []historyValue   //历史消息体处理时间，相同的消息体只取最后一个
	opt     *Option          //一些参数
	death   int64            //死亡时间
	action  Action           //消费消息体的方法
	err     error            //执行错误
	chb     chan bool        //关闭心跳 closeHeartbeat
}

// 当前消费账者消费的数据
type historyValue struct {
	msg   string
	date  string
	err   error
	msgId string
}

// ConsumeTask 创建消费者消费任务
func (t *RedisTask) ConsumeTask(opts ...SetConsumeTaskOption) (consume *Consume, doing bool) {
	option := &Option{
		ackMode:     AutoAck,
		nilWaitTime: defaultNilWaitTime,
		beforeAction: func(c *Data) {

		},
		afterAction: func(c *Data, err error) {

		},
	}

	for _, opt := range opts {
		opt(option)
	}

	//建造消费者
	consume = &Consume{
		RedisTask: t,
		AckName:   ackTaskSuffix(t.Key),
		exit:      make(chan error),
		opt:       option,
		chb:       make(chan bool),
	}

	//检查一下任务是否在进行中
	doing = t.TaskDoing(t.Ctx, t.Key)

	return
}

// Done 消费任务
func (c *Consume) Done(fc Action) {
	if !c.addLock() { //加锁不成功表示任务已经在进行中了，直接退出
		c.err = HasTaskDoing
		return
	}

	if c.opt.taskIdKey != "" {
		c.Ctx = context.WithValue(c.Ctx, c.opt.taskIdKey, c.tid)
	}

	//载入消费方法
	c.action = fc

	defer c.Release()
	defer c.closeHeartbeat() //任务完成停止心跳
	defer c.AfterConsume()   //消费完成处理

	//把没有被ack的消息体重新入队
	c.ReenterAckList()

	//监听消息体消费信息
	go func() {
		defer func() {
			if err := recover(); err != nil {
				c.err = errors.New(fmt.Sprintf("执行操作panic:%v", err))
				c.exit <- c.err
			}
		}()
		c.exec()
	}()

	//监听任务心跳，处理锁问题
	go c.checkHeartbeat()

	//监听程序退出
	select {
	case err := <-c.exit: //收到退出请求退出
		if !errors.Is(err, grds.Nil) && err != nil {
			if c.err == nil {
				c.err = AbnormalErr //异常错误信息
			}
		}
	}
}

// 消费消息体
func (c *Consume) exec() {
	for {
		select {
		case <-c.Ctx.Done():
			return
		default:
			res, err := c.getTaskData()

			//等待消亡
			if errors.Is(err, waitDeath) {
				continue
			}

			//有错误数据
			if err != nil {
				return
			}
			//数据处理
			c.dataProcessing(res)
		}
	}
}

// 获取数据
func (c *Consume) getTaskData() ([]byte, error) {
	res, err := c.Cli.RPop(c.Ctx, c.taskName).Bytes()
	if err == nil {
		return res, nil
	}

	if errors.Is(err, grds.Nil) { //如果是数据拿完了
		if c.death == 0 { //看看有没有设定消亡时间
			c.death = time.Now().Add(c.opt.nilWaitTime).UnixNano() / 1e6
		}

		if c.death < (time.Now().UnixNano() / 1e6) { //如果消亡时间小于当前时间，那么休息100ms再去取一次数据，直到到了消亡时间都没有数据进来，那就直接灭亡
			time.Sleep(100 * time.Millisecond)
			return nil, waitDeath
		}
	}

	c.exit <- err
	return nil, err
}

// 监听任务心跳，处理锁问题
func (c *Consume) checkHeartbeat() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)
	timer := time.NewTicker(heartbeatTime)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			c.lock.ExtendExp(lockExp)
		case <-ch: //程序退出，只要解锁redis，防止死锁
			return
		case <-c.chb: //任务完成停止心跳
			return
		}
	}
}

// 数据处理
func (c *Consume) dataProcessing(res []byte) {
	msgId := xid.New().String()

	ctx := c.Ctx
	if c.opt.traceIdKey != "" {
		ctx = context.WithValue(c.Ctx, c.opt.traceIdKey, msgId)
	}

	//任务数据
	taskData := &Data{
		c:     c,
		Ctx:   ctx,
		Msg:   res,
		MsgId: msgId,
		Begin: time.Now(),
	}

	//把消息体加入任务的ack表
	taskData.addAckTable()

	//处理前要做的事情
	c.opt.beforeAction(taskData)

	//执行处理方法
	fcErr := c.action(taskData)

	//处理后要做的事情
	c.opt.afterAction(taskData, fcErr)

	//ack
	if c.opt.ackMode == AutoAck || (c.opt.ackMode == DoneSuccessAck && fcErr == nil) {
		taskData.Ack()
	}

	//加入到历史数据
	taskData.addHistory(fcErr)
}

// ReenterAckList 把没有ack的消息体重新加入任务
func (c *Consume) ReenterAckList() {
	data := c.Cli.HGetAll(c.Ctx, c.AckName).Val()

	if len(data) == 0 {
		return
	}
	defer c.Cli.Del(c.Ctx, c.AckName)

	var msgList []string
	for _, msg := range data {
		msgList = append(msgList, msg)
	}

	c.Cli.RPush(c.Ctx, c.taskName, msgList)
}

// 消费错误
func (c *Consume) Error() error {
	return c.err
}

// 停止心跳
func (c *Consume) closeHeartbeat() {
	c.chb <- true
}

// Release 释放消费者占用资源
func (c *Consume) Release() {
	c.lock.Release() //解锁
	c.cancel()       //关闭ctx
	close(c.exit)    //关闭错误通道
	close(c.chb)     //关闭错误通道
}

// AfterConsume 消费完成后
func (c *Consume) AfterConsume() {
	if c.opt.afterDone == nil {
		return
	}
	c.opt.afterDone(c)
}

// 任务加锁 保证系统的这个消费任务只有一个
func (c *Consume) addLock() bool {
	c.tid = xid.New().String()
	lock := lock3.NewRedisLock(
		c.Ctx,
		c.Cli,
		taskLockSuffix(c.taskName),
		lock3.SetExpire(lockExp), //过期时间是心跳检测时间的一分钟
		lock3.SetValue(c.tid),
	)
	if !lock.SoftLock() { //没有设置成功，说明消费任务在运行着
		c.tid = c.Cli.Get(c.Ctx, taskLockSuffix(c.taskName)).Val()
		return false
	}
	c.lock = lock

	return true
}

// RemainLen 剩余长度
func (c *Consume) RemainLen() int64 {
	remainLen, _ := c.Cli.LLen(c.Ctx, c.taskName).Result() //获取剩余消息体长度
	return remainLen
}

// DoneLen 已经消费的消息体数
func (c *Consume) DoneLen() int {
	return len(c.history)
}

// GetFinishInfo 获取完成的信息
func (c *Consume) GetFinishInfo() (int, []string) {
	var infoStr []string
	for _, one := range c.history {
		infoStr = append(infoStr, fmt.Sprintf("消息Id: %s, 消费数据: %s, 消费完成时间: %s, 错误信息: %v", one.msgId, one.msg, one.date, one.err))
	}

	return len(c.history), infoStr
}

// GetTaskName 获取任务名称
func (c *Consume) GetTaskName() string {
	return c.taskName
}
