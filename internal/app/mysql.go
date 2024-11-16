package app

import (
	"jjsdapi/internal/config"
	"jjsdapi/internal/repository/dao"
	"jjsdapi/plugins/orm"
)

var DB orm.Plugin

func InitDB() (err error) {
	dbConfig := orm.PluginConfig{
		URL:             config.Conf.Database.URL,
		MaxIdleConn:     config.Conf.Database.MaxIdleConn,
		MaxOpenConn:     config.Conf.Database.MaxOpenConn,
		ConnMaxLifetime: config.Conf.Database.ConnMaxLifetime,
		Logger:          Logger,
		SlowThreshold:   config.Conf.Database.SlowThreshold,
	}
	DB, err = orm.NewPlugin(dbConfig)
	if err != nil {
		return
	}

	dao.SetDefault(DB.Client)
	return
}
