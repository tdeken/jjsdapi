package app

import (
	"jjsdapi/internal/config"
	"jjsdapi/plugins/cache"
)

var Redis cache.Plugin

func InitRedis() (err error) {
	cacheConfig := cache.PluginConfig{
		Host:               config.Conf.Redis.Addr,
		Password:           config.Conf.Redis.Password,
		DbNum:              config.Conf.Redis.DB,
		Sentinel:           config.Conf.Redis.Sentinel,
		SentinelMasterName: config.Conf.Redis.SentinelMasterName,
	}
	Redis, err = cache.NewPlugin(cacheConfig)
	return
}
