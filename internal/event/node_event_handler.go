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

package event

import (
	"context"
	"encoding/json"
	"github.com/quanxiang-cloud/flow/internal/event/node"
	"github.com/quanxiang-cloud/flow/internal/flow"
	"github.com/quanxiang-cloud/flow/internal/server/options"
	"github.com/quanxiang-cloud/flow/pkg"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/misc/logger"
	"github.com/quanxiang-cloud/flow/rpc/pb"
)

// NodeEventHandler node event handler
type NodeEventHandler interface {
	NodeStartEventHandler(ctx context.Context, eventReqData *pb.NodeEventReqData) (*pb.NodeEventRespData, error)
	NodeEndEventHandler(ctx context.Context, eventReqData *pb.NodeEventReqData) (*pb.NodeEventRespData, error)
	convertToEventData(ctx context.Context, eventReqData *pb.NodeEventReqData) (*node.EventData, error)
}

type nodeEventHandler struct {
	nodeFactory *NodeFactory
	Flow        flow.Flow
}

// NewNodeEventHandler new
func NewNodeEventHandler(conf *config.Configs, opts ...options.Options) (NodeEventHandler, error) {
	nodeFactory, err := NewNodeFactory(conf, opts...)
	if err != nil {
		return nil, nil
	}
	flow2, err := flow.NewFlow(conf, opts...)
	if err != nil {
		return nil, nil
	}

	c := &nodeEventHandler{
		nodeFactory: nodeFactory,
		Flow:        flow2,
	}

	return c, nil
}

func (neh *nodeEventHandler) convertToEventData(ctx context.Context, eventReqData *pb.NodeEventReqData) (*node.EventData, error) {
	eventData := &node.EventData{
		ProcessID:         eventReqData.ProcessID,
		NodeDefKey:        eventReqData.NodeDefKey,
		ProcessInstanceID: eventReqData.ProcessInstanceID,
		RequestID:         pkg.STDRequestID2(ctx),
		UserID:            pkg.STDUserID(ctx),
		TaskID:            eventReqData.TaskID,
		ExecutionID:       eventReqData.ExecutionID,
	}

	ctx = pkg.RPCCTXTransfer(eventData.RequestID, eventData.UserID)

	shape, err := neh.Flow.GetShapeByProcessID(ctx, eventData.ProcessID, eventData.NodeDefKey)
	if err != nil {
		return nil, err
	}
	eventData.Shape = shape
	return eventData, nil
}

// NodeStartEventHandler event init start event(not task event)
func (neh *nodeEventHandler) NodeStartEventHandler(ctx context.Context, eventReqData *pb.NodeEventReqData) (*pb.NodeEventRespData, error) {
	marshal, _ := json.Marshal(eventReqData)
	logger.Logger.Info("eventdata=,", string(marshal))

	eventData, err := neh.convertToEventData(ctx, eventReqData)
	if err != nil {
		return nil, err
	}

	node2 := neh.nodeFactory.GetNode(eventData.Shape.Type)

	if node2 != nil {
		return node2.InitBegin(ctx, eventData)
	}

	return nil, nil
}

// NodeEndEventHandler event init end event(not task event, is event event)，
func (neh *nodeEventHandler) NodeEndEventHandler(ctx context.Context, eventReqData *pb.NodeEventReqData) (*pb.NodeEventRespData, error) {
	eventData, err := neh.convertToEventData(ctx, eventReqData)
	if err != nil {
		return nil, err
	}

	node2 := neh.nodeFactory.GetNode(eventData.Shape.Type)

	if node2 != nil {
		return node2.InitEnd(ctx, eventData)
	}

	return nil, nil
}
