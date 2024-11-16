package task

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Action 执行操作函数
type Action func(td *Data) error

type Data struct {
	c       *Consume
	Ctx     context.Context
	MsgId   string    //当前消息体的唯一标识
	Msg     []byte    //数据结构体
	reenter bool      //是否重新进入队列
	Begin   time.Time //拿到数据的时间
}

// 把消息体加入任务的ack表
func (td *Data) addAckTable() {
	//把消息体存到hash表里面
	td.c.Cli.HSet(td.Ctx, td.c.AckName, td.MsgId, string(td.Msg))
}

// 加入历史数据
func (td *Data) addHistory(fcErr error) {
	if td.reenter { //重新入队的不算处理
		return
	}

	history := historyValue{
		msg:   string(td.Msg),
		date:  time.Now().Format("2006-01-02 15:04:05"),
		err:   fcErr,
		msgId: td.MsgId,
	}

	td.c.history = append(td.c.history, history)
}

// Ack 说明任务已经完成了
func (td *Data) Ack() {
	td.c.Cli.HDel(td.Ctx, td.c.AckName, td.MsgId)

	//如何设定了消亡时间，那要重置
	if td.c.death != 0 {
		td.c.death = 0
	}
}

// DealAgain 回滚数据，重新第一位重新处理
func (td *Data) DealAgain(times int64) error {
	msgRetryKey := fmt.Sprintf("%s_retry_%s", td.c.taskName, string(td.Msg))
	retryNum := td.c.Cli.Incr(td.Ctx, msgRetryKey).Val()
	if retryNum == 1 { //第一次重试的时候就要设定过期时间
		td.c.Cli.Expire(td.Ctx, msgRetryKey, 5*time.Minute)
	}

	if retryNum > times { // 重试次数大于给予的次数
		return ReenterLimit
	}

	return td.c.Cli.RPush(context.Background(), td.c.taskName, string(td.Msg)).Err()
}

// ReenterList 重新入队
func (td *Data) ReenterList() error {
	err := td.c.Cli.LPush(td.Ctx, td.c.taskName, string(td.Msg)).Err()
	if err != nil {
		return err
	}
	td.reenter = true
	return nil
}

// TaskExit 主动退出任务
func (td *Data) TaskExit(clear bool) {
	if td.c.RemainLen() > 0 && clear { //清楚数据，主动退出
		td.c.Cli.Del(td.c.Ctx, td.c.taskName)
	}
	td.c.exit <- nil
}

// ParsingMsg 解析数据
func (td *Data) ParsingMsg(v interface{}) error {
	return json.Unmarshal(td.Msg, &v)
}
