package delay

import (
	"context"
	"fmt"
	grds "github.com/go-redis/redis/v8"
	"reflect"
	"testing"
	"time"
)

var client *grds.Client

func connUp() {
	client = grds.NewClient(&grds.Options{
		Addr:     "8.135.23.37:6379",
		Password: "lylbdev@redis",
		DB:       8,
	})
}

func TestCustomer(t *testing.T) {
	connUp()

	d := NewRedisDelayQueue(context.Background(), client, "ttt")

	for i := 0; i < 3; i++ {
		num := i
		go func() {
			d.Listen(func(ctx context.Context, msg []byte) error {
				t.Log(num, string(msg))
				time.Sleep(time.Second)
				return nil
			})
		}()
	}

	select {}
}

type A struct {
	Name string `json:"name"`
}

type DelayQueueMsg struct {
	Scenes      string                 `json:"scenes"`
	PublishTime time.Time              `json:"publish_time"`
	Data        map[string]interface{} `json:"data"`
}

func TestProduct(t *testing.T) {
	connUp()

	d := NewRedisDelayQueue(context.Background(), client, "ttt")

	for i := 1; i <= 10; i++ {
		ti := time.Now().Add(time.Duration(i+(i-1)*5) * time.Second)

		err := d.Push(ti, DelayQueueMsg{
			Scenes:      "111",
			PublishTime: time.Now(),
			Data:        structToMap(A{Name: ti.Format("2006-01-02 15:04:05")}),
		})
		t.Log(err)
	}

}

func TestDelay_Consumer(t *testing.T) {
	connUp()

	d := NewRedisDelayQueue(context.Background(), client, "ttt1")

	go func() {
		time.Sleep(10 * time.Second)
		Stop()
	}()

	for i := 0; i < 3; i++ {
		num := i
		go func() {
			d.Consumer(func(ctx context.Context, d Delivery) error {
				t.Log(num, d.MsgId, d.Sec, d.Ttl, string(d.Body), time.Now().Sub(d.PublishTime).Seconds())
				return nil
			}).ConsumerListen()
		}()
	}

	select {}
}

type DelayProducerData struct {
	Scenes string                 `json:"scenes"`
	Data   map[string]interface{} `json:"data"`
}

func TestDelay_Producer(t *testing.T) {
	connUp()

	d := NewRedisDelayQueue(context.Background(), client, "ttt1")

	for i := 1; i <= 10; i++ {

		_, err := d.Producer().Push(i+(i-1)*5, DelayProducerData{
			Scenes: "111",
			Data:   structToMap(A{Name: fmt.Sprintf("测试%d", i)}),
		})
		t.Log(err)
	}
}

// structToMap 结构体转为Map[string]interface{}
func structToMap(in interface{}) map[string]interface{} {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct { // 非结构体返回错误提示
		return nil
	}

	t := v.Type()
	// 遍历结构体字段
	// 指定tagName值为map中key;字段值为map中value
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagValue := fi.Tag.Get("json"); tagValue != "" {
			out[tagValue] = v.Field(i).Interface()
		}
	}
	return out
}
