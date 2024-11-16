package orm

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

type PluginConfig struct {
	URL             string        // 数据库连接地址
	MaxIdleConn     int           // 最大空闲连接数
	MaxOpenConn     int           // 最大活动连接数
	ConnMaxLifetime time.Duration // 连接最大生命周期(单位:秒)
	Logger          Logger        //日志
	SlowThreshold   time.Duration //慢sql阈值(单位:秒)
}

type Plugin struct {
	//配置
	Config PluginConfig
	//gorm实例
	Client *gorm.DB
}

func NewPlugin(pluginConfig PluginConfig) (plugin Plugin, err error) {
	gormLogger := NewLogger(pluginConfig.Logger, logger.Config{
		SlowThreshold:             pluginConfig.SlowThreshold,
		Colorful:                  false,
		IgnoreRecordNotFoundError: true,
		LogLevel:                  logger.Info,
	})
	cfg := &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            false, //预编译模式
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: gormLogger,
	}

	client, err := gorm.Open(mysql.Open(pluginConfig.URL), cfg)
	if err == nil {
		sqlDB, _ := client.DB()
		sqlDB.SetMaxIdleConns(pluginConfig.MaxIdleConn)
		sqlDB.SetMaxOpenConns(pluginConfig.MaxOpenConn)
		sqlDB.SetConnMaxLifetime(pluginConfig.ConnMaxLifetime)
	}

	return Plugin{
		Config: pluginConfig,
		Client: client,
	}, err
}
