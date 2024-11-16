package main

import (
	"jjsdapi/internal/boot"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//加载配置，系统全局变量
	boot.Init()

	//启动服务
	boot.Run()

	//监听程序退出
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)

	select {
	case <-ch:
		defer close(ch)
	}

	//关闭程序
	boot.Shutdown()
}
