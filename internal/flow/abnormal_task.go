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

package flow

import (
	"context"
	"github.com/quanxiang-cloud/flow/internal/convert"
	"github.com/quanxiang-cloud/flow/internal/models"
	"github.com/quanxiang-cloud/flow/internal/models/mysql"
	"github.com/quanxiang-cloud/flow/internal/server/options"
	"github.com/quanxiang-cloud/flow/pkg"
	"github.com/quanxiang-cloud/flow/pkg/client"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/misc/error2"
	"github.com/quanxiang-cloud/flow/pkg/page"
	"gorm.io/gorm"
)

// AbnormalTask service
type AbnormalTask interface {
	List(ctx context.Context, req *models.AbnormalTaskReq) (*page.RespPage, error)
	AdminStepBack(ctx context.Context, req *AdminTaskReq, model *models.HandleTaskModel) (bool, error)
	AdminSendBack(ctx context.Context, req *AdminTaskReq) (bool, error)
	AdminAbandon(ctx context.Context, req *AdminTaskReq) (bool, error)
	AdminDeliverTask(ctx context.Context, req *AdminTaskReq, model *models.HandleTaskModel) (bool, error)
	AdminGetTaskForm(ctx context.Context, req *AdminTaskReq) (*InstanceDetailModel, error)
}

type abnormalTask struct {
	db               *gorm.DB
	conf             *config.Configs
	abnormalTaskRepo models.AbnormalTaskRepo
	processAPI       client.Process
	operationRecord  OperationRecord
	instanceRepo     models.InstanceRepo
	flowRepo         models.FlowRepo
	formAPI          client.Form
	appCenterAPI     client.AppCenter
}

// NewAbnormalTask init
func NewAbnormalTask(conf *config.Configs, opts ...options.Options) (AbnormalTask, error) {
	operationRecord, _ := NewOperationRecord(conf, opts...)
	t := &abnormalTask{
		conf:             conf,
		abnormalTaskRepo: mysql.NewAbnormalTaskRepo(),
		processAPI:       client.NewProcess(conf),
		operationRecord:  operationRecord,
		instanceRepo:     mysql.NewInstanceRepo(),
		flowRepo:         mysql.NewFlowRepo(),
		formAPI:          client.NewForm(conf),
		appCenterAPI:     client.NewAppCenter(conf),
	}

	for _, opt := range opts {
		opt(t)
	}
	return t, nil
}

// SetDB set db
func (a *abnormalTask) SetDB(db *gorm.DB) {
	a.db = db
}

func (a *abnormalTask) List(ctx context.Context, req *models.AbnormalTaskReq) (*page.RespPage, error) {
	appIDs, err := a.appCenterAPI.GetAdminAppIDs(ctx)
	if err != nil || len(appIDs) == 0 {
		return &page.RespPage{
			Data:       nil,
			TotalCount: 0,
		}, nil
	}

	req.AdminAppIDs = appIDs
	orderItem := page.OrderItem{
		Column:    "create_time",
		Direction: page.Asc,
	}
	req.Orders = []page.OrderItem{orderItem}
	tasks, count, _ := a.abnormalTaskRepo.Page(a.db, req)

	pages := page.RespPage{
		Data:       tasks,
		TotalCount: count,
	}

	return &pages, nil
}

func (a *abnormalTask) AdminStepBack(ctx context.Context, req *AdminTaskReq, model *models.HandleTaskModel) (bool, error) {
	if model.TaskDefKey == "" {
		return false, error2.NewErrorWithString(error2.Internal, "event id is nil")
	}
	resp, flag, err := a.check(ctx, req.ProcessInstanceID, req.TaskID)
	if err != nil {
		return flag, err
	}

	flowInstanceEntity, err := a.instanceRepo.GetEntityByProcessInstanceID(a.db, req.ProcessInstanceID)
	if err != nil {
		return false, err
	}

	toNode, err := a.processAPI.GetModelNode(ctx, req.ProcessInstanceID, model.TaskDefKey)
	if err != nil {
		return false, error2.NewErrorWithString(error2.Internal, "Can not find event ")
	}

	if err = a.processAPI.StepBack(ctx, req.ProcessInstanceID, req.TaskID, model.TaskDefKey); err != nil {
		return false, err
	}
	if err = a.abnormalTaskRepo.UpdateByProcessInstanceID(a.db, req.ProcessInstanceID, map[string]interface{}{
		"status": 1,
	}); err != nil {
		return false, err
	}

	// Add operation record
	model.HandleType = opStepBack
	model.HandleDesc = "将工作流回退至“" + toNode.Name + "”"
	a.operationRecord.AddOperationRecord(ctx, flowInstanceEntity, resp.Data[0], model)

	return true, nil
}

func (a *abnormalTask) AdminSendBack(ctx context.Context, req *AdminTaskReq) (bool, error) {
	resp, flag, err := a.check(ctx, req.ProcessInstanceID, req.TaskID)
	if err != nil {
		return flag, err
	}

	addTaskReq := client.AddTaskReq{
		InstanceID: req.ProcessInstanceID,
		TaskID:     req.TaskID,
		UserID:     pkg.STDUserID(ctx),
		Name:       "打回重填",
		Desc:       convert.SendBackTask,
	}
	err = a.processAPI.SendBack(ctx, addTaskReq)
	if err != nil {
		return false, err
	}

	flowInstanceEntity, err := a.instanceRepo.GetEntityByProcessInstanceID(a.db, req.ProcessInstanceID)
	if err != nil {
		return false, err
	}
	if flowInstanceEntity == nil {
		return false, error2.NewErrorWithString(error2.Internal, "Can not find flow instance data ")
	}

	updateMap := make(map[string]interface{}, 0)
	updateMap["modifier_id"] = pkg.STDUserID(ctx)
	updateMap["status"] = SendBack
	if err = a.instanceRepo.Update(a.db, flowInstanceEntity.ID, updateMap); err != nil {
		return false, err
	}
	a.instanceRepo.Update(a.db, flowInstanceEntity.ID, updateMap)
	if err = a.abnormalTaskRepo.UpdateByProcessInstanceID(a.db, req.ProcessInstanceID, map[string]interface{}{
		"status": 1,
	}); err != nil {
		return false, err
	}

	// Add operation record
	handleTaskModel := &models.HandleTaskModel{
		HandleType: opSendBack,
		HandleDesc: "将工作流打回至发起人",
	}
	a.operationRecord.AddOperationRecord(ctx, flowInstanceEntity, resp.Data[0], handleTaskModel)

	return true, nil
}

func (a *abnormalTask) AdminAbandon(ctx context.Context, req *AdminTaskReq) (bool, error) {
	if !a.isUnhandleAbnormalTask(req.ProcessInstanceID, req.TaskID) {
		return false, error2.NewErrorWithString(error2.Internal, "is not abnormal task")
	}
	instance, err := a.instanceRepo.GetEntityByProcessInstanceID(a.db, req.ProcessInstanceID)
	if err != nil {
		return false, err
	}
	if instance == nil {
		return false, error2.NewErrorWithString(error2.Internal, "Can not find flow instance data ")
	}

	flow, err := a.flowRepo.FindByID(a.db, instance.FlowID)
	if err != nil {
		return false, err
	}
	if flow == nil {
		return false, error2.NewErrorWithString(error2.Internal, "Can not find flow data ")
	}

	tasks, err := a.processAPI.GetTasksByInstanceID(ctx, req.ProcessInstanceID)
	if err != nil {
		return false, err
	}

	err = a.processAPI.AbendInstance(ctx, req.ProcessInstanceID)
	if err != nil {
		return false, err
	}

	updateMap := make(map[string]interface{}, 0)
	updateMap["modifier_id"] = pkg.STDUserID(ctx)
	updateMap["status"] = Abandon
	if err = a.instanceRepo.Update(a.db, instance.ID, updateMap); err != nil {
		return false, err
	}

	if err = a.abnormalTaskRepo.UpdateByProcessInstanceID(a.db, req.ProcessInstanceID, map[string]interface{}{
		"status": 1,
	}); err != nil {
		return false, err
	}

	// Add operation record
	handleTaskModel := &models.HandleTaskModel{
		HandleType: opAbandon,
		HandleDesc: "作废",
	}
	a.operationRecord.AddOperationRecord(ctx, instance, tasks[0], handleTaskModel)

	return true, nil
}

func (a *abnormalTask) AdminDeliverTask(ctx context.Context, req *AdminTaskReq, model *models.HandleTaskModel) (bool, error) {
	userID := pkg.STDUserID(ctx)

	if len(model.HandleUserIDs) != 1 {
		return false, error2.NewErrorWithString(error2.Internal, "Deliver user requried ")
	}
	if model.HandleUserIDs[0] == userID {
		return false, error2.NewErrorWithString(error2.Internal, "Can not deliver self ")
	}
	resp, f, err := a.check(ctx, req.ProcessInstanceID, req.TaskID)
	if err != nil {
		return f, err
	}
	flowInstanceEntity, err := a.instanceRepo.GetEntityByProcessInstanceID(a.db, req.ProcessInstanceID)
	if err != nil {
		return false, err
	}
	if flowInstanceEntity == nil {
		return false, error2.NewErrorWithString(error2.Internal, "Can not find flow instance data ")
	}

	task := resp.Data[0]

	addHistoryTaskReq := &client.AddHistoryTaskReq{
		Name:        task.Name,
		Desc:        task.Desc,
		Assignee:    userID,
		UserID:      userID,
		TaskID:      task.ID,
		NodeDefKey:  task.NodeDefKey,
		InstanceID:  task.ProcInstanceID,
		ExecutionID: task.ExecutionID,
	}
	_, err = a.processAPI.AddHistoryTask(ctx, addHistoryTaskReq)
	if err != nil {
		return false, err
	}

	err = a.processAPI.SetAssignee(ctx, req.ProcessInstanceID, req.TaskID, model.HandleUserIDs[0])
	if err != nil {
		return false, err
	}

	// Add operation record
	model.HandleType = opDeliver
	model.HandleDesc = "转交"
	a.operationRecord.AddOperationRecord(ctx, flowInstanceEntity, task, model)

	return true, nil
}

func (a *abnormalTask) AdminGetTaskForm(ctx context.Context, req *AdminTaskReq) (*InstanceDetailModel, error) {
	instance, err := a.instanceRepo.GetEntityByProcessInstanceID(a.db, req.ProcessInstanceID)
	if err != nil {
		return nil, err
	}
	if instance == nil {
		return nil, error2.NewErrorWithString(error2.Internal, "Can not find flow instance data ")
	}
	flow, err := a.flowRepo.FindByID(a.db, instance.FlowID)
	if err != nil {
		return nil, err
	}
	if flow == nil {
		return nil, error2.NewErrorWithString(error2.Internal, "Can not find flow data ")
	}

	if !a.isAbnormalTask(req.ProcessInstanceID, req.TaskID) {
		return nil, error2.NewErrorWithString(error2.Internal, "is not abnormal task")
	}

	taskDetailModels := make([]*TaskDetailModel, 0)
	formSchema, _ := a.formAPI.GetFormSchema(ctx, instance.AppID, instance.FormID)
	taskDetailModels = append(taskDetailModels, &TaskDetailModel{
		// FormData:   formData,
		FormSchema: formSchema,
	})

	return &InstanceDetailModel{
		FlowName:            flow.Name,
		CanMsg:              flow.CanMsg == 1,
		CanViewStatusAndMsg: flow.CanViewStatusMsg == 1,
		TaskDetailModels:    taskDetailModels,
	}, nil
}

func (a *abnormalTask) check(ctx context.Context, processInstanceID string, taskID string) (*client.GetTasksResp, bool, error) {
	if !a.isUnhandleAbnormalTask(processInstanceID, taskID) {
		return nil, false, error2.NewErrorWithString(error2.Internal, "is not abnormal task")
	}
	resp, err := a.processAPI.GetTasks(ctx, client.GetTasksReq{
		InstanceID: []string{processInstanceID},
		TaskID:     []string{taskID},
	})
	if err != nil {
		return nil, false, err
	}
	if resp == nil || resp.Data == nil || len(resp.Data) < 1 {
		return nil, false, error2.NewErrorWithString(error2.Internal, "Can not find task ")
	}
	return resp, true, nil
}

func (a *abnormalTask) isUnhandleAbnormalTask(processInstanceID string, taskID string) bool {
	rs, err := a.abnormalTaskRepo.Find(a.db, map[string]interface{}{
		"process_instance_id": processInstanceID,
		"task_id":             taskID,
		"status":              0,
	})
	if err != nil {
		return false
	}
	return rs != nil && len(rs) > 0
}

func (a *abnormalTask) isAbnormalTask(processInstanceID string, taskID string) bool {
	rs, err := a.abnormalTaskRepo.Find(a.db, map[string]interface{}{
		"process_instance_id": processInstanceID,
		"task_id":             taskID,
	})
	if err != nil {
		return false
	}
	return rs != nil && len(rs) > 0
}

// AdminTaskReq admin task req params
type AdminTaskReq struct {
	ProcessInstanceID string
	TaskID            string
}
