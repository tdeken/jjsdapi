package app

import (
	"context"
	"fmt"
	"go.uber.org/zap/zapcore"
	"jjsdapi/internal/config"
	"jjsdapi/internal/consts"
	"jjsdapi/plugins/logger"
)

var Logger logger.Plugin

func InitLogger() {
	loggerConfig := logger.PluginConfig{
		ModeProd:        config.Conf.Server.Env == consts.EnvProd,
		Level:           zapcore.InfoLevel,
		StdOut:          config.Conf.Logger.StdOut,
		FileOut:         config.Conf.Logger.FileOut,
		JsonEncode:      false,
		StackTraceLevel: zapcore.ErrorLevel,
		FileConfig: logger.FileConfig{
			Path:       config.Conf.Logger.Path,
			MaxSize:    config.Conf.Logger.MaxSize,
			MaxBackups: config.Conf.Logger.MaxBackups,
			MaxAge:     config.Conf.Logger.MaxAge,
			LocalTime:  true,
			Compress:   config.Conf.Logger.Compress,
		},
		AlarmWriter:     FeiShu.Alert,
		Env:             config.Conf.Server.Env,
		ContextTraceKey: consts.TraceIdKey,
		ContextKey:      []string{consts.CtxReqId, consts.CtxBootId, consts.CtxMsgId},
	}
	Logger = logger.NewPlugin(loggerConfig)

}

func TryCatch(ctx context.Context, msg string) {
	if err := recover(); err != nil {
		Logger.Error(ctx, fmt.Sprintf("备注信息：%s, 错误信息：%v", msg, err))
	}
}
