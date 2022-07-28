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
	"encoding/json"
	"github.com/quanxiang-cloud/flow/internal/convert"
	"github.com/quanxiang-cloud/flow/internal/models"
	"github.com/quanxiang-cloud/flow/internal/models/mysql"
	"github.com/quanxiang-cloud/flow/internal/server/options"
	"github.com/quanxiang-cloud/flow/pkg"
	"github.com/quanxiang-cloud/flow/pkg/client"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/misc/error2"
	"github.com/quanxiang-cloud/flow/pkg/misc/id2"
	"github.com/quanxiang-cloud/flow/pkg/misc/logger"
	"github.com/quanxiang-cloud/flow/pkg/misc/time2"
	"github.com/quanxiang-cloud/flow/pkg/utils"
	"gorm.io/gorm"
	"strings"
	"time"
)

const (
	// OpSubmit Operation
	opSubmit = "SUBMIT" // 提交
	// OpReSubmit Operation
	opReSubmit = "RE_SUBMIT" // 再次提交
	// OpCancel Operation
	opCancel = "CANCEL" // 撤回
	// OpAbandon Operation
	opAbandon = "ABANDON" // 作废
	// OpAbend Operation
	opAbend = "ABEND" // 异常结束

	// OpAgree Operation
	opAgree = "AGREE" // 通过
	// OpRefuse Operation
	opRefuse = "REFUSE" // 拒绝
	// OpFillIn Operation
	opFillIn = "FILL_IN" // 填写

	// OpDeliver Operation
	opDeliver = "DELIVER" // 转交
	// OpStepBack Operation
	opStepBack = "STEP_BACK" // 回退
	// OpSendBack Operation
	opSendBack = "SEND_BACK" // 打回重填
	// OpCC Operation
	opCC = "CC" // 抄送
	// OpRead Operation
	opRead = "READ" // 邀请阅示
	// OpHandleCC Operation
	opHandleCC = "HANDLE_CC" // 处理抄送
	// OpHandleRead Operation
	opHandleRead = "HANDLE_READ" // 处理阅示

	// OpAddSign Operation
	opAddSign = "ADD_SIGN" // 加签
	// OpAutoReview Operation
	opAutoReview = "AUTO_REVIEW" // 自动审批
	// OpAutoSkip Operation
	OpAutoSkip = "AUTO_SKIP" // 跳过
	// OpAutoCC Operation
	opAutoCC = "AUTO_CC" // 自动抄送
)

// OperationRecord service
type OperationRecord interface {
	// AddTaskComment(ctx context.Context, processInstanceID string, task *client.ProcessTask, handleTaskModel *models.HandleTaskModel, userID string) error
	GetAgreeUserIds(ctx context.Context, processInstanceID string) ([]string, error)
	AddOperationRecord(ctx context.Context, instance *models.Instance, task *client.ProcessTask, handleTaskModel *models.HandleTaskModel) error
	AddOperationRecords(ctx context.Context, instance *models.Instance, task *client.ProcessTask, handleTaskModel *models.HandleTaskModel, records []*models.OperationRecord) error
	ConvertOperationRecord(ctx context.Context, instance *models.Instance, task *client.ProcessTask, handleTaskModel *models.HandleTaskModel) []*models.OperationRecord
	UpdateByNodeInstanceID(ctx context.Context, flow *models.Flow, task *client.ProcessTask) error
}

type operationRecord struct {
	db                  *gorm.DB
	operationRecordRepo models.OperationRecordRepo
	instanceStepRepo    models.InstanceStepRepo
	instanceStep        InstanceStep
	processAPI          client.Process
	instanceRepo        models.InstanceRepo
	flowRepo            models.FlowRepo
	flow                Flow
}

// NewOperationRecord init
func NewOperationRecord(conf *config.Configs, opts ...options.Options) (OperationRecord, error) {
	instanceStep, err := NewInstanceStep(conf, opts...)
	if err != nil {
		return nil, err
	}
	flow, err := NewFlow(conf, opts...)
	if err != nil {
		return nil, err
	}
	or := &operationRecord{
		operationRecordRepo: mysql.NewOperationRecordRepo(),
		instanceStepRepo:    mysql.NewInstanceStepRepo(),
		instanceStep:        instanceStep,
		processAPI:          client.NewProcess(conf),
		instanceRepo:        mysql.NewInstanceRepo(),
		flowRepo:            mysql.NewFlowRepo(),
		flow:                flow,
	}

	for _, opt := range opts {
		opt(or)
	}
	return or, nil
}

// SetDB set db
func (or *operationRecord) SetDB(db *gorm.DB) {
	or.db = db
}

func (or *operationRecord) ConvertOperationRecord(ctx context.Context, instance *models.Instance, task *client.ProcessTask, handleTaskModel *models.HandleTaskModel) []*models.OperationRecord {
	var handleUserID string
	if handleTaskModel.HandleType == opAutoReview {
		handleUserID = handleTaskModel.AutoReviewUserID
	} else {
		handleUserID = pkg.STDUserID(ctx)
	}

	records := make([]*models.OperationRecord, 0)

	taskID := ""
	taskName := ""
	taskDefKey := ""
	if task != nil {
		taskID = task.ID
		taskName = task.Name
		taskDefKey = task.NodeDefKey
	}

	operation := &models.OperationRecord{
		ProcessInstanceID: instance.ProcessInstanceID,
		HandleUserID:      handleUserID,
		HandleType:        handleTaskModel.HandleType,
		HandleDesc:        handleTaskModel.HandleDesc,
		Remark:            handleTaskModel.Remark,
		Status:            Completed,
		TaskID:            taskID,
		TaskName:          taskName,
		TaskDefKey:        taskDefKey,
		CorrelationData:   strings.Join(handleTaskModel.HandleUserIDs, ","),
		RelNodeDefKey:     handleTaskModel.RelNodeDefKey,
		BaseModel: models.BaseModel{
			ID:         id2.GenID(),
			CreatorID:  handleUserID,
			CreateTime: time2.Now(),
			ModifyTime: time2.Now(),
		},
	}
	records = append(records, operation)

	return records
}

func (or *operationRecord) AddOperationRecords(ctx context.Context, instance *models.Instance, task *client.ProcessTask, handleTaskModel *models.HandleTaskModel, records []*models.OperationRecord) error {

	ID := or.processTaskStep(ctx, instance, task, handleTaskModel)

	var filter bool
	if handleTaskModel.HandleType == opCC || handleTaskModel.HandleType == opRead {
		derivationRecords, err := or.operationRecordRepo.FindRecords(or.db, instance.ProcessInstanceID, ID, []string{handleTaskModel.HandleType}, true)
		if err != nil {
			return err
		}
		if derivationRecords != nil && len(derivationRecords) > 0 {
			filter = true
		}
	}

	for _, record := range records {
		if filter && (record.HandleType == opCC || record.HandleType == opRead) {
			continue
		}
		record.InstanceStepID = ID
		or.operationRecordRepo.Create(or.db, record)
	}
	return nil

}

func (or *operationRecord) AddOperationRecord(ctx context.Context, instance *models.Instance, task *client.ProcessTask, handleTaskModel *models.HandleTaskModel) error {
	records := or.ConvertOperationRecord(ctx, instance, task, handleTaskModel)
	return or.AddOperationRecords(ctx, instance, task, handleTaskModel, records)
}

func (or *operationRecord) processTaskStep(ctx context.Context, instance *models.Instance, task *client.ProcessTask, handleTaskModel *models.HandleTaskModel) string {
	var ID string
	switch handleTaskModel.HandleType {
	case opSubmit:
		ID = or.processSubmitTaskStep(ctx, instance.ProcessInstanceID)
	case OpAutoSkip:
		ID = or.processAutoSkipTaskStep(ctx, instance, task)
	case opAgree:
		fallthrough
	case opRefuse:
		fallthrough
	case opFillIn:
		fallthrough
	case opStepBack:
		fallthrough
	case opCancel:
		fallthrough
	case opDeliver:
		fallthrough
	case opCC:
		fallthrough
	case opRead:
		fallthrough
	case opAutoReview:
		fallthrough
	case opAutoCC:
		fallthrough
	case opAbandon:
		fallthrough
	case opAbend:
		fallthrough
	case opSendBack:
		fallthrough
	case opAddSign:
		ID = or.processBaseTaskStep(ctx, instance, task, handleTaskModel)
	case opReSubmit:
		ID = or.processReSubmitTaskStep(ctx, instance.ProcessInstanceID)
	}
	return ID
}

func (or *operationRecord) processSubmitTaskStep(ctx context.Context, processInstanceID string) string {
	step := &models.InstanceStep{
		ProcessInstanceID: processInstanceID,
		TaskType:          convert.START,
		Status:            opSubmit,
		BaseModel: models.BaseModel{
			ID:         id2.GenID(),
			CreatorID:  pkg.STDUserID(ctx),
			CreateTime: time2.Now(),
			ModifyTime: time2.Now(),
		},
	}
	or.instanceStepRepo.Create(or.db, step)
	return step.ID
}

func (or *operationRecord) processAutoSkipTaskStep(ctx context.Context, instance *models.Instance, task *client.ProcessTask) string {
	logger.Logger.Info("processAutoSkipTaskStep")
	logger.Logger.Info(task)
	f, _ := or.flowRepo.FindByID(or.db, instance.FlowID)
	// 获取当前节点类型
	// 节点相关数据
	shape, _ := convert.GetShapeByTaskDefKey(f.BpmnText, task.NodeDefKey)
	if shape == nil {
		return ""
	}
	logger.Logger.Info(shape)
	basicConfig := convert.GetTaskBasicConfigModel(shape)

	step := &models.InstanceStep{
		ProcessInstanceID: instance.ProcessInstanceID,
		TaskType:          convert.GetCurrentNodeType(shape.Type, basicConfig.MultiplePersonWay),
		Status:            OpAutoSkip,
		TaskName:          task.Name,
		BaseModel: models.BaseModel{
			ID:         id2.GenID(),
			CreatorID:  pkg.STDUserID(ctx),
			CreateTime: time2.UnixToISO8601(time.Now().Unix() + 5),
			ModifyTime: time2.UnixToISO8601(time.Now().Unix() + 5),
		},
	}
	or.instanceStepRepo.Create(or.db, step)
	return step.ID
}

func (or *operationRecord) processReSubmitTaskStep(ctx context.Context, processInstanceID string) string {
	step := &models.InstanceStep{
		ProcessInstanceID: processInstanceID,
		TaskType:          convert.START,
		Status:            opSubmit,
		BaseModel: models.BaseModel{
			ID:         id2.GenID(),
			CreatorID:  pkg.STDUserID(ctx),
			CreateTime: time2.Now(),
			ModifyTime: time2.Now(),
		},
	}
	or.instanceStepRepo.Create(or.db, step)
	return step.ID
}

func (or *operationRecord) processBaseTaskStep(ctx context.Context, instance *models.Instance, task *client.ProcessTask, handleTaskModel *models.HandleTaskModel) string {
	if task == nil {
		return ""
	}

	processInstanceID := instance.ProcessInstanceID
	userID := pkg.STDUserID(ctx)

	f, _ := or.flowRepo.FindByID(or.db, instance.FlowID)
	// 获取当前节点类型
	// 节点相关数据
	shape, _ := convert.GetShapeByTaskDefKey(f.BpmnText, task.NodeDefKey)
	if shape == nil {
		return ""
	}
	basicConfig := convert.GetTaskBasicConfigModel(shape)

	var nodeInstanceID string
	if basicConfig.MultiplePersonWay == "and" {
		nodeInstanceID = task.NodeInstancePid
	} else {
		nodeInstanceID = task.NodeInstanceID
	}

	// 判断节点所有是否已经完成
	taskCondition := client.GetTasksReq{
		InstanceID: []string{processInstanceID},
		NodeDefKey: task.NodeDefKey,
		Desc:       []string{convert.ReviewTask, convert.WriteTask},
	}

	// 判断节点所有是否已经完成
	var nodeCompleted bool
	tasksResp, _ := or.processAPI.GetTasks(ctx, taskCondition)
	if len(tasksResp.Data) == 0 {
		nodeCompleted = true
	}

	step, _ := or.instanceStep.GetFlowInstanceStep(ctx, processInstanceID, nodeInstanceID)

	if step == nil {
		stepID := id2.GenID()
		var status string
		if nodeCompleted && handleTaskModel.HandleType != opAddSign {
			status = handleTaskModel.HandleType
		} else {
			status = InReview
		}

		var taskHandleUserIDStr string
		if task.TaskType == "TempModel" {
			relRecord, _ := or.operationRecordRepo.FindRecordByRelDefKey(or.db, processInstanceID, task.NodeDefKey)
			taskHandleUserIDStr = relRecord.CorrelationData
		} else {
			taskHandleUserIds, _ := or.flow.GetTaskHandleUserIDs(ctx, shape, instance)
			taskHandleUserIDStr = strings.Join(taskHandleUserIds, ",")
		}

		if handleTaskModel.HandleType == opDeliver {
			taskHandleUserIDStr = strings.Replace(taskHandleUserIDStr, userID, handleTaskModel.HandleUserIDs[0], 1)
		}

		// 该节点的类型
		var timeStr string
		if handleTaskModel.HandleType == opAutoReview {
			timeStr = time2.UnixToISO8601(time.Now().Unix() + 5)
		} else {
			timeStr = time2.Now()
		}
		step = &models.InstanceStep{
			ProcessInstanceID: processInstanceID,
			TaskType:          convert.GetCurrentNodeType(shape.Type, basicConfig.MultiplePersonWay),
			TaskDefKey:        task.NodeDefKey,
			TaskName:          task.Name,
			Status:            status,
			HandleUserIDs:     taskHandleUserIDStr,
			NodeInstanceID:    nodeInstanceID,
			BaseModel: models.BaseModel{
				ID:         stepID,
				CreatorID:  userID,
				CreateTime: timeStr,
				ModifyTime: timeStr,
			},
		}
		or.instanceStepRepo.Create(or.db, step)
	} else {

		if handleTaskModel.HandleType == opDeliver {
			step.HandleUserIDs = strings.Replace(step.HandleUserIDs, userID, handleTaskModel.HandleUserIDs[0], 1)
		}
		if !derivation(handleTaskModel.HandleType) {
			updateMap := map[string]interface{}{
				"status":          handleTaskModel.HandleType,
				"modify_time":     time2.Now(),
				"modifier_id":     userID,
				"handle_user_ids": step.HandleUserIDs,
			}
			or.instanceStepRepo.Update(or.db, step.ID, updateMap)
		}

	}
	return step.ID
}

// UpdateByNodeInstanceId func
func (or *operationRecord) UpdateByNodeInstanceID(ctx context.Context, flow *models.Flow, task *client.ProcessTask) error {
	shape, _ := convert.GetShapeByTaskDefKey(flow.BpmnText, task.NodeDefKey)
	if shape == nil {
		return nil
	}
	basicConfig := convert.GetTaskBasicConfigModel(shape)

	var nodeInstanceID string
	if basicConfig.MultiplePersonWay == "and" {
		nodeInstanceID = task.NodeInstancePid
	} else {
		nodeInstanceID = task.NodeInstanceID
	}
	updateMap := map[string]interface{}{
		"status":      Agree,
		"modify_time": time2.Now(),
	}
	err := or.instanceStepRepo.UpdateByNodeInstanceID(or.db, nodeInstanceID, updateMap)
	return err
}

func (or *operationRecord) AddTaskComment(ctx context.Context, processInstanceID string, task *client.ProcessTask, handleTaskModel *models.HandleTaskModel, userID string) error {
	switch handleTaskModel.HandleType {
	case opSubmit:
		or.submit(ctx, processInstanceID, handleTaskModel, userID)
	case opAgree:
		fallthrough
	case opRefuse:
		fallthrough
	case opFillIn:
		fallthrough
	case opStepBack:
		fallthrough
	case opCancel:
		fallthrough
	case opDeliver:
		fallthrough
	case opCC:
		fallthrough
	case opRead:
		fallthrough
	case opAutoReview:
		fallthrough
	case OpAutoSkip:
		or.agree(ctx, processInstanceID, task, handleTaskModel, userID)
	case opReSubmit:
		fallthrough
	case opAutoCC:
		or.resubmit(ctx, processInstanceID, handleTaskModel, userID)
	case opAbandon:
		fallthrough
	case opAbend:
		or.end(ctx, processInstanceID, handleTaskModel, userID)
	case opSendBack:
		or.sendBack(ctx, processInstanceID, task, handleTaskModel, userID)
	case opAddSign:
		// TODO
		break
	default:
		return error2.NewErrorWithString(error2.Internal, "no handletype match")
	}

	return nil
}

func (or *operationRecord) submit(ctx context.Context, processInstanceID string, handleTaskModel *models.HandleTaskModel, userID string) error {
	step := models.InstanceStep{
		ProcessInstanceID: processInstanceID,
		TaskType:          convert.Start,
		Status:            opSubmit,
		BaseModel: models.BaseModel{
			CreatorID:  userID,
			CreateTime: time2.Now(),
			ModifyTime: time2.Now(),
		},
	}
	or.instanceStepRepo.Create(or.db, &step)

	handleTaskModel.Remark = utils.UnicodeEmojiCode(handleTaskModel.Remark)
	handleTaskModel.FormData = nil
	correlationData, _ := json.Marshal(handleTaskModel)
	operation := &models.OperationRecord{
		HandleDesc:        handleTaskModel.HandleDesc,
		HandleType:        handleTaskModel.HandleType,
		ProcessInstanceID: processInstanceID,
		Remark:            handleTaskModel.Remark,
		CorrelationData:   string(correlationData),
		InstanceStepID:    step.ID,
		BaseModel: models.BaseModel{
			CreatorID:  userID,
			CreateTime: time2.Now(),
			ModifyTime: time2.Now(),
		},
	}
	or.operationRecordRepo.Create(or.db, operation)

	// or.linkNextStep(ctx, processInstanceID, handleTaskModel, userID)

	return nil
}

func (or *operationRecord) resubmit(ctx context.Context, processInstanceID string, handleTaskModel *models.HandleTaskModel, userID string) error {
	step, err := or.instanceStep.GetReStartFlowInstanceStep(ctx, processInstanceID)
	if err != nil {
		return err
	}
	if step == nil {
		return nil
	}

	step.Status = handleTaskModel.HandleType
	step.BaseModel.ModifierID = userID
	step.BaseModel.ModifyTime = time2.Now()
	or.instanceStepRepo.Create(or.db, step)

	handleTaskModel.Remark = utils.UnicodeEmojiCode(handleTaskModel.Remark)
	handleTaskModel.FormData = nil
	correlationData, _ := json.Marshal(handleTaskModel)
	operation := &models.OperationRecord{
		HandleDesc:        handleTaskModel.HandleDesc,
		HandleType:        handleTaskModel.HandleType,
		ProcessInstanceID: processInstanceID,
		Remark:            handleTaskModel.Remark,
		CorrelationData:   string(correlationData),
		InstanceStepID:    step.ID,
		BaseModel: models.BaseModel{
			CreatorID:  userID,
			CreateTime: time2.Now(),
			ModifyTime: time2.Now(),
		},
	}
	or.operationRecordRepo.Create(or.db, operation)

	or.linkNextStep(ctx, processInstanceID, handleTaskModel, userID)

	return nil
}

func (or *operationRecord) end(ctx context.Context, processInstanceID string, handleTaskModel *models.HandleTaskModel, userID string) error {
	steps, err := or.instanceStep.GetFlowInstanceSteps(ctx, processInstanceID)
	if err != nil {
		return err
	}

	if len(steps) == 0 {
		logger.Logger.Error("steps is null:[]"+processInstanceID, pkg.STDRequestID(ctx))
		return nil
	}

	for _, step := range steps {
		updateMap := map[string]interface{}{
			"status":      handleTaskModel.HandleType,
			"modify_time": time2.Now(),
			"modifier_id": userID,
		}
		or.instanceStepRepo.Update(or.db, step.ID, updateMap)

		handleTaskModel.Remark = utils.UnicodeEmojiCode(handleTaskModel.Remark)
		handleTaskModel.FormData = nil
		correlationData, _ := json.Marshal(handleTaskModel)
		operation := &models.OperationRecord{
			HandleDesc:        handleTaskModel.HandleDesc,
			HandleType:        handleTaskModel.HandleType,
			ProcessInstanceID: processInstanceID,
			Remark:            handleTaskModel.Remark,
			CorrelationData:   string(correlationData),
			InstanceStepID:    step.ID,
			BaseModel: models.BaseModel{
				CreatorID:  userID,
				CreateTime: time2.Now(),
				ModifyTime: time2.Now(),
			},
		}
		or.operationRecordRepo.Create(or.db, operation)
	}

	or.linkNextStep(ctx, processInstanceID, handleTaskModel, userID)

	return nil
}

func (or *operationRecord) sendBack(ctx context.Context, processInstanceID string, task *client.ProcessTask, handleTaskModel *models.HandleTaskModel, userID string) error {
	step, err := or.instanceStep.GetFlowInstanceStep(ctx, processInstanceID, task.NodeDefKey)
	if err != nil {
		return err
	}
	if step == nil {
		return nil
	}

	updateMap := map[string]interface{}{
		"status":      handleTaskModel.HandleType,
		"modify_time": time2.Now(),
		"modifier_id": userID,
	}
	or.instanceStepRepo.Update(or.db, step.ID, updateMap)

	handleTaskModel.Remark = utils.UnicodeEmojiCode(handleTaskModel.Remark)
	handleTaskModel.FormData = nil
	correlationData, _ := json.Marshal(handleTaskModel)
	operation := &models.OperationRecord{
		HandleDesc:        handleTaskModel.HandleDesc,
		HandleType:        handleTaskModel.HandleType,
		ProcessInstanceID: processInstanceID,
		Remark:            handleTaskModel.Remark,
		CorrelationData:   string(correlationData),
		InstanceStepID:    step.ID,
		BaseModel: models.BaseModel{
			CreatorID:  userID,
			CreateTime: time2.Now(),
			ModifyTime: time2.Now(),
		},
	}
	or.operationRecordRepo.Create(or.db, operation)

	or.linkNextStep(ctx, processInstanceID, handleTaskModel, userID)

	return nil
}

func (or *operationRecord) agree(ctx context.Context, processInstanceID string, task *client.ProcessTask, handleTaskModel *models.HandleTaskModel, userID string) error {

	step, err := or.instanceStep.GetFlowInstanceStep(ctx, processInstanceID, task.NodeDefKey)
	if err != nil {
		return err
	}
	if step == nil {
		flowInstance, _ := or.instanceRepo.GetEntityByProcessInstanceID(or.db, processInstanceID)
		flow, _ := or.flowRepo.FindByID(or.db, flowInstance.FlowID)
		// 获取当前节点类型
		// 节点相关数据
		shape, _ := convert.GetShapeByTaskDefKey(flow.BpmnText, task.NodeDefKey)
		if shape == nil {
			return nil
		}
		basicConfig := convert.GetTaskBasicConfigModel(shape)

		// 该节点的类型
		currentNodeType := convert.GetCurrentNodeType(shape.Type, basicConfig.MultiplePersonWay)
		step = &models.InstanceStep{
			ProcessInstanceID: processInstanceID,
			TaskType:          currentNodeType,
			TaskDefKey:        task.NodeDefKey,
			TaskName:          task.Name,
			Status:            Review,
			BaseModel: models.BaseModel{
				CreatorID:  userID,
				CreateTime: time2.Now(),
				ModifyTime: time2.Now(),
			},
		}
		or.instanceStepRepo.Create(or.db, step)

	} else {
		// linkNextNode := false // 是否需要添加下一节点
		taskCondition := client.GetTasksReq{
			InstanceID: []string{processInstanceID},
			NodeDefKey: task.NodeDefKey,
		}
		tasksResp, err := or.processAPI.GetTasks(ctx, taskCondition)
		if err != nil {
			return err
		}
		if len(tasksResp.Data) == 0 { // 当前节点任务已完成
			updateMap := map[string]interface{}{
				"status":      handleTaskModel.HandleType,
				"modify_time": time2.Now(),
				"modifier_id": userID,
			}
			or.instanceStepRepo.Update(or.db, step.ID, updateMap)
			// linkNextNode = true
		} else if step.Status == Review {
			// if handleTaskModel.HandleType == opDeliver {
			// 	step.HandleUserIDs = utils.StringJoin(handleTaskModel.HandleUserIds)
			// }

			updateMap := map[string]interface{}{
				"status":      handleTaskModel.HandleType,
				"modify_time": time2.Now(),
				"modifier_id": userID,
			}
			or.instanceStepRepo.Update(or.db, step.ID, updateMap)
		}
	}

	handleTaskModel.Remark = utils.UnicodeEmojiCode(handleTaskModel.Remark)
	handleTaskModel.FormData = nil
	correlationData, _ := json.Marshal(handleTaskModel)
	operation := &models.OperationRecord{
		HandleDesc:        handleTaskModel.HandleDesc,
		HandleType:        handleTaskModel.HandleType,
		ProcessInstanceID: processInstanceID,
		Remark:            handleTaskModel.Remark,
		CorrelationData:   string(correlationData),
		InstanceStepID:    step.ID,
		BaseModel: models.BaseModel{
			CreatorID:  userID,
			CreateTime: time2.Now(),
			ModifyTime: time2.Now(),
		},
	}
	or.operationRecordRepo.Create(or.db, operation)

	// if linkNextNode {
	// 	or.linkNextStep(ctx, processInstanceID, handleTaskModel, userID)
	// }

	return nil
}

func (or *operationRecord) linkNextStep(ctx context.Context, processInstanceID string, handleTaskModel *models.HandleTaskModel, userID string) error {
	taskCondition := client.GetTasksReq{
		InstanceID: []string{processInstanceID},
	}
	tasksResp, err := or.processAPI.GetTasks(ctx, taskCondition)
	if err != nil {
		return err
	}

	if len(tasksResp.Data) == 0 && handleTaskModel.HandleType != SendBack {
		or.linkEndStep(ctx, processInstanceID, userID)
		return nil
	}

	step := &models.InstanceStep{}
	if handleTaskModel.HandleType == SendBack {
		step.ProcessInstanceID = processInstanceID
		step.TaskType = convert.Start
		step.Status = opReSubmit
		step.BaseModel.CreatorID = userID
		step.BaseModel.CreateTime = time2.Now()
		step.BaseModel.ModifyTime = time2.Now()
	} else if handleTaskModel.HandleType == opCancel {
		or.linkEndStep(ctx, processInstanceID, userID)
	} else {
		flowInstance, _ := or.instanceRepo.GetEntityByProcessInstanceID(or.db, processInstanceID)
		flow, _ := or.flowRepo.FindByID(or.db, flowInstance.FlowID)
		currentTask := tasksResp.Data[0]

		// 获取当前节点类型
		// 节点相关数据
		shape, _ := convert.GetShapeByTaskDefKey(flow.BpmnText, currentTask.NodeDefKey)
		if shape == nil {
			return nil
		}
		basicConfig := convert.GetTaskBasicConfigModel(shape)

		// 该节点定义的可以审批的人
		taskHandleUserIds, _ := or.flow.GetTaskHandleUserIDs(ctx, shape, flowInstance)

		// 该节点的类型
		currentNodeType := convert.GetCurrentNodeType(shape.Type, basicConfig.MultiplePersonWay)
		var singleInstanceNode bool
		singleInstanceNode = currentNodeType != convert.AndApproval && currentNodeType != convert.AndFillIn

		step.ProcessInstanceID = processInstanceID
		if singleInstanceNode {
			step.TaskID = tasksResp.Data[0].ID
		}
		step.TaskType = currentNodeType
		step.TaskDefKey = currentTask.NodeDefKey
		step.TaskName = currentTask.Name
		step.HandleUserIDs = utils.StringJoin(taskHandleUserIds)
		step.Status = Review
		step.BaseModel.CreatorID = userID
		step.BaseModel.CreateTime = time2.Now()
		step.BaseModel.ModifyTime = time2.Now()
	}
	or.instanceStepRepo.Create(or.db, step)

	return nil
}

func (or *operationRecord) linkEndStep(ctx context.Context, processInstanceID string, userID string) error {
	step := &models.InstanceStep{
		ProcessInstanceID: processInstanceID,
		TaskType:          convert.End,
		TaskName:          "结束",
		Status:            convert.End,
		BaseModel: models.BaseModel{
			CreatorID:  userID,
			CreateTime: time2.Now(),
			ModifyTime: time2.Now(),
		},
	}
	or.instanceStepRepo.Create(or.db, step)

	return nil
}

func (or *operationRecord) GetAgreeUserIds(ctx context.Context, processInstanceID string) ([]string, error) {
	return or.operationRecordRepo.GetAgreeUserIds(or.db, processInstanceID)
}
