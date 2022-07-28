/*
Copyright 2022 QuanxiangCloud Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"github.com/quanxiang-cloud/flow/pkg/misc/logger"
	rpc2 "github.com/quanxiang-cloud/flow/rpc"
	"os"
	"os/signal"
	"syscall"

	"github.com/quanxiang-cloud/flow/pkg/redis"

	"github.com/quanxiang-cloud/flow/api/restful"
	"github.com/quanxiang-cloud/flow/pkg/config"
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
	rpc, err := rpc2.NewRPCServer(configPath, config.Config)
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
