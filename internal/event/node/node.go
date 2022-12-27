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
	"encoding/json"
	"github.com/quanxiang-cloud/flow/internal/flow"
	"github.com/quanxiang-cloud/flow/internal/flow/callback_tasks"
	"github.com/quanxiang-cloud/flow/internal/models"
	"github.com/quanxiang-cloud/flow/pkg/client"
	"github.com/quanxiang-cloud/flow/pkg/misc/logger"
	"gorm.io/gorm"
	"strings"
)

// Node struct
type Node struct {
	Db                      *gorm.DB
	FlowRepo                models.FlowRepo
	InstanceRepo            models.InstanceRepo
	InstanceVariablesRepo   models.InstanceVariablesRepo
	InstanceStepRepo        models.InstanceStepRepo
	AbnormalTaskRepo        models.AbnormalTaskRepo
	InstanceExecutionRepo   models.InstanceExecutionRepo
	FlowVariable            models.VariablesRepo
	DispatcherCallbackRepo  models.DispatcherCallbackRepo
	Urge                    callback_tasks.Urge
	Flow                    flow.Flow
	Instance                flow.Instance
	OperationRecord         flow.OperationRecord
	Task                    flow.Task
	FormAPI                 client.Form
	MessageCenterAPI        client.MessageCenter
	StructorAPI             client.Structor
	ProcessAPI              client.Process
	IdentityAPI             client.Identity
	Dispatcher              client.Dispatcher
	FlowProcessRelationRepo models.FlowProcessRelationRepo
}

// SetDB set db
func (n *Node) SetDB(db *gorm.DB) {
	n.Db = db
}

func (n *Node) CheckRefuse(ctx context.Context, db *gorm.DB, processInstanceID string) bool {
	instanceSteps, err := n.InstanceStepRepo.FindInstanceStepsByStatus(n.Db, processInstanceID, []string{"REFUSE"})
	if err != nil {
		return false
	}
	if len(instanceSteps) > 0 {
		return false
	}
	tasksReq := client.GetTasksReq{
		InstanceID: []string{processInstanceID},
	}
	taskResp, _ := n.ProcessAPI.GetHistoryTasks(ctx, tasksReq)
	tasks := taskResp.Data
	marshal, _ := json.Marshal(taskResp)
	logger.Logger.Debug("taskResp==", string(marshal))
	for k := range tasks {
		if strings.Contains(tasks[k].Comments, "REFUSE") {
			return false
		}
	}
	return true
}
