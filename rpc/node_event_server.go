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
	"context"
	"fmt"
	"github.com/quanxiang-cloud/flow/internal/event"
	"github.com/quanxiang-cloud/flow/pkg"
	"github.com/quanxiang-cloud/flow/rpc/pb"
	"google.golang.org/grpc/metadata"
)

// NodeEventServer node event server
type NodeEventServer struct {
	eventHandler event.NodeEventHandler
}

// NewNodeEventServer new
func NewNodeEventServer(eventHandler event.NodeEventHandler) *NodeEventServer {
	return &NodeEventServer{
		eventHandler: eventHandler,
	}
}

// Event event
func (nes *NodeEventServer) Event(ctx context.Context, in *pb.NodeEventReq) (*pb.NodeEventResp, error) {
	var resp *pb.NodeEventRespData
	var err error

	// request header
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		fmt.Println("get metadata error")
	}
	ctx = pkg.RPCCTXTransfer(md["request-id"][0], md["user-id"][0])

	switch in.EventType {
	case "nodeInitBeginEvent":
		resp, err = nes.eventHandler.NodeStartEventHandler(ctx, in.Data)
	case "nodeInitEndEvent":
		resp, err = nes.eventHandler.NodeEndEventHandler(ctx, in.Data)
	case "complete":
		//
	}

	if err != nil {
		return &pb.NodeEventResp{
			Code: pb.NodeEventResp_FAILURE,
			Msg:  err.Error(),
		}, err
	}

	return &pb.NodeEventResp{
		Code: pb.NodeEventResp_SUCCESS,
		Data: resp,
	}, nil
}
