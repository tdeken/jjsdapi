package task

import (
	"context"
	"encoding/json"
	grds "github.com/go-redis/redis/v8"
	"jjsdapi/plugins/cache/derivative"
	"reflect"
)

type RedisTask struct {
	derivative.Derivative
	cancel   context.CancelFunc
	taskName string //任务名称
}

// NewRedisTask 实例化一个redis任务实例
func NewRedisTask(ctx context.Context, cli *grds.Client, key string) *RedisTask {
	toCtx, cancel := context.WithCancel(ctx)
	task := &RedisTask{
		cancel:   cancel,
		taskName: taskPrefix(key),
	}

	task.Init(toCtx, cli, key)

	return task
}

// TaskDoing 任务是否在进行中
func (t *RedisTask) TaskDoing(ctx context.Context, key string) bool {
	if t.Cli.Exists(ctx, taskLockSuffix(key)).Val() > 0 {
		return true
	}
	return false
}

// PublicTask 发布消费任务
func (t *RedisTask) PublicTask(data ...interface{}) error {
	var jsonStr []string
	for _, item := range data {
		switch reflect.TypeOf(item).Kind() {
		case reflect.String:
			jsonStr = append(jsonStr, item.(string))
		default:
			str, _ := json.Marshal(item)
			jsonStr = append(jsonStr, string(str))
		}
	}
	return t.Cli.LPush(t.Ctx, t.taskName, jsonStr).Err()
}

// TaskPublicAndConsume 任务发布并消费
func (t *RedisTask) TaskPublicAndConsume(fc Action, data ...interface{}) (*Consume, error) {
	err := t.PublicTask(data...)
	if err != nil {
		return nil, err
	}

	consume, doing := t.ConsumeTask()
	if !doing {
		consume.Done(fc)
	}

	return consume, err
}
