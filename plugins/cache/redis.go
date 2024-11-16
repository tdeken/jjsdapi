package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
	"strings"
	"time"
)

var ErrCacheMiss = errors.New("cache miss")

type PluginConfig struct {
	Host               string //主机地址
	Password           string //连接密码
	DbNum              int    //连接库
	Sentinel           string //哨兵地址";"分隔
	SentinelMasterName string //master名称
}

type Plugin struct {
	//配置
	Config PluginConfig
	//redis实例
	Cache *redis.Client
}

func NewPlugin(pluginConfig PluginConfig) (plugin Plugin, err error) {
	var client *redis.Client
	if pluginConfig.Sentinel != "" {
		sentinelList := strings.Split(pluginConfig.Sentinel, ";")
		client = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    pluginConfig.SentinelMasterName,
			SentinelAddrs: sentinelList,
			Password:      pluginConfig.Password,
			DB:            pluginConfig.DbNum,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:     pluginConfig.Host,
			Password: pluginConfig.Password,
			DB:       pluginConfig.DbNum,
		})
	}
	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		err = fmt.Errorf("redis连接失败:%v", err)
		return
	}
	return Plugin{
		Config: pluginConfig,
		Cache:  client,
	}, nil
}

// Set 将键值对存储到Redis中，并设置过期时间
func (p *Plugin) Set(ctx context.Context, key string, value string, expire time.Duration) error {
	err := p.Cache.Set(ctx, key, value, expire).Err()
	if err != nil {
		return fmt.Errorf("failed to set value: %v", err)
	}

	return nil
}

// Get 从Redis中获取指定键的值，并将其反序列化为指定的类型
func (p *Plugin) Get(ctx context.Context, key string) (value string, err error) {
	value, err = p.Cache.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrCacheMiss
		}
		return "", fmt.Errorf("failed to get value: %v", err)
	}

	return
}

// SetObject 将对象对存储到Redis中，并设置过期时间
func (p *Plugin) SetObject(ctx context.Context, key string, value interface{}, expire time.Duration) error {
	// 将值序列化为JSON字符串
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %v", err)
	}

	err = p.Cache.Set(ctx, key, jsonValue, expire).Err()
	if err != nil {
		return fmt.Errorf("failed to set value: %v", err)
	}

	return nil
}

// GetObject 从Redis中获取指定键的对象，并将其反序列化为指定的类型
func (p *Plugin) GetObject(ctx context.Context, key string, value interface{}) error {
	jsonValue, err := p.Cache.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrCacheMiss
		}
		return fmt.Errorf("failed to get value: %v", err)
	}

	err = json.Unmarshal([]byte(jsonValue), value)
	if err != nil {
		return fmt.Errorf("failed to unmarshal value: %v", err)
	}

	return nil
}

// Del 从Redis中删除指定的键值对
func (p *Plugin) Del(ctx context.Context, key string) error {
	err := p.Cache.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete value: %v", err)
	}

	return nil
}

// Exists 检查指定的键是否存在于Redis中
func (p *Plugin) Exists(ctx context.Context, key string) (bool, error) {
	result, err := p.Cache.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check if key exists: %v", err)
	}

	return result == 1, nil
}

// Incr 对指定键所对应的值加1
func (p *Plugin) Incr(ctx context.Context, key string) (int64, error) {
	result, err := p.Cache.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment key: %v", err)
	}

	return result, nil
}

// Decr 对指定键所对应的值减1
func (p *Plugin) Decr(ctx context.Context, key string) (int64, error) {
	result, err := p.Cache.Decr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment key: %v", err)
	}

	return result, nil
}

type LockItem struct {
	key string
	val string
}

func (p *Plugin) UnlockItem(ctx context.Context, item LockItem) error {
	return p.Unlock(ctx, item.key, item.val)
}

// Lock 获取分布式锁
func (p *Plugin) Lock(ctx context.Context, key string, ttl time.Duration) (success bool, l LockItem, err error) {
	// 生成唯一随机字符串，用作锁值
	lockValue := xid.New().String()

	// 使用 SET key value NX PX 命令获取锁
	res, err := p.Cache.SetNX(ctx, key, lockValue, ttl).Result()
	if err != nil {
		return
	}
	if !res {
		return
	}

	// 加锁成功
	success = true
	l = LockItem{
		key: key,
		val: lockValue,
	}
	return
}

// Unlock 释放分布式锁
func (p *Plugin) Unlock(ctx context.Context, key string, lockValue string) error {
	// Lua脚本释放锁，只有锁值匹配才能释放锁
	script := "if redis.call(\"get\",KEYS[1]) == ARGV[1] then return redis.call(\"del\",KEYS[1]) else return 0 end"
	_, err := p.Cache.Eval(ctx, script, []string{key}, lockValue).Result()
	if err != nil {
		return err
	}

	// 释放锁成功
	return nil
}

// WaitAndAcquireLock 等待具有给定键和过期时间的分布式锁。
// 它尝试在间隔期内获取锁，直到成功或超时。
func (p *Plugin) WaitAndAcquireLock(ctx context.Context, key string, expiration time.Duration, timeout time.Duration) (success bool, l LockItem, err error) {
	startTime := time.Now()
	for {
		result, lockVal, er := p.Lock(ctx, key, expiration)
		if er != nil {
			err = er
			return
		}

		if result {
			// 成功获取到锁
			success = true
			l = lockVal
			return
		}

		// 检查是否超时
		if time.Since(startTime) >= timeout {
			err = errors.New("在指定的超时时间内未能获取到锁")
			return
		}

		// 等待一小段时间后重试
		time.Sleep(time.Millisecond * 50)
	}
}
