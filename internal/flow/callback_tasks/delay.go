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

package callback_tasks

import (
	"context"
	"encoding/json"
	"errors"
	flow2 "github.com/quanxiang-cloud/flow/internal/flow"
	"github.com/quanxiang-cloud/flow/internal/models/mysql"
	"github.com/quanxiang-cloud/flow/internal/server/options"
	"github.com/quanxiang-cloud/flow/pkg/client"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"gorm.io/gorm"
	"strings"
)

// Delay info
//type Delay interface {
//	UrgingExecute(ctx context.Context, code string) error
//
//	TaskUrge(ctx context.Context, req *TaskUrgeModel) error
//}

type Delay struct {
	C
}

// NewDelay init
func NewDelay(conf *config.Configs, opts ...options.Options) (*Delay, error) {
	var d Delay
	task, _ := flow2.NewTask(conf, opts...)
	flow, _ := flow2.NewFlow(conf, opts...)
	instance, _ := flow2.NewInstance(conf, opts...)
	operationRecord, _ := flow2.NewOperationRecord(conf, opts...)

	d.processAPI = client.NewProcess(conf)
	d.formAPI = client.NewForm(conf)
	d.dispatcherCallbackRepo = mysql.NewDispatcherCallbackRepo()
	d.dispatcherAPI = client.NewDispatcher(conf)
	d.urgeRepo = mysql.NewUrgeRepo()
	d.instanceRepo = mysql.NewInstanceRepo()
	d.flowRepo = mysql.NewFlowRepo()
	d.task = task
	d.flow = flow
	d.instance = instance
	d.operationRecord = operationRecord

	for _, opt := range opts {
		opt(&d)
	}
	return &d, nil
}

// SetDB set db
func (u *Delay) SetDB(db *gorm.DB) {
	u.db = db
}

// Execute 执行
// 处理回调
func (u *Delay) Execute(ctx context.Context, code *string) error {
	c := strings.Split(*code, "_")
	if len(c) != 3 {
		return errors.New("bad code fmt")
	}
	data, err := u.dispatcherCallbackRepo.FindByID(u.db, *code)
	if err != nil {
		return err
	}
	req := client.CompleteNodeReq{}
	if err := json.Unmarshal([]byte(data.OtherInfo), &req); err != nil {
		return err
	}
	if err := u.processAPI.CompleteNode(ctx, &req); err != nil {
		return err
	}
	return nil
}
