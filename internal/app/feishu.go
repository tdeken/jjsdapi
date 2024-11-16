package app

import (
	"fmt"
	"jjsdapi/internal/config"
	"jjsdapi/plugins/alert"
)

var FeiShu alert.Plugin

func InitFeiShu() (err error) {
	if config.Conf.Alert.FeiShu.Url != "" {
		alertConfig := alert.PluginConfig{
			Type:   alert.PluginConfigTypeFeiShu,
			Url:    config.Conf.Alert.FeiShu.Url,
			Secret: config.Conf.Alert.FeiShu.Secret,
		}

		for _, v := range config.Conf.Alert.Scenes {
			alertConfig.Scenes = append(alertConfig.Scenes, &alert.Conf{
				Type:   alert.PluginConfigTypeFeiShu,
				Url:    v.Url,
				Secret: v.Secret,
				Scene:  v.Scene,
			})
		}

		FeiShu, err = alert.NewPlugin(alertConfig)
	}
	return
}

const (
	normal = "normal"
	notice = "notice"
)

// FeishuNormalSend 发送普通消息，这个给自己监听用的
func FeishuNormalSend(format string, a ...any) {
	h, ok := FeiShu.Scenes[normal]
	if ok {
		h.SendText(fmt.Sprintf(format, a...))
	}
}

// FeishuNoticeSend 发送通知消息，这个是一些业务通知用的
func FeishuNoticeSend(format string, a ...any) {
	h, ok := FeiShu.Scenes[notice]
	if ok {
		h.SendText(fmt.Sprintf(format, a...))
	}
}
