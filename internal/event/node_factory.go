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
	"github.com/quanxiang-cloud/flow/internal/convert"
	"github.com/quanxiang-cloud/flow/internal/event/node"
	"github.com/quanxiang-cloud/flow/internal/flow"
	"github.com/quanxiang-cloud/flow/internal/flow/callback_tasks"
	"github.com/quanxiang-cloud/flow/internal/models/mysql"
	"github.com/quanxiang-cloud/flow/internal/server/options"
	"github.com/quanxiang-cloud/flow/pkg/client"
	"github.com/quanxiang-cloud/flow/pkg/config"
)

// NodeFactory struct
type NodeFactory struct {
	email          *node.Email
	letter         *node.Letter
	cc             *node.CC
	dataCreate     *node.DataCreate
	dataUpdate     *node.DataUpdate
	webHook        *node.WebHook
	variableUpdate *node.VariableUpdate
	userTask       *node.UserTask
	delay          *node.Delay
}

// NewNodeFactory new
func NewNodeFactory(conf *config.Configs, opts ...options.Options) (*NodeFactory, error) {
	instance, err := flow.NewInstance(conf, opts...)
	if err != nil {
		return nil, nil
	}
	urge, err := callback_tasks.NewUrge(conf, opts...)
	if err != nil {
		return nil, nil
	}
	operationRecord, err := flow.NewOperationRecord(conf, opts...)
	if err != nil {
		return nil, nil
	}
	task, err := flow.NewTask(conf, opts...)
	if err != nil {
		return nil, nil
	}
	flow, err := flow.NewFlow(conf, opts...)
	if err != nil {
		return nil, nil
	}

	n := &node.Node{
		FlowRepo:               mysql.NewFlowRepo(),
		FormAPI:                client.NewForm(conf),
		InstanceRepo:           mysql.NewInstanceRepo(),
		InstanceVariablesRepo:  mysql.NewInstanceVariablesRepo(),
		AbnormalTaskRepo:       mysql.NewAbnormalTaskRepo(),
		FlowVariable:           mysql.NewVariablesRepo(),
		DispatcherCallbackRepo: mysql.NewDispatcherCallbackRepo(),
		MessageCenterAPI:       client.NewMessageCenter(conf),
		StructorAPI:            client.NewStructor(conf),
		ProcessAPI:             client.NewProcess(conf),
		IdentityAPI:            client.NewIdentity(conf),
		Dispatcher:             client.NewDispatcher(conf),
		Urge:                   *urge,
		Flow:                   flow,
		Instance:               instance,
		OperationRecord:        operationRecord,
		InstanceExecutionRepo:  mysql.NewInstanceExecutionRepo(),
		Task:                   task,
	}
	for _, opt := range opts {
		opt(n)
	}

	return &NodeFactory{
		email:          node.NewEmail(conf, n),
		letter:         node.NewLetter(conf, n),
		cc:             node.NewCC(conf, n),
		dataCreate:     node.NewDataCreate(conf, n),
		dataUpdate:     node.NewDataUpdate(conf, n),
		webHook:        node.NewWebHook(conf, n),
		variableUpdate: node.NewVariableUpdate(conf, n),
		userTask:       node.NewUserTask(conf, n),
		delay:          node.NewDelay(conf, n),
	}, nil
}

// GetNode get
func (f *NodeFactory) GetNode(nodeName string) node.INode {
	switch nodeName {
	case convert.Email:
		return f.email
	case convert.Letter:
		return f.letter
	case convert.Autocc:
		return f.cc
	case convert.TableDataCreate:
		return f.dataCreate
	case convert.TableDataUpdate:
		return f.dataUpdate
	case convert.WebHook:
		return f.webHook
	case convert.ProcessVariableAssignment:
		return f.variableUpdate
	case convert.Approve:
		fallthrough
	case convert.FillIn:
		return f.userTask
	case convert.Delayed:
		return f.delay
	}

	return nil
}
