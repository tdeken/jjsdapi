package test

import (
	"fmt"
	"jjsdapi/internal/app"
	"jjsdapi/internal/config"
	"jjsdapi/internal/repository/dao"
	"os"
)

func init() {

	var idx = 5
	for _, arg := range os.Args {
		if arg == "-test.bench" {
			idx = 7
		}
	}

	var filePath = config.FilePath{
		ConfigName: "config-dev",
		ConfigType: "yaml",
		ConfigPath: fmt.Sprintf("%s/etc", os.Args[idx]),
	}

	err := config.LoadConfig(filePath)
	if err != nil {
		panic(err)
	}

	err = app.InitFeiShu()
	if err != nil {
		panic(err)
	}

	app.InitLogger()

	err = app.InitDB()
	if err != nil {
		panic(err)
	}

	// 设置dao的数据库
	dao.SetDefault(app.DB.Client)

	err = app.InitRedis()
	if err != nil {
		panic(err)
	}

}
