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

package convert

import (
	"encoding/json"
	"github.com/quanxiang-cloud/flow/internal"
	"github.com/quanxiang-cloud/flow/pkg/misc/error2"
	"github.com/quanxiang-cloud/flow/pkg/misc/logger"
	"github.com/quanxiang-cloud/flow/pkg/utils"
	"reflect"
	"strings"
)

const (
	// OrApproval type
	OrApproval = "OR_APPROVAL" // 或签
	// AndApproval type
	AndApproval = "AND_APPROVAL" // 会签
	// OrFillIn type
	OrFillIn = "OR_FILLIN" // 任填
	// AndFillIn type
	AndFillIn = "AND_FILLIN" // 全填
	// START type
	START = "START" // 开始
	// END type
	END = "END" // 结束

	// AssigneeList dynamic user list key
	AssigneeList = "assigneeList"
)

// Process Node Type
const (
	Start            = "Start"
	User             = "User"
	MultiUser        = "MultiUser"
	ParallelGateway  = "ParallelGateway"
	InclusiveGateway = "InclusiveGateway"
	End              = "End"
	Service          = "Service"
	Script           = "Script"
	SubProcess       = "SubProcess"
)

// 延时节点的延时类型
const (
	DelayTypeOfTableColumn = "tableColumn"
	DelayTypeOfaTime       = "aTime"
	DelayTypeOfSpecTime    = "specTime"
)

// 邮件节点类型
const (
	// EmailTypeOfField 人员选择字段(1.1.1版本前默认)
	EmailTypeOfField = "field"
	// EmailTypeOfMultipleField 多个字段
	EmailTypeOfMultipleField = "multipleField"
	// EmailTypeOfSuperior 上级领导
	EmailTypeOfSuperior = "superior"
	// EmailTypeOfLeadOfDepartment 部门负责人
	EmailTypeOfLeadOfDepartment = "leadOfDepartment"
	// EmailTypeOfProcessInitiator 发起人
	EmailTypeOfProcessInitiator = "processInitiator"
	// EmailTypeOfPerson 人员
	EmailTypeOfPerson = "person"
)

// 前端组件字段类型
const (
	// CompOfUserPicker 人员选择器
	CompOfUserPicker = "UserPicker"
	// CompOfRadioGroup 单选
	CompOfRadioGroup = "RadioGroup"
	// CompOfTextarea 多行文本
	CompOfTextarea = "Textarea"
	// CompOfCheckboxGroup 复选框
	CompOfCheckboxGroup = "CheckboxGroup"
	// CompOfInput input输入框
	CompOfInput = "Input"
	// CompOfSelect 下拉框
	CompOfSelect = "Select"
	// CompOfMultipleSelect 下拉复选框
	CompOfMultipleSelect = "MultipleSelect"
)

// dispatcher回调类型
const (
	// CallbackOfDelay 延时
	CallbackOfDelay = "delayed"
	// CallbackOfUrge 催单
	CallbackOfUrge = "urge"
	// CallbackOfCron 定时
	CallbackOfCron = "cron"
)

// LowCode flow Chart Type
const (
	FormData = "formData"

	FormTime = "FORM_TIME"
	// Approve task
	Approve = "approve"
	// FillIn task
	FillIn = "fillIn"
	// Email service task
	Email = "email"
	// Autocc service task
	Autocc = "autocc"
	// Letter service task
	Letter = "letter"
	// ProcessVariableAssignment service task
	ProcessVariableAssignment = "processVariableAssignment"
	// TableDataCreate service task
	TableDataCreate = "tableDataCreate"
	// TableDataUpdate service task
	TableDataUpdate = "tableDataUpdate"
	// WebHook service task
	WebHook = "webhook"
	// Delayed 延时节点
	Delayed = "delayed"

	end                 = "end"
	processBranchSource = "processBranchSource"
	processBranch       = "processBranch"
	processBranchTarget = "processBranchTarget"
	plus                = "plus"
	step                = "step"
	PauseExecution      = "pauseExecution"
)

// Process task type
const (
	// ReviewTask 审批任务
	ReviewTask = "REVIEW"
	// WriteTask 填写任务
	WriteTask = "WRITE"
	// ReadTask 阅示任务
	ReadTask = "READ"
	// CcTask 抄送任务
	CcTask = "CC"
	// SendBackTask 打回重填任务
	SendBackTask = "SEND_BACK"
)

// 人员'person' | 表单字段'field' | 岗位'position' | 上级领导'superior' | 部门负责人'leadOfDepartment'
const (
	Person           = "person"
	Field            = "field"
	Position         = "position"
	Superior         = "superior"
	LeadOfDepartment = "leadOfDepartment"
)

// webHook inputs type
const (
	Body   = "body"
	Path   = "path"
	Header = "header"
	Query  = "query"
)

// ToProcessModel func
func ToProcessModel(chartJSON string) (*ProcessModel, error) {
	p := &ProcessModel{}
	err := json.Unmarshal([]byte(chartJSON), p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// GetShapeByChartType func
func GetShapeByChartType(chartJSON string, chartType string) (*ShapeModel, error) {
	p, err := ToProcessModel(chartJSON)
	if err != nil {
		return nil, err
	}
	for _, s := range p.Shapes {
		if chartType == s.Type {
			return &s, nil
		}
	}
	return nil, err
}

// GetFormInfoFromShapes get form info
func GetFormInfoFromShapes(shapes []ShapeModel) string {
	var shape ShapeModel
	for _, s := range shapes {
		if "formData" == s.Type {
			shape = s
			break
		}
	}

	bd := shape.Data.BusinessData
	var form map[string]interface{}
	if v := bd["form"]; v != nil {
		form = v.(map[string]interface{})
	}
	if form == nil {
		return ""
	}
	// value = from的id
	return utils.Strval(form["value"])
}

// GetBusinessDataByType func
func GetBusinessDataByType(text string, nodeType string) (*map[string]interface{}, error) {
	s, err := GetShapeByChartType(text, nodeType)
	if err != nil {
		return nil, err
	}
	return &s.Data.BusinessData, nil
}

// GetShapeByTaskDefKey func
func GetShapeByTaskDefKey(flowJSON string, taskDefKey string) (*ShapeModel, error) {
	processModel := &ProcessModel{}
	err := json.Unmarshal([]byte(flowJSON), processModel)
	if err != nil {
		return nil, err
	}

	// 加签节点ID如下 currentNodeID+":"+newNodeID
	if len(taskDefKey) > 0 && strings.Contains(taskDefKey, ":") {
		taskDefKey = strings.Split(taskDefKey, ":")[0]
	}

	if len(processModel.Shapes) > 0 {
		for _, shapeModel := range processModel.Shapes {
			if shapeModel.ID == taskDefKey {
				return &shapeModel, nil
			}
		}
	}

	return nil, nil
}

// GetTaskBasicConfigModel get shape basic config model
func GetTaskBasicConfigModel(shape *ShapeModel) *TaskBasicConfigModel {
	businessData := shape.Data.BusinessData
	taskBasicConfigModel := &TaskBasicConfigModel{}
	basicConfigJSON, err := json.Marshal(businessData["basicConfig"])
	if err != nil {
		logger.Logger.Error(err)
		return nil
	}
	err = json.Unmarshal(basicConfigJSON, taskBasicConfigModel)
	if err != nil {
		logger.Logger.Error(err)
		return nil
	}
	return taskBasicConfigModel
}

// GetValueFromBusinessData func
func GetValueFromBusinessData(s ShapeModel, keyInfo string) interface{} {
	data := s.Data.BusinessData
	if len(data) == 0 {
		return ""
	}
	keys := strings.Split(keyInfo, ".")
	var tmp interface{}
	for _, key := range keys {
		tmp = data[key]
		if tmp == nil {
			return nil
		}
		if reflect.TypeOf(tmp).Kind() == reflect.Map {
			data = tmp.(map[string]interface{})
		}

	}
	return tmp
}

// SwitchToNodeType func
func switchToNodeType(s ShapeModel) string {
	switch s.Type {
	case FormData, FormTime:
		return Start
	case end:
		return End
	case processBranchSource:
		return ParallelGateway
	case processBranchTarget:
		return InclusiveGateway
	case Approve:
		fallthrough
	case FillIn:
		rt := GetValueFromBusinessData(s, "basicConfig.multiplePersonWay").(string)
		if rt == "or" {
			return User
		}
		return MultiUser
	case Email:
		fallthrough
	case Autocc:
		fallthrough
	case Letter:
		fallthrough
	case ProcessVariableAssignment:
		fallthrough
	case TableDataCreate:
		fallthrough
	case TableDataUpdate:
		fallthrough
	case Delayed:
		fallthrough
	case WebHook:
		return Service
	}

	return ""
}

func getDesc(chartType string) string {
	switch chartType {
	case Approve:
		return ReviewTask
	case FillIn:
		return WriteTask
	}
	return ""
}

// GenProcessNodes info
func GenProcessNodes(text string) ([]*ProcessNodeModel, interface{}, error) {
	formulaFields := make(map[string]interface{})
	p, err := ToProcessModel(text)
	if err != nil {
		return nil, nil, err
	}

	// {node_id: node_value, ......}
	nodeMap := make(map[string]ShapeModel)
	for _, s := range p.Shapes {
		nodeMap[s.ID] = s
	}

	formID := GetFormInfoFromShapes(p.Shapes)

	nodes := make([]*ProcessNodeModel, 0)
	// 确定每个节点的子节点
	for _, s := range p.Shapes {
		if processBranch == s.Type || plus == s.Type || step == s.Type {
			continue
		}
		nextNodes := make([]*ProcessNextNodeModel, 0)
		childrenIDs := s.Data.NodeData.ChildrenID
		for _, childrenID := range childrenIDs {
			condition := ""
			// 如果是分流节点，则该节点将会有多个子节点
			if processBranchSource == s.Type {
				conditionNode := nodeMap[childrenID]
				ignore := conditionNode.Data.BusinessData["ignore"]
				if v := conditionNode.Data.BusinessData["rule"]; v != nil {
					condition = v.(string)
					for k, v := range utils.ChangeObjectToMap(conditionNode.Data.BusinessData["formulaFields"]) {
						formulaFields[k] = v
					}
				} else if utils.Strval(ignore) == "true" {
					condition = "defaultBranch"
				} else {
					return nil, nil, error2.NewErrorWithString(error2.Internal, "Condition not set")
				}
				nextNode := nodeMap[conditionNode.Data.NodeData.ChildrenID[0]]
				condition = strings.Replace(condition, "$", "", -1)
				pnn := &ProcessNextNodeModel{
					ID:        nextNode.ID,
					Condition: condition,
				}
				nextNodes = append(nextNodes, pnn)
			} else {
				err = checkChart(s, formID)
				if err != nil {
					return nil, nil, err
				}

				nextNode := nodeMap[childrenID]
				if v := nextNode.Data.BusinessData["rule"]; v != nil {
					condition = v.(string)
					for k, v := range utils.ChangeObjectToMap(nextNode.Data.BusinessData["formulaFields"]) {
						formulaFields[k] = v
					}
				}
				if Email == s.Type {
					for k, v := range utils.ChangeObjectToMap(s.Data.BusinessData["formulaFields"]) {
						formulaFields[k] = v
					}
				}
				condition = strings.Replace(condition, "$", "", -1)
				pnn := &ProcessNextNodeModel{
					ID:        nextNode.ID,
					Condition: condition,
				}
				nextNodes = append(nextNodes, pnn)
			}
		}
		node := &ProcessNodeModel{
			ID:        s.ID,
			Name:      s.Data.NodeData.Name,
			Type:      switchToNodeType(s),
			NextNodes: nextNodes,
			Desc:      getDesc(s.Type),
			UserIds:   getUsers(s),
			GroupIds:  getGroups(s),
			Variable:  getVariable(s),
		}

		if node.Variable != "" { // 变量和固定的人只能二选一
			node.UserIds = nil
			node.GroupIds = nil
		}

		nodes = append(nodes, node)
	}
	return nodes, formulaFields, err
}

func checkChart(s ShapeModel, formID string) error {
	if FillIn == s.Type || Approve == s.Type {
		// 人员'person' | 表单字段'field' | 岗位'position' | 上级领导'superior' | 部门负责人'leadOfDepartment'
		t := GetValueFromBusinessData(s, "basicConfig.approvePersons.type")
		if t == "" {
			return error2.NewErrorWithString(error2.Internal, "Handler not set")
		}

		switch t {
		case Person:
			{
				users := GetValueFromBusinessData(s, "basicConfig.approvePersons.users")
				roles := GetValueFromBusinessData(s, "basicConfig.approvePersons.roles")
				departments := GetValueFromBusinessData(s, "basicConfig.approvePersons.departments")
				if (users == nil || reflect.ValueOf(users).Len() == 0) && (roles == nil || reflect.ValueOf(roles).Len() == 0) && (departments == nil || reflect.ValueOf(departments).Len() == 0) {
					return error2.NewErrorWithString(error2.Internal, "审批人未填写")
				}
			}
		case Field:
			{
				fieldsObj := GetValueFromBusinessData(s, "basicConfig.approvePersons.fields")
				if fieldsObj == nil || reflect.ValueOf(fieldsObj).Len() == 0 {
					return error2.NewErrorWithString(error2.Internal, "审批人未填写")
				}
				return nil
			}
		}
	} else if TableDataUpdate == s.Type {
		targetTableID := GetValueFromBusinessData(s, "targetTableId")
		if targetTableID != formID {
			conditions := GetValueFromBusinessData(s, "filterRule.conditions")
			if conditions == nil || reflect.ValueOf(conditions).Len() == 0 {
				return error2.NewErrorWithString(error2.Internal, "conditions not set")
			}
		}

		updateRule := GetValueFromBusinessData(s, "updateRule")
		if updateRule == nil || reflect.ValueOf(updateRule).Len() == 0 {
			return error2.NewErrorWithString(error2.Internal, "updateRule not set")
		}
	}

	return nil
}

func getUsers(s ShapeModel) []string {
	userIDs := make([]string, 0)
	users := GetValueFromBusinessData(s, "basicConfig.approvePersons.users")
	if users == nil || users == "" {
		return userIDs
	}
	for _, user := range users.([]interface{}) {
		id := user.(map[string]interface{})["id"]
		userIDs = append(userIDs, id.(string))
	}
	return userIDs
}

func getGroups(s ShapeModel) []string {
	groups := make([]string, 0)
	roles := GetValueFromBusinessData(s, "basicConfig.approvePersons.roles")
	if roles != nil && roles != "" {
		for _, role := range roles.([]interface{}) {
			roleID := role.(map[string]interface{})["id"]
			groups = append(groups, utils.StringJoins(internal.Role, "_", utils.Strval(roleID)))
		}
	}
	departments := GetValueFromBusinessData(s, "basicConfig.approvePersons.departments")
	if departments != nil && departments != "" {
		for _, department := range departments.([]interface{}) {
			departmentID := department.(map[string]interface{})["id"]
			groups = append(groups, utils.StringJoins(internal.Dep, "_", utils.Strval(departmentID)))
		}
	}
	return groups
}

func getVariable(s ShapeModel) string {
	//  上级领导'superior' | 部门负责人'leadOfDepartment'
	t := GetValueFromBusinessData(s, "basicConfig.approvePersons.type")
	if "superior" == t || "leadOfDepartment" == t || "field" == t || "processInitiator" == t {
		return AssigneeList
	}

	return ""
}

// GetCurrentNodeType func
func GetCurrentNodeType(t string, multiplePersonWay string) string {
	if len(t) == 0 || len(multiplePersonWay) == 0 {
		return OrApproval
	}

	if t == Approve {
		if "or" == multiplePersonWay {
			return OrApproval
		} else if "and" == multiplePersonWay {
			return AndApproval
		}
	} else if t == FillIn {
		if "or" == multiplePersonWay {
			return OrFillIn
		} else if "and" == multiplePersonWay {
			return AndFillIn
		}
	}

	return OrApproval
}

// FindNextNode func
func FindNextNode(currentNode *ShapeModel, chartJSON string) (string, bool) {
	nextNodeDefKey := currentNode.Data.NodeData.BranchTargetElementID
	if nextNodeDefKey == "" {
		endNode, _ := GetShapeByChartType(chartJSON, "end")
		return endNode.ID, true
	}
	nextNode, _ := GetShapeByTaskDefKey(chartJSON, nextNodeDefKey)
	// true : 任一分支拒绝结束流程 false：所有分支拒绝结束流程
	processBranchEndStrategy := GetValueFromBusinessData(*nextNode, "processBranchEndStrategy").(string) == "any"
	if !processBranchEndStrategy {
		return nextNode.ID, false
	}
	return FindNextNode(nextNode, chartJSON)
}

// TaskBasicConfigModel struct
type TaskBasicConfigModel struct {
	ApprovePersons    ApprovePersonsModel `json:"approvePersons"`    // 审批人
	MultiplePersonWay string              `json:"multiplePersonWay"` // 多人审批时
	// 无审批人时whenNoPerson：自动跳过该节点skip,转交给管理员transferAdmin
	WhenNoPerson string `json:"whenNoPerson"`
	// 自动审批通过规则autoRules：审批人为发起人时origin，审批人与上一节点审批人相同时parent，审批人与前置节点（非上一节点审批人相同时）previous
	AutoRules []string `json:"autoRules"`
	// 审批用时限制
	TimeRule TaskTimeRuleModel `json:"timeRule"`
}

// ApprovePersonsModel struct
type ApprovePersonsModel struct {
	// 人员'person' | 表单字段'field' | 岗位'position' | 上级领导'superior' | 部门负责人'leadOfDepartment'
	Type        string                   `json:"type"`
	Users       []map[string]interface{} `json:"users"`
	Roles       []map[string]interface{} `json:"roles"`
	Departments []map[string]interface{} `json:"departments"`
	Positions   []string                 `json:"positions"` // 岗位id集合
	Fields      []string                 `json:"fields"`    // 表单字段集合
}

// TaskTimeRuleModel struct
type TaskTimeRuleModel struct {
	Enabled     bool              `json:"enabled"`
	DeadLine    DeadLineModel     `json:"deadLine"`
	WhenTimeout TypeAndValueModel `json:"whenTimeout"`
}

// DeadLineModel struct
type DeadLineModel struct {
	DayHourMinuteModel
	// entry进入该节点后,firstEntry首次进入该节点后,flowWorked工作流开始后
	BreakPoint string    `json:"breakPoint"`
	Urge       UrgeModel `json:"urge"`
}

// TypeAndValueModel struct
type TypeAndValueModel struct {
	// noDealWith不处理,autoDealWith自动处理,jump跳转至其他节点
	Type  string `json:"type"`
	Value string `json:"value"`
}

// UrgeModel struct
type UrgeModel struct {
	DayHourMinuteModel
	Repeat DayHourMinuteModel `json:"repeat"`
}

// DayHourMinuteModel struct
type DayHourMinuteModel struct {
	Day     int `json:"day"`
	Hours   int `json:"hours"`
	Minutes int `json:"minutes"`
}

// FieldPermissionModel field permission model
type FieldPermissionModel struct {
	FieldName    string             `json:"fieldName"`
	XInternal    XInternalModel     `json:"x-internal"`
	InitialValue FieldValueSetModel `json:"initialValue"`
	SubmitValue  FieldValueSetModel `json:"submitValue"`
}

// XInternalModel permission
type XInternalModel struct {
	//  editable hidden write read 4位二进制表示，后端关注write和read
	Permission int8 `json:"permission"`
}

// FieldValueSetModel struct
type FieldValueSetModel struct {
	Variable    interface{} `json:"variable"`
	StaticValue interface{} `json:"staticValue"`
}

// OperatorPermissionModel struct
type OperatorPermissionModel struct {
	Custom []OperatorPermissionItemModel `json:"custom"`
	System []OperatorPermissionItemModel `json:"system"`
}

// OperatorPermissionItemModel struct
type OperatorPermissionItemModel struct {
	Enabled        bool   `json:"enabled"`
	Changeable     bool   `json:"changeable"`
	Name           string `json:"name"`
	DefaultText    string `json:"defaultText"`
	Text           string `json:"text"`
	Only           string `json:"only"`
	ReasonRequired bool   `json:"reasonRequired"`
	Value          string `json:"value"`
}

// ProcessModel struct
type ProcessModel struct {
	Version string       `json:"version"`
	Shapes  []ShapeModel `json:"shapes"`
}

// ShapeModel struct
type ShapeModel struct {
	ID     string         `json:"id"`
	Type   string         `json:"type"`
	Data   ShapeDataModel `json:"data"`
	Source string         `json:"source"`
	Target string         `json:"target"`
}

// ShapeDataModel info
type ShapeDataModel struct {
	NodeData     NodeDataModel          `json:"nodeData"`
	BusinessData map[string]interface{} `json:"businessData"`
}

// NodeDataModel info
type NodeDataModel struct {
	Name                  string   `json:"name"`
	BranchTargetElementID string   `json:"branchTargetElementID"`
	ChildrenID            []string `json:"childrenID"`
}

// ProcessNodeModel info
type ProcessNodeModel struct {
	ID        string                  `json:"id"`
	Name      string                  `json:"name"`
	Type      string                  `json:"type"`
	Desc      string                  `json:"desc"`
	Variable  string                  `json:"variable"`
	UserIds   []string                `json:"userIds"`
	GroupIds  []string                `json:"groupIds"`
	NextNodes []*ProcessNextNodeModel `json:"nextNodes"`
}

// ProcessNextNodeModel info
type ProcessNextNodeModel struct {
	ID        string `json:"id"`
	Condition string `json:"condition"`
}

// Input struct
type Input struct {
	Type      string      `json:"type"`
	FieldType string      `json:"fieldType"`
	Name      string      `json:"name"`
	Data      interface{} `json:"data"`
	In        string      `json:"in"`
	Title     string      `json:"title"`
	FieldName string      `json:"fieldName"`
	TableID   string      `json:"tableID"`
}
