package main

import (
	"flag"
	"github.com/quanxiang-cloud/flow/pkg/redis"
	"os"
	"os/signal"
	"syscall"

	"github.com/quanxiang-cloud/flow/api/restful"
	"github.com/quanxiang-cloud/flow/pkg/config"

	"github.com/quanxiang-cloud/flow/pkg/misc/logger"
)

var (
	configPath = flag.String("config", "configs/config.yml", "-config 配置文件地址")
)

func main() {
	flag.Parse()

	err := config.Init(*configPath)
	if err != nil {
		panic(err)
	}

	err = logger.New(&config.Config.Log)
	if err != nil {
		panic(err)
	}

	err = redis.Init()
	if err != nil {
		panic(err)
	}
	// 启动路由
	router, err := restful.NewRouter(config.Config)
	if err != nil {
		panic(err)
	}
	// start rpc server
	rpc, err := restful.NewRPCServer(configPath, config.Config)
	if err != nil {
		panic(err)
	}

	go rpc.Start()
	go router.Run()

	// Start kafka consumer
	//err = mq.NewkafkaConsumer(config.Config)
	//if err != nil {
	//	panic(err)
	//}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			router.Close()
			rpc.Stop()
			logger.Sync()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
