package delay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	grds "github.com/go-redis/redis/v8"
	"jjsdapi/internal/internalplugins/cache/derivative"
	rl "jjsdapi/plugins/cache/derivative/lock"
	"time"
)

var delayQueue = make(map[int]*Delay)

type Action func(ctx context.Context, msg []byte) error

type Delivery struct {
	AppId       string    `json:"app_id"`       //应用ID
	MsgId       string    `json:"msg_id"`       //消息ID
	PublishTime time.Time `json:"publish_time"` //发布时间
	Type        string    `json:"type"`         //消息类型
	Sec         int       `json:"sec"`          //延迟秒数
	Ttl         int       `json:"ttl"`          //生命周期
	Body        []byte    `json:"body"`         //发布的数据
}

type Delay struct {
	derivative.Derivative
	database string        //数据库
	exit     chan struct{} //退出
	idx      int
}

// NewRedisDelayQueue 实例一个redis延时队列
func NewRedisDelayQueue(ctx context.Context, cli *grds.Client, key string) *Delay {
	delay := &Delay{
		exit: make(chan struct{}),
	}

	//初始化
	delay.Init(ctx, cli, key)

	//设置一些相关的名称（如：数据库名称）
	delay.setName()

	return delay
}

// Stop 暂停监听
func (d *Delay) Stop() {
	d.exit <- struct{}{}
	delete(delayQueue, d.idx)
}

// Listen 监听延时队列
func (d *Delay) Listen(action Action) error {
	for {
		select {
		case <-d.exit:
			return nil
		default:
			wait, _ := d.consume(action)
			if wait {
				time.Sleep(time.Second)
			}
		}
	}
}

// 消费数据
func (d *Delay) consume(action Action) (wait bool, err error) {
	//拿到要消费的表名称
	tables, err := d.getDataTableNames("0", fmt.Sprintf("%d", time.Now().Unix()))
	if err != nil {
		return
	}

	//没有要运行的表，就等1s再运行
	if len(tables) == 0 {
		wait = true
		return
	}

	//处理要消费的表数据
	for _, table := range tables {
		for {
			//拿每一条数据出来
			row, getErr := d.Cli.LPop(d.Ctx, table).Bytes()

			//如果数据取空了，就直接把表删掉
			if errors.Is(getErr, grds.Nil) {
				d.delDataTable(table)
			}
			if getErr != nil {
				break
			}

			//协程消费消息
			go func() {
				_ = action(context.Background(), row)
			}()
		}
	}

	return
}

// 监听之前处理一些数据
func (d *Delay) beforeListen() {
	idx := len(delayQueue)
	d.idx, d.exit = idx, make(chan struct{})
	delayQueue[idx] = d
}

// Push 把数据压入延时表
func (d *Delay) Push(t time.Time, data ...interface{}) (err error) {
	tableName, err := d.dataTableName(float64(t.Unix()))
	if err != nil {
		return err
	}

	var pushData []string
	for _, v := range data {
		b, _ := json.Marshal(v)
		pushData = append(pushData, string(b))
	}

	err = d.Cli.RPush(d.Ctx, tableName, pushData).Err()
	return
}

// 删除数据表
func (d *Delay) delDataTable(table string) {
	d.Cli.ZRem(d.Ctx, d.database, table)
}

// 获取区间的所有数据表
func (d *Delay) getDataTableNames(min, max string) ([]string, error) {
	return d.Cli.ZRangeByScore(d.Ctx, d.database, &grds.ZRangeBy{
		Min: min,
		Max: max,
	}).Result()
}

// 获取某个时间的表名
func (d *Delay) dataTableName(score float64) (string, error) {
	sStr := fmt.Sprintf("%f", score)
	tableNames, err := d.getDataTableNames(sStr, sStr)
	if err != nil {
		return "", err
	}

	//存在表名就直接返回
	if len(tableNames) > 0 {
		return tableNames[0], nil
	}

	//创建一张表，上锁，防止重复创建
	lock := rl.NewRedisLock(d.Ctx, d.Cli, d.Key).Lock()
	defer lock.Release()
	if lock.Error() != nil {
		return "", lock.Error()
	}

	//可以创建的时候，先拿一次，看看是不是已经被创建了
	tableNames, err = d.getDataTableNames(sStr, sStr)
	if err != nil {
		return "", err
	}

	if len(tableNames) > 0 {
		return tableNames[0], nil
	}

	//创建一张表
	tableName := dataTableName(d.Key, score)
	_, err = d.Cli.ZAddNX(d.Ctx, d.database, &grds.Z{
		Score:  score,
		Member: tableName,
	}).Result()
	if err != nil {
		return "", err
	}

	return tableName, nil
}

// 设置一些redis键名称
func (d *Delay) setName() {
	d.database = fmt.Sprintf("delay_time_table:%s", d.Key)
}

// 获取一个数据表名称
func dataTableName(key string, score float64) string {
	return fmt.Sprintf("delay_data_table:%s_%f", key, score)
}

// Stop 暂停延时队列
func Stop() {
	if len(delayQueue) == 0 {
		return
	}
	for _, v := range delayQueue {
		go v.Stop()
	}
}
