package logic

import (
	"context"
	"fmt"
	"git.internal.yunify.com/qxp/misc/logger"
	"github.com/quanxiang-cloud/flow/internal/svc"
	"github.com/quanxiang-cloud/flow/pb"
	"github.com/quanxiang-cloud/flow/pkg"
)

// EventHandler handler
var EventHandler = map[string]handlerFunc{
	"test":           testEventHandler,
	"taskStartEvent": taskStartEventHandler,
	"taskEndEvent":   taskEndEventHandler,
}

type handlerFunc func(ctx context.Context, svcCtx *svc.ServiceContext, in *pb.EventReq) *pb.EventResp

// PublishLogic PublishLogic
type PublishLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewPublishLogic NewPublishLogic
func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishLogic {
	return &PublishLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Publish Publish
func (l *PublishLogic) Publish(in *pb.EventReq) (*pb.EventResp, error) {
	if handler, ok := EventHandler[in.EventName]; ok {
		r := handler(l.ctx, l.svcCtx, in)
		return r, nil
	}
	return &pb.EventResp{}, nil
}

func testEventHandler(ctx context.Context, svcCtx *svc.ServiceContext, in *pb.EventReq) *pb.EventResp {
	fmt.Println("==================")
	fmt.Println(in.EventData)
	return &pb.EventResp{
		Result: "success",
	}
}

func taskStartEventHandler(ctx context.Context, svcCtx *svc.ServiceContext, in *pb.EventReq) *pb.EventResp {
	fmt.Println("==================")
	fmt.Println(in.EventData)

	err := svcCtx.Callback.TaskStartEventHandler(ctx, in.EventData)
	if err != nil {
		logger.Logger.Errorw("rpc start event execute error:"+err.Error(), pkg.STDRequestID(ctx))
		return &pb.EventResp{
			Result: "fail",
		}
	}

	return &pb.EventResp{
		Result: "success",
	}
}

func taskEndEventHandler(ctx context.Context, svcCtx *svc.ServiceContext, in *pb.EventReq) *pb.EventResp {
	fmt.Println("==================")
	fmt.Println(in.EventData)

	err := svcCtx.Callback.TaskEndEventHandler(ctx, in.EventData)
	if err != nil {
		logger.Logger.Errorw("rpc end event execute error:"+err.Error(), pkg.STDRequestID(ctx))
		return &pb.EventResp{
			Result: "fail",
		}
	}
	return &pb.EventResp{
		Result: "success",
	}
}
