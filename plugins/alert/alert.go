package alert

import (
	"jjsdapi/plugins/alert/webhook"
)

const (
	PluginConfigTypeFeiShu = iota //飞书
	PluginConfigTypeWxWork        //企微
)

type Conf struct {
	Type   int    //类型
	Url    string //webhook地址
	Secret string //秘钥
	Scene  string //使用场景
}

type PluginConfig struct {
	Type   int    //类型
	Url    string //webhook地址
	Secret string //秘钥
	Scenes []*Conf
}

type Plugin struct {
	//配置
	Config PluginConfig
	//redis实例
	Alert webhook.WebHook

	//场景下的推送
	Scenes map[string]webhook.WebHook
}

func NewPlugin(pluginConfig PluginConfig) (plugin Plugin, err error) {
	var alert webhook.WebHook
	switch pluginConfig.Type {
	case PluginConfigTypeFeiShu:
		alert = webhook.NewFeiShu(pluginConfig.Url, pluginConfig.Secret)
	case PluginConfigTypeWxWork:
		alert = webhook.NewWxWork(pluginConfig.Url, pluginConfig.Secret)
	}

	plugin = Plugin{
		Config: pluginConfig,
		Alert:  alert,
		Scenes: map[string]webhook.WebHook{},
	}

	for _, conf := range pluginConfig.Scenes {
		switch conf.Type {
		case PluginConfigTypeFeiShu:
			plugin.Scenes[conf.Scene] = webhook.NewFeiShu(conf.Url, conf.Secret)
		case PluginConfigTypeWxWork:
			plugin.Scenes[conf.Scene] = webhook.NewWxWork(conf.Url, conf.Secret)
		}
	}

	return
}
