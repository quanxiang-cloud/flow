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
	"fmt"
	"strings"

	"github.com/quanxiang-cloud/flow/internal"
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
	"github.com/quanxiang-cloud/flow/pkg/page"
	"github.com/quanxiang-cloud/flow/pkg/utils"
	"gorm.io/gorm"
)

// Flow service
type Flow interface {
	SaveFlow(ctx context.Context, req *models.Flow, userID string) (*models.Flow, error)
	Info(ID string) (*models.Flow, error)
	CopyFlow(ctx context.Context, ID string) (*models.Flow, error)
	DeleteFlow(ctx context.Context, flowID string, userID string) (bool, error)
	FlowList(ctx context.Context, req *QueryFlowReq) (*page.RespPage, error)
	CorrelationFlowList(ctx context.Context, req *CorrelationFlowReq) ([]*models.Flow, error)
	DeleteApp(ctx context.Context, req *DeleteAppReq) *DeleteAppResp
	UpdateFlowStatus(ctx context.Context, req *PublishProcessReq, usrID string) (*UpdateFlowStatusResp, error)
	GetNodes(ctx context.Context, ID string) ([]*models.NodeModel, error)
	GetShapeByProcessID(ctx context.Context, processID, nodeDefKey string) (*convert.ShapeModel, error)

	GetVariableList(ctx context.Context, ID string) ([]*models.Variables, error)
	SaveFlowVariable(ctx context.Context, req *SaveVariablesReq, userID string) (*models.Variables, error)
	DeleteFlowVariable(ctx context.Context, ID string) (bool, error)
	GetInstanceVariableValues(ctx context.Context, instance *models.Instance) (map[string]interface{}, error)
	GetFlowVariableValues(ctx context.Context, instance *models.Instance) (map[string]interface{}, error)

	RefreshRule(ctx context.Context) (bool, error)

	GetTaskHandleUserIDs(ctx context.Context, shape *convert.ShapeModel, flowInstanceEntity *models.Instance) ([]string, []string)
	GetTaskHandleUserIDs2(ctx context.Context, approvePersonsStr interface{}, flowInstanceEntity *models.Instance) ([]string, []string)
	GetTaskHandleUsers(ctx context.Context, shape *convert.ShapeModel, flowInstanceEntity *models.Instance) []*client.UserInfoResp
	GetTaskHandleUsers2(ctx context.Context, approvePersonsStr interface{}, flowInstanceEntity *models.Instance) []*client.UserInfoResp
	GetFlowKeyFields(ctx context.Context, flowIDs []string) (map[string][]string, error)

	checkCanCancel(ctx context.Context, flow *models.Flow, flowInstance *models.Instance, tasks []*client.ProcessTask) bool
	FormatFormValue(instance *models.Instance, data map[string]interface{}) map[string]interface{}
	FormatFormValue2(formulaFields interface{}, data map[string]interface{}) map[string]interface{}
	AppReplicationExport(ctx context.Context, req *AppReplicationExportReq) (string, error)
	AppReplicationImport(ctx context.Context, req *AppReplicationImportReq, userID string) (bool bool, err error)

	suspendApp(ctx context.Context, appID string) error
	recoveryApp(ctx context.Context, appID string) error
}

type flow struct {
	db                    *gorm.DB
	flowRepo              models.FlowRepo
	triggerRuleRepo       models.TriggerRuleRepo
	variablesRepo         models.VariablesRepo
	instanceVariablesRepo models.InstanceVariablesRepo
	formFieldRepo         models.FormFieldRepo
	identityAPI           client.Identity
	formAPI               client.Form
	processAPI            client.Process
	instanceRepo          models.InstanceRepo
	abnormalTaskRepo      models.AbnormalTaskRepo
	commentRepo           models.CommentRepo
	stepRepo              models.InstanceStepRepo
	operationRecordRepo   models.OperationRecordRepo
	urgeRepo              models.UrgeRepo
	instanceExecutionRepo models.InstanceExecutionRepo
	dispatcher            client.Dispatcher
	conf                  *config.Configs
}

const (
	SystemAudit        = "SYS_AUDIT_BOOL"
	SystemType         = "SYSTEM"
	SystemFieldType    = "boolean"
	SystemDefaultValue = "True"
)

// NewFlow init
func NewFlow(conf *config.Configs, opts ...options.Options) (Flow, error) {
	f := &flow{
		flowRepo:              mysql.NewFlowRepo(),
		triggerRuleRepo:       mysql.NewTriggerRuleRepo(),
		variablesRepo:         mysql.NewVariablesRepo(),
		instanceVariablesRepo: mysql.NewInstanceVariablesRepo(),
		formFieldRepo:         mysql.NewFormFieldRepo(),
		identityAPI:           client.NewIdentity(conf),
		formAPI:               client.NewForm(conf),
		processAPI:            client.NewProcess(conf),
		instanceRepo:          mysql.NewInstanceRepo(),
		abnormalTaskRepo:      mysql.NewAbnormalTaskRepo(),
		commentRepo:           mysql.NewCommentRepo(),
		stepRepo:              mysql.NewInstanceStepRepo(),
		operationRecordRepo:   mysql.NewOperationRecordRepo(),
		urgeRepo:              mysql.NewUrgeRepo(),
		instanceExecutionRepo: mysql.NewInstanceExecutionRepo(),
		dispatcher:            client.NewDispatcher(conf),
		conf:                  conf,
	}

	for _, opt := range opts {
		opt(f)
	}
	return f, nil
}

// SetDB set db
func (f *flow) SetDB(db *gorm.DB) {
	f.db = db
}

func (f *flow) GetShapeByProcessID(ctx context.Context, processID, nodeDefKey string) (*convert.ShapeModel, error) {
	flow, err := f.flowRepo.FindByProcessID(f.db, processID)
	if err != nil {
		return nil, err
	}

	shape, err := convert.GetShapeByTaskDefKey(flow.BpmnText, nodeDefKey)
	if err != nil {
		return nil, err
	}

	return shape, nil
}

// SaveFlow save flow
func (f *flow) SaveFlow(ctx context.Context, req *models.Flow, userID string) (*models.Flow, error) {
	if utils.HasEmoji(req.Name) {
		return nil, error2.NewErrorWithString(error2.Internal, "工作流名称不能包含emoji")
	}
	if len(req.BpmnText) > 0 {
		req.BpmnText = utils.UnicodeEmojiCode(req.BpmnText)
	}
	tx := f.db.Begin()
	if len(req.ID) > 0 {
		flow, err := f.flowRepo.FindByID(f.db, req.ID)
		if err != nil {
			return nil, err
		}
		if flow == nil {
			return nil, error2.NewErrorWithString(error2.Internal, "Cannot find flow ")
		}

		if len(flow.SourceID) > 0 || flow.Status == models.ENABLE {
			err := error2.NewErrorWithString(error2.Internal, "Cannot edit this flow ")
			return nil, err
		}
		req.ModifierID = userID
		req.Status = models.DISABLE

		err = f.flowRepo.UpdateFlow(tx, req)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		vb := &models.Variables{
			FlowID:       req.ID,
			Name:         SystemAudit,
			Type:         SystemType,
			Code:         "flowVar_" + strings.Replace(id2.GenID(), "-", "", -1),
			FieldType:    SystemFieldType,
			DefaultValue: SystemDefaultValue,
			Desc:         "系统初始化变量，不可修改",
		}

		condition := make(map[string]interface{})
		condition["name"] = SystemAudit
		condition["flow_id"] = req.ID
		variables, err := f.variablesRepo.FindVariables(f.db, condition)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if len(variables) > 0 {
			tx.Commit()
			return req, nil
		}
		vb.CreatorID = userID
		vb.CreateTime = time2.Now()
		vb.ModifyTime = time2.Now()
		err = f.variablesRepo.Create(tx, vb)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		tx.Commit()
		return req, nil
	}
	if req.TriggerMode == "FORM_DATA" {
		shape, err := convert.GetShapeByChartType(req.BpmnText, "formData")
		if err != nil {
			return nil, err
		}
		bd := shape.Data.BusinessData
		var form map[string]interface{}
		if v := bd["form"]; v != nil {
			form = v.(map[string]interface{})
		}
		if form == nil {
			return nil, nil
		}
		req.FormID = utils.Strval(form["value"])
	}

	// req.FormID = bd["form"].(map[string]interface{})["value"].(string)
	req.CreatorID = userID
	req.ModifierID = userID
	req.ModifyTime = time2.Now()
	req.Status = models.DISABLE

	req.AppStatus = mysql.AppActiveStatus
	err := f.flowRepo.Create(tx, req)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	vb := &models.Variables{
		FlowID:       req.ID,
		Name:         SystemAudit,
		Type:         SystemType,
		Code:         "flowVar_" + strings.Replace(id2.GenID(), "-", "", -1),
		FieldType:    SystemFieldType,
		DefaultValue: SystemDefaultValue,
		Desc:         "系统初始化变量，不可修改",
	}

	condition := make(map[string]interface{})
	condition["name"] = SystemAudit
	condition["flow_id"] = req.ID
	variables, err := f.variablesRepo.FindVariables(f.db, condition)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if len(variables) > 0 {
		tx.Commit()
		return req, nil
	}
	vb.CreatorID = userID
	vb.CreateTime = time2.Now()
	vb.ModifyTime = time2.Now()
	err = f.variablesRepo.Create(f.db, vb)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return req, nil
}

// DeleteFlow delete flow
func (f *flow) DeleteFlow(ctx context.Context, flowID string, userID string) (bool, error) {
	tx := f.db.Begin()
	var err error
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	flow, err := f.flowRepo.FindByID(tx, flowID)
	if err != nil {
		return false, err
	}
	if flow.Status == models.ENABLE {
		return false, error2.NewErrorWithString(error2.Internal, "流程启用中，不允许删除")
	}
	err = f.flowRepo.Delete(tx, flowID)
	if err != nil {
		return false, err
	}
	err = f.updateFlowHistory(tx, flowID, userID)
	if err != nil {
		return false, err
	}

	err = f.variablesRepo.DeleteByFlowID(tx, flowID)
	if err != nil {
		return false, err
	}
	// 如果是定时的将任务注册或删除调度器
	if flow.TriggerMode == convert.FormTime {
		code := fmt.Sprintf("flow:%s_%s", convert.CallbackOfCron, flow.ID)
		req := client.UpdateTaskStateReq{
			Code:  code,
			State: 2,
		}
		err1 := f.dispatcher.UpdateState(ctx, req)
		if err1 != nil {
			logger.Logger.Error("update dispatcher state err,", err1)
		}
	}

	return true, nil
}

func (f *flow) updateFlowHistory(tx *gorm.DB, flowID string, userID string) error {
	condition := map[string]interface{}{"source_id": flowID}
	needDeleteFlows, err := f.flowRepo.FindFlows(tx, condition)
	if err != nil {
		return err
	}

	if len(needDeleteFlows) > 0 {
		updateMap := map[string]interface{}{"status": models.DELETED, "modifier_id": userID}
		err = f.flowRepo.UpdateFlows(tx, condition, updateMap)
		if err != nil {
			return err
		}

		// Delete flow trigger rules
		flowIDs := make([]string, 0)
		for _, value := range needDeleteFlows {
			flowIDs = append(flowIDs, value.ID)
		}
		f.triggerRuleRepo.DeleteByFlowIDS(tx, flowIDs)
	}
	return nil
}

func (f *flow) Info(ID string) (*models.Flow, error) {
	if ID == "" || len(ID) == 0 {
		err := error2.NewErrorWithString(error2.Internal, "ID cannot be empty")
		return nil, err
	}
	flowData, err := f.flowRepo.FindByID(f.db, ID)
	if flowData == nil {
		return nil, error2.NewErrorWithString(error2.Internal, " cannot find flowData")
	}
	if len(flowData.BpmnText) > 0 {
		flowData.BpmnText = utils.UnicodeEmojiDecode(flowData.BpmnText)
		// 检查需要回退的节点是否被删除
		if flowData.CanCancel == 1 && flowData.CanCancelType == 3 && flowData.CanCancelNodes != "" {
			oriCancelNodes := strings.Split(flowData.CanCancelNodes, ",")
			flag := true
			effCancelNodes := make([]string, 0)
			for _, node := range oriCancelNodes {
				shape, err := convert.GetShapeByTaskDefKey(flowData.BpmnText, node)
				if err != nil {
					return nil, err
				}
				if shape != nil {
					flag = false
					effCancelNodes = append(effCancelNodes, node)
				}
			}
			if flag {
				flowData.CanCancelType = 0
			} else {
				flowData.CanCancelNodes = strings.Join(effCancelNodes, ",")
			}
		}
	}
	if len(flowData.Name) > 0 {
		flowData.Name = utils.UnicodeEmojiDecode(flowData.Name)
	}
	return flowData, err
}

func (f *flow) CopyFlow(ctx context.Context, ID string) (*models.Flow, error) {
	userID := pkg.STDUserID(ctx)
	flowData, err := f.flowRepo.FindByID(f.db, ID)
	if err != nil {
		return nil, err
	}
	if flowData == nil {
		err := error2.NewErrorWithString(error2.Internal, "Can not find flow ")
		return nil, err
	}

	model, err := convert.ToProcessModel(flowData.BpmnText)
	if err != nil {
		return nil, err
	}

	for _, node := range model.Shapes {
		node.Data.BusinessData = nil
	}
	newJSON, err := json.Marshal(model)
	if err != nil {
		return nil, err
	}

	flowData.SourceID = ""
	flowData.Name = flowData.Name + "_副本1"
	flowData.FormID = ""
	flowData.Status = "DISABLE"
	flowData.CanCancelNodes = ""
	flowData.InstanceName = ""
	flowData.KeyFields = ""
	flowData.ProcessID = ""
	flowData.BpmnText = string(newJSON)
	flowData.ID = ""
	flowData.CreatorID = userID
	flowData.ModifierID = userID
	flowData.ModifyTime = time2.Now()
	err = f.flowRepo.Create(f.db, flowData)

	return flowData, err
}

func (f *flow) FlowList(ctx context.Context, req *QueryFlowReq) (*page.RespPage, error) {
	condition := make(map[string]interface{})
	if req.AppID != "" {
		condition["app_id"] = req.AppID
	}
	if req.Name != "" {
		condition["name"] = req.Name
	}
	if req.Status != "" {
		condition["status"] = req.Status
	}
	if req.TriggerMode != "" {
		condition["trigger_mode"] = req.TriggerMode
	}
	condition["source_id"] = ""
	flows, count := f.flowRepo.FindPageFlows(f.db, condition, req.Page, req.Size)
	if len(flows) > 0 {
		for _, e := range flows {
			user, err := f.identityAPI.FindUserByID(ctx, e.ModifierID)
			if err != nil {
				return nil, err
			}
			if user != nil {
				e.ModifierName = user.UserName
			}
		}
	}
	pages := page.RespPage{
		Data:       flows,
		TotalCount: count,
	}

	return &pages, nil
}

func (f *flow) CorrelationFlowList(ctx context.Context, req *CorrelationFlowReq) ([]*models.Flow, error) {
	condition := make(map[string]interface{})
	if req.AppID != "" {
		condition["app_id"] = req.AppID
	}
	condition["form_id"] = req.FormID
	condition["source_id"] = ""
	return f.flowRepo.FindFlowList(f.db, condition)
}

func (f *flow) GetVariableList(ctx context.Context, ID string) ([]*models.Variables, error) {
	variables, err := f.variablesRepo.FindVariablesByFlowID(f.db, ID)
	if err != nil {
		return nil, err
	}
	for _, e := range variables {
		e.Name = utils.UnicodeEmojiDecode(e.Name)
	}
	return variables, err
}

func (f *flow) GetNodes(ctx context.Context, ID string) ([]*models.NodeModel, error) {
	flow, err := f.flowRepo.FindByID(f.db, ID)
	if err != nil {
		return nil, err
	}
	if flow == nil {
		err := error2.NewErrorWithString(error2.Internal, "flow is not exist")
		return nil, err
	}
	if flow.BpmnText == "" {
		err := error2.NewErrorWithString(error2.Internal, "flow chart is not exist")
		return nil, err
	}
	p, err := convert.ToProcessModel(flow.BpmnText)
	if err != nil {
		return nil, err
	}
	if len(p.Shapes) <= 0 {
		err := error2.NewErrorWithString(error2.Internal, "flow chart is not exist")
		return nil, err
	}
	nodes := make([]*models.NodeModel, 0)
	for _, s := range p.Shapes {
		if convert.FillIn == s.Type || convert.Approve == s.Type {
			node := &models.NodeModel{
				TaskDefKey: s.ID,
				TaskName:   s.Data.NodeData.Name,
			}
			nodes = append(nodes, node)
		}
	}
	return nodes, err
}

func (f *flow) SaveFlowVariable(ctx context.Context, req *SaveVariablesReq, userID string) (model *models.Variables, err error) {
	if utils.HasEmoji(req.Name) {
		return nil, error2.NewErrorWithString(error2.Internal, "变量名中不能包含emoji")
	}

	model = &models.Variables{
		BaseModel: models.BaseModel{
			ID: req.ID,
		},
		FlowID:       req.FlowID,
		Name:         utils.UnicodeEmojiCode(req.Name),
		Type:         req.Type,
		Code:         req.Code,
		FieldType:    req.FieldType,
		Format:       req.Format,
		DefaultValue: utils.Strval(req.DefaultValue),
		Desc:         req.Desc,
	}

	tx := f.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	if model.FlowID == "" {
		err = error2.NewErrorWithString(error2.Internal, "no FlowId")
		return nil, err
	}

	condition := make(map[string]interface{})
	condition["name"] = model.Name
	condition["flow_id"] = model.FlowID
	variables, err := f.variablesRepo.FindVariables(tx, condition)
	if err != nil {
		return nil, err
	}
	if model.ID == "" {
		if len(variables) > 0 {
			err = error2.NewErrorWithString(error2.Internal, "Process variable already exists ")
			return nil, err
		}
		model.CreatorID = userID
		model.CreateTime = time2.Now()
		model.ModifyTime = time2.Now()
		model.Code = "flowVar_" + strings.Replace(id2.GenID(), "-", "", -1)
		model.Type = "CUSTOM"
		err = f.variablesRepo.Create(tx, model)
		if err != nil {
			return nil, err
		}
	} else {
		if len(variables) > 1 {
			err = error2.NewErrorWithString(error2.Internal, "Process variable already exists ")
			return nil, err
		}
		variable, err := f.variablesRepo.FindByID(tx, model.ID)
		if err != nil {
			return nil, err
		}
		if variable == nil {
			err = error2.NewErrorWithString(error2.Internal, "Process variable is not exists ")
			return nil, err
		}
		if variable.Type == SystemType {
			err = error2.NewErrorWithString(error2.Internal, "Process system variable,user can not modify ")
			return nil, err
		}
		if len(variables) > 0 && variables[0].ID != model.ID {
			err = error2.NewErrorWithString(error2.Internal, "Process variable already exists ")
			return nil, err
		}
		model.ModifyTime = time2.Now()
		model.ModifierID = userID
		err = f.variablesRepo.Update(tx, model.ID, map[string]interface{}{
			"name":          model.Name,
			"field_type":    model.FieldType,
			"default_value": model.DefaultValue,
		})
		if err != nil {
			return nil, err
		}
	}
	tx.Commit()
	return model, nil
}

func (f *flow) getAppFlowIds(ctx context.Context, appID string) ([]string, error) {
	condition := map[string]interface{}{
		"app_id":    appID,
		"source_id": "",
	}

	flows, err := f.flowRepo.FindFlowList(f.db, condition)
	if err != nil {
		return nil, err
	}

	IDs := make([]string, 0)
	for _, value := range flows {
		IDs = append(IDs, value.ID)
	}

	return IDs, nil
}

func (f *flow) suspendApp(ctx context.Context, appID string) error {
	// Suspend flow and instance
	err := f.flowRepo.UpdateAppStatus(f.db, appID, mysql.AppSuspendStatus)
	if err != nil {
		return err
	}
	err = f.instanceRepo.UpdateAppStatus(f.db, appID, mysql.AppSuspendStatus)
	if err != nil {
		return err
	}

	defKeys, err := f.getAppFlowIds(ctx, appID)
	if err != nil {
		return err
	}
	err = f.processAPI.UpdateAppStatus(ctx, defKeys, "suspend")
	if err != nil {
		return err
	}

	return nil
}

func (f *flow) recoveryApp(ctx context.Context, appID string) error {
	// Suspend flow and instance
	err := f.flowRepo.UpdateAppStatus(f.db, appID, mysql.AppActiveStatus)
	if err != nil {
		return err
	}
	err = f.instanceRepo.UpdateAppStatus(f.db, appID, mysql.AppActiveStatus)
	if err != nil {
		return err
	}

	defKeys, err := f.getAppFlowIds(ctx, appID)
	if err != nil {
		return err
	}
	err = f.processAPI.UpdateAppStatus(ctx, defKeys, "reactive")
	if err != nil {
		return err
	}

	return nil
}

// DeleteApp delete app handle
func (f *flow) DeleteApp(ctx context.Context, req *DeleteAppReq) *DeleteAppResp {
	errs := make([]DeleteAppError, 0)

	switch req.Status {
	case mysql.AppPreDelete:
		err := f.suspendApp(ctx, req.AppID)
		if err != nil {
			e := DeleteAppError{
				DB:    "flow",
				Table: "flow,flow_instance",
				SQL:   "",
				Error: err,
			}
			errs = append(errs, e)
		}
	case mysql.AppDelete:
		// Delete flow data
		err := f.deleteFlowByAppID(ctx, req.AppID)
		if err != nil {
			e := DeleteAppError{
				DB:    "flow",
				Table: "flow,flow_instance",
				SQL:   "",
				Error: err,
			}
			errs = append(errs, e)
		}
	case mysql.AppRecovery:
		err := f.recoveryApp(ctx, req.AppID)
		if err != nil {
			e := DeleteAppError{
				DB:    "flow",
				Table: "flow,flow_instance",
				SQL:   "",
				Error: err,
			}
			errs = append(errs, e)
		}
	}

	return &DeleteAppResp{
		Errors: errs,
	}
}

func (f *flow) deleteFlowByAppID(ctx context.Context, appID string) error {
	tx := f.db.Begin()

	defKeys, err := f.getAppFlowIds(ctx, appID)
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	if err != nil {
		return err
	}
	err = f.processAPI.UpdateAppStatus(ctx, defKeys, "dump")
	if err != nil {
		return err
	}

	// 查询flowID和flowInstanceID
	condition := map[string]interface{}{"app_id": appID}
	flows, err := f.flowRepo.GetFlows(tx, condition)

	if err != nil {
		return err
	}
	if len(flows) == 0 {
		return nil
	}
	instances, err := f.instanceRepo.GetInstances(tx, condition)
	if err != nil {
		return err
	}
	flowIDs := make([]string, len(flows))
	processIDs := make([]string, 0)
	for index, flow := range flows {
		flowIDs[index] = flow.ID
		if flow.SourceID == "" {
			processIDs = append(processIDs, flow.ID)
		}
		// 如果是定时的将任务注册或删除调度器
		if flow.TriggerMode == convert.FormTime {
			code := fmt.Sprintf("flow:%s_%s", convert.CallbackOfCron, flow.ID)
			req := client.UpdateTaskStateReq{
				Code:  code,
				State: 2,
			}
			err1 := f.dispatcher.UpdateState(ctx, req)
			if err1 != nil {
				logger.Logger.Error("update dispatcher state err,", err1)
			}
		}

	}
	instanceIDs := make([]string, len(instances))
	processInstanceIDs := make([]string, len(instances))
	for index, instance := range instances {
		instanceIDs[index] = instance.ID
		processInstanceIDs[index] = instance.ProcessInstanceID
	}
	// 删除flow表
	if err = f.flowRepo.DeleteByIDs(tx, flowIDs); err != nil {
		return err
	}
	// 删除flowInstance
	if err = f.instanceRepo.DeleteByFlowIDs(tx, flowIDs); err != nil {
		return err
	}
	// 删除flow_abnormal_task
	if err = f.abnormalTaskRepo.DeleteByInstanceIDs(tx, instanceIDs); err != nil {
		return err
	}
	// 删除flow_comment
	if err = f.commentRepo.DeleteByInstanceIDs(tx, instanceIDs); err != nil {
		return err
	}
	// 删除flow_form_field
	if err = f.formFieldRepo.DeleteByFlowIDs(tx, flowIDs); err != nil {
		return err
	}
	// 删除flow_instance_step
	if err = f.stepRepo.DeleteByProcessInstanceIDs(tx, processInstanceIDs); err != nil {
		return err
	}
	// 删除flow_instance_variable
	if err = f.instanceVariablesRepo.DeleteByProcessInstanceIDs(tx, processInstanceIDs); err != nil {
		return err
	}
	// 删除flow_operation_record
	if err = f.operationRecordRepo.DeleteByProcessInstanceIDs(tx, processInstanceIDs); err != nil {
		return err
	}
	// 删除flow_trigger_rule
	if err = f.triggerRuleRepo.DeleteByFlowIDS(tx, flowIDs); err != nil {
		return err
	}
	// 删除flow_urge
	if err = f.urgeRepo.DeleteByProcessInstanceIDs(tx, processInstanceIDs); err != nil {
		return err
	}
	// 删除flow_variables
	if err = f.variablesRepo.DeleteByFlowIDs(tx, flowIDs); err != nil {
		return err
	}
	// 删除instance_execution
	if err = f.instanceExecutionRepo.DeleteByProcessInstanceIDs(tx, processInstanceIDs); err != nil {
		return err
	}
	// 调process接口
	if err = f.processAPI.UpdateAppStatus(ctx, processIDs, mysql.AppDeleteStatus); err != nil {
		return err
	}
	return nil
}

func (f *flow) DeleteFlowVariable(ctx context.Context, ID string) (bool, error) {
	tx := f.db.Begin()
	err := f.variablesRepo.Delete(tx, ID)
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (f *flow) RefreshRule(ctx context.Context) (bool bool, err error) {
	tx := f.db.Begin()
	condition := make(map[string]interface{})
	condition["source_id"] = ""
	condition["status"] = models.ENABLE

	flows, err := f.flowRepo.FindFlows(tx, condition)
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	if err != nil {
		return false, err
	}
	if len(flows) > 0 {
		for _, s := range flows {
			b, err := convert.GetBusinessDataByType(s.BpmnText, convert.FormData)
			if err != nil {
				return false, err
			}
			rule, err := json.Marshal(b)
			if err != nil {
				return false, err
			}
			ruleOne, err := f.triggerRuleRepo.FindByFormIDAndDFlowID(f.db, s.FormID, s.ID)
			if err != nil {
				return false, err
			}
			if ruleOne == nil {
				ftr := &models.TriggerRule{
					Rule:   string(rule),
					FlowID: s.ID,
					FormID: s.FormID,
				}
				f.triggerRuleRepo.Create(tx, ftr)
			}

		}

	}
	tx.Commit()
	return true, err
}

type UpdateFlowStatusResp struct {
	Flag        bool
	TriggerMode string
}

func (f *flow) UpdateFlowStatus(ctx context.Context, req *PublishProcessReq, usrID string) (resp *UpdateFlowStatusResp, err error) {
	tx := f.db.Begin()

	fl, err := f.flowRepo.FindByID(tx, req.ID)
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	flowStatusResp := &UpdateFlowStatusResp{}
	flowStatusResp.TriggerMode = fl.TriggerMode
	if err != nil {
		flowStatusResp.Flag = false
		return flowStatusResp, err
	}
	if fl == nil {
		err = error2.NewErrorWithString(error2.Internal, "Process is not exists ")
		flowStatusResp.Flag = false
		return flowStatusResp, err
	}
	fl.ModifierID = usrID
	fl.ModifyTime = time2.Now()
	fl.Status = req.Status

	if models.ENABLE == req.Status {
		if fl.BpmnText == "" {
			err = error2.NewErrorWithString(error2.Internal, "Process chart is nil ")
			flowStatusResp.Flag = false
			return flowStatusResp, err
		}
		if fl.TriggerMode == "FORM_DATA" {
			// 找到form节点对应的数据
			s, err := convert.GetShapeByChartType(fl.BpmnText, convert.FormData)
			if err != nil {
				flowStatusResp.Flag = false
				return flowStatusResp, err
			}
			if err = checkChartJSON(s); err != nil {
				flowStatusResp.Flag = false
				return flowStatusResp, err
			}
			formID := s.Data.BusinessData["form"].(map[string]interface{})["value"].(string)
			fl.FormID = formID

			// add new trigger rule
			rule, err := json.Marshal(s.Data.BusinessData)
			if err != nil {
				flowStatusResp.Flag = false
				return flowStatusResp, err
			}
			ruleOne, err := f.triggerRuleRepo.FindByFormIDAndDFlowID(f.db, fl.FormID, fl.ID)
			if err != nil {
				flowStatusResp.Flag = false
				return flowStatusResp, err
			}
			if ruleOne == nil {
				ftr := &models.TriggerRule{
					Rule:   string(rule),
					FlowID: fl.ID,
					FormID: fl.FormID,
				}
				if err = f.triggerRuleRepo.Create(tx, ftr); err != nil {
					flowStatusResp.Flag = false
					return flowStatusResp, err
				}
			} else {
				err := f.triggerRuleRepo.Update(tx, ruleOne.ID, map[string]interface{}{
					"rule": string(rule),
				})
				if err != nil {
					flowStatusResp.Flag = false
					return flowStatusResp, err
				}
			}

		}

		// deploy
		processID, formulaFields, err := f.deploy(ctx, fl)
		if err != nil {
			flowStatusResp.Flag = false
			return flowStatusResp, err
		}
		if processID == "" {
			flowStatusResp.Flag = false
			return flowStatusResp, error2.NewErrorWithString(error2.Internal, "Deployment failed")
		}

		fl.ProcessID = processID
		err = f.flowRepo.UpdateFlow(tx, fl)
		if err != nil {
			flowStatusResp.Flag = false
			return flowStatusResp, err
		}
		////// add new record
		////fl.SourceID = fl.ID
		//originalID := fl.ID
		////fl.ID = id2.GenID()
		////if err = f.flowRepo.Create(f.db, fl); err != nil {
		////	flowStatusResp.Flag = false
		////	return flowStatusResp, err
		////}
		//
		//// sync variable list
		//variables, err := f.variablesRepo.FindVariables(f.db, map[string]interface{}{
		//	"flow_id": originalID,
		//})
		//if err != nil {
		//	flowStatusResp.Flag = false
		//	return flowStatusResp, err
		//}
		//if len(variables) > 0 {
		//	for _, v := range variables {
		//		v.FlowID = fl.ID
		//		v.ID = id2.GenID()
		//		if err = f.variablesRepo.Create(f.db, v); err != nil {
		//			flowStatusResp.Flag = false
		//			return flowStatusResp, err
		//		}
		//	}
		//}
		// 如果是定时的将任务注册或删除调度器
		if fl.TriggerMode == convert.FormTime {
			code := fmt.Sprintf("flow:%s_%s", convert.CallbackOfCron, req.ID)
			req := client.TaskPostReq{
				Code:    code,
				Title:   "cron_" + fl.Name,
				Type:    2,
				TimeBar: fl.Cron,
				State:   1,
			}
			err1 := f.dispatcher.TakePost(ctx, req)
			if err1 != nil {
				logger.Logger.Error("register dispatcher err,", err1)
				code := fmt.Sprintf("flow:%s_%s", convert.CallbackOfCron, fl.ID)
				dispReq := client.UpdateTaskStateReq{
					Code:  code,
					State: 1,
				}
				err1 := f.dispatcher.UpdateState(ctx, dispReq)
				if err1 != nil {
					logger.Logger.Error("unable flow update dispatcher state err,", err1)
				}
			}
		}

		// save form field with path
		f.formFieldRepo.DeleteByFlowID(f.db, fl.ID)
		if formulaFields != nil {
			formulaFieldsMap := utils.ChangeObjectToMap(formulaFields)
			if formulaFieldsMap != nil {
				for key, value := range formulaFieldsMap {
					formField := &models.FormField{
						FlowID:         fl.ID,
						FormID:         fl.FormID,
						FieldName:      key,
						FieldValuePath: value.(string),
					}
					f.formFieldRepo.Create(f.db, formField)
				}
			}
		}

	} else {
		err := f.flowRepo.UpdateFlow(tx, fl)
		if err != nil {
			flowStatusResp.Flag = false
			return flowStatusResp, err
		}
		f.updateFlowHistory(tx, fl.ID, usrID)
		// 如果是定时的将任务注册或删除调度器
		if fl.TriggerMode == convert.FormTime {
			code := fmt.Sprintf("flow:%s_%s", convert.CallbackOfCron, fl.ID)
			dispReq := client.UpdateTaskStateReq{
				Code:  code,
				State: 2,
			}
			err1 := f.dispatcher.UpdateState(ctx, dispReq)
			if err1 != nil {
				logger.Logger.Error("unable flow update dispatcher state err,", err1)
			}
		}
	}
	tx.Commit()
	flowStatusResp.Flag = true
	return flowStatusResp, err
}

func checkChartJSON(s *convert.ShapeModel) error {
	data := s.Data.BusinessData
	triggerWay := data["triggerWay"].([]interface{})
	triggerCondition := data["triggerCondition"].(map[string]interface{})
	expr := triggerCondition["expr"].([]interface{})
	if len(triggerWay) <= 0 && len(expr) <= 0 {
		err := error2.NewErrorWithString(error2.Internal, "not set trigger rule ")
		return err
	}
	return nil
}

func (f *flow) deploy(ctx context.Context, fl *models.Flow) (string, interface{}, error) {
	nodes, formulaFields, err := convert.GenProcessNodes(fl.BpmnText)
	if err != nil {
		return "", nil, err
	}
	if nodes == nil || len(nodes) < 1 {
		return "", nil, error2.NewErrorWithString(error2.Internal, "The drawing has no nodes ")
	}
	reqMap := &map[string]interface{}{
		"id":    fl.ID,
		"name":  fl.Name,
		"nodes": nodes,
	}
	model, err := json.Marshal(reqMap)
	if err != nil {
		return "", nil, err
	}
	resp, err := f.processAPI.DeployProcess(ctx, client.DeployReq{
		Model:     string(model),
		CreatorID: pkg.STDUserID(ctx),
	})
	if err != nil {
		return "", nil, err
	}

	return resp.ID, formulaFields, nil
}

func (f *flow) FormatFormValue(instance *models.Instance, data map[string]interface{}) map[string]interface{} {
	formFields, _ := f.formFieldRepo.FindByFlowID(f.db, instance.FlowID)
	if len(formFields) > 0 && data != nil {
		for _, field := range formFields {
			paths := strings.Split(field.FieldValuePath, ".")
			data[field.FieldName] = f.formAPI.GetValueByFieldFormat(paths, data[field.FieldName])
		}
	}
	return data
}

func (f *flow) FormatFormValue2(formulaFieldsObj interface{}, data map[string]interface{}) map[string]interface{} {
	formulaFields := utils.ChangeObjectToStringMap(formulaFieldsObj)
	if len(formulaFields) > 0 {
		for fieldName, fieldValuePath := range formulaFields {
			paths := strings.Split(fieldValuePath, ".")
			data[fieldName] = f.formAPI.GetValueByFieldFormat(paths, data[fieldName])
		}
	}
	return data
}

func (f *flow) GetTaskHandleUserIDs(ctx context.Context, shape *convert.ShapeModel, flowInstanceEntity *models.Instance) ([]string, []string) {
	bd := shape.Data.BusinessData
	if bd == nil {
		return nil, nil
	}

	return f.GetTaskHandleUserIDs2(ctx, convert.GetValueFromBusinessData(*shape, "basicConfig.approvePersons"), flowInstanceEntity)
}

func (f *flow) GetTaskHandleUserIDs2(ctx context.Context, approvePersonsStr interface{}, flowInstanceEntity *models.Instance) ([]string, []string) {

	approvePersonsJSON, err := json.Marshal(approvePersonsStr)
	if err != nil {
		logger.Logger.Error(err)
		return nil, nil
	}
	var approvePersons convert.ApprovePersonsModel
	err = json.Unmarshal(approvePersonsJSON, &approvePersons)
	if err != nil {
		logger.Logger.Error(err)
		return nil, nil
	}

	// 人员'person' | 表单字段'field' | 岗位'position' | 上级领导'superior' | 部门负责人'leadOfDepartment' | 发起人processInitiator
	handleUserIds := make([]string, 0) // all handle users
	assigneeList := make([]string, 0)  // dynamic handle users
	if convert.EmailTypeOfField == approvePersons.Type || convert.EmailTypeOfMultipleField == approvePersons.Type {
		if len(approvePersons.Fields) > 0 {
			condition := client.FormDataConditionModel{
				AppID:   flowInstanceEntity.AppID,
				TableID: flowInstanceEntity.FormID,
				DataID:  flowInstanceEntity.FormInstanceID,
			}
			formData, err := f.formAPI.GetFormData(ctx, condition)
			if err != nil {
				return nil, nil
			}

			for _, field := range approvePersons.Fields {
				if formData != nil && formData[field] != nil {
					value, ok := formData[field].(string)
					if ok {
						assigneeList = append(assigneeList, value)
					} else {
						fieldValue := utils.ChangeObjectToMapList(formData[field])
						for _, mm := range fieldValue {
							assigneeList = append(assigneeList, utils.GetAsString(mm["value"]))
						}
					}
				}
			}
		}
	} else if convert.EmailTypeOfSuperior == approvePersons.Type {
		userID, err := f.identityAPI.GetSuperior(ctx, flowInstanceEntity.ApplyUserID)
		if err != nil {
			return nil, nil
		}

		if len(userID) > 0 {
			assigneeList = append(assigneeList, userID)
		}
	} else if convert.EmailTypeOfLeadOfDepartment == approvePersons.Type {
		userID, err := f.identityAPI.GetLeadOfDepartment(ctx, flowInstanceEntity.ApplyUserID)
		if err != nil {
			return nil, nil
		}

		if len(userID) > 0 {
			assigneeList = append(assigneeList, userID)
		}
	} else if convert.EmailTypeOfProcessInitiator == approvePersons.Type {
		assigneeList = append(assigneeList, flowInstanceEntity.ApplyUserID)
	} else if convert.EmailTypeOfPerson == approvePersons.Type {
		if len(approvePersons.Users) > 0 {
			for _, user := range approvePersons.Users {
				handleUserIds = append(handleUserIds, user["id"].(string))
			}

			handleUserIds, _ = f.identityAPI.ValidateUserIDs(ctx, handleUserIds)
		}

		if len(approvePersons.Departments) > 0 {
			for _, depart := range approvePersons.Departments {
				users, _ := f.identityAPI.FindUsersByDepartID(ctx, depart["id"].(string))
				if len(users) > 0 {
					handleUserIds = append(handleUserIds, users...)
				}
			}
		}

	}

	handleUserIds = append(handleUserIds, assigneeList...)

	return utils.RemoveReplicaSliceString(handleUserIds), utils.RemoveReplicaSliceString(assigneeList)
}

func (f *flow) GetTaskHandleUsers(ctx context.Context, shape *convert.ShapeModel, flowInstanceEntity *models.Instance) []*client.UserInfoResp {
	handleUserIds, _ := f.GetTaskHandleUserIDs(ctx, shape, flowInstanceEntity)
	return f.getTaskHandleUsersByIDs(ctx, handleUserIds)
}

func (f *flow) GetTaskHandleUsers2(ctx context.Context, approvePersonsStr interface{}, flowInstanceEntity *models.Instance) []*client.UserInfoResp {
	handleUserIds, _ := f.GetTaskHandleUserIDs2(ctx, approvePersonsStr, flowInstanceEntity)
	return f.getTaskHandleUsersByIDs(ctx, handleUserIds)
}

func (f *flow) getTaskHandleUsersByIDs(ctx context.Context, handleUserIds []string) []*client.UserInfoResp {
	userMap, err := f.identityAPI.FindUsersByIDs(ctx, handleUserIds)
	if err != nil {
		return nil
	}
	handleUsers := make([]*client.UserInfoResp, 0)
	for _, user := range userMap {
		handleUsers = append(handleUsers, user)
	}
	return handleUsers
}

func (f *flow) GetFlowKeyFields(ctx context.Context, flowIDs []string) (map[string][]string, error) {
	list, err := f.flowRepo.FindByIDs(f.db, flowIDs)
	if err != nil {
		return nil, err
	}

	dataMap := make(map[string][]string, 0)
	if len(list) > 0 {
		for _, flowEntity := range list {
			if len(flowEntity.KeyFields) > 0 {
				dataMap[flowEntity.ID] = strings.Split(flowEntity.KeyFields, ",")
			}
		}
	}
	return dataMap, nil
}

func (f *flow) checkCanCancel(ctx context.Context, flow *models.Flow, flowInstance *models.Instance, tasks []*client.ProcessTask) bool {
	if flow.CanCancel != 1 {
		return false
	}

	switch flow.CanCancelType {
	case 1:
		return flowInstance.Status == Review
	case 2:
		return true
	case 3:
		if len(flow.CanCancelNodes) > 0 {
			canCancelNodes := strings.Split(flow.CanCancelNodes, ",")
			currentNodes := make([]string, 0)
			for _, task := range tasks {
				currentNodes = append(currentNodes, task.NodeDefKey)
			}
			return utils.IntersectString(canCancelNodes, currentNodes)
		}

	}
	return false
}

func (f *flow) GetInstanceVariableValues(ctx context.Context, instance *models.Instance) (map[string]interface{}, error) {
	// default variables
	valueMap := make(map[string]interface{})
	valueMap["flowVar_instanceCreatorName"] = instance.ApplyUserName
	valueMap["flowVar_instanceCreateTime"] = utils.ChangeISO8601ToBjTime(instance.CreateTime)
	valueMap["flowVar_instanceStatus"] = GetStatusName(instance.Status)

	// self define variables
	variablesEntities, err := f.instanceVariablesRepo.FindVariablesByProcessInstanceID(f.db, instance.ProcessInstanceID)
	if err != nil {
		return nil, err
	}
	for _, value := range variablesEntities {
		valueMap[value.Code] = utils.FormatValue(value.Value, value.FieldType)
	}

	// event variables
	orders := make([]client.QueryOrder, 0)
	order := client.QueryOrder{
		OrderType: internal.OrderTypeAsc,
		Column:    "end_time",
	}
	orders = append(orders, order)
	req := client.GetTasksReq{
		InstanceID: []string{instance.ProcessInstanceID},
		Status:     Completed,
		Page:       0,
		Limit:      10000,
		Order:      orders,
	}
	tasks, err := f.processAPI.GetHistoryTasks(ctx, req)

	if err != nil {
		return nil, err
	}
	if tasks != nil && tasks.Data != nil {
		assignees := make([]string, 0)
		for _, task := range tasks.Data {
			assignees = append(assignees, task.Assignee)
		}
		users, _ := f.identityAPI.FindUsersByIDs(ctx, assignees)
		for _, task := range tasks.Data {
			valueMap["$"+task.NodeDefKey+".handleUserId"] = task.Assignee
			if users != nil && task.Assignee != "" {
				valueMap["$"+task.NodeDefKey+".handleUserName"] = users[task.Assignee].UserName
			}
			valueMap["$"+task.NodeDefKey+".handleTime"] = task.EndTime

			if len(task.Comments) > 0 {
				var comments map[string]interface{}
				err := json.Unmarshal([]byte(task.Comments), &comments)
				if err != nil {
					return nil, err
				}
				valueMap["$"+task.NodeDefKey+".reviewResult"] = comments["reviewResult"]
				valueMap["$"+task.NodeDefKey+".reviewRemark"] = comments["reviewRemark"]
			}
		}
	}

	// form data variables
	formReq := client.FormDataConditionModel{
		AppID:   instance.AppID,
		TableID: instance.FormID,
		DataID:  instance.FormInstanceID,
	}
	formData, err := f.formAPI.GetFormData(ctx, formReq)
	if err == nil && formData != nil {
		for k, v := range formData {
			valueMap[k] = v
		}
	}

	return valueMap, nil
}

func (f *flow) GetFlowVariableValues(ctx context.Context, instance *models.Instance) (map[string]interface{}, error) {
	valueMap := make(map[string]interface{})

	variablesEntities, err := f.variablesRepo.FindVariablesByFlowID(f.db, instance.FlowID)
	if err != nil {
		return nil, err
	}
	for _, value := range variablesEntities {
		valueMap[value.Code] = utils.FormatValue(value.DefaultValue, value.FieldType)
	}

	valueMap["flowVar_instanceCreatorName"] = instance.ApplyUserName
	valueMap["flowVar_instanceCreateTime"] = utils.ChangeISO8601ToBjTime(instance.CreateTime)
	valueMap["flowVar_instanceStatus"] = GetStatusName(instance.Status)

	return valueMap, nil
}

func (f *flow) AppReplicationExport(ctx context.Context, req *AppReplicationExportReq) (string, error) {
	if req.AppID == "" {
		return "", error2.NewErrorWithString(error2.ErrParams, "AppID is nil")
	}
	condition := make(map[string]interface{})
	condition["source_id"] = ""
	condition["app_id"] = req.AppID
	flows, err := f.flowRepo.FindFlows(f.db, condition)
	if err != nil {
		return "", err
	}
	for _, flow := range flows {
		variables, err := f.variablesRepo.FindVariablesByFlowID(f.db, flow.ID)
		if err != nil {
			return "", err
		}
		flow.Variables = variables
	}
	marshal, err := json.Marshal(flows)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}

func (f *flow) AppImportFlow(ctx context.Context, req *models.Flow) (*models.Flow, error) {
	// req.CreatorID = userID
	// req.ModifierID = userID
	req.ID = id2.GenID()
	req.ModifyTime = time2.Now()
	req.AppStatus = mysql.AppActiveStatus
	err := f.flowRepo.Create2(f.db, req)

	return req, err
}

func (f *flow) AppReplicationImport(ctx context.Context, req *AppReplicationImportReq, userID string) (bool bool, err error) {
	tx := f.db.Begin()

	if req.AppID == "" {
		return false, error2.NewErrorWithString(error2.ErrParams, "AppID is nil")
	}
	if req.Flows == "" {
		return false, error2.NewErrorWithString(error2.ErrParams, "flows is nil")
	}

	var flows []models.Flow
	err = json.Unmarshal([]byte(req.Flows), &flows)
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	if err != nil {
		return false, err
	}
	logger.Logger.Info(req)
	for _, flow := range flows {
		// replace variables
		flow.BpmnText = strings.Replace(flow.BpmnText, flow.AppID, req.AppID, -1)
		// for k, v := range req.FormID {
		// 	flow.BpmnText = strings.Replace(flow.BpmnText, k, v, -1)
		// }

		s, err := convert.GetShapeByChartType(flow.BpmnText, convert.FormData)
		if err != nil {
			return false, err
		}
		formID := s.Data.BusinessData["form"].(map[string]interface{})["value"].(string)

		// save flow
		flow.ID = id2.GenID()
		flow.ModifyTime = time2.Now()
		flow.AppStatus = mysql.AppActiveStatus
		flow.AppID = req.AppID
		flow.FormID = formID
		flow.CreatorID = userID
		flow.ModifierID = userID
		err = f.flowRepo.Create2(f.db, &flow)
		if err != nil {
			return false, err
		}

		// deploy flow
		if flow.Status == models.ENABLE {
			// 自动发布流程
			processID, formulaFields, err := f.deploy(ctx, &flow)
			if err != nil {
				continue
			}
			if processID == "" {
				continue
			}
			flow.ProcessID = processID
			err = f.flowRepo.UpdateFlow(tx, &flow)
			if err != nil {
				return false, err
			}

			flow.SourceID = flow.ID
			flow.ID = id2.GenID()
			if err := f.flowRepo.Create2(tx, &flow); err != nil {
				continue
			}

			if err = checkChartJSON(s); err != nil {
				continue
			}
			// add new trigger rule
			rule, err := json.Marshal(s.Data.BusinessData)
			if err != nil {
				continue
			}
			ruleOne, err := f.triggerRuleRepo.FindByFormIDAndDFlowID(f.db, flow.FormID, flow.ID)
			if err != nil {
				continue
			}
			if ruleOne == nil {
				ruleStr := string(rule)
				ftr := &models.TriggerRule{
					Rule:   ruleStr,
					FlowID: flow.ID,
					FormID: formID,
				}
				ftr.CreatorID = userID
				ftr.CreateTime = time2.Now()
				if err = f.triggerRuleRepo.Create(tx, ftr); err != nil {
					continue
				}
			} else {
				err := f.triggerRuleRepo.Update(tx, ruleOne.ID, map[string]interface{}{
					"rule": string(rule),
				})
				if err != nil {
					continue
				}
			}
			// sync variable list
			for _, variable := range flow.Variables {
				if variable.Type == "CUSTOM" {
					variable.FlowID = flow.ID
					variable.CreatorID = userID
					variable.CreateTime = time2.Now()
					if err := f.variablesRepo.Create(tx, variable); err != nil {
						continue
					}
				}
			}

			// save form field with path
			formulaFieldsMap := utils.ChangeObjectToMap(formulaFields)
			if formulaFieldsMap != nil {
				for key, value := range formulaFieldsMap {
					formField := &models.FormField{
						FlowID:         flow.ID,
						FormID:         formID,
						FieldName:      key,
						FieldValuePath: value.(string),
					}
					formField.CreatorID = userID
					formField.CreateTime = time2.Now()
					f.formFieldRepo.Create(tx, formField)
				}
			}
		}

	}
	return true, nil
}

// DeleteAppResp delete app resp
type DeleteAppResp struct {
	Errors []DeleteAppError `json:"errors"`
}

// DeleteAppError delete app error
type DeleteAppError struct {
	DB    string      `json:"db"`
	Table string      `json:"table"`
	SQL   string      `json:"sql"`
	Error interface{} `json:"err"`
}

// SaveVariablesReq save variables req
type SaveVariablesReq struct {
	ID           string      `json:"id"`
	FlowID       string      `json:"flowId"`
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	Code         string      `json:"code"`
	FieldType    string      `json:"fieldType"`
	Format       string      `json:"format"`
	DefaultValue interface{} `json:"defaultValue"`
	Desc         string      `json:"desc"`
}

// PublishProcessReq deploy process
type PublishProcessReq struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// QueryFlowReq flow request
type QueryFlowReq struct {
	Page        int    `json:"page"`
	Size        int    `json:"size"`
	TriggerMode string `json:"triggerMode"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	AppID       string `json:"appId"`
}

// CorrelationFlowReq req
type CorrelationFlowReq struct {
	AppID  string `json:"appID"`
	FormID string `json:"formID"`
}

// DeleteAppReq delete app req
type DeleteAppReq struct {
	AppID  string `json:"appID" binding:"required"`
	Status string `json:"status" binding:"required"` // preDelete、delete，recovery
}

// AppReplicationExportReq 应用复制导出入参
type AppReplicationExportReq struct {
	AppID string `json:"appID"`
}

// AppReplicationImportReq 应用复制导出入参
type AppReplicationImportReq struct {
	AppID string `json:"appID"`
	// FormID map[string]string `json:"formID"`
	Flows string `json:"flows"`
}
