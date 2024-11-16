package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type PluginConfig struct {
	ModeProd        bool                                                         //是否线上模式
	Level           zapcore.Level                                                //日志等级
	StdOut          bool                                                         //是否输出命令行
	FileOut         bool                                                         //是否输出文件日志
	JsonEncode      bool                                                         //是否json形式输出
	StackTraceLevel zapcore.Level                                                //达到该等级的日志输出链路信息
	FileConfig      FileConfig                                                   //文件配置
	AlarmWriter     io.Writer                                                    //告警写入
	Env             string                                                       //环境
	ContextTraceKey string                                                       //trace key
	ContextKey      []string                                                     //ctx key
	BeforePrint     func(ctx context.Context, lv zapcore.Level, msg string) bool //打印前
}

type Plugin struct {
	//配置
	Config PluginConfig
	//logger输出
	LoggerPrinter *zap.Logger
}

type FileConfig struct {
	Path       string //文件输出路径
	MaxSize    int    //一个文件最大大小(单位:MB)
	MaxBackups int    //保留旧文件个数
	MaxAge     int    //保留旧文件最大天数
	LocalTime  bool   //是否使用本地时间
	Compress   bool   //是否压缩/归档旧文件
}

func NewPlugin(pluginConfig PluginConfig) Plugin {
	//trace
	if pluginConfig.ContextTraceKey == "" {
		pluginConfig.ContextTraceKey = "trace_id"
	}
	//writer
	var writer zapcore.WriteSyncer      //普通写入
	var alarmWriter zapcore.WriteSyncer //告警写入
	var syncer []zapcore.WriteSyncer
	var alarmSyncer []zapcore.WriteSyncer

	if pluginConfig.StdOut {
		syncer = append(syncer, zapcore.AddSync(os.Stdout))
	}
	if pluginConfig.FileOut {
		umberJackLogger := &lumberjack.Logger{
			Filename:   pluginConfig.FileConfig.Path,       //日志文件的位置
			MaxSize:    pluginConfig.FileConfig.MaxSize,    //在进行切割之前，日志文件的最大大小（以MB为单位）
			MaxAge:     pluginConfig.FileConfig.MaxAge,     //保留旧文件的最大天数
			MaxBackups: pluginConfig.FileConfig.MaxBackups, //保留旧文件的最大个数
			LocalTime:  pluginConfig.FileConfig.LocalTime,
			Compress:   pluginConfig.FileConfig.Compress, //是否压缩/归档旧文件
		}
		syncer = append(syncer, zapcore.AddSync(umberJackLogger))
	}
	if pluginConfig.AlarmWriter != nil {
		alarmSyncer = append(alarmSyncer, zapcore.AddSync(pluginConfig.AlarmWriter))
	}
	writer = zapcore.NewMultiWriteSyncer(syncer...)
	alarmWriter = zapcore.NewMultiWriteSyncer(alarmSyncer...)
	//encoder
	var encoder zapcore.Encoder
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	if pluginConfig.ModeProd {
		encoderConfig = zap.NewProductionEncoderConfig()
	}
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006/01/02 15:04:05.000")
	encoder = zapcore.NewConsoleEncoder(encoderConfig)
	if pluginConfig.JsonEncode {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}
	//core
	var coreTree []zapcore.Core
	coreTree = append(coreTree, zapcore.NewCore(encoder, writer, zap.NewAtomicLevelAt(pluginConfig.Level)))
	if alarmWriter != nil {
		coreTree = append(coreTree, zapcore.NewCore(encoder, alarmWriter, zap.NewAtomicLevelAt(zapcore.ErrorLevel)))
	}
	core := zapcore.NewTee(coreTree...)
	//logger
	loggerPrinter := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2), zap.AddStacktrace(pluginConfig.StackTraceLevel))
	return Plugin{
		Config:        pluginConfig,
		LoggerPrinter: loggerPrinter,
	}
}

//统一输出

func (p Plugin) Print(ctx context.Context, lv zapcore.Level, msg string, opName ...string) {
	finalMsg := msg
	if len(opName) > 0 {
		finalMsg = fmt.Sprintf("%v (%v)", finalMsg, strings.Join(opName, "|"))
	}

	if len(p.Config.ContextKey) > 0 {
		for _, key := range p.Config.ContextKey {
			val := ctx.Value(key)
			if val != nil {
				finalMsg = fmt.Sprintf("[%s.%v] %v", key, fmt.Sprint(val), finalMsg)
			}
		}
	}

	if p.Config.ContextTraceKey != "" {
		traceId := ctx.Value(p.Config.ContextTraceKey)
		if traceId != nil {
			finalMsg = fmt.Sprintf("[%v] %v", fmt.Sprint(traceId), finalMsg)
		}
	}

	if p.Config.Env != "" {
		finalMsg = fmt.Sprintf("[%v] %v", fmt.Sprint(p.Config.Env), finalMsg)
	}

	ce := p.LoggerPrinter.Check(lv, finalMsg)
	if ce != nil {
		ce.Write()
	}
}

func (p Plugin) Debug(ctx context.Context, msg string, opName ...string) {
	lv := zap.DebugLevel
	p.Print(ctx, lv, msg, opName...)
}

func (p Plugin) Info(ctx context.Context, msg string, opName ...string) {
	lv := zap.InfoLevel
	p.Print(ctx, lv, msg, opName...)
}
func (p Plugin) Warn(ctx context.Context, msg string, opName ...string) {
	lv := zap.WarnLevel
	p.Print(ctx, lv, msg, opName...)
}
func (p Plugin) Error(ctx context.Context, msg string, opName ...string) {
	lv := zap.ErrorLevel
	p.Print(ctx, lv, msg, opName...)
}
func (p Plugin) ErrorWithStack(ctx context.Context, msg string, opName ...string) {
	msg += stackInfo()
	lv := zap.ErrorLevel
	p.Print(ctx, lv, msg, opName...)
}
func (p Plugin) Panic(ctx context.Context, msg string, opName ...string) {
	lv := zap.PanicLevel
	p.Print(ctx, lv, msg, opName...)
}

func stackInfo() (info string) {
	//获取计数器
	callers := make([]uintptr, 6)
	n := runtime.Callers(3, callers)
	frames := runtime.CallersFrames(callers[:n])
	info = "\n堆栈信息:\n"
	for {
		frame, more := frames.Next()
		info += fmt.Sprintf("%s:%d %s\n", frame.File, frame.Line, strings.TrimPrefix(frame.Function, "."))
		if !more {
			break
		}
	}
	return
}
