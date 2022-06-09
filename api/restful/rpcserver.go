package restful

import (
	"github.com/quanxiang-cloud/flow/internal/callback"
	"github.com/quanxiang-cloud/flow/internal/server"
	"github.com/quanxiang-cloud/flow/internal/server/options"
	"github.com/quanxiang-cloud/flow/internal/svc"
	"github.com/quanxiang-cloud/flow/pb"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/misc/logger"
	"github.com/quanxiang-cloud/flow/pkg/misc/mysql2"
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
	var c svc.Config
	conf.MustLoad(*configPath, &c)
	c.Timeout = 600000

	// init db
	db, err := mysql2.New(config.Mysql, logger.Logger)
	if err != nil {
		return nil, err
	}

	optDB := options.WithDB(db)
	callback, err := callback.NewCallback(config, optDB)
	if err != nil {
		return nil, err
	}

	ctx := svc.NewServiceContext(c, callback)
	srv := server.NewEventServer(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterEventServer(grpcServer, srv)
	})
	return s, nil
}
