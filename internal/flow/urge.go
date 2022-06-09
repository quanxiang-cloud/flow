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
	"github.com/quanxiang-cloud/flow/pkg/utils"
	"gorm.io/gorm"
)

const (
	noDealWith   = "noDealWith"
	autoDealWith = "autoDealWith"
	jump         = "jump"
)

const (
	firstEntry = "firstEntry"
	entry      = "entry"
	flowWorked = "flowWorked"
)

// urge type
const (
	DEADLINE = "DEADLINE"
	URGE     = "URGE"
)

// Urge info
type Urge interface {
	UrgingExecute(ctx context.Context, code string) error

	TaskUrge(ctx context.Context, req *TaskUrgeModel) error
}

type urge struct {
	db                     *gorm.DB
	processAPI             client.Process
	formAPI                client.Form
	dispatcherCallbackRepo models.DispatcherCallbackRepo
	dispatcherAPI          client.Dispatcher
	urgeRepo               models.UrgeRepo
	flowRepo               models.FlowRepo
	instanceRepo           models.InstanceRepo
	task                   Task
	flow                   Flow
	instance               Instance
	operationRecord        OperationRecord
}

// NewUrge init
func NewUrge(conf *config.Configs, opts ...options.Options) (Urge, error) {
	task, _ := NewTask(conf, opts...)
	flow, _ := NewFlow(conf, opts...)
	instance, _ := NewInstance(conf, opts...)
	operationRecord, _ := NewOperationRecord(conf, opts...)
	u := &urge{
		processAPI:             client.NewProcess(conf),
		formAPI:                client.NewForm(conf),
		dispatcherCallbackRepo: mysql.NewDispatcherCallbackRepo(),
		dispatcherAPI:          client.NewDispatcher(conf),
		urgeRepo:               mysql.NewUrgeRepo(),
		instanceRepo:           mysql.NewInstanceRepo(),
		flowRepo:               mysql.NewFlowRepo(),
		task:                   task,
		flow:                   flow,
		instance:               instance,
		operationRecord:        operationRecord,
	}
	for _, opt := range opts {
		opt(u)
	}
	return u, nil
}

// SetDB set db
func (u *urge) SetDB(db *gorm.DB) {
	u.db = db
}

func (u *urge) Transaction() *gorm.DB {
	return u.db.Begin()
}

func (u *urge) TaskUrge(ctx context.Context, req *TaskUrgeModel) error {
	userID := pkg.STDUserID(ctx)
	flowInstanceEntity, err := u.instanceRepo.GetEntityByProcessInstanceID(u.db, req.ProcessInstanceID)
	if err != nil {
		return err
	}
	if flowInstanceEntity == nil {
		return error2.NewErrorWithString(error2.Internal, "Can not find flow instance data ")
	}

	flowEntity, err := u.flowRepo.FindByID(u.db, flowInstanceEntity.FlowID)
	if err != nil {
		return err
	}
	if flowEntity == nil {
		return error2.NewErrorWithString(error2.Internal, "Can not find flow data ")
	}
	if !(flowEntity.CanUrge == 1 && isFlowOngoing(flowInstanceEntity.Status)) {
		return error2.NewErrorWithString(error2.Internal, "Can not urge this instance ")
	}

	tasks, err := u.processAPI.GetTasksByInstanceID(ctx, req.ProcessInstanceID)
	if err != nil {
		return err
	}
	if len(tasks) == 0 {
		return error2.NewErrorWithString(error2.Internal, "Can not urge this instance ")
	}

	for _, task := range tasks {
		// urges, err := u.urgeRepo.FindByTaskID(u.db, task.ID)
		// if err != nil {
		// 	return err
		// }

		// if len(urges) == 0 {
		err = u.urgeRepo.Create(u.db, &models.Urge{
			BaseModel: models.BaseModel{
				CreatorID: userID,
			},
			TaskID:            task.ID,
			ProcessInstanceID: req.ProcessInstanceID,
		})
		if err != nil {
			return err
		}

		// return error2.NewErrorWithString(code.CannotRepeatUrge, code.CodeTable[code.CannotRepeatUrge])
		// }
	}
	return nil
}

// DispatcherCallOtherInfo struct
type DispatcherCallOtherInfo struct {
	FlowInstanceID string `json:"flowInstanceId"`
	TaskID         string `json:"TaskID"`
	TaskDefKey     string `json:"taskDefKey"`
}

func (u *urge) UrgingExecute(ctx context.Context, code string) error {
	ctx = pkg.RPCCTXTransfer("", "")
	if code == "" {
		return error2.NewErrorWithString(error2.Internal, "code is nil")
	}
	d, err := u.dispatcherCallbackRepo.FindByID(u.db, code)
	if err != nil {
		return err
	}
	switch d.Type {
	case URGE:
		err = u.urge(ctx, *d)
	case DEADLINE:
		err = u.dealLine(ctx, *d)
	}
	return err
}

func (u *urge) urge(ctx context.Context, callback models.DispatcherCallback) error {
	tx := u.Transaction()
	otherInfo := &DispatcherCallOtherInfo{}
	err := json.Unmarshal([]byte(callback.OtherInfo), otherInfo)
	if err != nil {
		return err
	}
	resp, err := u.processAPI.GetHistoryTasks(ctx, client.GetTasksReq{
		TaskID: []string{otherInfo.TaskID},
	})
	if err != nil {
		return err
	}
	if resp != nil && len(resp.Data) > 0 {
		return nil
	}

	err = u.urgeRepo.Create(u.db, &models.Urge{
		BaseModel: models.BaseModel{
			ID: id2.GenID(),
		},
		TaskID:            otherInfo.TaskID,
		ProcessInstanceID: otherInfo.FlowInstanceID,
	})
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (u *urge) dealLine(ctx context.Context, callback models.DispatcherCallback) error {
	tx := u.Transaction()
	otherInfo := &DispatcherCallOtherInfo{}
	err := json.Unmarshal([]byte(callback.OtherInfo), otherInfo)
	if err != nil {
		return err
	}
	resp, err := u.processAPI.GetHistoryTasks(ctx, client.GetTasksReq{
		TaskID: []string{otherInfo.TaskID},
	})
	if err != nil {
		return err
	}
	if resp != nil && len(resp.Data) > 0 {
		return nil
	}

	instance, err := u.instanceRepo.FindByID(u.db, otherInfo.FlowInstanceID)
	if err != nil {
		return err
	}
	flow, err := u.flowRepo.FindByID(u.db, instance.FlowID)
	if err != nil {
		return err
	}

	s, err := convert.GetShapeByTaskDefKey(flow.BpmnText, otherInfo.TaskDefKey)
	if err != nil {
		return err
	}
	tmp := convert.GetValueFromBusinessData(*s, "basicConfig.timeRule")
	tempJSON, err := json.Marshal(tmp)
	if err != nil {
		return err
	}
	var timeRule convert.TaskTimeRuleModel
	err = json.Unmarshal(tempJSON, &timeRule)
	if err != nil {
		return err
	}

	timeOut := timeRule.WhenTimeout
	timeOutType := timeOut.Type
	if timeOutType == "" {
		timeOutType = noDealWith
	}

	taskCondition := client.GetTasksReq{
		InstanceID: []string{instance.ProcessInstanceID},
		NodeDefKey: otherInfo.TaskDefKey,
	}
	tasksResp, err := u.processAPI.GetTasks(ctx, taskCondition)
	if timeOutType == autoDealWith {
		params, err := u.flow.GetInstanceVariableValues(ctx, instance)
		formReq := client.FormDataConditionModel{
			AppID:   instance.AppID,
			TableID: instance.FormID,
			DataID:  instance.FormInstanceID,
		}
		formData, err := u.formAPI.GetFormData(ctx, formReq)
		params = utils.MergeMap(params, formData)
		if err != nil {
			return err
		}

		for _, task := range tasksResp.Data {
			u.task.AutoReviewTask(ctx, flow, instance, task, "", params)
		}
		processInstance, err := u.processAPI.GetInstanceByID(ctx, instance.ProcessInstanceID)
		if err != nil {
			return err
		}
		if processInstance != nil && processInstance.Status != Active {
			// 说明此时流程结束了
			dataMap := make(map[string]interface{})
			dataMap["modifier_id"] = pkg.STDUserID(ctx)
			dataMap["status"] = Agree
			err = u.instanceRepo.Update(u.db, instance.ID, dataMap)
			if err != nil {
				return err
			}
		}
	} else if timeOutType == jump {
		params, err := u.instance.GetInstanceVariableValues(ctx, instance)
		if err != nil {
			return err
		}

		comments := map[string]interface{}{
			"reviewResult": Agree,
			"reviewRemark": "",
		}
		u.processAPI.CompleteTaskToNode(ctx, instance.ProcessInstanceID, tasksResp.Data[0].ID, params, timeOut.Value, comments)

		model := &models.HandleTaskModel{
			HandleType: OpAutoSkip,
			HandleDesc: "该节点超时未处理，已跳转到指定节点",
		}

		// 增加操作日志
		task := &client.ProcessTask{
			NodeDefKey: s.ID,
			Name:       s.Data.NodeData.Name,
		}
		u.operationRecord.AddOperationRecord(ctx, instance, task, model)
	}
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// TaskUrgeModel task urge model
type TaskUrgeModel struct {
	ProcessInstanceID string `json:"processInstanceID" binding:"required"`
}
