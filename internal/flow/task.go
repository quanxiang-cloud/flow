package flow

import (
	"context"
	"encoding/json"
	"git.internal.yunify.com/qxp/misc/id2"
	"git.internal.yunify.com/qxp/misc/logger"
	"git.internal.yunify.com/qxp/misc/time2"
	"github.com/quanxiang-cloud/flow/internal/convert"
	"github.com/quanxiang-cloud/flow/internal/models"
	"github.com/quanxiang-cloud/flow/internal/models/mysql"
	"github.com/quanxiang-cloud/flow/internal/server/options"
	"github.com/quanxiang-cloud/flow/pkg/client"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/utils"
	"gorm.io/gorm"
)

// Task service
type Task interface {
	TaskInitHandle(ctx context.Context, flowEntity *models.Flow, flowInstanceEntity *models.Instance, task *client.ProcessTask, currentUserID string) error
	noUserHandle(ctx context.Context, flowInstanceEntity *models.Instance, task *client.ProcessTask, taskBasicConfigModel *convert.TaskBasicConfigModel) error
	TaskUrging(ctx context.Context, rule convert.TaskTimeRuleModel, nodeID string, instance models.Instance, TaskID string) error

	// TaskCheck(ctx context.Context, flowEntity *models.Flow, flowInstanceEntity *models.Instance, currentUserID string) error
	FilterCanEditFormData(ctx context.Context, flowEntity *models.Flow, flowInstanceEntity *models.Instance,
		taskDefKey string, formData map[string]interface{}) map[string]interface{}
	FilterCanReadFormData(ctx context.Context, flowInstanceEntity *models.Instance, fieldPermissionObj interface{}, formData interface{}) map[string]interface{}
	GetCurrentNodes(ctx context.Context, processInstanceID string) ([]models.NodeModel, error)
	AutoReviewTask(ctx context.Context, flowEntity *models.Flow, flowInstanceEntity *models.Instance, task *client.ProcessTask, userID string, params map[string]interface{}) error
	PermissionConvertWriteRead(permission int8) (bool, bool)
}

type task struct {
	db                     *gorm.DB
	operationRecordRepo    models.OperationRecordRepo
	formAPI                client.Form
	identityAPI            client.Identity
	processAPI             client.Process
	flow                   Flow
	abnormalTaskRepo       models.AbnormalTaskRepo
	operationRecord        OperationRecord
	dispatcherCallbackRepo models.DispatcherCallbackRepo
	dispatcherAPI          client.Dispatcher
	flowInstanceRepo       models.InstanceRepo
}

// NewTask init
func NewTask(conf *config.Configs, opts ...options.Options) (Task, error) {
	flow, err := NewFlow(conf, opts...)
	if err != nil {
		return nil, err
	}

	operationRecord, err := NewOperationRecord(conf, opts...)
	if err != nil {
		return nil, err
	}
	t := &task{
		operationRecordRepo:    mysql.NewOperationRecordRepo(),
		formAPI:                client.NewForm(conf),
		identityAPI:            client.NewIdentity(conf),
		processAPI:             client.NewProcess(conf),
		flow:                   flow,
		abnormalTaskRepo:       mysql.NewAbnormalTaskRepo(),
		operationRecord:        operationRecord,
		dispatcherCallbackRepo: mysql.NewDispatcherCallbackRepo(),
		dispatcherAPI:          client.NewDispatcher(conf),
		flowInstanceRepo:       mysql.NewInstanceRepo(),
	}

	for _, opt := range opts {
		opt(t)
	}
	return t, nil
}

// SetDB set db
func (t *task) SetDB(db *gorm.DB) {
	t.db = db
}

// TaskInitHandle auto review, no user handle, time limit handle and urge handle
func (t *task) TaskInitHandle(ctx context.Context, flowEntity *models.Flow, flowInstanceEntity *models.Instance, task *client.ProcessTask, currentUserID string) error {
	taskAutoHandle := false
	shape, err := convert.GetShapeByTaskDefKey(flowEntity.BpmnText, task.NodeDefKey)
	if err != nil {
		return err
	}
	taskBasicConfigModel := convert.GetTaskBasicConfigModel(shape)
	if taskBasicConfigModel == nil {
		return nil
	}

	handleUserIds, _ := t.flow.GetTaskHandleUserIDs(ctx, shape, flowInstanceEntity)
	if len(handleUserIds) == 0 { // 无审批人时
		logger.Logger.Info("事件：加载handleUserIds，handleUserIds=" + utils.ChangeStringArrayToString(handleUserIds))
		t.noUserHandle(ctx, flowInstanceEntity, task, taskBasicConfigModel)
	} else { // 有人审批时,审批人为发起人时origin，审批人与上一节点审批人相同时parent，审批人与前置节点（非上一节点审批人相同时）previous
		autoRules := taskBasicConfigModel.AutoRules
		if len(autoRules) > 0 { // 自动审批通过
			params, err := t.flow.GetInstanceVariableValues(ctx, flowInstanceEntity)
			if err != nil {
				logger.Logger.Error(err)
			}
			formReq := client.FormDataConditionModel{
				AppID:   flowInstanceEntity.AppID,
				TableID: flowInstanceEntity.FormID,
				DataID:  flowInstanceEntity.FormInstanceID,
			}
			formData, err := t.formAPI.GetFormData(ctx, formReq)
			if err != nil {
				logger.Logger.Error(err)
			}
			params = utils.MergeMap(params, formData)

			if utils.Contain(autoRules, "origin") && utils.Contain(handleUserIds, flowInstanceEntity.CreatorID) {
				t.AutoReviewTask(ctx, flowEntity, flowInstanceEntity, task, flowInstanceEntity.CreatorID, params)
				taskAutoHandle = true
			} else if utils.Contain(autoRules, "parent") && utils.Contain(handleUserIds, currentUserID) {
				t.AutoReviewTask(ctx, flowEntity, flowInstanceEntity, task, currentUserID, params)
				taskAutoHandle = true
			} else if utils.Contain(autoRules, "previous") {
				agreeUserIds, _ := t.operationRecord.GetAgreeUserIds(ctx, flowInstanceEntity.ProcessInstanceID)
				if len(agreeUserIds) > 0 && utils.IntersectString(agreeUserIds, handleUserIds) {
					t.AutoReviewTask(ctx, flowEntity, flowInstanceEntity, task, currentUserID, params)
					taskAutoHandle = true
				}
			}
		}
	}

	if !taskAutoHandle {
		// 审批用时限制
		err = t.setTimeRule(ctx, flowInstanceEntity, task, taskBasicConfigModel.TimeRule)
		if err != nil {
			logger.Logger.Error(err)
		}
		// 催办
		err = t.TaskUrging(ctx, taskBasicConfigModel.TimeRule, task.NodeDefKey, *flowInstanceEntity, task.ID)
		if err != nil {
			logger.Logger.Error(err)
		}
	}

	// 检查流程实例状态，如果流程实例结束，则flow要同步状态
	processInstance, err := t.processAPI.GetInstanceByID(ctx, flowInstanceEntity.ProcessInstanceID)
	if err != nil {
		logger.Logger.Error(err)
	}
	if processInstance.Status != Active {
		dataMap := make(map[string]interface{})
		dataMap["status"] = opAgree
		dataMap["modifier_id"] = currentUserID
		err = t.flowInstanceRepo.Update(t.db, flowInstanceEntity.ID, dataMap)
		if err != nil {
			logger.Logger.Error(err)
		}
	}

	return nil
}

// noUserHandle no user handle
func (t *task) noUserHandle(ctx context.Context, flowInstanceEntity *models.Instance, task *client.ProcessTask, taskBasicConfigModel *convert.TaskBasicConfigModel) error {
	params, err := t.flow.GetInstanceVariableValues(ctx, flowInstanceEntity)
	if err != nil {
		return err
	}

	if "skip" == taskBasicConfigModel.WhenNoPerson { // 自动跳过节点
		err := t.processAPI.CompleteTask(ctx, flowInstanceEntity.ProcessInstanceID, task.ID, params, nil)
		if err != nil {
			return err
		}

		model := &models.HandleTaskModel{
			HandleType: OpAutoSkip,
			HandleDesc: "该节点下无相关负责人，已自动跳过",
		}

		// 增加操作日志
		task := &client.ProcessTask{
			NodeDefKey: task.NodeDefKey,
			Name:       task.Name,
		}
		t.operationRecord.AddOperationRecord(ctx, flowInstanceEntity, task, model)
	} else { // 无审批人时，转交给管理员
		t.processAPI.SetAssignee(ctx, flowInstanceEntity.ProcessInstanceID, task.ID, "0") // 锁定为无人处理
		// add to exception task list
		flowAbnormalTaskEntity := &models.AbnormalTask{
			FlowInstanceID:    flowInstanceEntity.ID,
			ProcessInstanceID: flowInstanceEntity.ProcessInstanceID,
			TaskID:            task.ID,
			TaskName:          task.Name,
			TaskDefKey:        task.NodeDefKey,
			Reason:            "节点无处理人",
			Status:            0,
			BaseModel: models.BaseModel{
				CreatorID:  "",
				ModifyTime: time2.Now(),
			},
		}
		err := t.abnormalTaskRepo.Create(t.db, flowAbnormalTaskEntity)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *task) TaskUrging(ctx context.Context, rule convert.TaskTimeRuleModel, nodeID string, instance models.Instance, TaskID string) error {
	if !rule.Enabled {
		return nil
	}
	deadLine := rule.DeadLine
	if deadLine.BreakPoint == firstEntry {
		resp, err := t.processAPI.GetHistoryTasks(ctx, client.GetTasksReq{
			InstanceID: []string{instance.ProcessInstanceID},
			NodeDefKey: nodeID,
		})
		if err != nil {
			return err
		}
		if resp != nil && len(resp.Data) > 0 {
			return nil
		}
	}
	if deadLine.Day == 0 && deadLine.Hours == 0 && deadLine.Minutes == 0 {
		return nil
	}

	// eval deadline
	now := time2.Now()
	deadLineTime := utils.AddDaysHoursMinutes(now, deadLine.Day, deadLine.Hours, deadLine.Minutes)
	otherInfo, _ := json.Marshal(&DispatcherCallOtherInfo{
		FlowInstanceID: instance.ID,
		TaskID:         TaskID,
		TaskDefKey:     nodeID,
	})
	id := id2.GenID()
	dc := models.DispatcherCallback{
		BaseModel: models.BaseModel{
			ID: id,
		},
		Type:      DEADLINE,
		OtherInfo: string(otherInfo),
	}
	err := t.dispatcherCallbackRepo.Create(t.db, &dc)
	if err != nil {
		return err
	}
	err = t.dispatcherAPI.TakePost(ctx, client.TaskPostReq{
		Code:    "flow:" + id,
		Type:    1,
		TimeBar: deadLineTime,
		State:   1,
	})
	if err != nil {
		return err
	}

	// urge
	urgeInfo := deadLine.Urge
	if urgeInfo.Day == 0 && urgeInfo.Hours == 0 && urgeInfo.Minutes == 0 {
		return nil
	}
	urgeTime := utils.AddDaysHoursMinutes(now, -urgeInfo.Day, -urgeInfo.Hours, -urgeInfo.Minutes)
	id = id2.GenID()
	dc = models.DispatcherCallback{
		BaseModel: models.BaseModel{
			ID: id,
		},
		Type:      URGE,
		OtherInfo: string(otherInfo),
	}
	err = t.dispatcherCallbackRepo.Create(t.db, &dc)
	if err != nil {
		return err
	}
	err = t.dispatcherAPI.TakePost(ctx, client.TaskPostReq{
		Code:    "flow:" + id,
		Type:    1,
		TimeBar: urgeTime,
		State:   1,
	})
	if err != nil {
		return err
	}
	return nil
}

func (t *task) AutoReviewTask(ctx context.Context, flowEntity *models.Flow, flowInstanceEntity *models.Instance, task *client.ProcessTask, userID string, params map[string]interface{}) error {
	// 保存数据到表单接口

	formData := t.FilterCanEditFormData(ctx, flowEntity, flowInstanceEntity, task.NodeDefKey, make(map[string]interface{}))
	if len(formData) > 0 {
		saveFormDataReq := client.UpdateEntity{
			Entity: formData["entity"],
			// Ref: formData["ref"],
		}

		err := t.formAPI.UpdateData(ctx, flowInstanceEntity.AppID, flowInstanceEntity.FormID, flowInstanceEntity.FormInstanceID, saveFormDataReq, false)
		if err != nil {
			return err
		}
	}

	comments := map[string]interface{}{
		"reviewResult": Agree,
		"reviewRemark": "",
	}

	err := t.processAPI.CompleteTask(ctx, flowInstanceEntity.ProcessInstanceID, task.ID, t.flow.FormatFormValue(flowInstanceEntity, params), comments)
	if err != nil {
		return err
	}

	model := &models.HandleTaskModel{
		HandleType: opAutoReview,
		HandleDesc: "自动审批通过",
	}

	// 增加操作日志
	model.AutoReviewUserID = userID
	return t.operationRecord.AddOperationRecord(ctx, flowInstanceEntity, task, model)
}

func (t *task) setTimeRule(ctx context.Context, flowInstanceEntity *models.Instance, task *client.ProcessTask, timeRule convert.TaskTimeRuleModel) error {
	// entry进入该节点后,firstEntry首次进入该节点后,flowWorked工作流开始后
	if len(timeRule.DeadLine.BreakPoint) > 0 {
		if "entry" == timeRule.DeadLine.BreakPoint { // 进入该节点后
			dueDate := utils.AddDaysHoursMinutes(task.CreateTime, int(timeRule.DeadLine.Day), int(timeRule.DeadLine.Hours), int(timeRule.DeadLine.Minutes))
			t.processAPI.SetDueDate(ctx, task.ID, dueDate)
		} else if "firstEntry" == timeRule.DeadLine.BreakPoint { // 首次进入该节点后
			taskCondition := client.GetTasksReq{
				InstanceID: []string{flowInstanceEntity.ProcessInstanceID},
				NodeDefKey: task.NodeDefKey,
			}
			actRuTaskEntity, err := t.processAPI.GetEarliestEntryThisTask(ctx, taskCondition)
			if err != nil {
				return err
			}

			if actRuTaskEntity != nil {
				dueDate := utils.AddDaysHoursMinutes(actRuTaskEntity.CreateTime, int(timeRule.DeadLine.Day), int(timeRule.DeadLine.Hours), int(timeRule.DeadLine.Minutes))
				t.processAPI.SetDueDate(ctx, task.ID, dueDate)
			}
		} else if "flowWorked" == timeRule.DeadLine.BreakPoint { // 工作流开始后
			dueDate := utils.AddDaysHoursMinutes(flowInstanceEntity.CreateTime, int(timeRule.DeadLine.Day), int(timeRule.DeadLine.Hours), int(timeRule.DeadLine.Minutes))
			t.processAPI.SetDueDate(ctx, task.ID, dueDate)
		}
	}
	return nil
}

// editable hidden write read 4位二进制表示，后端关注write和read,返回write,read
func (t *task) PermissionConvertWriteRead(permission int8) (bool, bool) {
	write := permission&2 == 2
	read := permission&1 == 1
	return write, read
}

// 保存数据时过滤可以编辑的form字段
func (t *task) FilterCanEditFormData(ctx context.Context, flowEntity *models.Flow, flowInstanceEntity *models.Instance,
	taskDefKey string, formData map[string]interface{}) map[string]interface{} {
	if len(formData) == 0 {
		formData = make(map[string]interface{}, 0)
	}

	shapeModel, err := convert.GetShapeByTaskDefKey(flowEntity.BpmnText, taskDefKey)
	if err != nil {
		return formData
	}
	fieldPermissionObj := shapeModel.Data.BusinessData["fieldPermission"]
	if fieldPermissionObj == nil {
		return formData
	}

	fieldPermissionModel := utils.ChangeObjectToMap(fieldPermissionObj)
	if len(fieldPermissionModel) == 0 {
		return formData
	}

	variableValues, err := t.flow.GetInstanceVariableValues(ctx, flowInstanceEntity)
	if err != nil {
		return formData
	}

	entityObj, ok := formData["entity"]
	if !ok {
		entityObj = new(interface{})
	}
	formData["entity"] = t.filterEntityCanEditFormData(entityObj, fieldPermissionModel, variableValues)

	refObj, ok := formData["ref"]
	if ok {
		refMap := utils.ChangeObjectToMap(refObj)
		for _, fieldValue := range refMap {
			refFieldOperationModel := &client.RefData{}
			fieldValueJSON, err := json.Marshal(fieldValue)
			if err != nil {
				continue
			}
			err = json.Unmarshal(fieldValueJSON, refFieldOperationModel)
			if err != nil {
				continue
			}

			if len(refFieldOperationModel.Updated) > 0 {
				for _, model := range refFieldOperationModel.Updated {
					model.Entity = t.filterEntityCanEditFormData(model.Entity, fieldPermissionModel, variableValues)
				}
			}
			if "associated_records" != refFieldOperationModel.Type {
				var newEntity []client.CreateEntity
				newJSON, err := json.Marshal(refFieldOperationModel.New)
				if err != nil {
					continue
				}
				err = json.Unmarshal(newJSON, &newEntity)
				if err != nil {
					continue
				}

				if len(newEntity) > 0 {
					for _, model := range newEntity {
						model.Entity = t.filterEntityCanEditFormData(model.Entity, fieldPermissionModel, variableValues)
					}
				}
				refFieldOperationModel.New = newEntity
			}
		}
		formData["ref"] = refMap
	}

	return formData
}

// filterEntityCanEditFormData filter entity can edit data by permission
func (t *task) filterEntityCanEditFormData(entity interface{}, fieldPermissionMap map[string]interface{}, variableMap map[string]interface{}) map[string]interface{} {
	entityMap := utils.ChangeObjectToMap(entity)
	if entityMap == nil {
		entityMap = make(map[string]interface{}, 0)
	}
	for pKey, pValue := range fieldPermissionMap {
		permissionModel := &convert.FieldPermissionModel{}
		permissionModelJSON, err := json.Marshal(pValue)
		if err != nil {
			continue
		}
		err = json.Unmarshal(permissionModelJSON, permissionModel)
		if err != nil {
			continue
		}

		write, _ := t.PermissionConvertWriteRead(permissionModel.XInternal.Permission)
		if write {
			eFieldValue := entityMap[pKey]
			if utils.IsNotNil(permissionModel.SubmitValue.Variable) || utils.IsNotNil(permissionModel.SubmitValue.StaticValue) {
				if utils.IsNotNil(permissionModel.SubmitValue.Variable) {
					entityMap[pKey] = variableMap[permissionModel.SubmitValue.Variable.(string)]
				} else if utils.IsNotNil(permissionModel.SubmitValue.StaticValue) {
					entityMap[pKey] = permissionModel.SubmitValue.StaticValue
				}
			} else if utils.IsNil(eFieldValue) && (utils.IsNotNil(permissionModel.InitialValue.Variable) || utils.IsNotNil(permissionModel.InitialValue.StaticValue)) {
				if utils.IsNotNil(permissionModel.InitialValue.Variable) {
					entityMap[pKey] = variableMap[permissionModel.InitialValue.Variable.(string)]
				} else if utils.IsNotNil(permissionModel.InitialValue.StaticValue) {
					entityMap[pKey] = permissionModel.InitialValue.StaticValue
				}
			}
		}
	}

	for fieldKey := range entityMap {
		fieldPermission := fieldPermissionMap[fieldKey]
		permissionModel := &convert.FieldPermissionModel{}
		permissionModelJSON, err := json.Marshal(fieldPermission)
		if err != nil {
			continue
		}
		err = json.Unmarshal(permissionModelJSON, permissionModel)
		if err != nil {
			continue
		}

		write, _ := t.PermissionConvertWriteRead(permissionModel.XInternal.Permission)
		if !write {
			delete(entityMap, fieldKey)
		}
	}
	return entityMap
}

// FilterCanReadFormData filter can read from data
func (t *task) FilterCanReadFormData(ctx context.Context, flowInstanceEntity *models.Instance, fieldPermissionObj interface{}, formData interface{}) map[string]interface{} {
	// fieldPermissionObj := model.FieldPermission
	if fieldPermissionObj == nil {
		return nil
	}

	fieldPermissionModel := utils.ChangeObjectToMap(fieldPermissionObj)
	if len(fieldPermissionModel) == 0 {
		return nil
	}

	variableValues, err := t.flow.GetInstanceVariableValues(ctx, flowInstanceEntity)
	if err != nil {
		return nil
	}

	return t.filterEntityCanReadFormData(formData, fieldPermissionModel, variableValues)

	// todo 子表单字段过滤,子表单权限暂时不校验，后期跟form一起加
}

// filterEntityCanReadFormData filter entity can read data by permission
func (t *task) filterEntityCanReadFormData(entity interface{}, fieldPermissionMap map[string]interface{}, variableMap map[string]interface{}) map[string]interface{} {
	entityMap := utils.ChangeObjectToMap(entity)
	for pKey, pValue := range fieldPermissionMap {
		permissionModel := &convert.FieldPermissionModel{}
		permissionModelJSON, err := json.Marshal(pValue)
		if err != nil {
			continue
		}
		err = json.Unmarshal(permissionModelJSON, permissionModel)
		if err != nil {
			continue
		}

		_, read := t.PermissionConvertWriteRead(permissionModel.XInternal.Permission)
		if read {
			eFieldValue := entityMap[pKey]
			if utils.IsNil(eFieldValue) && (utils.IsNotNil(permissionModel.InitialValue.Variable) || utils.IsNotNil(permissionModel.InitialValue.StaticValue)) {
				if utils.IsNotNil(permissionModel.InitialValue.Variable) {
					entityMap[pKey] = variableMap[permissionModel.InitialValue.Variable.(string)]
				} else if utils.IsNotNil(permissionModel.InitialValue.StaticValue) {
					entityMap[pKey] = permissionModel.InitialValue.StaticValue
				}
			}
		}
	}

	for fieldKey := range entityMap {
		fieldPermission := fieldPermissionMap[fieldKey]
		permissionModel := &convert.FieldPermissionModel{}
		permissionModelJSON, err := json.Marshal(fieldPermission)
		if err != nil {
			continue
		}
		err = json.Unmarshal(permissionModelJSON, permissionModel)
		if err != nil {
			continue
		}
		_, read := t.PermissionConvertWriteRead(permissionModel.XInternal.Permission)
		if !read {
			delete(entityMap, fieldKey)
		}
	}
	return entityMap
}

// GetCurrentNodes
func (t *task) GetCurrentNodes(ctx context.Context, processInstanceID string) ([]models.NodeModel, error) {
	taskCondition := client.GetTasksReq{
		InstanceID: []string{processInstanceID},
	}
	tasksResp, err := t.processAPI.GetTasks(ctx, taskCondition)
	if err != nil {
		return nil, err
	}
	if len(tasksResp.Data) > 0 {
		nodes := make([]models.NodeModel, 0)
		for _, value := range tasksResp.Data {
			if value.TaskType == "NON_MODEL" {
				continue
			}
			node := models.NodeModel{
				TaskDefKey: value.NodeDefKey,
				TaskName:   value.Name,
			}
			nodes = append(nodes, node)
		}

		return nodes, nil
	}
	return nil, nil
}
