package callback

import (
	"context"
	"encoding/json"
	"fmt"
	"git.internal.yunify.com/qxp/misc/logger"
	"github.com/quanxiang-cloud/flow/internal/convert"
	"github.com/quanxiang-cloud/flow/internal/flow"
	"github.com/quanxiang-cloud/flow/internal/models"
	"github.com/quanxiang-cloud/flow/internal/models/mysql"
	"github.com/quanxiang-cloud/flow/internal/server/options"
	"github.com/quanxiang-cloud/flow/pkg"
	"github.com/quanxiang-cloud/flow/pkg/client"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/utils"
	"gorm.io/gorm"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// Callback Callback
type Callback interface {
	TaskStartEventHandler(ctx context.Context, data map[string]string) error
	TaskEndEventHandler(ctx context.Context, data map[string]string) error
}

type callback struct {
	db                    *gorm.DB
	flowRepo              models.FlowRepo
	formAPI               client.Form
	instanceRepo          models.InstanceRepo
	flowVariable          models.VariablesRepo
	instanceVariablesRepo models.InstanceVariablesRepo
	abnormalTaskRepo      models.AbnormalTaskRepo
	messageCenterAPI      client.MessageCenter
	structorAPI           client.Structor
	processAPI            client.Process
	identityAPI           client.Identity
	urge                  flow.Urge
	flow                  flow.Flow
	instance              flow.Instance
	operationRecord       flow.OperationRecord
	instanceExecutionRepo models.InstanceExecutionRepo
	task                  flow.Task
	formFieldRepo         models.FormFieldRepo
	polyAPI               client.PolyAPI
}

// NewCallback init
func NewCallback(conf *config.Configs, opts ...options.Options) (Callback, error) {
	instance, err := flow.NewInstance(conf, opts...)
	if err != nil {
		return nil, nil
	}
	urge, err := flow.NewUrge(conf, opts...)
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

	c := &callback{
		flowRepo:              mysql.NewFlowRepo(),
		formAPI:               client.NewForm(conf),
		instanceRepo:          mysql.NewInstanceRepo(),
		instanceVariablesRepo: mysql.NewInstanceVariablesRepo(),
		abnormalTaskRepo:      mysql.NewAbnormalTaskRepo(),
		flowVariable:          mysql.NewVariablesRepo(),
		messageCenterAPI:      client.NewMessageCenter(conf),
		structorAPI:           client.NewStructor(conf),
		processAPI:            client.NewProcess(conf),
		identityAPI:           client.NewIdentity(conf),
		urge:                  urge,
		flow:                  flow,
		instance:              instance,
		operationRecord:       operationRecord,
		instanceExecutionRepo: mysql.NewInstanceExecutionRepo(),
		task:                  task,
		formFieldRepo:         mysql.NewFormFieldRepo(),
		polyAPI:               client.NewPolyAPI(conf),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

// SetDB set db
func (c *callback) SetDB(db *gorm.DB) {
	c.db = db
}

func (c *callback) Transaction() *gorm.DB {
	return c.db.Begin()
}

// TaskStartEventHandler node init start event(not task event)
func (c *callback) TaskStartEventHandler(ctx context.Context, data map[string]string) error {
	processID := data["processID"]
	nodeDefKey := data["nodeDefKey"]
	requestID := data["requestID"]
	userID := data["userID"]
	processInstanceID := data["processInstanceID"]
	ctx = pkg.RPCCTXTransfer(requestID, userID)
	flow, err := c.flowRepo.FindByProcessID(c.db, processID)
	if err != nil {
		return err
	}

	shape, err := convert.GetShapeByTaskDefKey(flow.BpmnText, nodeDefKey)
	if err != nil {
		return err
	}
	if shape == nil {
		return nil
	}

	if shape.Type == convert.ProcessVariableAssignment {
		formShape, err := convert.GetShapeByChartType(flow.BpmnText, convert.FormData)
		if err != nil {
			return err
		}
		return c.processVariableChangeServiceTask(ctx, shape, processInstanceID, formShape.ID)
	} else if shape.Type == convert.WebHook {
		formShape, err := convert.GetShapeByChartType(flow.BpmnText, convert.FormData)
		if err != nil {
			return err
		}
		return c.webHookServiceTask(ctx, shape, processInstanceID, formShape.ID)
	}

	logger.Logger.Info("事件：加载AssigneeList，节点为" + nodeDefKey)
	// add dynamic assignee user list
	if shape.Type == convert.Approve || shape.Type == convert.FillIn {
		approvePersonsObj := convert.GetValueFromBusinessData(*shape, "basicConfig.approvePersons")
		if approvePersonsObj != nil {
			err = c.addAssignee(ctx, shape, data)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// TaskEndEventHandler node init end event(not task event, is node event)，
func (c *callback) TaskEndEventHandler(ctx context.Context, data map[string]string) error {
	processID := data["processID"]
	nodeDefKey := data["nodeDefKey"]
	processInstanceID := data["processInstanceID"]
	requestID := data["requestID"]
	userID := data["userID"]

	ctx = pkg.RPCCTXTransfer(requestID, userID)
	flow, err := c.flowRepo.FindByProcessID(c.db, processID)
	if err != nil {
		return err
	}
	instance, err := c.instanceRepo.GetEntityByProcessInstanceID(c.db, processInstanceID)
	if err != nil {
		return err
	}
	shape, err := convert.GetShapeByTaskDefKey(flow.BpmnText, nodeDefKey)
	if err != nil {
		return err
	}
	formShape, err := convert.GetShapeByChartType(flow.BpmnText, convert.FormData)
	if err != nil {
		return err
	}
	switch shape.Type {
	case convert.Email:
		return c.emailServiceTask(ctx, shape, processInstanceID)
	case convert.Autocc:
		return c.autoCCServiceTask(ctx, shape, processInstanceID)
	case convert.Letter:
		return c.letterServiceTask(ctx, shape, processInstanceID)
	case convert.TableDataCreate:
		return c.tableDataCreateServiceTask(ctx, shape, processInstanceID, formShape.ID)
	case convert.TableDataUpdate:
		return c.tableDataUpdateServiceTask(ctx, shape, processInstanceID, formShape.ID)
	case convert.Approve:
		fallthrough
	case convert.FillIn:
		taskIDs := data["taskID"] // 如果是会签则taskID是数组，逗号间隔的数组
		if len(taskIDs) > 0 {
			req := client.GetTasksReq{
				TaskID: strings.Split(taskIDs, ","),
			}
			resp, err := c.processAPI.GetTasks(ctx, req)
			if err != nil {
				return err
			}
			if len(resp.Data) == 0 {
				return nil
			}
			for _, value := range resp.Data {
				c.task.TaskInitHandle(ctx, flow, instance, value, userID)
			}
		}
	}
	return nil
}

func (c *callback) addAssignee(ctx context.Context, shape *convert.ShapeModel, data map[string]string) error {
	processID := data["processID"]
	processInstanceID := data["processInstanceID"]
	nodeDefKey := data["nodeDefKey"]

	flowInstanceEntity, err := c.instanceRepo.GetEntityByProcessInstanceID(c.db, processInstanceID)
	if err != nil {
		return err
	}

	// assigneeList dynamic handle users
	_, assigneeList := c.flow.GetTaskHandleUserIDs(ctx, shape, flowInstanceEntity)

	logger.Logger.Info("事件：加载AssigneeList，AssigneeList=" + utils.ChangeStringArrayToString(assigneeList))
	if len(assigneeList) > 0 {
		return c.processAPI.SetProcessVariables(ctx, processID, processInstanceID, nodeDefKey, convert.AssigneeList, assigneeList)
	}

	return nil
}

// Send email node
func (c *callback) emailServiceTask(ctx context.Context, shape *convert.ShapeModel, processInstanceID string) error {
	instance, err := c.instanceRepo.GetEntityByProcessInstanceID(c.db, processInstanceID)
	if err != nil {
		return err
	}

	bd := shape.Data.BusinessData
	if bd == nil {
		return nil
	}

	dataReq := client.FormDataConditionModel{
		AppID:   instance.AppID,
		TableID: instance.FormID,
		DataID:  instance.FormInstanceID,
	}

	dataResp, err := c.formAPI.GetFormData(ctx, dataReq)
	if err != nil {
		return err
	}
	if dataResp == nil {
		return nil
	}

	// replace content
	content := utils.Strval(bd["content"])
	value := c.flow.FormatFormValue(instance, dataResp)
	var fieldType map[string]interface{}
	if v := bd["fieldType"]; v != nil {
		fieldType = v.(map[string]interface{})
	}
	for k, v := range value {
		t := fieldType[k]
		if t == "datepicker" {
			vt := v.(string)
			if strings.Contains(vt, ".000Z") {
				vt = strings.Replace(vt, ".000Z", "+0000", 1)
			}
			v = utils.ChangeISO8601ToBjTime(vt)
		}
		content = strings.Replace(content, "${"+k+"}", utils.Strval(v), 1)
	}

	handleUsers := c.flow.GetTaskHandleUsers2(ctx, bd["approvePersons"], instance)

	// gen req params
	mesAttachments := make([]map[string]interface{}, 0)
	if v := bd["mes_attachment"]; v != nil {
		arr := v.([]interface{})
		for _, e := range arr {
			tmp := e.(map[string]interface{})
			mesAttachment := make(map[string]interface{})
			mesAttachment["name"] = utils.Strval(tmp["file_name"])
			mesAttachment["path"] = utils.Strval(tmp["file_url"])
			mesAttachments = append(mesAttachments, mesAttachment)
		}
	}

	emailAddr := make([]string, 0)
	for _, user := range handleUsers {
		emailAddr = append(emailAddr, user.Email)
	}
	email := client.Email{
		To: emailAddr,
		Contents: client.Contents{
			Content: content,
		},
		Title: utils.Strval(bd["title"]),
		Files: mesAttachments,
	}
	msgReq := client.MsgReq{
		Email: email,
	}
	// post msg
	err = c.messageCenterAPI.MessageCreate(ctx, msgReq)
	if err != nil {
		return err
	}
	return nil
}

// Send website letter node
func (c *callback) letterServiceTask(ctx context.Context, shape *convert.ShapeModel, processInstanceID string) error {
	instance, err := c.instanceRepo.GetEntityByProcessInstanceID(c.db, processInstanceID)
	if err != nil {
		return err
	}

	bd := shape.Data.BusinessData
	if bd == nil {
		return nil
	}

	handleUsers := c.flow.GetTaskHandleUsers2(ctx, bd["approvePersons"], instance)
	// gen req params
	var recivers []client.Receivers
	for _, user := range handleUsers {
		reciver := client.Receivers{
			Type: 1,
			ID:   user.ID,
			Name: user.UserName,
		}
		recivers = append(recivers, reciver)
	}
	types, err := strconv.Atoi(bd["sort"].(string))
	if err != nil {
		return err
	}
	web := client.Web{
		IsSend: true,
		Title:  utils.Strval(bd["title"]),
		Contents: client.Contents{
			Content: utils.Strval(bd["content"]),
		},
		Receivers: recivers,
		Types:     types,
	}
	m := client.Mail{
		Web: web,
	}
	// post msg
	err = c.messageCenterAPI.MessageCreateff(ctx, m)
	if err != nil {
		return err
	}
	return nil
}

// Auto cc node
func (c *callback) autoCCServiceTask(ctx context.Context, shape *convert.ShapeModel, processInstanceID string) error {
	instance, err := c.instanceRepo.GetEntityByProcessInstanceID(c.db, processInstanceID)
	if err != nil {
		return err
	}

	ccUsers, _ := c.flow.GetTaskHandleUserIDs(ctx, shape, instance)

	if len(ccUsers) == 0 {
		return nil
	}

	req := &client.AddTaskReq{
		InstanceID: processInstanceID,
		NodeDefKey: shape.ID,
		Name:       shape.Data.NodeData.Name,
		Desc:       convert.CcTask,
		Assignee:   ccUsers,
	}
	_, err = c.processAPI.AddNonModelTask(ctx, req)
	return err
}

func (c *callback) processVariableChangeServiceTask(ctx context.Context, shape *convert.ShapeModel, processInstanceID string, formDefKey string) error {
	instance, err := c.instanceRepo.GetEntityByProcessInstanceID(c.db, processInstanceID)
	if err != nil {
		return err
	}

	bd := shape.Data.BusinessData
	if bd == nil {
		return nil
	}

	var assignmentRules []map[string]interface{}
	if v := bd["assignmentRules"]; v != nil {
		arr := v.([]interface{})
		for _, e := range arr {
			assignmentRules = append(assignmentRules, e.(map[string]interface{}))
		}
	}
	if assignmentRules == nil {
		return nil
	}
	variables, err := c.flow.GetInstanceVariableValues(ctx, instance)
	if err != nil {
		return err
	}
	for _, e := range assignmentRules {
		variableName := utils.Strval(e["variableName"])
		valueFrom := utils.Strval(e["valueFrom"])
		valueOf := e["valueOf"]

		value, err := c.cal(ctx, valueFrom, valueOf, nil, instance, variables, nil, formDefKey)
		if err != nil {
			return err
		}

		// c.instanceVariablesRepo.f
		fieldValue, fieldType := utils.StrvalAndType(value)
		err = c.instanceVariablesRepo.UpdateTypeAndValue(c.db, processInstanceID, variableName, fieldType, fieldValue)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *callback) tableDataCreateServiceTask(ctx context.Context, shape *convert.ShapeModel, processInstanceID string, formDefKey string) error {
	instance, err := c.instanceRepo.GetEntityByProcessInstanceID(c.db, processInstanceID)
	if err != nil {
		return err
	}

	bd := shape.Data.BusinessData
	if bd == nil {
		return nil
	}

	var createRules map[string]interface{}
	if v := bd["createRule"]; v != nil {
		createRules = v.(map[string]interface{})
	}
	if createRules == nil {
		return nil
	}

	variables, err := c.instance.GetInstanceVariableValues(ctx, instance)
	if err != nil {
		return err
	}

	// master form
	filedValueReq := make(map[string]interface{})
	for k, v := range createRules {
		tmp := v.(map[string]interface{})
		valueFrom := utils.Strval(tmp["valueFrom"])
		valueOf := tmp["valueOf"]

		value, err := c.cal(ctx, valueFrom, valueOf, nil, instance, variables, nil, formDefKey)
		if err != nil {
			return err
		}
		filedValueReq[k] = value
	}

	// child form
	refDataReq := make(map[string]client.RefData)
	var refs map[string]interface{}
	if v := bd["ref"]; v != nil {
		refs = v.(map[string]interface{})
	}
	if refs != nil {
		for k, v := range refs {
			tmp := v.(map[string]interface{})
			refData := client.RefData{
				AppID:   instance.AppID,
				TableID: utils.Strval(tmp["tableId"]),
				Type:    utils.Strval(tmp["type"]),
			}

			var subTableCreateRules []map[string]interface{}
			if t := tmp["createRules"]; t != nil {
				arr := t.([]interface{})
				for _, e := range arr {
					subTableCreateRules = append(subTableCreateRules, e.(map[string]interface{}))
				}
			}

			new := make([]client.CreateEntity, 0)
			if subTableCreateRules != nil {
				for _, e := range subTableCreateRules {
					record := client.CreateEntity{}
					recordEntity := make(map[string]interface{})
					for k1, v1 := range e {
						valueFrom := utils.Strval(v1.(map[string]interface{})["valueFrom"])
						valueOf := v1.(map[string]interface{})["valueOf"]
						value, err := c.cal(ctx, valueFrom, valueOf, nil, instance, variables, nil, formDefKey)
						if err != nil {
							return err
						}
						recordEntity[k1] = value
					}
					record.Entity = recordEntity
					new = append(new, record)
				}
			}
			refData.New = new
			refDataReq[k] = refData
		}
	}

	err = c.formAPI.CreateData(ctx, instance.AppID, utils.Strval(bd["targetTableId"]), client.CreateEntity{
		Entity: filedValueReq,
		Ref:    refDataReq,
	}, bd["silent"].(bool))
	if err != nil {
		return err
	}

	return nil
}

func (c *callback) tableDataUpdateServiceTask(ctx context.Context, shape *convert.ShapeModel, processInstanceID string, formDefKey string) error {
	instance, err := c.instanceRepo.GetEntityByProcessInstanceID(c.db, processInstanceID)
	if err != nil {
		return err
	}

	bd := shape.Data.BusinessData
	if bd == nil {
		return nil
	}

	variables, err := c.instance.GetInstanceVariableValues(ctx, instance)
	if err != nil {
		return err
	}

	targetTableID := utils.Strval(bd["targetTableId"])
	formQueryRef := bd["formQueryRef"]
	triggerAgain := bd["silent"].(bool)

	var updateIDs []string
	if instance.FormID != targetTableID { // 非本表支持过滤条件，本表只能更新当前记录
		var filterRule map[string]interface{}
		if v := bd["filterRule"]; v != nil {
			filterRule = v.(map[string]interface{})
		}
		if filterRule == nil {
			return nil
		}

		var conditions []map[string]interface{}
		if v := filterRule["conditions"]; v != nil {
			arr := v.([]interface{})
			for _, e := range arr {
				conditions = append(conditions, e.(map[string]interface{}))
			}
		}
		if conditions == nil {
			return nil
		}

		// reqConditions := make([]map[string]interface{}, 0)
		boolMap := make(map[string]interface{})
		queryMap := map[string]interface{}{
			"bool": boolMap,
		}
		terms := make([]map[string]interface{}, 0)
		if utils.Strval(filterRule["tag"]) == "or" {
			boolMap["should"] = &terms
		} else {
			boolMap["must"] = &terms
		}

		for _, v := range conditions {

			valueOf := v["value"]
			value, err := c.cal(ctx, "currentFormValue", valueOf, nil, instance, variables, formQueryRef, formDefKey)
			if err != nil {
				return err
			}

			// 等于eq, 不等于neq，包含in，不包含nin
			if utils.Strval(v["operator"]) == "eq" {
				term := map[string]interface{}{
					"term": map[string]interface{}{
						utils.Strval(v["fieldName"]): value,
					},
				}
				terms = append(terms, term)
			} else if utils.Strval(v["operator"]) == "neq" {
				mustNot := make([]map[string]interface{}, 0)
				mustNot = append(mustNot, map[string]interface{}{
					"term": map[string]interface{}{
						utils.Strval(v["fieldName"]): value,
					},
				})
				term := map[string]interface{}{
					"bool": map[string]interface{}{
						"mustNot": mustNot,
					},
				}
				terms = append(terms, term)
			} else if utils.Strval(v["operator"]) == "in" {
				// todo 如果是数组格式的需要修改in判断
				term := map[string]interface{}{
					"term": map[string]interface{}{
						utils.Strval(v["fieldName"]): value,
					},
				}
				terms = append(terms, term)
			} else if utils.Strval(v["operator"]) == "nin" {
				// todo 如果是数组格式的需要修改in判断
				mustNot := make([]map[string]interface{}, 0)
				mustNot = append(mustNot, map[string]interface{}{
					"term": map[string]interface{}{
						utils.Strval(v["fieldName"]): value,
					},
				})
				term := map[string]interface{}{
					"bool": map[string]interface{}{
						"mustNot": mustNot,
					},
				}
				terms = append(terms, term)
			}
		}

		updateIDs, err = c.formAPI.GetIDs(ctx, instance.AppID, targetTableID, queryMap)
		if err != nil {
			return err
		}
	} else { // if target form is current, can not trigger flow again
		triggerAgain = false

		updateIDs = []string{instance.FormInstanceID} // if target form is current , only update current data
	}

	var updateRules []map[string]interface{}
	if v := bd["updateRule"]; v != nil {
		arr := v.([]interface{})
		for _, e := range arr {
			updateRules = append(updateRules, e.(map[string]interface{}))
		}
	}

	updateReq := make(map[string]interface{}, 0)
	for _, updateRule := range updateRules {
		fieldName := utils.Strval(updateRule["fieldName"])
		valueFrom := utils.Strval(updateRule["valueFrom"])
		valueOf := updateRule["valueOf"]
		formulaFields := updateRule["formulaFields"]

		val, err := c.cal(ctx, valueFrom, valueOf, formulaFields, instance, variables, formQueryRef, formDefKey)
		if err != nil {
			return err
		}
		updateReq[fieldName] = val
	}

	selectField := utils.Strval(bd["selectField"])               // 普通组件为空，高级组件为字段名
	selectFieldType := utils.Strval(bd["selectFieldType"])       // 高级组件类型
	selectFieldTableID := utils.Strval(bd["selectFieldTableId"]) // 高级组件涉及的tableId
	ref := make(map[string]client.RefData, 0)
	if len(selectField) > 0 && selectField != "normal" {
		if instance.FormID == targetTableID { // 本表
			dataReq := client.FormDataConditionModel{
				AppID:   instance.AppID,
				TableID: instance.FormID,
				DataID:  instance.FormInstanceID,
			}
			if selectFieldType == "associated_records" || selectFieldType == "foreign_table" || selectFieldType == "sub_table" {
				dataReq.Ref = map[string]interface{}{
					selectField: map[string]interface{}{
						"appID":   instance.AppID,
						"tableID": selectFieldTableID,
						"type":    selectFieldType,
					},
				}
			}

			dataResp, err := c.formAPI.GetFormData(ctx, dataReq)
			if err != nil {
				return err
			}
			if dataResp == nil {
				return nil
			}

			fmt.Println("updateNode formData:" + utils.Strval(dataResp))
			fmt.Println("updateNode selectField:" + utils.Strval(dataResp[selectField]))
			fmt.Println("updateNode updateReq:" + utils.Strval(updateReq))
			if selectFieldType == "associated_records" || selectFieldType == "foreign_table" { // 外表
				targetTableID = selectFieldTableID
				updateIDs = utils.ChangeInterfaceToIDArray(dataResp[selectField])
			} else if selectFieldType == "associated_data" { // 外表
				targetTableID = selectFieldTableID
				// updateIDs = utils.ChangeInterfaceToValueArray(dataResp[selectField])
				associatedData := utils.ChangeObjectToMap(dataResp[selectField])
				if associatedData != nil {
					updateIDs = append(updateIDs, utils.Strval(associatedData["value"]))
				}
			} else if selectFieldType == "sub_table" { // 本表
				selectFieldDatas := utils.ChangeObjectToMapList(dataResp[selectField])
				newArr := make([]client.UpdateEntity, 0)
				for _, selectData := range selectFieldDatas {
					newArr = append(newArr, client.UpdateEntity{
						Entity: updateReq,
						Query: map[string]interface{}{
							"term": map[string]interface{}{
								"_id": selectData["_id"],
							},
						},
					})
				}

				ref[selectField] = client.RefData{
					AppID:   instance.AppID,
					TableID: selectFieldTableID,
					Type:    selectFieldType,
					Updated: newArr,
				}

				updateReq = make(map[string]interface{}, 0)
			}

			fmt.Println("updateNode updateIDs:" + utils.Strval(updateIDs))
			fmt.Println("updateNode updateReq:" + utils.Strval(updateReq))
		} else { // 非本表
			fmt.Println("updateNode filter updateIDs:" + utils.Strval(updateIDs))
			if len(updateIDs) > 0 {
				for _, updateID := range updateIDs {
					updateReq2 := make(map[string]interface{}, 0)
					updateIDs2 := make([]string, 0)
					ref2 := make(map[string]client.RefData, 0)
					targetTableID2 := targetTableID
					dataReq := client.FormDataConditionModel{
						AppID:   instance.AppID,
						TableID: targetTableID,
						DataID:  updateID,
					}
					if selectFieldType == "associated_records" || selectFieldType == "foreign_table" || selectFieldType == "sub_table" {
						dataReq.Ref = map[string]interface{}{
							selectField: map[string]interface{}{
								"appID":   instance.AppID,
								"tableID": selectFieldTableID,
								"type":    selectFieldType,
							},
						}
					}
					dataResp, err := c.formAPI.GetFormData(ctx, dataReq)
					if err != nil {
						return err
					}
					if dataResp == nil {
						return nil
					}

					fmt.Println("updateNode formData:" + utils.Strval(dataResp))
					fmt.Println("updateNode selectField:" + utils.Strval(dataResp[selectField]))
					fmt.Println("updateNode updateReq:" + utils.Strval(updateReq))

					if selectFieldType == "associated_records" || selectFieldType == "foreign_table" { // 外表
						targetTableID2 = selectFieldTableID
						updateIDs2 = utils.ChangeInterfaceToIDArray(dataResp[selectField])
						updateReq2 = updateReq
					} else if selectFieldType == "associated_data" { // 外表
						targetTableID2 = selectFieldTableID
						// updateIDs2 = utils.ChangeInterfaceToValueArray(dataResp[selectField])
						associatedData := utils.ChangeObjectToMap(dataResp[selectField])
						if associatedData != nil {
							updateIDs2 = append(updateIDs2, utils.Strval(associatedData["value"]))
						}
						updateReq2 = updateReq
					} else if selectFieldType == "sub_table" { // 本表
						selectFieldDatas := utils.ChangeObjectToMapList(dataResp[selectField])
						newArr := make([]client.UpdateEntity, 0)
						for _, selectData := range selectFieldDatas {
							newArr = append(newArr, client.UpdateEntity{
								Entity: updateReq,
								Query: map[string]interface{}{
									"term": map[string]interface{}{
										"_id": selectData["_id"],
									},
								},
							})
						}
						updateIDs2 = []string{updateID}

						ref2[selectField] = client.RefData{
							AppID:   instance.AppID,
							TableID: selectFieldTableID,
							Type:    selectFieldType,
							Updated: newArr,
						}
					}

					fmt.Println("updateNode updateIDs:" + utils.Strval(updateIDs2))
					fmt.Println("updateNode updateReq:" + utils.Strval(updateReq2))

					err = c.formAPI.UpdateData(ctx, instance.AppID, targetTableID2, "", client.UpdateEntity{
						Entity: updateReq2,
						Query: map[string]interface{}{
							"terms": map[string]interface{}{
								"_id": updateIDs2,
							},
						},
						Ref: ref2,
					}, triggerAgain)
					if err != nil {
						return err
					}
				}
			}
			return nil
		}
	}

	err = c.formAPI.UpdateData(ctx, instance.AppID, targetTableID, "", client.UpdateEntity{
		Entity: updateReq,
		Query: map[string]interface{}{
			"terms": map[string]interface{}{
				"_id": updateIDs,
			},
		},
		Ref: ref,
	}, triggerAgain)
	if err != nil {
		return err
	}
	return err
}

func (c *callback) webHookServiceTask(ctx context.Context, shape *convert.ShapeModel, processInstanceID string, formDefKey string) error {
	instance, err := c.instanceRepo.GetEntityByProcessInstanceID(c.db, processInstanceID)
	if err != nil {
		return err
	}

	variables, err := c.instance.GetInstanceVariableValues(ctx, instance)
	if err != nil {
		return err
	}

	dataReq := client.FormDataConditionModel{
		AppID:   instance.AppID,
		TableID: instance.FormID,
		DataID:  instance.FormInstanceID,
	}
	dataResp, err := c.formAPI.GetFormData(ctx, dataReq)
	if err != nil {
		return err
	}
	if dataResp == nil {
		return nil
	}

	for k, v := range dataResp {
		variables[k] = v
	}

	bd := shape.Data.BusinessData
	if bd == nil {
		return nil
	}

	hookType := utils.Strval(bd["type"])
	var conf map[string]interface{}
	if v := bd["config"]; v != nil {
		conf = v.(map[string]interface{})
	}
	var inputs []convert.Input
	if v := conf["inputs"]; v != nil {
		arr := v.([]interface{})
		for _, e := range arr {
			marshal, err := json.Marshal(e)
			if err != nil {
				return err
			}

			input := convert.Input{}
			err = json.Unmarshal(marshal, &input)
			if err != nil {
				return err
			}
			inputs = append(inputs, input)
		}
	}

	// gen req
	requestBody := make(map[string]interface{})
	requestHeader := make(map[string]string) // req header
	path := ""
	method := utils.Strval(conf["method"])
	if hookType == "request" {
		var api map[string]interface{}
		if v := conf["api"]; v != nil {
			api = v.(map[string]interface{})
		}
		path = utils.Strval(api["value"])
	} else {
		path = utils.Strval(conf["sendUrl"])
	}

	// do req
	queryFlag := false
	v := url.Values{}
	for _, e := range inputs {
		if e.In == convert.Header {
			val, err := c.webHookCal(ctx, e, variables, formDefKey)
			if err != nil {
				return err
			}
			requestHeader[e.Name] = utils.Strval(val)
		} else if e.In == convert.Path {
			val, err := c.webHookCal(ctx, e, variables, formDefKey)
			if err != nil {
				return err
			}
			path = strings.Replace(path, e.Name, utils.Strval(val), 1)
		} else if e.In == convert.Query {
			val, err := c.webHookCal(ctx, e, variables, formDefKey)
			if err != nil {
				return err
			}
			queryFlag = true
			v.Add(e.Name, utils.Strval(val))
			// path += e.Name + "=" + utils.Strval(val) + "&"
		} else if e.In == convert.Body {
			if v := e.Data; v != nil {
				if reflect.TypeOf(v).Kind() == reflect.Slice {
					arr := v.([]interface{})
					for _, e := range arr {
						marshal, err := json.Marshal(e)
						if err != nil {
							return err
						}

						input := convert.Input{}
						err = json.Unmarshal(marshal, &input)
						if err != nil {
							return err
						}
						val, err := c.exchangeParam(ctx, input, variables, formDefKey)
						if err != nil {
							return err
						}
						requestBody[input.Name] = val
					}
				} else {
					// requestBody[e.Name] = e.Data

					val, err := c.webHookCal(ctx, e, variables, formDefKey)
					if err != nil {
						return err
					}
					requestBody[e.Name] = utils.Strval(val)
				}

			}

		}

	}
	if queryFlag {
		path += "?"
		path += v.Encode()
	}
	if hookType == "request" {

		resp, err := c.polyAPI.InnerRequest(ctx, path, requestBody, requestHeader, method)
		if err != nil {
			return err
		}
		if resp == nil {
			return nil
		}
		apiMap := c.apiMap(resp)
		// save resp
		for k, v := range apiMap {
			code := "$" + shape.ID + "." + k
			variable, err := c.instanceVariablesRepo.FindVariablesByCode(c.db, processInstanceID, code)
			if err != nil {
				return err
			}

			fieldValue, fieldType := utils.StrvalAndType(v)
			if variable.ID != "" {
				if err := c.instanceVariablesRepo.UpdateTypeAndValue(c.db, processInstanceID, code, fieldType, fieldValue); err != nil {
					return err
				}
			} else {
				variable := models.InstanceVariables{
					ProcessInstanceID: processInstanceID,
					Code:              code,
					FieldType:         fieldType,
					Value:             fieldValue,
				}
				if err := c.instanceVariablesRepo.Create(c.db, &variable); err != nil {
					return err
				}
			}
		}
	} else if hookType == "send" {
		requestHeader["Content-Type"] = utils.Strval(conf["contentType"])
		_, err := c.polyAPI.SendRequest(ctx, path, requestBody, requestHeader, method)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *callback) apiMap(apiMap map[string]interface{}) map[string]interface{} {
	ret := make(map[string]interface{})
	for k, v := range apiMap {
		typeOf := reflect.TypeOf(v).Kind()
		if typeOf == reflect.Map {
			m := c.apiMap(v.(map[string]interface{}))
			for k1, v1 := range m {
				ret[k+"."+k1] = v1
			}
		} else {
			ret[k] = v
		}
	}
	return ret
}

// exchangeParam exchange body val
func (c *callback) exchangeParam(ctx context.Context, input convert.Input, variables map[string]interface{}, formDefKey string) (interface{}, error) {
	if input.Type == "object" {
		var inputs []convert.Input
		if v := input.Data; v != nil {
			arr := v.([]interface{})
			for _, e := range arr {
				marshal, err := json.Marshal(e)
				if err != nil {
					return nil, err
				}

				input := convert.Input{}
				err = json.Unmarshal(marshal, &input)
				if err != nil {
					return nil, err
				}
				inputs = append(inputs, input)
			}
		}
		param0 := make(map[string]interface{})
		for _, e := range inputs {
			val, err := c.exchangeParam(ctx, e, variables, formDefKey)
			if err != nil {
				return nil, err
			}
			param0[e.Name] = val
		}
		return param0, nil
	}
	val, err := c.webHookCal(ctx, input, variables, formDefKey)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (c *callback) webHookCal(ctx context.Context, input convert.Input, variables map[string]interface{}, formDefKey string) (interface{}, error) {
	if input.Type == "direct_expr" {
		// 判断有没有运算符
		if strings.Contains(utils.Strval(input.Data), "$") {

			expression := strings.TrimSpace(utils.Strval(input.Data))
			expression = strings.Replace(expression, "$variable.", "", -1)
			expression = strings.Replace(expression, "$"+formDefKey+".", "", -1)

			for k, v := range variables {
				// if strings.HasPrefix(k, "$") && strings.Contains(expression, k) {
				if strings.Contains(expression, k) {
					expression = strings.Replace(expression, k, utils.Strval(v), -1)
				}
			}

			expression = strings.Replace(expression, "$", "", -1)
			ret, err := c.structorAPI.CalExpression(ctx, map[string]interface{}{
				"expression": expression,
				"parameter":  variables,
			})
			if err != nil { // 公式计算，不能计算字符串，只能计算数值
				return expression, nil
			}
			return ret, nil
		}
		return input.Data, nil
	}
	return input.Data, nil
}

func (c *callback) cal(ctx context.Context, valueFrom string, valueOf interface{}, formulaFields interface{}, instance *models.Instance, variables map[string]interface{}, formQueryRef interface{}, formDefKey string) (interface{}, error) {

	var value interface{}
	if "fixedValue" == valueFrom {
		value = valueOf
	} else if "currentFormValue" == valueFrom {
		data, err := c.formAPI.GetFormData(ctx, client.FormDataConditionModel{
			AppID:   instance.AppID,
			DataID:  instance.FormInstanceID,
			TableID: instance.FormID,
			Ref:     formQueryRef,
		})
		if err != nil {
			return nil, err
		}
		// value = data[utils.Strval(valueOf)]
		value, _ = c.formAPI.GetValue(data, utils.Strval(valueOf), nil)
	} else if "formula" == valueFrom {
		data, err := c.formAPI.GetFormData(ctx, client.FormDataConditionModel{
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
			ret, err := c.structorAPI.CalExpression(ctx, map[string]interface{}{
				"expression": expression,
				"parameter":  c.flow.FormatFormValue2(formulaFields, variables),
			})
			if err != nil {
				logger.Logger.Error(err)
				return "", nil
			}
			return ret, nil
		}
		return valueOf, nil

	} else if "processVariable" == valueFrom {
		data, err := c.formAPI.GetFormData(ctx, client.FormDataConditionModel{
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
