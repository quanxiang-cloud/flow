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
	"errors"
	"fmt"
	"github.com/quanxiang-cloud/flow/internal/convert"
	"github.com/quanxiang-cloud/flow/pkg/client"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/utils"
	"github.com/quanxiang-cloud/flow/rpc/pb"
)

// UserTask struct
type UserTask struct {
	*Node
}

// NewUserTask new
func NewUserTask(conf *config.Configs, node *Node) *UserTask {
	return &UserTask{
		Node: node,
	}
}

// InitBegin event
func (n *UserTask) InitBegin(ctx context.Context, eventData *EventData) (*pb.NodeEventRespData, error) {
	fmt.Println("事件：加载AssigneeList，节点为" + eventData.NodeDefKey)
	// add dynamic assignee user list
	flowInstanceEntity, err := n.InstanceRepo.GetEntityByProcessInstanceID(n.Db, eventData.ProcessInstanceID)
	if err != nil {
		return nil, err
	}

	// assigneeList dynamic handle users
	_, assigneeList := n.Flow.GetTaskHandleUserIDs(ctx, eventData.Shape, flowInstanceEntity)

	fmt.Println("事件：加载AssigneeList，AssigneeList=" + utils.ChangeStringArrayToString(assigneeList))
	if len(assigneeList) > 0 {
		return nil, n.ProcessAPI.SetProcessVariables(ctx, eventData.ProcessID, eventData.ProcessInstanceID, eventData.NodeDefKey, convert.AssigneeList, assigneeList)
	}

	return nil, nil
}

// InitEnd event
func (n *UserTask) InitEnd(ctx context.Context, eventData *EventData) (*pb.NodeEventRespData, error) {
	flow, err := n.FlowRepo.FindByProcessID(n.Db, eventData.ProcessID)
	if err != nil {
		return nil, err
	}
	if flow == nil {
		flowProcessRelation, err := n.FlowProcessRelationRepo.FindByProcessID(n.Db, eventData.ProcessID)
		if err != nil {
			return nil, err
		}
		flow, err = n.FlowRepo.FindByID(n.Db, flowProcessRelation.FlowID)
		if err != nil {
			return nil, err
		}
		if flow == nil {
			return nil, errors.New("user node not match flow")
		}
	}
	instance, err := n.InstanceRepo.GetEntityByProcessInstanceID(n.Db, eventData.ProcessInstanceID)
	if err != nil {
		return nil, err
	}

	if len(eventData.TaskID) > 0 {
		apiReq := client.GetTasksReq{
			TaskID: eventData.TaskID,
		}
		resp, err := n.ProcessAPI.GetTasks(ctx, apiReq)
		if err != nil {
			return nil, err
		}
		if len(resp.Data) == 0 {
			return nil, nil
		}
		for _, value := range resp.Data {
			n.Task.TaskInitHandle(ctx, flow, instance, value, eventData.UserID)
		}
	}
	return nil, nil
}
