package fiber

import (
	"context"
	"fmt"
	"jjsdapi/internal/config"
	"jjsdapi/internal/fiber/route"
	"jjsdapi/internal/fiber/server"
	"time"
)

func Run() (err error) {
	//注册路由
	route.Route()

	// do custom route

	//启动项目
	if err = server.Web.Server.Listen(fmt.Sprintf(":%d", config.Conf.Server.Port)); err != nil {
		return
	}
	return
}

func Shutdown() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = server.Web.Server.ShutdownWithContext(ctx); err != nil {
		return
	}

	// 给3秒时间，处理剩余程序未处理内容
	time.Sleep(3 * time.Second)
	return
}
