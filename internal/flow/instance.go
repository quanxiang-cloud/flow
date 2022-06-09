package flow

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/quanxiang-cloud/flow/internal"
	"github.com/quanxiang-cloud/flow/internal/convert"
	"github.com/quanxiang-cloud/flow/internal/models"
	"github.com/quanxiang-cloud/flow/internal/models/mysql"
	"github.com/quanxiang-cloud/flow/internal/server/options"
	"github.com/quanxiang-cloud/flow/pkg"
	"github.com/quanxiang-cloud/flow/pkg/client"
	"github.com/quanxiang-cloud/flow/pkg/code"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/misc/error2"
	"github.com/quanxiang-cloud/flow/pkg/misc/id2"
	"github.com/quanxiang-cloud/flow/pkg/misc/logger"
	"github.com/quanxiang-cloud/flow/pkg/misc/time2"
	"github.com/quanxiang-cloud/flow/pkg/page"
	"github.com/quanxiang-cloud/flow/pkg/utils"
	"gorm.io/gorm"
	"sort"
	"strings"
)

const (
	// Completed status
	Completed = "COMPLETED"
	// Active status
	Active = "ACTIVE"
)

const (
	// Review instance status
	Review = "REVIEW" // 待处理
	// InReview instance status
	InReview = "IN_REVIEW" // 处理中
	// SendBack instance status
	SendBack = "SEND_BACK" // 打回重填
	// Agree instance status
	Agree = "AGREE" // 通过
	// Refuse instance status
	Refuse = "REFUSE" // 拒绝
	// Cancel instance status
	Cancel = "CANCEL" // 撤销
	// Abandon instance status
	Abandon = "ABANDON" // 作废
	// Abend instance status
	Abend = "ABEND" // 异常结束

	// UnTreated 还未处理
	UnTreated = "UNTREATED"
)

// WaitHandleTasks wait handle task type
func WaitHandleTasks() []string {
	return []string{convert.ReviewTask, convert.WriteTask, convert.ReadTask, convert.SendBackTask}
}

// GetStatusName get status name
func GetStatusName(status string) string {
	switch status {
	case "REVIEW":
		{
			return "待审批"
		}
	case "IN_REVIEW":
		{
			return "审批中"
		}
	case "SEND_BACK":
		{
			return "待补充"
		}
	case "AGREE":
		{
			return "通过"
		}
	case "REFUSE":
		{
			return "拒绝"
		}
	case "CANCEL":
		{
			return "撤销"
		}
	case "ABANDON":
		{
			return "作废"
		}
	case "ABEND":
		{
			return "异常结束"
		}
	}

	return ""
}

func isFlowOngoing(status string) bool {
	return Review == status || InReview == status || SendBack == status
}

// GetOngoingStatus get going status
func GetOngoingStatus() []string {
	return []string{Review, InReview, SendBack}
}

func isReviewStatus(status string) bool {
	return opAgree == status || opRefuse == status || opFillIn == status
}

// Instance service
type Instance interface {
	StartFlow(ctx context.Context, req *StartFlowModel) (string, error)
	FormDataDeleted(ctx context.Context, formDataIDs []string, userID string) error
	AbendFlowInstance(ctx context.Context, flowInstanceEntity *models.Instance, remark string, userID string) error

	MyApplyList(ctx context.Context, req *MyApplyReq) (*page.RespPage, error)
	WaitReviewList(ctx context.Context, req *TaskListReq) (*page.RespPage, error)
	ReviewedList(ctx context.Context, req *TaskListReq) (*page.RespPage, error)
	CcToMeList(ctx context.Context, req *CcListReq) (*page.RespPage, error)
	AllList(ctx context.Context, req *TaskListReq) (*page.RespPage, error)
	InstanceAddFormData(ctx context.Context, data []map[string]interface{}) error

	ConvertTask(ctx context.Context, tasks []*client.ProcessTask) ([]*models.ActTaskEntity, []string, []string)
	ConvertInstance(ctx context.Context, instances []*client.ProcessInstance) ([]*models.Instance, []string)
	FlowInstanceMapAddFormData(ctx context.Context, processInstanceIDs []string) map[string]interface{}

	Cancel(ctx context.Context, processInstanceID string) (bool, error)
	Resubmit(ctx context.Context, processInstanceID string, req *ResubmitReq) (bool, error)
	FlowInfo(ctx context.Context, processInstanceID string) (*models.Flow, error)
	SendBack(ctx context.Context, processInstanceID string, taskID string, model *models.HandleTaskModel) (bool, error)
	StepBack(ctx context.Context, processInstanceID string, taskID string, model *models.HandleTaskModel) (bool, error)
	CcFlow(ctx context.Context, processInstanceID string, taskID string, model *models.HandleTaskModel) (bool, error)
	ReadFlow(ctx context.Context, processInstanceID string, taskID string, model *models.HandleTaskModel) (bool, error)
	HandleCc(ctx context.Context, taskIDs []string) (bool, error)
	HandleRead(ctx context.Context, processInstanceID string, taskID string, model *models.HandleTaskModel) (bool, error)

	AddSign(ctx context.Context, processInstanceID string, taskID string, model *models.AddSignModel) (bool, error)
	DeliverTask(ctx context.Context, processInstanceID string, taskID string, model *models.HandleTaskModel) (bool, error)
	StepBackNodes(ctx context.Context, processInstanceID string) ([]*models.NodeModel, error)
	FlowInstanceCount(ctx context.Context) (*InstanceCountModel, error)

	ReviewTask(ctx context.Context, processInstanceID string, taskID string, model *models.HandleTaskModel) (bool, error)
	GetFlowInstanceForm(ctx context.Context, processInstanceID string, taskTypeDetailModel *TaskTypeDetailModel) (*InstanceDetailModel, error)
	GetFormData(ctx context.Context, processInstanceID string, taskID string, req *GetFormDataReq) (interface{}, error)
	ProcessHistories(ctx context.Context, processInstanceID string) ([]*models.InstanceStep, error)

	GetInstanceVariableValues(ctx context.Context, instance *models.Instance) (map[string]interface{}, error)
	Cal(ctx context.Context, valueFrom string, valueOf interface{}, formulaFields interface{}, instance *models.Instance, variables map[string]interface{}, formQueryRef interface{}, formDefKey string) (interface{}, error)
}

type instance struct {
	db                    *gorm.DB
	conf                  *config.Configs
	instanceRepo          models.InstanceRepo
	flowRepo              models.FlowRepo
	operationRecordRepo   models.OperationRecordRepo
	operationRecord       OperationRecord
	flow                  Flow
	task                  Task
	processAPI            client.Process
	identityAPI           client.Identity
	formAPI               client.Form
	structorAPI           client.Structor
	urgeRepo              models.UrgeRepo
	stepRepo              models.InstanceStepRepo
	recordRepo            models.OperationRecordRepo
	instanceVariablesRepo models.InstanceVariablesRepo
	variablesRepo         models.VariablesRepo
	instanceExecutionRepo models.InstanceExecutionRepo
	abnormalTaskRepo      models.AbnormalTaskRepo
}

// NewInstance init
func NewInstance(conf *config.Configs, opts ...options.Options) (Instance, error) {
	flow, _ := NewFlow(conf, opts...)
	operationRecord, _ := NewOperationRecord(conf, opts...)
	task, _ := NewTask(conf, opts...)
	i := &instance{
		conf:                  conf,
		instanceRepo:          mysql.NewInstanceRepo(),
		flowRepo:              mysql.NewFlowRepo(),
		flow:                  flow,
		operationRecord:       operationRecord,
		task:                  task,
		processAPI:            client.NewProcess(conf),
		identityAPI:           client.NewIdentity(conf),
		structorAPI:           client.NewStructor(conf),
		formAPI:               client.NewForm(conf),
		urgeRepo:              mysql.NewUrgeRepo(),
		operationRecordRepo:   mysql.NewOperationRecordRepo(),
		stepRepo:              mysql.NewInstanceStepRepo(),
		recordRepo:            mysql.NewOperationRecordRepo(),
		variablesRepo:         mysql.NewVariablesRepo(),
		instanceVariablesRepo: mysql.NewInstanceVariablesRepo(),
		instanceExecutionRepo: mysql.NewInstanceExecutionRepo(),
		abnormalTaskRepo:      mysql.NewAbnormalTaskRepo(),
	}

	for _, opt := range opts {
		opt(i)
	}
	return i, nil
}

// SetDB set db
func (i *instance) SetDB(db *gorm.DB) {
	i.db = db
}

// StartFlow start flow
func (i *instance) StartFlow(ctx context.Context, req *StartFlowModel) (string, error) {
	flowEntity, err := i.flowRepo.FindByID(i.db, req.FlowID)
	if err != nil {
		return "", err
	}

	appCenter := client.NewAppCenter(i.conf)
	appName, err := appCenter.GetAppName(ctx, flowEntity.AppID)

	identity := client.NewIdentity(i.conf)
	myUserEntity, err := identity.FindUserByID(ctx, req.UserID)
	if err != nil || myUserEntity == nil {
		return "", error2.NewErrorWithString(error2.Internal, "Can not find apply user info ")
	}

	formDataID, ok := req.FormData["_id"]
	if !ok {
		return "", nil
	}
	formReq := client.FormDataConditionModel{
		AppID:   flowEntity.AppID,
		TableID: flowEntity.FormID,
		DataID:  formDataID.(string),
	}
	formData, err := i.formAPI.GetFormData(ctx, formReq)
	if err != nil {
		return "", error2.NewErrorWithString(error2.Internal, "Get form data error ")
	}

	req.FormData = formData

	process := client.NewProcess(i.conf)
	startProcessResp, err := process.StartProcessInstance(ctx, client.StartProcessReq{
		ProcessID: flowEntity.ProcessID,
		UserID:    req.UserID,
		Params:    req.FormData,
	})
	if err != nil {
		return "", err
	}

	flowInstanceEntity := &models.Instance{
		AppID:             flowEntity.AppID,
		AppName:           appName,
		FlowID:            req.FlowID,
		ProcessInstanceID: startProcessResp.InstanceID,
		FormID:            flowEntity.FormID,
		FormInstanceID:    req.FormData["_id"].(string),
		ApplyUserID:       req.UserID,
		ApplyUserName:     myUserEntity.UserName,
		AppStatus:         mysql.AppActiveStatus,
		BaseModel: models.BaseModel{
			CreatorID:  req.UserID,
			ModifyTime: time2.Now(),
			CreateTime: time2.Now(),
		},
		Status: Review,
	}

	if len(flowEntity.InstanceName) > 0 {
		flowInstanceEntity.Name = i.instanceNameConverter(nil, flowEntity.InstanceName, flowInstanceEntity)
	} else {
		flowInstanceEntity.Name = flowEntity.Name
	}
	err = i.instanceRepo.Create(i.db, flowInstanceEntity)

	// add instance variable value
	err = i.addInstanceVariableValues(ctx, flowInstanceEntity)
	params, err := i.GetInstanceVariableValues(ctx, flowInstanceEntity)
	params = utils.MergeMap(params, req.FormData)
	if err != nil {
		return "", err
	}
	_, err = process.InitProcessInstance(ctx, client.InitInstanceReq{
		InstanceID: startProcessResp.InstanceID,
		UserID:     req.UserID,
		Params:     params,
	})
	if err != nil {
		return "", err
	}

	// Add operation record
	handleTaskModel := &models.HandleTaskModel{
		HandleType: opSubmit,
		HandleDesc: "发起流程",
	}
	i.operationRecord.AddOperationRecord(ctx, flowInstanceEntity, nil, handleTaskModel)

	// i.task.TaskCheck(ctx, flowEntity, flowInstanceEntity, req.UserID)
	processInstance, err := i.processAPI.GetInstanceByID(ctx, startProcessResp.InstanceID)
	if err != nil {
		return "", err
	}
	if processInstance != nil && processInstance.Status != Active {
		// 说明此时流程结束了
		dataMap := make(map[string]interface{})
		dataMap["modifier_id"] = pkg.STDUserID(ctx)
		dataMap["status"] = Agree
		err = i.instanceRepo.Update(i.db, flowInstanceEntity.ID, dataMap)
		if err != nil {
			return "", err
		}
	}

	return flowInstanceEntity.ID, nil
}

func (i *instance) addInstanceVariableValues(ctx context.Context, instance *models.Instance) error {
	flowVariablesEntities, err := i.variablesRepo.FindVariablesByFlowID(i.db, instance.FlowID)
	if err != nil {
		return err
	}

	instanceVariables := make([]*models.InstanceVariables, 0)
	for _, value := range flowVariablesEntities {
		if "CUSTOM" == value.Type {
			instanceVariable := &models.InstanceVariables{
				ProcessInstanceID: instance.ProcessInstanceID,
				Name:              value.Name,
				Type:              value.Type,
				Code:              value.Code,
				FieldType:         value.FieldType,
				Format:            value.Format,
				Value:             value.DefaultValue,
				Desc:              value.Desc,
				BaseModel: models.BaseModel{
					ID:         id2.GenID(),
					CreateTime: time2.Now(),
				},
			}
			instanceVariables = append(instanceVariables, instanceVariable)
		}
	}
	return i.instanceVariablesRepo.BatchCreate(i.db, instanceVariables)
}

func (i *instance) GetInstanceVariableValues(ctx context.Context, instance *models.Instance) (map[string]interface{}, error) {
	return i.flow.GetInstanceVariableValues(ctx, instance)
}

func (i *instance) instanceNameConverter(ctx context.Context, instanceName string, instance *models.Instance) string {
	valueMap, _ := i.flow.GetFlowVariableValues(ctx, instance)
	for key, value := range valueMap {
		instanceName = strings.Replace(instanceName, "$"+key, utils.Strval(value), -1)
		instanceName = strings.Replace(instanceName, key, utils.Strval(value), -1)
	}
	return instanceName
}

// FormDataDeleted form data deleted abend flow instance
func (i *instance) FormDataDeleted(ctx context.Context, formDataIDs []string, userID string) error {
	if len(formDataIDs) > 0 {
		for _, value := range formDataIDs {
			i.AbendFormDataDeleted(ctx, value, userID)
		}
	}

	return nil
}

func (i *instance) AbendFormDataDeleted(ctx context.Context, formDataID string, userID string) error {
	instances, err := i.instanceRepo.FindInstances(i.db, formDataID, GetOngoingStatus())
	if err != nil {
		return err
	}

	if len(instances) > 0 {
		for _, instance := range instances {
			tasks, _ := i.processAPI.GetTasksByInstanceID(ctx, instance.ProcessInstanceID)
			if len(tasks) > 0 {
				for _, task := range tasks {
					// add to exception task list
					flowAbnormalTaskEntity := &models.AbnormalTask{
						FlowInstanceID:    instance.ID,
						ProcessInstanceID: instance.ProcessInstanceID,
						TaskID:            task.ID,
						TaskName:          task.Name,
						TaskDefKey:        task.NodeDefKey,
						Reason:            "表单数据被删除，任务异常结束",
						Status:            2,
						BaseModel: models.BaseModel{
							CreatorID:  "",
							ModifyTime: time2.Now(),
						},
					}
					err := i.abnormalTaskRepo.Create(i.db, flowAbnormalTaskEntity)
					if err != nil {
						logger.Logger.Error(err)
					}
				}
			}

			i.AbendFlowInstance(ctx, instance, "表单数据被删除，流程异常结束", userID)
		}
	}
	return nil
}

func (i *instance) AbendFlowInstance(ctx context.Context, flowInstanceEntity *models.Instance, remark string, userID string) error {
	processInstanceID := flowInstanceEntity.ProcessInstanceID

	err := i.processAPI.AbendInstance(ctx, processInstanceID)
	if err != nil {
		return err
	}

	updateMap := map[string]interface{}{
		"status": opAbend,
	}
	err = i.instanceRepo.Update(i.db, flowInstanceEntity.ID, updateMap)
	if err != nil {
		return err
	}

	handleTaskModel := &models.HandleTaskModel{
		HandleType: opAbend,
		HandleDesc: "异常结束",
		Remark:     remark,
	}
	i.operationRecord.AddOperationRecord(ctx, flowInstanceEntity, nil, handleTaskModel)

	return nil
}

func (i *instance) InstanceAddFormData(ctx context.Context, data []map[string]interface{}) error {
	formDataCondition := make([]client.FormDataConditionModel, 0)
	formSchemaCondition := make([]client.FormSchemaConditionModel, 0)
	if len(data) > 0 {
		flowIDs := make([]string, 0)
		for _, value := range data {
			flowIDs = append(flowIDs, value["flowId"].(string))
		}

		flowKeyFieldMap, _ := i.flow.GetFlowKeyFields(ctx, flowIDs)

		for _, flowInstanceEntity := range data {
			c := client.FormDataConditionModel{
				AppID:   flowInstanceEntity["appId"].(string),
				TableID: flowInstanceEntity["formId"].(string),
				DataID:  flowInstanceEntity["formInstanceId"].(string),
			}
			formDataCondition = append(formDataCondition, c)

			s := client.FormSchemaConditionModel{
				AppID:   flowInstanceEntity["appId"].(string),
				TableID: flowInstanceEntity["formId"].(string),
			}
			formSchemaCondition = append(formSchemaCondition, s)
		}

		formDataMap, _ := i.formAPI.BatchGetFormData(ctx, formDataCondition)
		formSchemaMap, _ := i.formAPI.BatchGetFormSchema(ctx, formSchemaCondition)

		for _, flowInstanceEntity := range data {
			if formDataMap != nil {
				flowInstanceEntity["formData"] = formDataMap[flowInstanceEntity["formInstanceId"].(string)]
			}
			if formSchemaMap != nil {
				flowInstanceEntity["formSchema"] = formSchemaMap[flowInstanceEntity["formId"].(string)]
			}

			flowInstanceEntity["keyFields"] = flowKeyFieldMap[flowInstanceEntity["flowId"].(string)]
		}
	}

	return nil
}

func (i *instance) MyApplyList(ctx context.Context, req *MyApplyReq) (*page.RespPage, error) {
	userID := pkg.STDUserID(ctx)

	if len(req.OrderType) == 0 {
		req.OrderType = page.Asc
	}
	orderItem := page.OrderItem{
		Column:    "create_time",
		Direction: req.OrderType,
	}
	req.ReqPage.Orders = []page.OrderItem{orderItem}

	queryWrapper := &models.PageInstancesReq{
		ReqPage:         req.ReqPage,
		ApplyUserID:     userID,
		Status:          req.Status,
		CreateTimeBegin: req.BeginDate,
		CreateTimeEnd:   req.EndDate,
		Keyword:         req.Keyword,
	}

	dataList, count, err := i.instanceRepo.PageInstances(i.db, queryWrapper)
	if err != nil {
		return nil, err
	}

	resp := &page.RespPage{}
	if len(dataList) > 0 {
		data := utils.ChangeObjectToMapList(dataList)
		i.identityAPI.AddUserInfo(ctx, data)
		i.InstanceAddFormData(ctx, data)

		for _, flowInstanceEntity := range data {
			nodes, _ := i.task.GetCurrentNodes(ctx, flowInstanceEntity["processInstanceId"].(string))
			flowInstanceEntity["nodes"] = nodes
		}

		resp.TotalCount = count
		resp.Data = data
	}

	if resp.Data == nil {
		resp.Data = make([]map[string]interface{}, 0)
	}

	return resp, nil
}

func (i *instance) WaitReviewList(ctx context.Context, req *TaskListReq) (*page.RespPage, error) {
	userID := pkg.STDUserID(ctx)
	resp := &page.RespPage{}

	tasksReq := client.GetTasksReq{}
	tasksReq.Assignee = userID
	if len(req.TagType) > 0 {
		if req.TagType == "OVERTIME" {
			tasksReq.DueTime = time2.Now()
		} else if req.TagType == "URGE" {
			tasksReq.TaskID, _ = i.urgeRepo.FindTaskIDs(i.db)
			if len(tasksReq.TaskID) == 0 {
				return &page.RespPage{
					TotalCount: 0,
					Data:       make([]*models.ActTaskEntity, 0),
				}, nil
			}
		}
	}

	// handle type : REVIEW, WRITE， READ， OTHER
	if len(req.HandleType) > 0 {
		tasksReq.Desc = []string{req.HandleType}
	} else {
		tasksReq.Desc = WaitHandleTasks()
	}

	tasksReq.Page = req.Page
	tasksReq.Limit = req.Size
	if len(req.OrderType) == 0 {
		req.OrderType = internal.OrderTypeDesc
	}
	order := client.QueryOrder{
		Column:    "create_time",
		OrderType: req.OrderType,
	}
	orders := []client.QueryOrder{order}
	tasksReq.Order = orders

	if len(req.AppID) > 0 || len(req.Keyword) > 0 {
		queryWrapper := &models.PageInstancesReq{
			AppID:   req.AppID,
			Keyword: req.Keyword,
		}
		instances, _, _ := i.instanceRepo.PageInstances(i.db, queryWrapper)
		if len(instances) > 0 {
			instanceIDs := make([]string, 0)
			for _, value := range instances {
				instanceIDs = append(instanceIDs, value.ProcessInstanceID)
			}
			tasksReq.InstanceID = instanceIDs
		}
	}

	tasksResp, err := i.processAPI.GetTasks(ctx, tasksReq)
	if err != nil {
		return resp, err
	}

	// TaskModel 需要转换成 ActTaskEntity
	list, processInstanceIDs, taskIDs := i.ConvertTask(ctx, tasksResp.Data)

	// 催办次数
	urgeNums, _ := i.urgeRepo.GetUrgeNums(i.db, taskIDs)
	if len(urgeNums) > 0 {
		for _, value := range list {
			value.UrgeNum = urgeNums[value.ID]
		}
	}

	flowInstanceMap := i.FlowInstanceMapAddFormData(ctx, processInstanceIDs)
	for _, task := range list {
		flowInstanceEntity := flowInstanceMap[task.ProcInstID]
		task.FlowInstanceEntity = flowInstanceEntity
	}

	resp.TotalCount = tasksResp.TotalCount
	resp.Data = list
	return resp, nil
}

func (i *instance) ReviewedList(ctx context.Context, req *TaskListReq) (*page.RespPage, error) {
	userID := pkg.STDUserID(ctx)
	resp := &page.RespPage{}
	if req.Agent == 1 {
		return resp, nil
	}

	tasksReq := client.GetTasksReq{}
	tasksReq.Assignee = userID
	tasksReq.Page = req.Page
	tasksReq.Limit = req.Size
	tasksReq.Desc = WaitHandleTasks()
	if len(req.OrderType) == 0 {
		req.OrderType = internal.OrderTypeDesc
	}
	order := client.QueryOrder{
		Column:    "modify_time",
		OrderType: req.OrderType,
	}
	orders := []client.QueryOrder{order}
	tasksReq.Order = orders

	if len(req.AppID) > 0 || len(req.Keyword) > 0 {
		queryWrapper := &models.PageInstancesReq{
			AppID:   req.AppID,
			Keyword: req.Keyword,
		}
		instances, _, _ := i.instanceRepo.PageInstances(i.db, queryWrapper)
		if len(instances) > 0 {
			instanceIDs := make([]string, 0)
			for _, value := range instances {
				instanceIDs = append(instanceIDs, value.ProcessInstanceID)
			}
			tasksReq.InstanceID = instanceIDs
		}
	}

	tasksResp, err := i.processAPI.GetDoneInstances(ctx, tasksReq)
	if err != nil {
		return resp, err
	}

	// ProcessInstanceModel 需要转换成 FlowInstanceEntity
	list, _ := i.ConvertInstance(ctx, tasksResp.Data)

	if len(list) > 0 {
		data := utils.ChangeObjectToMapList(list)
		if err = i.identityAPI.AddUserInfo(ctx, data); err != nil {
			return nil, err
		}
		if err = i.InstanceAddFormData(ctx, data); err != nil {
			return nil, err
		}

		for _, flowInstanceEntity := range data {
			nodes, _ := i.task.GetCurrentNodes(ctx, flowInstanceEntity["processInstanceId"].(string))
			flowInstanceEntity["nodes"] = nodes
		}

		resp.TotalCount = tasksResp.TotalCount
		resp.Data = data
		return resp, nil
	}

	if resp.Data == nil {
		resp.Data = make([]map[string]interface{}, 0)
	}
	return resp, nil
}

func (i *instance) CcToMeList(ctx context.Context, req *CcListReq) (*page.RespPage, error) {
	userID := pkg.STDUserID(ctx)
	resp := &page.RespPage{}

	tasksReq := client.GetTasksReq{}
	tasksReq.Desc = []string{convert.CcTask}
	if req.Status == 0 {
		tasksReq.Status = Active
	} else if req.Status == 1 {
		tasksReq.Status = Completed
	}

	tasksReq.Assignee = userID
	tasksReq.Page = req.Page
	tasksReq.Limit = req.Size
	if len(req.OrderType) == 0 {
		req.OrderType = internal.OrderTypeDesc
	}
	order := client.QueryOrder{
		Column:    "create_time",
		OrderType: req.OrderType,
	}
	orders := []client.QueryOrder{order}
	tasksReq.Order = orders

	if len(req.AppID) > 0 || len(req.Keyword) > 0 {
		queryWrapper := &models.PageInstancesReq{
			AppID:   req.AppID,
			Keyword: req.Keyword,
		}
		instances, _, _ := i.instanceRepo.PageInstances(i.db, queryWrapper)
		if len(instances) > 0 {
			instanceIDs := make([]string, 0)
			for _, value := range instances {
				instanceIDs = append(instanceIDs, value.ProcessInstanceID)
			}
			tasksReq.InstanceID = instanceIDs
		}
	}

	tasksResp, err := i.processAPI.GetAllTasks(ctx, tasksReq)
	if err != nil {
		return resp, err
	}

	// TaskModel 需要转换成 ActTaskEntity
	list, processInstanceIDs, _ := i.ConvertTask(ctx, tasksResp.Data)

	taskIDs := make([]string, 0)
	for _, value := range tasksResp.Data {
		taskIDs = append(taskIDs, value.ID)
	}

	flowInstanceMap := i.FlowInstanceMapAddFormData(ctx, processInstanceIDs)
	for _, task := range list {
		flowInstanceEntity := flowInstanceMap[task.ProcInstID]
		task.FlowInstanceEntity = flowInstanceEntity
	}

	resp.TotalCount = tasksResp.TotalCount
	resp.Data = list

	if resp.Data == nil {
		resp.Data = make([]*models.ActTaskEntity, 0)
	}
	return resp, nil
}

func (i *instance) AllList(ctx context.Context, req *TaskListReq) (*page.RespPage, error) {
	userID := pkg.STDUserID(ctx)
	resp := &page.RespPage{}

	tasksReq := client.GetTasksReq{}
	tasksReq.Assignee = userID

	tasksReq.Page = req.Page
	tasksReq.Limit = req.Size
	if len(req.OrderType) == 0 {
		req.OrderType = internal.OrderTypeDesc
	}
	order := client.QueryOrder{
		Column:    "modify_time",
		OrderType: req.OrderType,
	}
	orders := []client.QueryOrder{order}
	tasksReq.Order = orders

	if len(req.AppID) > 0 || len(req.Keyword) > 0 {
		queryWrapper := &models.PageInstancesReq{
			AppID:   req.AppID,
			Keyword: req.Keyword,
		}
		instances, _, _ := i.instanceRepo.PageInstances(i.db, queryWrapper)
		if len(instances) > 0 {
			instanceIDs := make([]string, 0)
			for _, value := range instances {
				instanceIDs = append(instanceIDs, value.ProcessInstanceID)
			}
			tasksReq.InstanceID = instanceIDs
		}
	}

	tasksResp, err := i.processAPI.GetInstances(ctx, tasksReq)
	if err != nil {
		return resp, err
	}

	list, _ := i.ConvertInstance(ctx, tasksResp.Data)
	if len(list) > 0 {
		data := utils.ChangeObjectToMapList(list)
		i.identityAPI.AddUserInfo(ctx, data)
		i.InstanceAddFormData(ctx, data)

		for _, flowInstanceEntity := range data {
			nodes, _ := i.task.GetCurrentNodes(ctx, flowInstanceEntity["processInstanceId"].(string))
			flowInstanceEntity["nodes"] = nodes
		}

		resp.TotalCount = tasksResp.TotalCount
		resp.Data = data
	}

	if resp.Data == nil {
		resp.Data = make([]map[string]interface{}, 0)
	}

	return resp, nil
}

func (i *instance) FlowInstanceMapAddFormData(ctx context.Context, processInstanceIDs []string) map[string]interface{} {
	flowInstanceMap := make(map[string]interface{}, 0)
	list, err := i.instanceRepo.FindByProcessInstanceIDs(i.db, processInstanceIDs)
	if err != nil {
		return nil
	}

	data := utils.ChangeObjectToMapList(list)

	i.identityAPI.AddUserInfo(ctx, data)

	if len(data) > 0 {
		i.InstanceAddFormData(ctx, data)

		for _, flowInstanceEntity := range data {
			flowInstanceMap[flowInstanceEntity["processInstanceId"].(string)] = flowInstanceEntity
		}
	}

	return flowInstanceMap
}

func (i *instance) Cancel(ctx context.Context, processInstanceID string) (bool, error) {
	userID := pkg.STDUserID(ctx)

	flowInstanceEntity, err := i.instanceRepo.GetEntityByProcessInstanceID(i.db, processInstanceID)
	if err != nil {
		return false, err
	}
	if flowInstanceEntity == nil {
		return false, error2.NewErrorWithString(error2.Internal, "Can not find flow instance data ")
	}

	flowEntity, err := i.flowRepo.FindByID(i.db, flowInstanceEntity.FlowID)
	if err != nil {
		return false, err
	}
	if flowEntity == nil {
		return false, error2.NewErrorWithString(error2.Internal, "Can not find flow data ")
	}

	tasks, err := i.processAPI.GetTasksByInstanceID(ctx, processInstanceID)
	if err != nil {
		return false, err
	}

	if !(isFlowOngoing(flowInstanceEntity.Status) && flowInstanceEntity.CreatorID == userID && i.flow.checkCanCancel(ctx, flowEntity, flowInstanceEntity, tasks)) {
		return false, error2.NewErrorWithString(error2.Internal, "Can not cancel ")
	}

	err = i.processAPI.AbendInstance(ctx, processInstanceID)
	if err != nil {
		return false, err
	}

	updateMap := make(map[string]interface{}, 0)
	updateMap["modifier_id"] = userID
	updateMap["status"] = Cancel
	i.instanceRepo.Update(i.db, flowInstanceEntity.ID, updateMap)

	// Add operation record
	handleTaskModel := &models.HandleTaskModel{
		HandleType: opCancel,
		HandleDesc: "撤销",
	}
	i.operationRecord.AddOperationRecord(ctx, flowInstanceEntity, tasks[0], handleTaskModel)

	return true, nil
}

func (i *instance) Resubmit(ctx context.Context, processInstanceID string, params *ResubmitReq) (bool, error) {
	userID := pkg.STDUserID(ctx)

	flowInstanceEntity, err := i.instanceRepo.GetEntityByProcessInstanceID(i.db, processInstanceID)
	if err != nil {
		return false, err
	}
	if flowInstanceEntity == nil {
		return false, error2.NewErrorWithString(error2.Internal, "Can not find flow instance data ")
	}

	flowEntity, err := i.flowRepo.FindByID(i.db, flowInstanceEntity.FlowID)
	if err != nil {
		return false, err
	}
	if flowEntity == nil {
		return false, error2.NewErrorWithString(error2.Internal, "Can not find flow data ")
	}

	if !(flowInstanceEntity.Status == SendBack && flowInstanceEntity.CreatorID == userID) {
		return false, error2.NewErrorWithString(error2.Internal, "Can not resubmit ")
	}

	saveFormDataReq := &client.UpdateEntity{}
	formDataJSON, err := json.Marshal(params.FormData)
	if err == nil {
		err = json.Unmarshal(formDataJSON, saveFormDataReq)
	}
	if err != nil {
		return false, err
	}

	err = i.formAPI.UpdateData(ctx, flowInstanceEntity.AppID, flowInstanceEntity.FormID, flowInstanceEntity.FormInstanceID, *saveFormDataReq, false)
	if err != nil {
		return false, err
	}

	req := client.GetTasksReq{
		InstanceID: []string{processInstanceID},
		Desc:       []string{convert.SendBackTask},
		Status:     Active,
	}
	getTasksResp, err := i.processAPI.GetTasks(ctx, req)
	if err != nil {
		return false, err
	}
	if len(getTasksResp.Data) > 0 {
		params, err := i.GetInstanceVariableValues(ctx, flowInstanceEntity)
		if err != nil {
			return false, err
		}
		i.processAPI.CompleteTask(ctx, processInstanceID, getTasksResp.Data[0].ID, i.flow.FormatFormValue(flowInstanceEntity, params), nil)
	} else {
		return false, error2.NewErrorWithString(error2.Internal, "Can not resubmit task ")
	}

	updateMap := make(map[string]interface{}, 0)
	updateMap["modifier_id"] = userID
	updateMap["status"] = Review
	i.instanceRepo.Update(i.db, flowInstanceEntity.ID, updateMap)

	// Add operation record
	handleTaskModel := &models.HandleTaskModel{
		HandleType: opReSubmit,
		HandleDesc: "重新提交",
	}
	i.operationRecord.AddOperationRecord(ctx, flowInstanceEntity, nil, handleTaskModel)

	return true, nil
}

func (i *instance) FlowInfo(ctx context.Context, processInstanceID string) (*models.Flow, error) {
	flowInstanceEntity, err := i.instanceRepo.GetEntityByProcessInstanceID(i.db, processInstanceID)
	if err != nil {
		return nil, err
	}
	if flowInstanceEntity == nil {
		return nil, error2.NewErrorWithString(error2.Internal, "Can not find flow instance data ")
	}

	flowEntity, err := i.flowRepo.FindByID(i.db, flowInstanceEntity.FlowID)
	if err != nil {
		return nil, err
	}
	return flowEntity, nil
}

func (i *instance) SendBack(ctx context.Context, processInstanceID string, taskID string, model *models.HandleTaskModel) (bool, error) {
	userID := pkg.STDUserID(ctx)

	task, err := i.processAPI.CheckActiveTask(ctx, processInstanceID, taskID, userID)
	if err != nil {
		return false, err
	}
	if task == nil {
		return false, error2.NewErrorWithString(error2.Internal, "Can not find task ")
	}

	req := client.AddTaskReq{
		InstanceID: processInstanceID,
		TaskID:     taskID,
		UserID:     pkg.STDUserID(ctx),
		Name:       "打回重填",
		Desc:       convert.SendBackTask,
	}
	err = i.processAPI.SendBack(ctx, req)
	if err != nil {
		return false, err
	}

	flowInstanceEntity, err := i.instanceRepo.GetEntityByProcessInstanceID(i.db, processInstanceID)
	if err != nil {
		return false, err
	}
	if flowInstanceEntity == nil {
		return false, error2.NewErrorWithString(error2.Internal, "Can not find flow instance data ")
	}

	updateMap := make(map[string]interface{}, 0)
	updateMap["modifier_id"] = userID
	updateMap["status"] = SendBack
	i.instanceRepo.Update(i.db, flowInstanceEntity.ID, updateMap)

	// Add operation record
	handleTaskModel := &models.HandleTaskModel{
		HandleType: opSendBack,
		HandleDesc: "将工作流打回至发起人",
	}
	i.operationRecord.AddOperationRecord(ctx, flowInstanceEntity, task, handleTaskModel)

	return true, nil
}

func (i *instance) StepBack(ctx context.Context, processInstanceID string, taskID string, model *models.HandleTaskModel) (bool, error) {
	if len(model.TaskDefKey) == 0 {
		return false, error2.NewErrorWithString(error2.Internal, "Step back event is null ")
	}

	userID := pkg.STDUserID(ctx)

	flowInstanceEntity, err := i.instanceRepo.GetEntityByProcessInstanceID(i.db, processInstanceID)
	if err != nil {
		return false, err
	}

	task, err := i.processAPI.CheckActiveTask(ctx, processInstanceID, taskID, userID)
	if err != nil {
		return false, err
	}
	if task == nil {
		return false, error2.NewErrorWithString(error2.Internal, "Can not find task ")
	}

	toNode, err := i.processAPI.GetModelNode(ctx, task.ProcID, model.TaskDefKey)
	if err != nil {
		return false, error2.NewErrorWithString(error2.Internal, "Can not find event ")
	}

	err = i.processAPI.StepBack(ctx, processInstanceID, taskID, model.TaskDefKey)
	if err != nil {
		return false, err
	}

	// Add operation record
	model.HandleType = opStepBack
	model.HandleDesc = "将工作流回退至“" + toNode.Name + "”"

	i.operationRecord.AddOperationRecord(ctx, flowInstanceEntity, task, model)

	return true, nil
}

func (i *instance) CcFlow(ctx context.Context, processInstanceID string, taskID string, model *models.HandleTaskModel) (bool, error) {
	if len(model.HandleUserIDs) == 0 {
		return false, error2.NewErrorWithString(error2.Internal, "Handle user is requried ")
	}

	userID := pkg.STDUserID(ctx)
	if utils.Contain(model.HandleUserIDs, userID) {
		return false, error2.NewErrorWithString(code.CannotCcToSelf, code.CodeTable[code.CannotCcToSelf])
	}

	task, err := i.processAPI.CheckActiveTask(ctx, processInstanceID, taskID, userID)
	if err != nil {
		return false, err
	}
	if task == nil {
		return false, error2.NewErrorWithString(error2.Internal, "Can not find task ")
	}

	flowInstanceEntity, err := i.instanceRepo.GetEntityByProcessInstanceID(i.db, processInstanceID)
	if err != nil {
		return false, err
	}
	if flowInstanceEntity == nil {
		return false, error2.NewErrorWithString(error2.Internal, "Can not find flow instance data ")
	}

	dataMap, err := i.operationRecordRepo.GetHandleUserIDs(i.db, processInstanceID, taskID, opCC)
	if err != nil {
		return false, err
	}

	// tasks := make([]client.ProcessTask, 0)
	ccUsers := make([]string, 0)
	for _, handleUser := range model.HandleUserIDs {
		if !utils.Contain(dataMap, handleUser) {
			ccUsers = append(ccUsers, handleUser)
		}
	}
	if len(ccUsers) == 0 {
		return false, error2.NewErrorWithString(error2.Internal, "Can not repeat cc ")
	}

	req := &client.AddTaskReq{
		InstanceID: processInstanceID,
		NodeDefKey: task.NodeDefKey,
		Name:       task.Name,
		Desc:       convert.CcTask,
		TaskID:     taskID,
		Assignee:   ccUsers,
		UserID:     userID,
	}
	tasks, err := i.processAPI.AddNonModelTask(ctx, req)
	if err != nil {
		return false, err
	}

	model.HandleType = opCC
	model.HandleDesc = "抄送"

	records := make([]*models.OperationRecord, 0)
	operation := &models.OperationRecord{
		ProcessInstanceID: flowInstanceEntity.ProcessInstanceID,
		HandleUserID:      userID,
		HandleType:        opCC,
		HandleDesc:        model.HandleDesc,
		Remark:            model.Remark,
		Status:            Completed,
		TaskID:            task.ID,
		TaskName:          task.Name,
		TaskDefKey:        task.NodeDefKey,
		CorrelationData:   utils.ChangeStringArrayToString(model.HandleUserIDs),
		BaseModel: models.BaseModel{
			ID:         id2.GenID(),
			CreatorID:  userID,
			CreateTime: time2.Now(),
			ModifyTime: time2.Now(),
		},
	}
	records = append(records, operation)
	for _, newTask := range tasks {
		operation := &models.OperationRecord{
			ProcessInstanceID: flowInstanceEntity.ProcessInstanceID,
			HandleUserID:      newTask.Assignee,
			HandleType:        opHandleCC,
			HandleDesc:        model.HandleDesc,
			Status:            Active,
			TaskID:            newTask.ID,
			TaskName:          newTask.Name,
			TaskDefKey:        newTask.NodeDefKey,
			BaseModel: models.BaseModel{
				ID:         id2.GenID(),
				CreatorID:  userID,
				CreateTime: time2.Now(),
				ModifyTime: time2.Now(),
			},
		}
		records = append(records, operation)
	}

	i.operationRecord.AddOperationRecords(ctx, flowInstanceEntity, task, model, records)

	return true, nil
}

func (i *instance) ReadFlow(ctx context.Context, processInstanceID string, taskID string, model *models.HandleTaskModel) (bool, error) {
	userID := pkg.STDUserID(ctx)
	if utils.Contain(model.HandleUserIDs, userID) {
		return false, error2.NewErrorWithString(code.CannotInviteReadToSelf, code.CodeTable[code.CannotInviteReadToSelf])
	}

	task, err := i.processAPI.CheckActiveTask(ctx, processInstanceID, taskID, userID)
	if err != nil {
		return false, err
	}
	if task == nil {
		return false, error2.NewErrorWithString(error2.Internal, "Can not find task ")
	}

	flowInstanceEntity, err := i.instanceRepo.GetEntityByProcessInstanceID(i.db, processInstanceID)
	if err != nil {
		return false, err
	}
	if flowInstanceEntity == nil {
		return false, error2.NewErrorWithString(error2.Internal, "Can not find flow instance data ")
	}

	dataMap, err := i.operationRecordRepo.GetHandleUserIDs(i.db, processInstanceID, taskID, opRead)
	if err != nil {
		return false, err
	}

	// tasks := make([]client.ProcessTask, 0)

	ccUsers := make([]string, 0)
	for _, handleUser := range model.HandleUserIDs {
		if !utils.Contain(dataMap, handleUser) {
			ccUsers = append(ccUsers, handleUser)
		}
	}
	if len(ccUsers) == 0 {
		return false, error2.NewErrorWithString(code.CannotRepeatInviteRead, code.CodeTable[code.CannotRepeatInviteRead])
	}

	req := &client.AddTaskReq{
		InstanceID: processInstanceID,
		NodeDefKey: task.NodeDefKey,
		Name:       task.Name,
		Desc:       convert.ReadTask,
		TaskID:     taskID,
		Assignee:   ccUsers,
		UserID:     userID,
	}
	tasks, err := i.processAPI.AddNonModelTask(ctx, req)
	if err != nil {
		return false, err
	}

	model.HandleType = opRead
	model.HandleDesc = "邀请阅示"

	records := make([]*models.OperationRecord, 0)
	operation := &models.OperationRecord{
		ProcessInstanceID: flowInstanceEntity.ProcessInstanceID,
		HandleUserID:      userID,
		HandleType:        opRead,
		HandleDesc:        model.HandleDesc,
		Remark:            model.Remark,
		Status:            Completed,
		TaskID:            task.ID,
		TaskName:          task.Name,
		TaskDefKey:        task.NodeDefKey,
		CorrelationData:   utils.ChangeStringArrayToString(model.HandleUserIDs),
		BaseModel: models.BaseModel{
			ID:         id2.GenID(),
			CreatorID:  userID,
			CreateTime: time2.Now(),
			ModifyTime: time2.Now(),
		},
	}
	records = append(records, operation)
	for _, newTask := range tasks {
		operation := &models.OperationRecord{
			ProcessInstanceID: flowInstanceEntity.ProcessInstanceID,
			HandleUserID:      newTask.Assignee,
			HandleType:        opHandleRead,
			HandleDesc:        model.HandleDesc,
			Status:            Active,
			TaskID:            newTask.ID,
			TaskName:          newTask.Name,
			TaskDefKey:        newTask.NodeDefKey,
			BaseModel: models.BaseModel{
				ID:         id2.GenID(),
				CreatorID:  userID,
				CreateTime: time2.Now(),
				ModifyTime: time2.Now(),
			},
		}
		records = append(records, operation)
	}

	i.operationRecord.AddOperationRecords(ctx, flowInstanceEntity, task, model, records)

	return true, nil
}

func (i *instance) HandleCc(ctx context.Context, taskIDs []string) (bool, error) {
	if len(taskIDs) == 0 {
		return false, error2.NewErrorWithString(error2.Internal, "Task id requried ")
	}

	userID := pkg.STDUserID(ctx)

	err := i.processAPI.CompleteNonModelTask(ctx, taskIDs)
	if err != nil {
		return false, err
	}

	records, err := i.operationRecordRepo.FindRecordsByTaskIDs(i.db, opHandleCC, taskIDs)
	recordIDs := make([]string, 0)
	if len(records) > 0 {
		for _, value := range records {
			recordIDs = append(recordIDs, value.ID)
		}
		i.operationRecordRepo.UpdateStatus(i.db, recordIDs, Completed, userID)
	}

	return true, nil
}

func (i *instance) HandleRead(ctx context.Context, processInstanceID string, taskID string, model *models.HandleTaskModel) (bool, error) {
	if len(processInstanceID) == 0 {
		return false, error2.NewErrorWithString(error2.Internal, "Instance id requried ")
	}
	if len(taskID) == 0 {
		return false, error2.NewErrorWithString(error2.Internal, "Task id requried ")
	}

	userID := pkg.STDUserID(ctx)

	err := i.processAPI.CompleteNonModelTask(ctx, []string{taskID})
	if err != nil {
		return false, err
	}

	records, err := i.operationRecordRepo.FindRecordsByTaskIDs(i.db, opHandleRead, []string{taskID})
	recordIDs := make([]string, 0)
	if len(records) > 0 {
		for _, value := range records {
			recordIDs = append(recordIDs, value.ID)
		}
		i.operationRecordRepo.UpdateStatus2(i.db, recordIDs, Completed, userID, model.Remark)
	}

	return true, nil
}

func (i *instance) DeliverTask(ctx context.Context, processInstanceID string, taskID string, model *models.HandleTaskModel) (bool, error) {
	userID := pkg.STDUserID(ctx)
	if len(model.HandleUserIDs) != 1 {
		return false, error2.NewErrorWithString(error2.Internal, "Deliver user requried ")
	}
	if model.HandleUserIDs[0] == userID {
		return false, error2.NewErrorWithString(error2.Internal, "Can not deliver self ")
	}

	task, err := i.processAPI.CheckActiveTask(ctx, processInstanceID, taskID, userID)
	if err != nil {
		return false, err
	}
	if task == nil {
		return false, error2.NewErrorWithString(error2.Internal, "Can not find task ")
	}

	flowInstanceEntity, err := i.instanceRepo.GetEntityByProcessInstanceID(i.db, processInstanceID)
	if err != nil {
		return false, err
	}
	if flowInstanceEntity == nil {
		return false, error2.NewErrorWithString(error2.Internal, "Can not find flow instance data ")
	}

	req := &client.AddHistoryTaskReq{
		Name:     task.Name,
		Desc:     task.Desc,
		Assignee: userID,
		UserID:   userID,
		// TaskID:      task.ID,
		NodeDefKey:  task.NodeDefKey,
		InstanceID:  task.ProcInstanceID,
		ExecutionID: task.ExecutionID,
	}
	_, err = i.processAPI.AddHistoryTask(ctx, req)
	if err != nil {
		return false, err
	}

	err = i.processAPI.SetAssignee(ctx, processInstanceID, taskID, model.HandleUserIDs[0])
	if err != nil {
		return false, err
	}

	// Add operation record
	model.HandleType = opDeliver
	model.HandleDesc = "转交"
	i.operationRecord.AddOperationRecord(ctx, flowInstanceEntity, task, model)

	return true, nil
}

func (i *instance) StepBackNodes(ctx context.Context, processInstanceID string) ([]*models.NodeModel, error) {
	tasksReq := client.GetTasksReq{
		InstanceID: []string{processInstanceID},
	}
	getTasksResp, err := i.processAPI.GetTasks(ctx, tasksReq)
	if err != nil {
		return nil, err
	}
	if len(getTasksResp.Data) > 0 {

		preNodesResp, err := i.processAPI.GetPreNode(ctx, getTasksResp.Data[0].ID)
		if err != nil {
			return nil, err
		}

		nodes := make([]*models.NodeModel, 0)
		for _, value := range preNodesResp.Nodes {
			if value.NodeType != convert.User && value.NodeType != convert.MultiUser {
				continue
			}
			node := &models.NodeModel{
				TaskDefKey: value.NodeDefKey,
				TaskName:   value.Name,
			}
			nodes = append(nodes, node)
		}
		return nodes, nil
	}
	return nil, nil
}

func (i *instance) FlowInstanceCount(ctx context.Context) (*InstanceCountModel, error) {
	userID := pkg.STDUserID(ctx)

	// 待办数量
	tasksReq := client.GetTasksReq{
		Assignee: userID,
		Desc:     WaitHandleTasks(),
	}
	waitHandleCount, _ := i.processAPI.GetTasksCount(ctx, tasksReq)

	// 抄送给我数量
	tasksReq = client.GetTasksReq{
		Assignee: userID,
		Desc:     []string{opCC},
	}
	ccToMeCount, _ := i.processAPI.GetTasksCount(ctx, tasksReq)

	// 已超时数量
	tasksReq = client.GetTasksReq{
		Assignee: userID,
		Desc:     WaitHandleTasks(),
		DueTime:  time2.Now(),
	}
	overTimeCount, _ := i.processAPI.GetTasksCount(ctx, tasksReq)

	// 催办数量
	taskIDs, _ := i.urgeRepo.FindTaskIDs(i.db)
	urgeCount := int64(0)
	if len(taskIDs) > 0 {
		tasksReq = client.GetTasksReq{
			Assignee: userID,
			Desc:     WaitHandleTasks(),
			TaskID:   taskIDs,
		}
		urgeCount, _ = i.processAPI.GetTasksCount(ctx, tasksReq)
	}

	model := &InstanceCountModel{
		OverTimeCount:   overTimeCount,
		UrgeCount:       urgeCount,
		WaitHandleCount: waitHandleCount,
		CcToMeCount:     ccToMeCount,
	}

	return model, nil
}

func (i *instance) AddSign(ctx context.Context, processInstanceID string, taskID string, model *models.AddSignModel) (bool, error) {
	userID := pkg.STDUserID(ctx)

	if len(model.Assignee) == 0 {
		return false, error2.NewErrorWithString(error2.Internal, "Add sign user is empty ")
	}

	flowInstanceEntity, err := i.instanceRepo.GetEntityByProcessInstanceID(i.db, processInstanceID)
	if err != nil {
		return false, err
	}

	task, err := i.processAPI.CheckActiveTask(ctx, processInstanceID, taskID, userID)
	if err != nil {
		return false, err
	}
	if task == nil {
		return false, error2.NewErrorWithString(error2.Internal, "Can not find task ")
	}

	assignees := make([]string, 0)
	for _, value := range model.Assignee {
		valueMap := utils.ChangeObjectToMap(value)
		assignees = append(assignees, valueMap["id"].(string))
	}

	relNodeDefKey := task.NodeDefKey + ":" + id2.GenID() // currentNodeID+":"+newNodeID
	req := &client.AddTaskReq{
		TaskID: taskID,
		UserID: userID,
		Node: &client.NodeData{
			DefKey:  relNodeDefKey,
			Name:    "加签",
			Type:    convert.User,
			Desc:    convert.ReviewTask,
			UserIDs: assignees,
		},
	}
	if model.MultiplePersonWay == "or" {
		req.Node.Type = convert.User
	} else if model.MultiplePersonWay == "and" {
		req.Node.Type = convert.MultiUser
	}

	if model.Type == "BEFORE" {
		_, err := i.processAPI.AddBeforeModelTask(ctx, req)
		if err != nil {
			return false, err
		}
	} else if model.Type == "AFTER" {
		_, err := i.processAPI.AddAfterModelTask(ctx, req)
		if err != nil {
			return false, err
		}
	} else {
		return false, error2.NewErrorWithString(error2.Internal, "Add type is error ")
	}

	// Add operation record
	handleTaskModel := &models.HandleTaskModel{
		HandleType:    opAddSign,
		HandleDesc:    "加签",
		RelNodeDefKey: relNodeDefKey,
		HandleUserIDs: assignees,
	}
	relMap := make(map[string]string)
	relMap["multiplePersonWay"] = model.MultiplePersonWay
	relMap["assignees"] = strings.Join(assignees, ",")
	relMapByte, err := json.Marshal(relMap)
	if err != nil {
		return false, err
	}
	records := make([]*models.OperationRecord, 0)
	operation := &models.OperationRecord{
		ProcessInstanceID: flowInstanceEntity.ProcessInstanceID,
		HandleUserID:      userID,
		HandleType:        opAddSign,
		Status:            Completed,
		TaskID:            task.ID,
		TaskName:          task.Name,
		TaskDefKey:        task.NodeDefKey,
		RelNodeDefKey:     relNodeDefKey,
		CorrelationData:   string(relMapByte),
		BaseModel: models.BaseModel{
			ID:         id2.GenID(),
			CreatorID:  userID,
			CreateTime: time2.Now(),
			ModifyTime: time2.Now(),
		},
	}
	records = append(records, operation)
	i.operationRecord.AddOperationRecords(ctx, flowInstanceEntity, task, handleTaskModel, records)
	return true, nil
}

func (i *instance) ReviewTask(ctx context.Context, processInstanceID string, taskID string, model *models.HandleTaskModel) (bool, error) {
	userID := pkg.STDUserID(ctx)

	if !isReviewStatus(model.HandleType) {
		return false, error2.NewErrorWithString(error2.Internal, "Handle type must be agree、refuse、fillIn ")
	}

	task, err := i.processAPI.CheckActiveTask(ctx, processInstanceID, taskID, userID)
	if err != nil {
		return false, err
	}
	if task == nil {
		return false, error2.NewError(code.TaskCannotFind)
		// return false, error2.NewErrorWithString(error2.Internal, "Can not find task ")
	}

	flowInstanceEntity, err := i.instanceRepo.GetEntityByProcessInstanceID(i.db, processInstanceID)
	if err != nil {
		return false, err
	}
	if flowInstanceEntity == nil {
		return false, error2.NewErrorWithString(error2.Internal, "Can not find flow instance data ")
	}

	entity, err := i.flowRepo.FindByID(i.db, flowInstanceEntity.FlowID)
	if err != nil {
		return false, err
	}
	if entity == nil {
		return false, error2.NewErrorWithString(error2.Internal, "Can not find flow data ")
	}

	comments := map[string]interface{}{
		"reviewResult": model.HandleType,
		"reviewRemark": model.Remark,
	}

	// 保存数据到表单接口,（将表达式替换成真实值，权限判断和数值拼装）
	if model.FormData != nil {
		formData := i.task.FilterCanEditFormData(ctx, entity, flowInstanceEntity, task.NodeDefKey, model.FormData)
		if len(formData) > 0 {

			saveFormDataReq := &client.UpdateEntity{}
			formDataJSON, err := json.Marshal(formData)
			if err == nil {
				err = json.Unmarshal(formDataJSON, saveFormDataReq)
			}
			if err != nil {
				return false, err
			}

			err = i.formAPI.UpdateData(ctx, flowInstanceEntity.AppID, flowInstanceEntity.FormID, flowInstanceEntity.FormInstanceID, *saveFormDataReq, false)
			if err != nil {
				return false, err
			}
		}
	}

	params, err := i.GetInstanceVariableValues(ctx, flowInstanceEntity)
	if err != nil {
		return false, err
	}

	status := ""
	if model.HandleType == Refuse { // 拒绝
		shapeModel, err := convert.GetShapeByTaskDefKey(entity.BpmnText, task.NodeDefKey)
		if err != nil {
			return false, err
		}
		branchTargetElementID := shapeModel.Data.NodeData.BranchTargetElementID // 合流节点id
		if len(branchTargetElementID) == 0 {                                    // 不在分支中，则直接结束流程
			err = i.processAPI.CompleteTask(ctx, processInstanceID, taskID, i.flow.FormatFormValue(flowInstanceEntity, params), comments)
			if err != nil {
				return false, err
			}
			err = i.processAPI.AbendInstance(ctx, processInstanceID)
			if err != nil {
				return false, err
			}
			status = Refuse
		} else { // 在分支中，需要判断合流设置的逻辑

			status, err = i.BranchReject0(ctx, shapeModel, task, entity, processInstanceID)

			if err != nil {
				return false, err
			}
		}

	} else {
		err = i.processAPI.CompleteTask(ctx, processInstanceID, taskID, i.flow.FormatFormValue(flowInstanceEntity, params), comments)
		if err != nil {
			return false, err
		}
	}

	// 增加操作日志
	i.operationRecord.AddOperationRecord(ctx, flowInstanceEntity, task, model)

	dataMap := make(map[string]interface{})
	if len(status) > 0 {
		dataMap["status"] = status
	} else {
		// tasks, err := i.processAPI.GetTasksByInstanceID(ctx, processInstanceID)
		processInstance, err := i.processAPI.GetInstanceByID(ctx, processInstanceID)
		if err != nil {
			return false, err
		}
		if processInstance != nil && processInstance.Status != Active {
			if opFillIn == model.HandleType {
				dataMap["status"] = opAgree
			} else {
				dataMap["status"] = model.HandleType
			}
		} else {
			// i.task.TaskCheck(ctx, entity, flowInstanceEntity, userID)
			processInstance, err = i.processAPI.GetInstanceByID(ctx, processInstanceID)
			if err != nil {
				return false, err
			}
			if processInstance != nil && processInstance.Status != Active {
				// 更新flow_instance_step
				err = i.operationRecord.UpdateByNodeInstanceID(ctx, entity, task)
				if err != nil {
					return false, err
				}
				if opFillIn == model.HandleType {
					dataMap["status"] = opAgree
				} else {
					dataMap["status"] = model.HandleType
				}
			} else {
				dataMap["status"] = InReview
			}
		}
	}

	dataMap["modifier_id"] = userID
	err = i.instanceRepo.Update(i.db, flowInstanceEntity.ID, dataMap)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (i *instance) BranchReject0(ctx context.Context, currentNode *convert.ShapeModel, task *client.ProcessTask,
	flowEntity *models.Flow, processInstanceID string) (string, error) {

	nextNodeDefKey, isProcessCompleted := convert.FindNextNode(currentNode, flowEntity.BpmnText)

	// 获取当前任务的同级 executionIDs
	gateWayExecutionReq := &client.GateWayExecutionReq{
		TaskID: task.ID,
	}
	executionResp, _ := i.processAPI.GetExecution(ctx, gateWayExecutionReq)

	// 判断合流节点的设置
	branchTargetElementID := currentNode.Data.NodeData.BranchTargetElementID // 合流节点id
	branchTargetElement, _ := convert.GetShapeByTaskDefKey(flowEntity.BpmnText, branchTargetElementID)

	// true : 任一分支拒绝结束流程 false：所有分支拒绝结束流程
	processBranchEndStrategy := convert.GetValueFromBusinessData(*branchTargetElement, "processBranchEndStrategy").(string) == "any"
	gatewayExecutionIDs := make([]string, 0)
	if processBranchEndStrategy {
		gatewayExecutionIDs = append(executionResp.Executions, executionResp.CurrentExecutionID)
	} else {
		gatewayExecutionIDs = append(gatewayExecutionIDs, executionResp.CurrentExecutionID)
		instanceExecutions, _ := i.instanceExecutionRepo.FindByExecutionIDs(i.db, executionResp.Executions)
		flag := len(instanceExecutions) == len(executionResp.Executions)
		if flag {
			nextNodeDefKey = currentNode.Data.NodeData.BranchTargetElementID
			nextNode, _ := convert.GetShapeByTaskDefKey(flowEntity.BpmnText, nextNodeDefKey)
			nextNodeDefKey, isProcessCompleted = convert.FindNextNode(nextNode, flowEntity.BpmnText)
		}

	}

	resp, err := i.processAPI.InclusiveExecution(ctx, client.ParentExecutionReq{
		TaskID: task.ID,
		DefKey: nextNodeDefKey,
	})
	if err != nil {
		return "", err
	}

	comments := map[string]interface{}{
		"reviewResult": Refuse,
		"reviewRemark": "",
	}
	commentsJSON, err := json.Marshal(comments)
	if err != nil {
		return "", err
	}
	commentsStr := string(commentsJSON)

	req := &client.CompleteExecutionReq{
		ExecutionID: gatewayExecutionIDs,
		NextDefKey:  nextNodeDefKey,
		TaskID:      task.ID,
		UserID:      pkg.STDUserID(ctx),
		Comments:    commentsStr,
	}
	_, err = i.processAPI.CompleteExecution(ctx, req)
	if err != nil {
		return "", err
	}

	if isProcessCompleted {
		return Refuse, nil
	}
	// 记录分支的结果
	instanceExecution := &models.InstanceExecution{
		ProcessInstanceID: processInstanceID,
		ExecutionID:       resp.ExecutionID, // 如果是会签需要使用父级executionID
		Result:            Refuse,
	}
	i.instanceExecutionRepo.Create(i.db, instanceExecution)
	return "", nil
}

// BranchReject 审核拒绝检查合流配置规则
func (i *instance) BranchReject(ctx context.Context, shapeModel *convert.ShapeModel, task *client.ProcessTask,
	flowEntity *models.Flow, processInstanceID string, params map[string]interface{}) (string, error) {
	branchTargetElementID := shapeModel.Data.NodeData.BranchTargetElementID // 合流节点id

	gateWayExecutionReq := &client.GateWayExecutionReq{
		TaskID: task.ID,
	}
	executionResp, _ := i.processAPI.GetExecution(ctx, gateWayExecutionReq)
	currentTaskBranchExecution := executionResp.CurrentExecutionID // 如果是会签需要使用父级executionID
	instanceExecution := &models.InstanceExecution{
		ProcessInstanceID: processInstanceID,
		ExecutionID:       currentTaskBranchExecution,
		Result:            Refuse,
	}
	i.instanceExecutionRepo.Create(i.db, instanceExecution)

	gatewayExecutionIDs := append(executionResp.Executions, executionResp.CurrentExecutionID) // 同级分支所有executionID

	branchTargetElement, _ := convert.GetShapeByTaskDefKey(flowEntity.BpmnText, branchTargetElementID)
	outBranchTargetElementID := branchTargetElement.Data.NodeData.BranchTargetElementID

	comments := map[string]interface{}{
		"reviewResult": Refuse,
		"reviewRemark": "",
	}

	commentsJSON, err := json.Marshal(comments)
	if err != nil {
		return "", err
	}
	commentsStr := string(commentsJSON)

	// true : 任一分支拒绝结束流程 false：所有分支拒绝结束流程
	processBranchEndStrategy := convert.GetValueFromBusinessData(*branchTargetElement, "processBranchEndStrategy").(string) == "any"
	if processBranchEndStrategy {
		if len(outBranchTargetElementID) == 0 { // 已经是最外层合流
			i.processAPI.CompleteTask(ctx, processInstanceID, task.ID, params, comments)
			i.processAPI.AbendInstance(ctx, processInstanceID)
			return Refuse, nil
		}

		// 结束其他分支，直接合流
		req := &client.CompleteExecutionReq{
			ExecutionID: executionResp.Executions,
			Comments:    commentsStr,
		}
		i.processAPI.CompleteExecution(ctx, req)
		i.processAPI.CompleteTaskToNode(ctx, processInstanceID, task.ID, params, branchTargetElementID, comments)
		return "", nil

	}

	instanceExecutions, _ := i.instanceExecutionRepo.FindByExecutionIDs(i.db, gatewayExecutionIDs)
	flag := len(instanceExecutions) == len(gatewayExecutionIDs)

	// 其他execution都是拒绝，则结束
	if flag && len(outBranchTargetElementID) == 0 { // 其他execution都是拒绝, 并且已经是最外层合流，则结束流程
		i.processAPI.CompleteTask(ctx, processInstanceID, task.ID, params, comments)
		i.processAPI.AbendInstance(ctx, processInstanceID)
		return Refuse, nil
	}
	// 结束分支
	i.processAPI.CompleteTaskToNode(ctx, processInstanceID, task.ID, params, branchTargetElementID, comments)
	return "", nil
}

func (i *instance) GetFormData(ctx context.Context, processInstanceID string, taskID string, req *GetFormDataReq) (interface{}, error) {
	userID := pkg.STDUserID(ctx)
	flowInstanceEntity, err := i.instanceRepo.GetEntityByProcessInstanceID(i.db, processInstanceID)
	if err != nil {
		return nil, err
	}
	if taskID == "0" || (flowInstanceEntity.CreatorID == userID && flowInstanceEntity.Status == opSendBack) {
		getFormDataReq := client.FormDataConditionModel{
			AppID:   flowInstanceEntity.AppID,
			TableID: flowInstanceEntity.FormID,
			DataID:  flowInstanceEntity.FormInstanceID,
			Ref:     req.Ref,
		}

		formData, err := i.formAPI.GetFormData(ctx, getFormDataReq)
		return formData, err
	}

	flowEntity, err := i.flowRepo.FindByID(i.db, flowInstanceEntity.FlowID)
	if err != nil {
		return nil, err
	}

	abnormalTasks, err := i.abnormalTaskRepo.Find(i.db, map[string]interface{}{
		"task_id": taskID,
		"status":  0,
	})
	if len(abnormalTasks) > 0 {
		userID = ""
	}
	tasksReq := client.GetTasksReq{
		InstanceID: []string{flowInstanceEntity.ProcessInstanceID},
		Assignee:   userID,
		TaskID:     []string{taskID},
	}
	if len(taskID) > 0 && taskID != "0" {
		tasksReq.TaskID = []string{taskID}
	}
	taskResp, _ := i.processAPI.GetAllTasks(ctx, tasksReq)
	tasks := taskResp.Data

	var task *client.ProcessTask
	if len(tasks) > 0 {
		task = tasks[0]
	} else {
		return nil, nil
	}

	taskDetailModel := &TaskDetailModel{
		TaskDefKey: task.NodeDefKey,
	}
	var fieldPermissionModel interface{}

	shapeModel, err := convert.GetShapeByTaskDefKey(flowEntity.BpmnText, taskDetailModel.TaskDefKey)
	if err != nil {
		return nil, err
	}
	if shapeModel != nil {
		fieldPermissionModel = shapeModel.Data.BusinessData["fieldPermission"]
		taskDetailModel.FieldPermission = fieldPermissionModel
	}
	getFormDataReq := client.FormDataConditionModel{
		AppID:   flowInstanceEntity.AppID,
		TableID: flowInstanceEntity.FormID,
		DataID:  flowInstanceEntity.FormInstanceID,
		Ref:     req.Ref,
	}

	formData, _ := i.formAPI.GetFormData(ctx, getFormDataReq)

	// taskDetailModel.FormData = formData
	return i.task.FilterCanReadFormData(ctx, flowInstanceEntity, taskDetailModel.FieldPermission, formData), nil
}

func (i *instance) GetFlowInstanceForm(ctx context.Context, processInstanceID string, taskTypeDetailModel *TaskTypeDetailModel) (*InstanceDetailModel, error) {
	flowInstanceEntity, err := i.instanceRepo.GetEntityByProcessInstanceID(i.db, processInstanceID)
	if err != nil {
		return nil, err
	}
	if flowInstanceEntity == nil {
		return nil, error2.NewErrorWithString(error2.Internal, "Can not find flow instance data ")
	}

	flowEntity, err := i.flowRepo.FindByID(i.db, flowInstanceEntity.FlowID)
	if err != nil {
		return nil, err
	}
	if flowEntity == nil {
		return nil, error2.NewErrorWithString(error2.Internal, "Can not find flow data ")
	}

	flowInstanceDetailModel := &InstanceDetailModel{
		FlowName:            flowEntity.Name,
		CanMsg:              flowEntity.CanMsg == 1,
		CanViewStatusAndMsg: flowEntity.CanViewStatusMsg == 1,
		AppID:               flowEntity.AppID,
		TableID:             flowEntity.FormID,
	}

	taskDetailModels := make([]*TaskDetailModel, 0)
	if "APPLY_PAGE" == taskTypeDetailModel.Type { // 我发起的流程详情页面
		taskDetailModel := i.getApplyPageTaskDetailModel(ctx, flowEntity, flowInstanceEntity)
		if taskDetailModel == nil {
			return nil, error2.NewErrorWithString(error2.Internal, "This flow is not apply by yourself ")
		}

		taskDetailModels = append(taskDetailModels, taskDetailModel)
		flowInstanceDetailModel.TaskDetailModels = taskDetailModels
	} else if "WAIT_HANDLE_PAGE" == taskTypeDetailModel.Type { // 待我处理
		tempTaskDetailModels := i.getWaitHandlePageTaskDetailModel(ctx, flowEntity, flowInstanceEntity, taskTypeDetailModel.TaskID)

		if len(tempTaskDetailModels) == 0 {
			return nil, error2.NewErrorWithString(error2.Internal, "Can not find task data ")
		}

		flowInstanceDetailModel.TaskDetailModels = tempTaskDetailModels
	} else if "HANDLED_PAGE" == taskTypeDetailModel.Type { // 我已处理的流程详情页面
		tempTaskDetailModels := i.getHandledPageTaskDetailModel(ctx, flowEntity, flowInstanceEntity, taskTypeDetailModel.TaskID)

		if len(tempTaskDetailModels) == 0 {
			return nil, error2.NewErrorWithString(error2.Internal, "Can not find task data ")
		}

		flowInstanceDetailModel.TaskDetailModels = tempTaskDetailModels
	} else if "CC_PAGE" == taskTypeDetailModel.Type { // 抄送给我
		tempTaskDetailModels := i.getCcPageTaskDetailModel(ctx, flowEntity, flowInstanceEntity, taskTypeDetailModel.TaskID)

		if len(tempTaskDetailModels) == 0 {
			return nil, error2.NewErrorWithString(error2.Internal, "Can not find task data ")
		}

		flowInstanceDetailModel.TaskDetailModels = tempTaskDetailModels
	} else if "ALL_PAGE" == taskTypeDetailModel.Type {
		tempTaskDetailModels := i.getAllPageTaskDetailModel(ctx, flowEntity, flowInstanceEntity, taskTypeDetailModel.TaskID)
		if len(tempTaskDetailModels) == 0 {
			return nil, error2.NewErrorWithString(error2.Internal, "Can not find task data ")
		}

		flowInstanceDetailModel.TaskDetailModels = tempTaskDetailModels
	}

	return flowInstanceDetailModel, nil
}

func (i *instance) getApplyPageTaskDetailModel(ctx context.Context, flowEntity *models.Flow, flowInstanceEntity *models.Instance) *TaskDetailModel {
	userID := pkg.STDUserID(ctx)

	if flowInstanceEntity.CreatorID != userID {
		return nil
	}

	taskDetailModel := &TaskDetailModel{
		TaskID: "0",
	}

	formSchema, _ := i.formAPI.GetFormSchema(ctx, flowInstanceEntity.AppID, flowInstanceEntity.FormID)
	taskDetailModel.FormSchema = formSchema

	if opSendBack == flowInstanceEntity.Status {
		taskDetailModel.HasResubmitBtn = true
	}

	tasks, _ := i.processAPI.GetTasksByInstanceID(ctx, flowInstanceEntity.ProcessInstanceID)
	if i.flow.checkCanCancel(ctx, flowEntity, flowInstanceEntity, tasks) && isFlowOngoing(flowInstanceEntity.Status) {
		taskDetailModel.HasCancelBtn = true
	}
	if flowEntity.CanUrge == 1 && isFlowOngoing(flowInstanceEntity.Status) {
		taskDetailModel.HasUrgeBtn = true
	}
	return taskDetailModel
}

func (i *instance) getWaitHandlePageTaskDetailModel(ctx context.Context, flowEntity *models.Flow, flowInstanceEntity *models.Instance, taskID string) []*TaskDetailModel {
	userID := pkg.STDUserID(ctx)

	taskDetailModels := make([]*TaskDetailModel, 0)

	// 待办
	waitTasks := make([]*client.ProcessTask, 0)
	if len(taskID) > 0 {
		task, _ := i.processAPI.CheckActiveTask(ctx, flowInstanceEntity.ProcessInstanceID, taskID, userID)
		waitTasks = append(waitTasks, task)
	} else {
		tasksReq := client.GetTasksReq{
			InstanceID: []string{flowInstanceEntity.ProcessInstanceID},
			Assignee:   userID,
			Desc:       WaitHandleTasks(),
		}
		taskResp, _ := i.processAPI.GetTasks(ctx, tasksReq)
		waitTasks = append(waitTasks, taskResp.Data...)
	}

	taskMap := make(map[string]*TaskDetailModel, 0)
	if len(waitTasks) > 0 {
		for _, task := range waitTasks {
			taskDetailModel := &TaskDetailModel{
				TaskID:     task.ID,
				TaskName:   task.Name,
				TaskDefKey: task.NodeDefKey,
				TaskType:   task.Desc,
			}
			if taskDetailModel.TaskType == opRead {
				taskDetailModel.HasReadHandleBtn = true
			}

			taskDetailModels = append(taskDetailModels, taskDetailModel)
			taskMap[task.ID] = taskDetailModel
		}
	}

	if len(taskDetailModels) > 0 {
		var taskDetailModel *TaskDetailModel
		if len(taskID) > 0 {
			taskDetailModel = taskMap[taskID]
		} else {
			taskDetailModel = taskDetailModels[0]
		}

		hasBtn := taskDetailModel.TaskType == convert.ReviewTask || taskDetailModel.TaskType == convert.WriteTask
		i.getTaskDetailModel(ctx, flowInstanceEntity, flowEntity, taskDetailModel, hasBtn)
	}

	return taskDetailModels
}

func (i *instance) getHandledPageTaskDetailModel(ctx context.Context, flowEntity *models.Flow, flowInstanceEntity *models.Instance, taskID string) []*TaskDetailModel {
	userID := pkg.STDUserID(ctx)

	taskDetailModels := make([]*TaskDetailModel, 0)

	tasksReq := client.GetTasksReq{
		InstanceID: []string{flowInstanceEntity.ProcessInstanceID},
		Assignee:   userID,
		Desc:       []string{convert.ReviewTask, convert.WriteTask, convert.ReadTask},
	}

	taskResp, _ := i.processAPI.GetHistoryTasks(ctx, tasksReq)
	tasks := taskResp.Data

	taskMap := make(map[string]*TaskDetailModel, 0)
	if len(tasks) > 0 {
		for _, task := range tasks {
			taskDetailModel := &TaskDetailModel{
				TaskID:     task.ID,
				TaskName:   task.Name,
				TaskDefKey: task.NodeDefKey,
				TaskType:   task.Desc,
			}

			taskDetailModels = append(taskDetailModels, taskDetailModel)
			taskMap[task.ID] = taskDetailModel
		}
	}

	if len(taskDetailModels) > 0 {
		var taskDetailModel *TaskDetailModel
		if len(taskID) > 0 {
			taskDetailModel = taskMap[taskID]
		} else {
			taskDetailModel = taskDetailModels[0]
		}

		i.getTaskDetailModel(ctx, flowInstanceEntity, flowEntity, taskDetailModel, false)
	}

	return taskDetailModels
}

func (i *instance) getCcPageTaskDetailModel(ctx context.Context, flowEntity *models.Flow, flowInstanceEntity *models.Instance, taskID string) []*TaskDetailModel {
	userID := pkg.STDUserID(ctx)

	taskDetailModels := make([]*TaskDetailModel, 0)
	tasksReq := client.GetTasksReq{
		InstanceID: []string{flowInstanceEntity.ProcessInstanceID},
		Assignee:   userID,
		Desc:       []string{convert.CcTask},
	}
	if len(taskID) > 0 {
		tasksReq.TaskID = []string{taskID}
	}
	taskResp, _ := i.processAPI.GetAllTasks(ctx, tasksReq)
	tasks := taskResp.Data

	taskMap := make(map[string]*TaskDetailModel, 0)
	if len(tasks) > 0 {
		for _, task := range tasks {
			taskDetailModel := &TaskDetailModel{
				TaskID:     task.ID,
				TaskName:   task.Name,
				TaskDefKey: task.NodeDefKey,
				TaskType:   task.Desc,
			}
			taskDetailModel.HasCcHandleBtn = task.Status == Active

			taskDetailModels = append(taskDetailModels, taskDetailModel)
			taskMap[task.ID] = taskDetailModel
		}
	}

	if len(taskDetailModels) > 0 {
		var taskDetailModel *TaskDetailModel
		if len(taskID) > 0 {
			taskDetailModel = taskMap[taskID]
		} else {
			taskDetailModel = taskDetailModels[0]
		}

		i.getTaskDetailModel(ctx, flowInstanceEntity, flowEntity, taskDetailModel, false)
	}

	return taskDetailModels
}

func (i *instance) getAllPageTaskDetailModel(ctx context.Context, flowEntity *models.Flow, flowInstanceEntity *models.Instance, taskID string) []*TaskDetailModel {
	userID := pkg.STDUserID(ctx)

	processModel := &convert.ProcessModel{}
	err := json.Unmarshal([]byte(flowEntity.BpmnText), processModel)
	if err != nil {
		return nil
	}

	taskDetailModels := make([]*TaskDetailModel, 0)
	taskMap := make(map[string]*TaskDetailModel, 0)

	tasksReq := client.GetTasksReq{
		InstanceID: []string{flowInstanceEntity.ProcessInstanceID},
		Assignee:   userID,
	}

	taskResp, _ := i.processAPI.GetAllTasks(ctx, tasksReq)
	tasks := taskResp.Data

	if len(tasks) > 0 {
		for _, task := range tasks {
			taskDetailModel := &TaskDetailModel{
				TaskID:     task.ID,
				TaskName:   task.Name,
				TaskDefKey: task.NodeDefKey,
				TaskType:   task.Desc,
			}
			if task.EndTime == "" {
				if task.Desc == opCC {
					taskDetailModel.HasCcHandleBtn = true
				} else if task.Desc == opRead {
					taskDetailModel.HasReadHandleBtn = true
				}
			}

			taskDetailModels = append(taskDetailModels, taskDetailModel)
			taskMap[task.ID] = taskDetailModel
		}
	}

	// 发起
	if flowInstanceEntity.CreatorID == userID {
		nodeDataModel := processModel.Shapes[0].Data.NodeData
		taskDetailModel := &TaskDetailModel{
			TaskID:     "0",
			TaskName:   nodeDataModel.Name,
			TaskDefKey: processModel.Shapes[0].ID,
		}
		taskDetailModels = append(taskDetailModels, taskDetailModel)
		taskMap[taskDetailModel.TaskID] = taskDetailModel
	}

	if len(taskDetailModels) > 0 {
		var taskDetailModel *TaskDetailModel
		if len(taskID) > 0 {
			taskDetailModel = taskMap[taskID]
		} else {
			taskDetailModel = taskDetailModels[0]
		}

		hasBtn := taskDetailModel.TaskType == convert.ReviewTask || taskDetailModel.TaskType == convert.WriteTask
		i.getTaskDetailModel(ctx, flowInstanceEntity, flowEntity, taskDetailModel, hasBtn)
	}

	return taskDetailModels
}

func (i *instance) getTaskDetailModel(ctx context.Context, flowInstanceEntity *models.Instance, flowEntity *models.Flow, taskDetailModel *TaskDetailModel, hasBtn bool) {
	userID := pkg.STDUserID(ctx)
	if taskDetailModel.TaskID == "0" || (flowInstanceEntity.CreatorID == userID && flowInstanceEntity.Status == opSendBack) { // 开始节点
		formSchema, _ := i.formAPI.GetFormSchema(ctx, flowInstanceEntity.AppID, flowInstanceEntity.FormID)
		taskDetailModel.FormSchema = formSchema

		if flowInstanceEntity.Status == opSendBack {
			taskDetailModel.HasResubmitBtn = true
		}
		tasks, _ := i.processAPI.GetTasksByInstanceID(ctx, flowInstanceEntity.ProcessInstanceID)

		if i.flow.checkCanCancel(ctx, flowEntity, flowInstanceEntity, tasks) && isFlowOngoing(flowInstanceEntity.Status) {
			taskDetailModel.HasCancelBtn = true
		}
		if flowEntity.CanUrge == 1 && isFlowOngoing(flowInstanceEntity.Status) {
			taskDetailModel.HasUrgeBtn = true
		}
		return
	}

	operatorPermissionModel := &convert.OperatorPermissionModel{}
	var fieldPermissionModel interface{}

	shapeModel, err := convert.GetShapeByTaskDefKey(flowEntity.BpmnText, taskDetailModel.TaskDefKey)
	if err != nil {
		return
	}
	if shapeModel != nil {
		fieldPermissionModel = shapeModel.Data.BusinessData["fieldPermission"]
		taskDetailModel.FieldPermission = fieldPermissionModel

		if hasBtn {
			operatorPermissionObj := shapeModel.Data.BusinessData["operatorPermission"]
			operatorPermission, err := json.Marshal(operatorPermissionObj)
			err = json.Unmarshal(operatorPermission, operatorPermissionModel)
			if err != nil {
				return
			}
		}
	}

	formSchema, _ := i.formAPI.GetFormSchema(ctx, flowInstanceEntity.AppID, flowInstanceEntity.FormID)

	// 将节点字段权限赋值到schema中去
	formSchema = i.setPermissionToSchema(&formSchema, fieldPermissionModel)

	taskDetailModel.FormSchema = formSchema
	// i.task.FilterCanReadFormData(ctx, flowInstanceEntity, taskDetailModel.FieldPermission,)

	if hasBtn && operatorPermissionModel != nil {
		// 判断按钮权限
		if len(operatorPermissionModel.Custom) > 0 {
			for index := 0; index < len(operatorPermissionModel.Custom); {
				if !operatorPermissionModel.Custom[index].Enabled {
					operatorPermissionModel.Custom = append(operatorPermissionModel.Custom[:index], operatorPermissionModel.Custom[index+1:]...)
				} else {
					if operatorPermissionModel.Custom[index].Value == opStepBack {
						// 判断是不是首节点
						nodes, err := i.StepBackNodes(ctx, flowInstanceEntity.ProcessInstanceID)
						if err != nil {
							return
						}
						if nodes == nil || len(nodes) == 0 {
							operatorPermissionModel.Custom = append(operatorPermissionModel.Custom[:index], operatorPermissionModel.Custom[index+1:]...)
						}
					}
					index++
				}
			}
		}
		if len(operatorPermissionModel.System) > 0 {
			for i := 0; i < len(operatorPermissionModel.System); {
				if !operatorPermissionModel.System[i].Enabled {
					operatorPermissionModel.System = append(operatorPermissionModel.System[:i], operatorPermissionModel.System[i+1:]...)
				} else {
					i++
				}
			}
		}

		if isFlowOngoing(flowInstanceEntity.Status) {
			taskDetailModel.OperatorPermission = operatorPermissionModel
		}

	}
}

// setPermissionToSchema 将节点字段权限赋值到schema中去, 子表单的字段权限设置也在fieldPermission的同级
func (i *instance) setPermissionToSchema(formSchema *interface{}, fieldPermissionObj interface{}) interface{} {
	if formSchema == nil {
		return nil
	}
	fieldPermissionModel := utils.ChangeObjectToMap(fieldPermissionObj)
	formSchemaMap := utils.ChangeObjectToMap(formSchema)
	if formSchemaMap != nil && formSchemaMap["properties"] != nil {
		propertiesMap := utils.ChangeObjectToMap(formSchemaMap["properties"])
		propertiesMap = i.propertiesFieldPermissionUpdate(propertiesMap, fieldPermissionModel)
		formSchemaMap["properties"] = propertiesMap
		// formSchemaMap["table"] = formSchemaTableMap
	}

	return formSchemaMap
}

func (i *instance) propertiesFieldPermissionUpdate(propertiesMap map[string]interface{}, fieldPermissionModel map[string]interface{}) map[string]interface{} {
	if len(propertiesMap) > 0 {
		for key, fieldObj := range propertiesMap {
			if fieldObj != nil {
				fieldMap := utils.ChangeObjectToMap(fieldObj)
				xinternalMap := utils.ChangeObjectToMap(fieldMap["x-internal"])
				if xinternalMap != nil {
					permissionModel := &convert.FieldPermissionModel{}
					permissionModelObj := fieldPermissionModel[key]
					permissionModelJSON, err := json.Marshal(permissionModelObj)
					if err == nil {
						json.Unmarshal(permissionModelJSON, permissionModel)
					}
					var permission int8
					if utils.IsNotNil(permissionModel.XInternal.Permission) {
						permission = permissionModel.XInternal.Permission
					} else {
						permission = 0
					}
					xinternalMap["permission"] = permission

					if fieldMap["properties"] == nil { // 布局组件字段不校验权限
						_, read := i.task.PermissionConvertWriteRead(permission)
						if !read {
							delete(propertiesMap, key)
						}
					}
				}
				fieldMap["x-internal"] = xinternalMap

				if fieldMap["items"] != nil { // 子表单字段
					itemsMap := utils.ChangeObjectToMap(fieldMap["items"])
					itemsPropertiesMap := utils.ChangeObjectToMap(itemsMap["properties"])
					if len(itemsPropertiesMap) > 0 {
						itemsPropertiesMap = i.propertiesFieldPermissionUpdate(itemsPropertiesMap, fieldPermissionModel)
						itemsMap["properties"] = itemsPropertiesMap
						fieldMap["items"] = itemsMap
					}
				} else if fieldMap["properties"] != nil { // 布局中的字段
					layoutPropertiesMap := utils.ChangeObjectToMap(fieldMap["properties"])
					if len(layoutPropertiesMap) > 0 {
						layoutPropertiesMap = i.propertiesFieldPermissionUpdate(layoutPropertiesMap, fieldPermissionModel)
						fieldMap["properties"] = layoutPropertiesMap
					}
				}
				propertiesMap[key] = fieldMap
			}
		}
	}
	return propertiesMap
}

func (i *instance) ProcessHistories(ctx context.Context, processInstanceID string) ([]*models.InstanceStep, error) {
	instance, err := i.instanceRepo.GetEntityByProcessInstanceID(i.db, processInstanceID)
	if err != nil {
		return nil, err
	}
	if instance == nil {
		return nil, error2.NewErrorWithString(code.InvalidInstanceID, "flowInstance is nil")
	}
	flow, err := i.flowRepo.FindByID(i.db, instance.FlowID)
	if err != nil {
		return nil, err
	}
	if flow == nil {
		return nil, error2.NewErrorWithString(code.InvalidProcessID, "flow is nil")
	}
	launchUser, err := i.identityAPI.FindUserByID(ctx, instance.ApplyUserID)
	if err != nil {
		return nil, err
	}
	if launchUser == nil {
		launchUser = &client.UserInfoResp{}
	}
	steps, err := i.stepRepo.FindInstanceSteps(i.db, &models.InstanceStep{
		ProcessInstanceID: processInstanceID,
	})
	if err != nil {
		return make([]*models.InstanceStep, 0), err
	}
	if steps == nil || len(steps) == 0 {
		return make([]*models.InstanceStep, 0), nil
	}

	derivationSteps := make([]*models.InstanceStep, 0)
	tailStep := steps[0]
	for _, step := range steps {
		if step.Status == Cancel {
			step.CreatorName = launchUser.UserName
			step.CreatorAvatar = launchUser.Avatar
			step.FlowName = flow.Name
			continue
		}
		records, err := i.recordRepo.FindRecords(i.db, processInstanceID, step.ID, []string{opDeliver, opRead, opCC, opHandleRead, opHandleCC}, false)
		if err != nil {
			return make([]*models.InstanceStep, 0), err
		}
		if len(records) == 1 && records[0].HandleType == opAddSign { // 当前节点只有加签操作
			step.Status = opAddSign

			var tempMap map[string]string
			err := json.Unmarshal([]byte(records[0].CorrelationData), &tempMap)
			if err != nil {
				return nil, err
			}
			if tempMap == nil || tempMap["assignees"] == "" {
				return nil, nil
			}
			addSignUserIDs := strings.Split(tempMap["assignees"], ",")
			tmpRecords := make([]*models.OperationRecord, 0)
			for _, addSignUserID := range addSignUserIDs {
				addSignUser, _ := i.identityAPI.FindUserByID(ctx, addSignUserID)
				replenish := &models.OperationRecord{
					HandleType: opAddSign,
					BaseModel: models.BaseModel{
						CreatorID:     addSignUser.ID,
						CreatorName:   addSignUser.UserName,
						CreatorAvatar: addSignUser.Avatar,
					},
				}
				tmpRecords = append(tmpRecords, replenish)
			}
			step.OperationRecords = tmpRecords
			continue
		}
		for i, r := range records { // 去除加签的record
			if r.HandleType == opAddSign {
				records = append(records[:i], records[i+1:]...)
			}
		}

		taskHandleUserIDs := make([]string, 0)
		if step.HandleUserIDs != "" {
			taskHandleUserIDs = strings.Split(step.HandleUserIDs, ",")
		}
		for i := 0; i < len(records); {
			taskHandleUserIDs = utils.SliceRemoveElement(taskHandleUserIDs, records[i].CreatorID)
			i++
		}

		if len(taskHandleUserIDs) > 0 {
			for _, taskHandleUserID := range taskHandleUserIDs {
				replenish := &models.OperationRecord{
					HandleType: UnTreated,
					BaseModel: models.BaseModel{
						CreatorID: taskHandleUserID,
					},
				}
				records = append(records, replenish)
			}
		}

		if step.TaskType == convert.START {
			step.CreatorName = launchUser.UserName
			step.CreatorAvatar = launchUser.Avatar
			step.FlowName = flow.Name
		}
		i.fillComment(ctx, records)
		step.OperationRecords = records

		derivationRecords, err := i.operationRecordRepo.FindRecords(i.db, processInstanceID, step.ID, []string{opDeliver, opRead, opCC, opAddSign}, true)
		for _, derivationRecord := range derivationRecords {
			user, _ := i.identityAPI.FindUserByID(ctx, derivationRecord.CreatorID)
			if derivationRecord.HandleType == opDeliver {
				deliverUser, _ := i.identityAPI.FindUserByID(ctx, derivationRecord.CorrelationData)
				replenish := &models.OperationRecord{
					HandleType: opDeliver,
					BaseModel: models.BaseModel{
						CreatorID:     deliverUser.ID,
						CreatorName:   deliverUser.UserName,
						CreatorAvatar: deliverUser.Avatar,
					},
				}
				derivationStep := &models.InstanceStep{
					FlowName: "转交",
					Status:   opDeliver,
					TaskType: opDeliver,
					BaseModel: models.BaseModel{
						CreateTime:    derivationRecord.CreateTime,
						ModifyTime:    derivationRecord.ModifyTime,
						CreatorName:   user.UserName,
						CreatorAvatar: user.Avatar,
					},
					OperationRecords: []*models.OperationRecord{replenish},
				}
				derivationSteps = append(derivationSteps, derivationStep)
			}
			if derivationRecord.HandleType == opRead {
				operationRecords, _ := i.operationRecordRepo.FindRecords(i.db, processInstanceID, step.ID, []string{opHandleRead}, true)
				for _, r := range operationRecords {
					if r.CreatorID == "" {
						continue
					}
					user, _ := i.identityAPI.FindUserByID(ctx, r.HandleUserID)
					r.CreatorName = user.UserName
					r.CreatorAvatar = user.Avatar
				}

				derivationStep := &models.InstanceStep{
					FlowName: "邀请阅示",
					Status:   opRead,
					TaskType: opRead,
					Reason:   derivationRecord.Remark,
					BaseModel: models.BaseModel{
						CreateTime:    derivationRecord.CreateTime,
						ModifyTime:    derivationRecord.ModifyTime,
						CreatorName:   user.UserName,
						CreatorAvatar: user.Avatar,
					},
					OperationRecords: operationRecords,
				}
				derivationSteps = append(derivationSteps, derivationStep)
			} else if derivationRecord.HandleType == opCC {
				operationRecords, _ := i.operationRecordRepo.FindRecords(i.db, processInstanceID, step.ID, []string{opHandleCC}, true)
				for _, r := range operationRecords {
					if r.CreatorID == "" {
						continue
					}
					user, _ := i.identityAPI.FindUserByID(ctx, r.HandleUserID)
					r.CreatorName = user.UserName
					r.CreatorAvatar = user.Avatar
				}
				derivationStep := &models.InstanceStep{
					FlowName: "抄送",
					Status:   opCC,
					TaskType: opCC,
					Reason:   derivationRecord.Remark,
					BaseModel: models.BaseModel{
						CreateTime:    derivationRecord.CreateTime,
						ModifyTime:    derivationRecord.ModifyTime,
						CreatorName:   user.UserName,
						CreatorAvatar: user.Avatar,
					},
					OperationRecords: operationRecords,
				}
				derivationSteps = append(derivationSteps, derivationStep)
			}
			if derivationRecord.HandleType == opAddSign {
				relMap := make(map[string]string)
				err := json.Unmarshal([]byte(derivationRecord.CorrelationData), &relMap)
				if err != nil {
					return nil, err
				}
				addSignUserIDs := strings.Split(relMap["assignees"], ",")
				tmpRecords := make([]*models.OperationRecord, 0)
				for _, addSignUserID := range addSignUserIDs {
					addSignUser, _ := i.identityAPI.FindUserByID(ctx, addSignUserID)
					replenish := &models.OperationRecord{
						HandleType: opAddSign,
						BaseModel: models.BaseModel{
							CreatorID:     addSignUser.ID,
							CreatorName:   addSignUser.UserName,
							CreatorAvatar: addSignUser.Avatar,
						},
					}
					tmpRecords = append(tmpRecords, replenish)
				}
				derivationStep := &models.InstanceStep{
					FlowName: "加签",
					Status:   opAddSign,
					TaskType: opAddSign,
					BaseModel: models.BaseModel{
						CreateTime:    derivationRecord.CreateTime,
						ModifyTime:    derivationRecord.ModifyTime,
						CreatorName:   user.UserName,
						CreatorAvatar: user.Avatar,
					},
					OperationRecords: tmpRecords,
				}
				derivationSteps = append(derivationSteps, derivationStep)
			}
		}

	}
	steps = append(steps, derivationSteps...)
	sort.Sort(StepSlice(steps))

	// 拼接当前节点
	tasksReq := client.GetTasksReq{
		InstanceID: []string{processInstanceID},
	}
	tasks, err := i.processAPI.GetTasks(ctx, tasksReq)
	if err != nil {
		return make([]*models.InstanceStep, 0), err
	}
	var processCompleted bool
	processInstance, _ := i.processAPI.GetInstanceByID(ctx, processInstanceID)
	if processInstance == nil || processInstance.Status != Active {
		processCompleted = true
	}

	currentStep := models.InstanceStep{}
	if instance.Status == SendBack {

		user, _ := i.identityAPI.FindUserByID(ctx, instance.ApplyUserID)
		currentStep.TaskType = convert.START
		currentStep.ProcessInstanceID = processInstanceID
		currentStep.Status = opReSubmit
		currentStep.CreatorName = user.UserName
		currentStep.CreatorID = user.ID
		currentStep.CreateTime = instance.ModifyTime
		currentStep.ModifyTime = instance.ModifyTime
		rear := append([]*models.InstanceStep{}, steps[0:]...)
		steps = append(append(steps[:0], &currentStep), rear...)
	} else if !processCompleted && len(tasks.Data) > 0 && (tailStep.Status == opAddSign || tasks.Data[0].NodeDefKey != tailStep.TaskDefKey) {

		currentStep.ProcessInstanceID = processInstanceID
		currentStep.TaskName = tasks.Data[0].Name
		currentStep.Status = Review
		currentStep.CreatorID = pkg.STDUserID(ctx)
		currentStep.CreateTime = instance.ModifyTime
		currentStep.ModifyTime = instance.ModifyTime

		taskHandleUserIDs := make([]string, 0)
		if tasks.Data[0].TaskType == "TEMP_MODEL" {
			relRecord, err := i.operationRecordRepo.FindRecordByRelDefKey(i.db, processInstanceID, tasks.Data[0].NodeDefKey)
			if err != nil {
				return nil, err
			}
			var tempMap map[string]string
			err = json.Unmarshal([]byte(relRecord.CorrelationData), &tempMap)
			if err != nil {
				return nil, err
			}
			multiplePersonWay := tempMap["multiplePersonWay"]
			if multiplePersonWay == "or" {
				currentStep.TaskType = convert.OrApproval
			} else {
				currentStep.TaskType = convert.AndApproval
			}
			taskHandleUserIDs = strings.Split(tempMap["assignees"], ",")
		} else {
			// 获取当前节点类型
			// 节点相关数据
			shape, _ := convert.GetShapeByTaskDefKey(flow.BpmnText, tasks.Data[0].NodeDefKey)
			basicConfig := convert.GetTaskBasicConfigModel(shape)
			if basicConfig != nil {
				currentStep.TaskType = convert.GetCurrentNodeType(shape.Type, basicConfig.MultiplePersonWay)
				taskHandleUserIDs, _ = i.flow.GetTaskHandleUserIDs(ctx, shape, instance)
			}
		}

		records := make([]*models.OperationRecord, 0)
		if len(taskHandleUserIDs) == 0 { // 没有处理人
			rs, err := i.abnormalTaskRepo.Find(i.db, map[string]interface{}{
				"process_instance_id": processInstanceID,
				"task_id":             tasks.Data[0].ID,
			})
			if err != nil {
				return nil, err
			}
			if rs != nil && len(rs) > 0 { // 异常任务
				replenish := &models.OperationRecord{
					HandleDesc: "该节点下无相关负责人，已交给管理员处理",
				}
				records = append(records, replenish)
			}
		} else {
			for _, taskHandleUserID := range taskHandleUserIDs {
				user, _ := i.identityAPI.FindUserByID(ctx, taskHandleUserID)
				replenish := &models.OperationRecord{
					HandleType: UnTreated,
					BaseModel: models.BaseModel{
						CreatorID:     taskHandleUserID,
						CreatorName:   user.UserName,
						CreatorAvatar: user.Avatar,
					},
				}
				records = append(records, replenish)
			}
		}

		currentStep.OperationRecords = records
		rear := append([]*models.InstanceStep{}, steps[0:]...)
		steps = append(append(steps[:0], &currentStep), rear...)

	} else if processCompleted {
		currentStep.TaskType = convert.END
		currentStep.ProcessInstanceID = processInstanceID
		currentStep.TaskName = "结束"
		currentStep.Status = convert.END
		currentStep.CreatorID = pkg.STDUserID(ctx)
		currentStep.CreateTime = instance.ModifyTime
		currentStep.ModifyTime = instance.ModifyTime
		rear := append([]*models.InstanceStep{}, steps[0:]...)
		steps = append(append(steps[:0], &currentStep), rear...)
	}

	return steps, nil
}

// StepSlice sort
type StepSlice []*models.InstanceStep

func (s StepSlice) Len() int {
	return len(s)
}
func (s StepSlice) Less(i, j int) bool {
	return s[i].CreateTime > s[j].CreateTime
}
func (s StepSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (i *instance) fillComment(ctx context.Context, records []*models.OperationRecord) error {
	for _, record := range records {
		if record.CreatorID == "" {
			continue
		}
		user, err := i.identityAPI.FindUserByID(ctx, record.CreatorID)
		if err != nil {
			return err
		}
		record.CreatorName = user.UserName
		record.CreatorAvatar = user.Avatar
	}
	return nil
}

func derivation(opType string) bool {
	return opType != "" && (opCC == opType || opRead == opType || opDeliver == opType)
}

// ProcessTask 需要转换成 ActTaskEntity
func (i *instance) ConvertTask(ctx context.Context, tasks []*client.ProcessTask) ([]*models.ActTaskEntity, []string, []string) {
	if len(tasks) > 0 {
		list := make([]*models.ActTaskEntity, 0)
		processInstanceIDs := make([]string, 0)
		taskIDs := make([]string, 0)
		for _, value := range tasks {
			item := &models.ActTaskEntity{
				ID:          value.ID,
				TaskDefKey:  value.NodeDefKey,
				ProcInstID:  value.ProcInstanceID,
				Name:        value.Name,
				Description: value.Desc,
				Assignee:    value.Assignee,
				StartTime:   value.CreateTime,
				EndTime:     value.EndTime,
				DueDate:     value.DueTime,
				Handled:     value.Status,
			}
			list = append(list, item)

			processInstanceIDs = append(processInstanceIDs, value.ProcInstanceID)
			taskIDs = append(taskIDs, value.ID)
		}
		return list, processInstanceIDs, taskIDs
	}
	return make([]*models.ActTaskEntity, 0), make([]string, 0), make([]string, 0)
}

// ProcessInstance 需要转换成 FlowInstanceEntity
func (i *instance) ConvertInstance(ctx context.Context, instances []*client.ProcessInstance) ([]*models.Instance, []string) {
	if len(instances) > 0 {
		var list []*models.Instance
		processInstanceIDs := make([]string, 0)
		for _, value := range instances {
			processInstanceIDs = append(processInstanceIDs, value.ID)
		}

		list, _ = i.instanceRepo.FindByProcessInstanceIDs(i.db, processInstanceIDs)
		retMap := make(map[string]*models.Instance)
		for _, v := range list {
			retMap[v.ProcessInstanceID] = v
		}
		ret := make([]*models.Instance, 0)
		for _, v := range instances {
			e := retMap[v.ID]
			e.ModifyTime = v.ModifyTime
			ret = append(ret, e)
		}
		return ret, processInstanceIDs
	}
	return make([]*models.Instance, 0), make([]string, 0)
}

func (i *instance) Cal(ctx context.Context, valueFrom string, valueOf interface{}, formulaFields interface{}, instance *models.Instance, variables map[string]interface{}, formQueryRef interface{}, formDefKey string) (interface{}, error) {

	var value interface{}
	if "fixedValue" == valueFrom {
		value = valueOf
	} else if "currentFormValue" == valueFrom {
		data, err := i.formAPI.GetFormData(ctx, client.FormDataConditionModel{
			AppID:   instance.AppID,
			DataID:  instance.FormInstanceID,
			TableID: instance.FormID,
			Ref:     formQueryRef,
		})
		if err != nil {
			return nil, err
		}
		// value = data[utils.Strval(valueOf)]
		value, _ = i.formAPI.GetValue(data, utils.Strval(valueOf), nil)
	} else if "formula" == valueFrom {
		data, err := i.formAPI.GetFormData(ctx, client.FormDataConditionModel{
			AppID:   instance.AppID,
			DataID:  instance.FormInstanceID,
			TableID: instance.FormID,
			Ref:     formQueryRef,
		})
		if err != nil {
			return nil, err
		}
		for k, v := range data {
			variables[k] = v
		}

		// 判断有没有运算符
		if strings.Contains(utils.Strval(valueOf), "$") {
			expression := utils.Strval(valueOf)
			expression = strings.Replace(expression, "$variable.", "", -1)
			expression = strings.Replace(expression, "$"+formDefKey+".", "", -1)

			expression = strings.Replace(expression, "$", "", -1)
			ret, err := i.structorAPI.CalExpression(ctx, map[string]interface{}{
				"expression": expression,
				"parameter":  i.flow.FormatFormValue2(formulaFields, variables),
			})
			if err != nil {
				logger.Logger.Error(err)
				return "", nil
			}
			return ret, nil
		}
		return valueOf, nil

	} else if "processVariable" == valueFrom {
		data, err := i.formAPI.GetFormData(ctx, client.FormDataConditionModel{
			AppID:   instance.AppID,
			DataID:  instance.FormInstanceID,
			TableID: instance.FormID,
			Ref:     formQueryRef,
		})
		if err != nil {
			return nil, err
		}
		for k, v := range data {
			variables[k] = v
		}
		key := utils.Strval(valueOf)
		key = strings.Replace(key, "$variable.", "", 1)
		key = strings.Replace(key, "$"+formDefKey+".", "", 1)
		value = utils.GetFieldValue(variables, key)

		paramByte, err := json.Marshal(variables)
		fmt.Println("GetFieldValue " + string(paramByte))
		fmt.Println("GetFieldValue     key:" + key + "    value:" + utils.Strval(value))
	}

	return value, nil
}

// StartFlowModel start flow model
type StartFlowModel struct {
	UserID   string                 `json:"userId" binding:"required"`
	FormData map[string]interface{} `json:"formData" binding:"required"`
	FlowID   string                 `json:"flowId" binding:"required"`
}

// MyApplyReq my apply request
type MyApplyReq struct {
	page.ReqPage

	Status    string `json:"status"` // handle type : ALL，REVIEW, WRITE， READ， OTHER
	BeginDate string `json:"beginDate"`
	EndDate   string `json:"endDate"`
	Keyword   string `json:"keyword"`
	OrderType string `json:"orderType"` // orderType: ASC|DESC
}

// ResubmitReq resubmit req
type ResubmitReq struct {
	FormData map[string]interface{} `json:"formData"`
}

// TaskListReq task list request
type TaskListReq struct {
	page.ReqPage

	TagType    string `json:"tagType"`    // tag type : OVERTIME, URGE
	HandleType string `json:"handleType"` // handle type : REVIEW, WRITE， READ， OTHER
	Agent      int8   `json:"agent"`      // only see agent :1 or 0
	Keyword    string `json:"keyword"`
	AppID      string `json:"appId"`
	OrderType  string `json:"orderType"` // orderType: ASC|DESC
	Status     int8   `json:"status"`
}

// CcListReq cc list request
type CcListReq struct {
	page.ReqPage

	AppID     string `json:"appId"`
	Status    int8   `json:"status"` // status: 1 ，0, -1查全部
	Keyword   string `json:"keyword"`
	OrderType string `json:"orderType"` // orderType: ASC|DESC
}

// GetFormDataReq get form data req
type GetFormDataReq struct {
	Ref interface{} `json:"ref"`
}

// InstanceCountModel instance count info
type InstanceCountModel struct {
	OverTimeCount   int64 `json:"overTimeCount"`
	UrgeCount       int64 `json:"urgeCount"`
	WaitHandleCount int64 `json:"waitHandleCount"`
	CcToMeCount     int64 `json:"ccToMeCount"`
}

// TaskTypeDetailModel detail info request
type TaskTypeDetailModel struct {
	TaskID string `json:"taskId"`
	Type   string `json:"type"` // Type:APPLY_PAGE我发起的,WAIT_HANDLE_PAGE待办,HANDLED_PAGE已办,CC_PAGE抄送,ALL_PAGE全部
}

// InstanceDetailModel instance detail
type InstanceDetailModel struct {
	FlowName            string             `json:"flowName"`
	CanViewStatusAndMsg bool               `json:"canViewStatusAndMsg"`
	CanMsg              bool               `json:"canMsg"`
	TaskDetailModels    []*TaskDetailModel `json:"taskDetailModels"`
	AppID               string             `json:"appId"`
	TableID             string             `json:"tableId"`
}

// TaskDetailModel task detail
type TaskDetailModel struct {
	TaskID             string      `json:"taskId"`
	TaskName           string      `json:"taskName"`
	TaskType           string      `json:"taskType"` // review、correlation
	TaskDefKey         string      `json:"taskDefKey"`
	FormSchema         interface{} `json:"formSchema"`
	FieldPermission    interface{} `json:"fieldPermission"` // map[string]convert.FieldPermissionModel
	OperatorPermission interface{} `json:"operatorPermission"`
	HasCancelBtn       bool        `json:"hasCancelBtn"`
	HasResubmitBtn     bool        `json:"hasResubmitBtn"`
	HasReadHandleBtn   bool        `json:"hasReadHandleBtn"`
	HasCcHandleBtn     bool        `json:"hasCcHandleBtn"`
	HasUrgeBtn         bool        `json:"hasUrgeBtn"`
}
