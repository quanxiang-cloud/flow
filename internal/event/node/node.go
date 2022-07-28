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
	"github.com/quanxiang-cloud/flow/internal/flow"
	"github.com/quanxiang-cloud/flow/internal/flow/callback_tasks"
	"github.com/quanxiang-cloud/flow/internal/models"
	"github.com/quanxiang-cloud/flow/pkg/client"
	"gorm.io/gorm"
)

// Node struct
type Node struct {
	Db                     *gorm.DB
	FlowRepo               models.FlowRepo
	InstanceRepo           models.InstanceRepo
	InstanceVariablesRepo  models.InstanceVariablesRepo
	AbnormalTaskRepo       models.AbnormalTaskRepo
	InstanceExecutionRepo  models.InstanceExecutionRepo
	FlowVariable           models.VariablesRepo
	DispatcherCallbackRepo models.DispatcherCallbackRepo
	Urge                   callback_tasks.Urge
	Flow                   flow.Flow
	Instance               flow.Instance
	OperationRecord        flow.OperationRecord
	Task                   flow.Task
	FormAPI                client.Form
	MessageCenterAPI       client.MessageCenter
	StructorAPI            client.Structor
	ProcessAPI             client.Process
	IdentityAPI            client.Identity
	Dispatcher             client.Dispatcher
}

// SetDB set db
func (n *Node) SetDB(db *gorm.DB) {
	n.Db = db
}
