package delay

import (
	"encoding/json"
	"github.com/rs/xid"
	"time"
)

type Producer struct {
	Delay
}

// Producer 生产者
func (d Delay) Producer() *Producer {
	return &Producer{d}
}

// Push 用生产者压入队列
func (p *Producer) Push(sec int, data interface{}, opts ...PushOption) (msgId string, err error) {
	t := time.Now()

	//执行的时间戳
	doT := t.Add(time.Duration(sec) * time.Second).Unix()

	tableName, err := p.dataTableName(float64(doT))
	if err != nil {
		return
	}

	b, err := json.Marshal(data)
	if err != nil {
		return
	}

	//传输结构
	var delivery = Delivery{
		PublishTime: t,
		Sec:         sec,
		Body:        b,
		MsgId:       xid.New().String(),
	}

	//处理参数
	for _, opt := range opts {
		opt(&delivery)
	}

	msgId = delivery.MsgId

	d, err := json.Marshal(delivery)
	if err != nil {
		return
	}

	err = p.Cli.RPush(p.Ctx, tableName, string(d)).Err()
	return
}
