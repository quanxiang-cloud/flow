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

package rpc

import (
	"github.com/quanxiang-cloud/flow/internal/event"
	"github.com/quanxiang-cloud/flow/internal/server/options"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/misc/logger"
	"github.com/quanxiang-cloud/flow/pkg/misc/mysql2"
	"github.com/quanxiang-cloud/flow/rpc/pb"
	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/zrpc"
	"google.golang.org/grpc"
)

// RPC rpc
type RPC struct {
	configPath string
}

// NewRPCServer NewRpcServer
func NewRPCServer(configPath *string, config *config.Configs) (*zrpc.RpcServer, error) {
	var rpcServerConf zrpc.RpcServerConf
	conf.MustLoad(*configPath, &rpcServerConf)
	rpcServerConf.Timeout = 0

	// init db
	db, err := mysql2.New(config.Mysql, logger.Logger)
	if err != nil {
		return nil, err
	}

	optDB := options.WithDB(db)

	nodeEventHandler, err := event.NewNodeEventHandler(config, optDB)
	if err != nil {
		return nil, err
	}

	nodeEventsrv := NewNodeEventServer(nodeEventHandler)

	s := zrpc.MustNewServer(rpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterNodeEventServer(grpcServer, nodeEventsrv)
	})
	return s, nil
}
