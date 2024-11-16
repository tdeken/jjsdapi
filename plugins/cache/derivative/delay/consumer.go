package delay

import (
	"context"
	"encoding/json"
	"time"
)

type ConsumerAction func(ctx context.Context, d Delivery) error

type Logger interface {
	Trace(ctx context.Context, diy Delivery) context.Context //消费者链路跟踪
	Error(ctx context.Context, err error, diy Delivery)      //执行错误
}

type Consumers []*Consumer

// RegisterConsumer 注册消费者
func RegisterConsumer(consumers ...*Consumer) Consumers {
	var cs Consumers
	if len(consumers) > 0 {
		cs = consumers
	}

	return cs
}

// Listen 监听
func (c Consumers) Listen() (err error) {
	l := len(c)
	for i := 0; i < l; i++ {
		go func(i int) {
			re := c[i].ConsumerListen()
			if re != nil {
				err = re
			}
		}(i)
	}
	return
}

type Consumer struct {
	Delay
	action ConsumerAction
	logger Logger
	death  func(ctx context.Context, diy Delivery) error
}

// Consumer 消费者
func (d Delay) Consumer(action ConsumerAction, opts ...ConsumeOption) *Consumer {
	c := &Consumer{
		Delay:  d,
		action: action,
	}
	for _, opt := range opts {
		opt(c)
	}

	return c
}

// ConsumerListen 消费者监听
func (c Consumer) ConsumerListen() error {
	delivery, err := c.GetDelivery()
	if err != nil {
		return err
	}
	for {
		select {
		case d, ok := <-delivery:
			if !ok {
				return nil
			}
			go c.actionDone(d)
		}
	}
}

// GetDelivery 获取数据
func (c *Consumer) GetDelivery() (<-chan Delivery, error) {
	ch := make(chan Delivery)

	go func(ch chan Delivery) {
		_ = c.Listen(func(ctx context.Context, msg []byte) error {
			var data Delivery

			_ = json.Unmarshal(msg, &data)

			//过了生命周期，就直接销毁
			if data.Ttl > 0 && data.PublishTime.Add(time.Duration(data.Ttl)*time.Second).Before(time.Now()) {
				return c.deathDone(ctx, data)
			}

			ch <- data
			return nil
		})
	}(ch)

	return ch, nil
}

// 处理日志ctx
func (c *Consumer) ctxWithLogger(d Delivery) context.Context {
	ctx := context.Background()
	if c.logger != nil {
		ctx = c.logger.Trace(ctx, d)
	}

	return ctx
}

// 处理日志ctx
func (c *Consumer) deathDone(ctx context.Context, d Delivery) error {
	if c.death == nil {
		return nil
	}

	if c.logger != nil {
		ctx = c.logger.Trace(ctx, d)
	}

	return c.death(ctx, d)
}

// 执行数据处理
func (c Consumer) actionDone(d Delivery) {
	ctx := c.ctxWithLogger(d)
	err := c.action(ctx, d)
	if c.logger != nil {
		c.logger.Error(ctx, err, d)
	}
}
