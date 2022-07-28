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

package node

import (
	"context"
	"github.com/quanxiang-cloud/flow/internal/convert"
	"github.com/quanxiang-cloud/flow/pkg/client"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/rpc/pb"
)

// CC struct
type CC struct {
	*Node
}

// NewCC new
func NewCC(conf *config.Configs, node *Node) *CC {
	return &CC{
		Node: node,
	}
}

// InitBegin event
func (n *CC) InitBegin(ctx context.Context, eventData *EventData) (*pb.NodeEventRespData, error) {
	return nil, nil
}

// InitEnd event
func (n *CC) InitEnd(ctx context.Context, eventData *EventData) (*pb.NodeEventRespData, error) {
	instance, err := n.InstanceRepo.GetEntityByProcessInstanceID(n.Db, eventData.ProcessInstanceID)
	if err != nil {
		return nil, err
	}

	ccUsers, _ := n.Flow.GetTaskHandleUserIDs(ctx, eventData.Shape, instance)

	if len(ccUsers) == 0 {
		return nil, nil
	}

	apiReq := &client.AddTaskReq{
		InstanceID: eventData.ProcessInstanceID,
		NodeDefKey: eventData.Shape.ID,
		Name:       eventData.Shape.Data.NodeData.Name,
		Desc:       convert.CcTask,
		Assignee:   ccUsers,
	}
	_, err = n.ProcessAPI.AddNonModelTask(ctx, apiReq)
	return nil, err
}
