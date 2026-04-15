package main

import (
	"flag"
	"fmt"
	"log"
)

var envString string

func init() {
	flag.StringVar(&envString, "envString", "dev", "环境配置：dev | prod")
}

func main() {
	flag.Parse()
	fmt.Printf("▶ 启动环境: %s\n", envString)

	app, err := InitApp(envString)
	if err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	if err = app.Run(); err != nil {
		log.Fatalf("服务异常退出: %v", err)
	}
}
