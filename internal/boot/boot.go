package boot

import (
	"flag"
	"fmt"
	"jjsdapi/internal/app"
	"jjsdapi/internal/config"
	"jjsdapi/internal/fiber"
	"log"
)

// Init 初始化项目配置
func Init() {
	var env string
	flag.StringVar(&env, "env", "dev", "启动环境")
	flag.Parse()

	err := config.LoadConfig(config.FilePath{
		ConfigName: fmt.Sprintf("config-%s", env),
		ConfigType: "yaml",
		ConfigPath: "etc",
	})
	if err != nil {
		log.Fatalf("配置文件加载失败:%v", err)
	}

	//报警推送
	err = app.InitFeiShu()
	if err != nil {
		return
	}

	//日志
	app.InitLogger()

	//mysql
	err = app.InitDB()
	if err != nil {
		return
	}

	//redis
	err = app.InitRedis()
	if err != nil {
		return
	}
}

func Run() {
	go fiber.Run()

	return
}

// Shutdown 关闭运行程序
func Shutdown() {
	if err := fiber.Shutdown(); err != nil {
		panic(err)
	}
}
