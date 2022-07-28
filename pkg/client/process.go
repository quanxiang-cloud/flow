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

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/quanxiang-cloud/flow/pkg"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"net/http"
)

// Process process client
type Process interface {
	StartProcessInstance(ctx context.Context, req StartProcessReq) (*StartProcessResp, error)
	InitProcessInstance(ctx context.Context, req InitInstanceReq) (*StartProcessResp, error)

	// get all wait handle tasks
	GetTasksByInstanceID(ctx context.Context, instanceID string) ([]*ProcessTask, error)
	GetInstanceByID(ctx context.Context, processInstanceID string) (*ProcessInstance, error)
	GetInstanceByIDs(ctx context.Context, req *ListProcessReq) (*ListProcessResp, error)
	CheckActiveTask(ctx context.Context, instanceID string, taskID string, userID string) (*ProcessTask, error)
	GetTaskByID(ctx context.Context, taskID string) (*ProcessTask, error)

	GetTasks(ctx context.Context, req GetTasksReq) (*GetTasksResp, error)
	GetTasksCount(ctx context.Context, req GetTasksReq) (int64, error)
	GetHistoryTasks(ctx context.Context, req GetTasksReq) (*GetTasksResp, error)
	GetAllTasks(ctx context.Context, req GetTasksReq) (*GetTasksResp, error)

	GetDoneInstances(ctx context.Context, req GetTasksReq) (*GetInstancesResp, error)
	GetInstances(ctx context.Context, req GetTasksReq) (*GetInstancesResp, error)

	SendBack(ctx context.Context, req AddTaskReq) error
	StepBack(ctx context.Context, processInstanceID string, taskID string, nodeDefKey string) error
	SetAssignee(ctx context.Context, processInstanceID string, taskID string, assignee string) error
	SetDueDate(ctx context.Context, taskID string, dueDate string) error
	GetEarliestEntryThisTask(ctx context.Context, req GetTasksReq) (*ProcessTask, error)
	CompleteTask(ctx context.Context, processInstanceID string, taskID string, params map[string]interface{}, comments map[string]interface{}) error
	CompleteTaskToNode(ctx context.Context, processInstanceID string, taskID string, params map[string]interface{}, toNodeDefKey string, comments map[string]interface{}) error
	CompleteNonModelTask(ctx context.Context, taskIDs []string) error
	AbendInstance(ctx context.Context, processInstanceID string) error
	CompleteExecution(ctx context.Context, req *CompleteExecutionReq) (*CompleteExecutionResp, error)
	GetExecution(ctx context.Context, req *GateWayExecutionReq) (*GateWayExecutionResp, error)

	GetModelNode(ctx context.Context, processID string, nodeDefKey string) (*Node, error)
	AddNonModelTask(ctx context.Context, req *AddTaskReq) ([]ProcessTask, error)
	AddBeforeModelTask(ctx context.Context, req *AddTaskReq) (*AddTaskResp, error)
	AddAfterModelTask(ctx context.Context, req *AddTaskReq) (*AddTaskResp, error)
	AddHistoryTask(ctx context.Context, req *AddHistoryTaskReq) (string, error)
	SetProcessVariables(ctx context.Context, processID string, processInstanceID string, nodeDefKey string, k string, v interface{}) error
	GetProcessVariable(ctx context.Context, processInstanceID string, keys []string) (map[string]interface{}, error)

	DeployProcess(ctx context.Context, deployReq DeployReq) (*DeployResp, error)
	GetPreNode(ctx context.Context, taskID string) (*TaskPreNodeResp, error)
	InclusiveExecution(ctx context.Context, req ParentExecutionReq) (*ParentExecutionResp, error)

	UpdateAppStatus(ctx context.Context, processDefKeys []string, action string) error
	CompleteNode(ctx context.Context, req *CompleteNodeReq) error
	NodeInstanceList(ctx context.Context, req *NodeInstanceListReq) ([]*ProcessNodeInstanceVO, error)
}

type process struct {
	conf   *config.Configs
	client http.Client
}

// NewProcess new
func NewProcess(conf *config.Configs) Process {
	return &process{
		conf:   conf,
		client: NewClient(conf.InternalNet),
	}
}

// DeployProcess deploy process model
func (p *process) DeployProcess(ctx context.Context, deployReq DeployReq) (*DeployResp, error) {
	resp := &DeployResp{}
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/deploy")
	err := POST(ctx, &p.client, url, deployReq, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// StartProcessInstance start process instance
func (p *process) StartProcessInstance(ctx context.Context, req StartProcessReq) (*StartProcessResp, error) {
	resp := &StartProcessResp{}
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/startInstance")
	err := POST(ctx, &p.client, url, req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// InitProcessInstance init process instance
func (p *process) InitProcessInstance(ctx context.Context, req InitInstanceReq) (*StartProcessResp, error) {
	resp := &StartProcessResp{}
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/initInstance")
	err := POST(ctx, &p.client, url, req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetInstanceByIDs get instance by ids
func (p *process) GetInstanceByIDs(ctx context.Context, req *ListProcessReq) (*ListProcessResp, error) {
	resp := &ListProcessResp{}
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/listInstance")
	err := POST(ctx, &p.client, url, req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetInstanceByID get instance by id
func (p *process) GetInstanceByID(ctx context.Context, processInstanceID string) (*ProcessInstance, error) {
	req := &ListProcessReq{
		InstanceID: []string{processInstanceID},
	}

	resp, err := p.GetInstanceByIDs(ctx, req)
	if err != nil {
		return nil, err
	}

	if resp != nil && len(resp.Instances) > 0 {
		return resp.Instances[0], nil
	}

	return nil, nil
}

// GetTaskByID get active task
func (p *process) GetTaskByID(ctx context.Context, taskID string) (*ProcessTask, error) {
	req := GetTasksReq{
		TaskID: []string{taskID},
	}
	tasksResp, err := p.GetTasks(ctx, req)
	if err != nil {
		return nil, err
	}
	if len(tasksResp.Data) > 0 {
		return tasksResp.Data[0], nil
	}
	return nil, nil
}

// GetTasksByInstanceID get instace active tasks
func (p *process) GetTasksByInstanceID(ctx context.Context, instanceID string) ([]*ProcessTask, error) {
	req := GetTasksReq{
		InstanceID: []string{instanceID},
	}
	tasksResp, err := p.GetTasks(ctx, req)
	if err != nil {
		return nil, err
	}
	return tasksResp.Data, nil
}

// GetActiveTask check active task exsit
func (p *process) CheckActiveTask(ctx context.Context, instanceID string, taskID string, userID string) (*ProcessTask, error) {
	req := GetTasksReq{
		InstanceID: []string{instanceID},
		TaskID:     []string{taskID},
		Assignee:   userID,
	}
	tasksResp, err := p.GetTasks(ctx, req)
	if err != nil {
		return nil, err
	}
	if len(tasksResp.Data) > 0 {
		return tasksResp.Data[0], nil
	}
	return nil, nil
}

// GetInstances 查询流程实例，包含待办实例和已办实例
func (p *process) GetInstances(ctx context.Context, req GetTasksReq) (*GetInstancesResp, error) {
	resp := &GetInstancesResp{}
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/wholeInstance")
	err := POST(ctx, &p.client, url, req, resp)
	if err != nil {
		return nil, err
	}
	if resp.Data == nil {
		resp.Data = make([]*ProcessInstance, 0)
	}

	return resp, nil
}

// GetDoneInstances 查询已处理的流程实例
func (p *process) GetDoneInstances(ctx context.Context, req GetTasksReq) (*GetInstancesResp, error) {
	resp := &GetInstancesResp{}
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/doneInstance")
	err := POST(ctx, &p.client, url, req, resp)
	if err != nil {
		return nil, err
	}

	if resp.Data == nil {
		resp.Data = make([]*ProcessInstance, 0)
	}

	return resp, nil
}

// GetTasks 查询待办
func (p *process) GetTasks(ctx context.Context, req GetTasksReq) (*GetTasksResp, error) {
	resp := &GetTasksResp{}
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/agencyTask")
	err := POST(ctx, &p.client, url, req, resp)
	if err != nil {
		return nil, err
	}
	if resp.Data == nil {
		resp.Data = make([]*ProcessTask, 0)
	}

	return resp, nil
}

// GetTasksCount 查询待办数量
func (p *process) GetTasksCount(ctx context.Context, req GetTasksReq) (int64, error) {
	resp := &GetTasksCountResp{}
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/agencyTaskTotal")
	err := POST(ctx, &p.client, url, req, resp)
	if err != nil {
		return 0, err
	}

	return resp.Total, nil
}

// GetHistoryTasks 查询已办任务
func (p *process) GetHistoryTasks(ctx context.Context, req GetTasksReq) (*GetTasksResp, error) {
	resp := &GetTasksResp{}
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/doneTask")
	err := POST(ctx, &p.client, url, req, resp)
	if err != nil {
		return nil, err
	}
	if resp.Data == nil {
		resp.Data = make([]*ProcessTask, 0)
	}

	return resp, nil
}

// GetAllTasks 查询所有任务，包括待办和已办
func (p *process) GetAllTasks(ctx context.Context, req GetTasksReq) (*GetTasksResp, error) {
	resp := &GetTasksResp{}
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/wholeTask")
	err := POST(ctx, &p.client, url, req, resp)
	if err != nil {
		return nil, err
	}
	if resp.Data == nil {
		resp.Data = make([]*ProcessTask, 0)
	}

	return resp, nil
}

// GetEarliestEntryThisTask 最早进入该节点的任务
func (p *process) GetEarliestEntryThisTask(ctx context.Context, req GetTasksReq) (*ProcessTask, error) {
	resp := &GetTasksResp{}
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/wholeTask")
	err := POST(ctx, &p.client, url, req, resp)
	if err != nil {
		return nil, err
	}

	if len(resp.Data) > 0 {
		return resp.Data[0], nil
	}

	return nil, nil
}

// SendBack 打回重填
func (p *process) SendBack(ctx context.Context, req AddTaskReq) error {
	resp := new(interface{})
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/refillTask")
	err := POST(ctx, &p.client, url, req, resp)
	if err != nil {
		return err
	}
	return nil
}

// StepBack 可以回退到任意节点
func (p *process) StepBack(ctx context.Context, processInstanceID string, taskID string, nodeDefKey string) error {
	var resp interface{}
	req := AddTaskReq{
		TaskID:     taskID,
		UserID:     pkg.STDUserID(ctx),
		NodeDefKey: nodeDefKey,
	}
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/fallbackTask")
	err := POST(ctx, &p.client, url, req, resp)
	if err != nil {
		return err
	}
	return nil
}

// SetAssignee 设置任务处理人
func (p *process) SetAssignee(ctx context.Context, processInstanceID string, taskID string, assignee string) error {
	resp := new(interface{})
	req := AddTaskConditionReq{
		InstanceID: processInstanceID,
		TaskID:     taskID,
		UserID:     pkg.STDUserID(ctx),
		Assignee:   assignee,
	}
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/assigneeTask")
	err := POST(ctx, &p.client, url, req, resp)
	if err != nil {
		return err
	}
	return nil
}

// SetDueDate 设置任务有效期
func (p *process) SetDueDate(ctx context.Context, taskID string, dueDate string) error {
	resp := new(interface{})
	req := AddTaskConditionReq{
		TaskID:  taskID,
		UserID:  pkg.STDUserID(ctx),
		DueTime: dueDate,
	}
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/dueTimeTask")
	err := POST(ctx, &p.client, url, req, resp)
	if err != nil {
		return err
	}
	return nil
}

// CompleteTask 完成任务
func (p *process) CompleteTask(ctx context.Context, processInstanceID string, taskID string, params map[string]interface{}, comments map[string]interface{}) error {
	commentsStr := ""

	if comments != nil {
		commentsJSON, err := json.Marshal(comments)
		if err != nil {
			return err
		}
		commentsStr = string(commentsJSON)
	}

	resp := new(interface{})
	req := CompleteTaskReq{
		InstanceID: processInstanceID,
		TaskID:     taskID,
		UserID:     pkg.STDUserID(ctx),
		Params:     params,
		Comments:   commentsStr,
	}
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/completeTask")
	err := POST(ctx, &p.client, url, req, resp)
	if err != nil {
		return err
	}
	return nil
}

// CompleteTaskToNode 完成任务并且跳转到指定节点
func (p *process) CompleteTaskToNode(ctx context.Context, processInstanceID string, taskID string, params map[string]interface{}, toNodeDefKey string, comments map[string]interface{}) error {
	commentsStr := ""

	if comments != nil {
		commentsJSON, err := json.Marshal(comments)
		if err != nil {
			return err
		}
		commentsStr = string(commentsJSON)
	}

	resp := new(interface{})
	req := CompleteTaskReq{
		InstanceID:     processInstanceID,
		TaskID:         taskID,
		UserID:         pkg.STDUserID(ctx),
		Params:         params,
		NextNodeDefKey: toNodeDefKey,
		Comments:       commentsStr,
	}
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/completeTask")
	err := POST(ctx, &p.client, url, req, resp)
	if err != nil {
		return err
	}
	return nil
}

// CompleteNonModelTask 完成非模型任务
func (p *process) CompleteNonModelTask(ctx context.Context, taskIDs []string) error {
	resp := new(CompleteTaskResp)
	req := CompleteNonModelTaskReq{
		TaskID: taskIDs,
		UserID: pkg.STDUserID(ctx),
	}
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/completeNonModelTasks")
	err := POST(ctx, &p.client, url, req, resp)
	if err != nil {
		return err
	}
	return nil
}

// GetModelNode 获取流程模型中的节点
func (p *process) GetModelNode(ctx context.Context, processID string, nodeDefKey string) (*Node, error) {
	var resp []Node
	req := map[string]interface{}{
		"processID":  processID,
		"nodeDefKey": nodeDefKey,
	}
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/ListProcessNode")
	err := POST(ctx, &p.client, url, req, &resp)
	if err != nil {
		return nil, err
	}
	if resp == nil || len(resp) < 1 {
		return nil, nil
	}

	return &resp[0], nil
}

// AbendInstance 异常结束流程
func (p *process) AbendInstance(ctx context.Context, processInstanceID string) error {
	resp := new(interface{})
	req := DeleteProcessReq{
		InstanceID: processInstanceID,
		UserID:     pkg.STDUserID(ctx),
	}
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/terminatedInstance")
	err := POST(ctx, &p.client, url, req, resp)
	if err != nil {
		return err
	}
	return nil
}

// AddNonModelTask 新增非模型任务
func (p *process) AddNonModelTask(ctx context.Context, req *AddTaskReq) ([]ProcessTask, error) {
	resp := &AddTaskResp{}
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/addTask")
	err := POST(ctx, &p.client, url, req, resp)
	if err != nil {
		return nil, err
	}

	return resp.Tasks, nil
}

// AddBeforeModelTask 新增前置临时模型任务，用于前加签
func (p *process) AddBeforeModelTask(ctx context.Context, req *AddTaskReq) (*AddTaskResp, error) {
	resp := &AddTaskResp{}
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/addFrondTask")
	err := POST(ctx, &p.client, url, req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// AddAfterModelTask 新增后置临时模型任务，用于后加签
func (p *process) AddAfterModelTask(ctx context.Context, req *AddTaskReq) (*AddTaskResp, error) {
	resp := &AddTaskResp{}
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/addBackTask")
	err := POST(ctx, &p.client, url, req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// AddHistoryTask 添加一个历史任务（转交业务需要在已办中查询）
func (p *process) AddHistoryTask(ctx context.Context, req *AddHistoryTaskReq) (string, error) {
	resp := ""
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/addDoneTask")
	err := POST(ctx, &p.client, url, req, resp)
	if err != nil {
		return "", err
	}
	return resp, nil
}

//  SetProcessVariables set process variables
func (p *process) SetProcessVariables(ctx context.Context, processID string, processInstanceID string, nodeDefKey string, k string, v interface{}) error {
	resp := &SaveVariablesResp{}
	req := SaveVariablesReq{
		ProcessID:  processID,
		InstanceID: processInstanceID,
		UserID:     pkg.STDUserID(ctx),
		NodeDefKey: nodeDefKey,
		Key:        k,
		Value:      v,
	}
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/saveVariables")
	return POST(ctx, &p.client, url, req, resp)
}

// GetProcessVariable get process variable
func (p *process) GetProcessVariable(ctx context.Context, processInstanceID string, keys []string) (map[string]interface{}, error) {
	resp := map[string]interface{}{}
	req := GetVariablesReq{
		InstanceID: processInstanceID,
		Key:        keys,
	}
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/getVariables")
	err := POST(ctx, &p.client, url, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CompleteExecution CompleteExecution
func (p *process) CompleteExecution(ctx context.Context, req *CompleteExecutionReq) (*CompleteExecutionResp, error) {
	resp := &CompleteExecutionResp{}
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/completeExecution")
	err := POST(ctx, &p.client, url, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetExecution get execution info
func (p *process) GetExecution(ctx context.Context, req *GateWayExecutionReq) (*GateWayExecutionResp, error) {
	resp := &GateWayExecutionResp{}
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/gatewayExecution")
	err := POST(ctx, &p.client, url, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetPreNode func
func (p *process) GetPreNode(ctx context.Context, taskID string) (*TaskPreNodeResp, error) {
	req := make(map[string]interface{})
	req["taskID"] = taskID
	resp := &TaskPreNodeResp{}
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/preNode")
	err := POST(ctx, &p.client, url, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// InclusiveExecution func
func (p *process) InclusiveExecution(ctx context.Context, req ParentExecutionReq) (*ParentExecutionResp, error) {
	resp := &ParentExecutionResp{}
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/inclusiveExecution")
	err := POST(ctx, &p.client, url, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdateAppStatus update process model app status
func (p *process) UpdateAppStatus(ctx context.Context, processDefKeys []string, action string) error {
	req := AppDelReq{
		AppDefKey: processDefKeys,
		Action:    action,
	}
	resp := &AppDelResp{}
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/appOperation")
	return POST(ctx, &p.client, url, req, resp)
}

// CompleteNode 延时节点触发调用
func (p *process) CompleteNode(ctx context.Context, req *CompleteNodeReq) error {
	resp := &CompleteNodeResp{}
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/completeNode")
	return POST(ctx, &p.client, url, req, resp)
}

func (p *process) NodeInstanceList(ctx context.Context, req *NodeInstanceListReq) ([]*ProcessNodeInstanceVO, error) {
	var resp []*ProcessNodeInstanceVO
	url := fmt.Sprintf("%s%s", p.conf.APIHost.ProcessHost, "api/v1/process/nodeInstanceList")
	err := POST(ctx, &p.client, url, req, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// AppDelReq add delete update app status
type AppDelReq struct {
	AppDefKey []string `json:"defKey"`
	Action    string   `json:"action"`
}

// AppDelResp AppDelResp
type AppDelResp struct{}

// ListProcessReq list process request
type ListProcessReq struct {
	// Page             int      `json:"page"`
	// Limit            int      `json:"limit"`
	Name             string   `json:"name"`
	ProcessID        []string `json:"processID"`
	InstanceID       []string `json:"instanceID"`
	ParentInstanceID []string `json:"parentInstanceID"`
	ProcessStatus    string   `json:"processStatus"`
}

// ListProcessResp list process request
type ListProcessResp struct {
	Instances []*ProcessInstance `json:"instances"`
}

// GateWayExecutionReq GateWayExecutionReq
type GateWayExecutionReq struct {
	TaskID string `json:"taskID" binding:"required"`
	// UserID string `json:"userID"`
}

// GateWayExecutionResp GateWayExecutionResp
type GateWayExecutionResp struct {
	Executions         []string `json:"executionIds"`
	CurrentExecutionID string   `json:"currentExecutionID"`
}

// CompleteExecutionReq CompleteExecutionReq
type CompleteExecutionReq struct {
	ExecutionID []string `json:"executionID"`
	UserID      string   `json:"userID"`
	NextDefKey  string   `json:"nextDefKey"`
	TaskID      string   `json:"taskID"`
	Comments    string   `json:"comments"`
}

// CompleteExecutionResp CompleteExecutionResp
type CompleteExecutionResp struct{}

// InitInstanceReq init instance request
type InitInstanceReq struct {
	InstanceID string                 `json:"instanceID"`
	UserID     string                 `json:"userID"`
	Params     map[string]interface{} `json:"params"`
}

// GetVariablesReq get variables request
type GetVariablesReq struct {
	InstanceID string   `json:"instanceID"`
	Key        []string `json:"keys"`
}

// SaveVariablesReq save variables request
type SaveVariablesReq struct {
	ProcessID  string      `json:"processID"`
	InstanceID string      `json:"instanceID"`
	NodeDefKey string      `json:"nodeDefKey"`
	UserID     string      `json:"userID"`
	Key        string      `json:"key"`
	Value      interface{} `json:"value"`
}

// SaveVariablesResp save variables
type SaveVariablesResp struct{}

// CompleteNonModelTaskReq CompleteNonModelTaskReq
type CompleteNonModelTaskReq struct {
	TaskID []string `json:"tasks" binding:"required"`
	UserID string   `json:"userID"`
}

// CompleteTaskResp CompleteTaskResp
type CompleteTaskResp struct{}

// AddHistoryTaskReq AddHistoryTaskReq
type AddHistoryTaskReq struct {
	Name        string `json:"name"`
	Desc        string `json:"desc"`
	UserID      string `json:"userID"`
	Assignee    string `json:"assignee"`
	TaskID      string `json:"taskID"`
	NodeDefKey  string `json:"nodeDefKey"`
	InstanceID  string `json:"instanceID"`
	ExecutionID string `json:"executionID"`
}

// DeleteProcessReq delete process request
type DeleteProcessReq struct {
	InstanceID string `json:"instanceID"`
	UserID     string `json:"userID"`
}

// AddTaskReq AddTaskReq
type AddTaskReq struct {
	Name       string    `json:"name"`
	Desc       string    `json:"desc"`
	UserID     string    `json:"userID"`
	Assignee   []string  `json:"assignee"`
	TaskID     string    `json:"taskID" binding:"required"`
	NodeDefKey string    `json:"nodeDefKey"`
	InstanceID string    `json:"instanceID"`
	Node       *NodeData `json:"event"`
}

// NodeData event json data
type NodeData struct {
	DefKey     string          `json:"id" binding:"required"`
	Name       string          `json:"name"`
	Type       string          `json:"type" binding:"required"`
	SubModel   *ModelData      `json:"subModel"`
	Desc       string          `json:"desc"`
	UserIDs    []string        `json:"userIds"`
	GroupIDs   []string        `json:"groupIds"`
	Variable   string          `json:"variable"`
	NextNodes  []*NodeLinkData `json:"nextNodes"`
	InstanceID string          `json:"-"`
}

// NodeLinkData event link json data
type NodeLinkData struct {
	NodeID    string `json:"id"`
	Condition string `json:"condition"`
}

// ModelData model json data
type ModelData struct {
	DefKey string      `json:"id"`
	Name   string      `json:"name"`
	Nodes  []*NodeData `json:"nodes"`
}

// AddTaskResp AddTaskResp
type AddTaskResp struct {
	Tasks []ProcessTask `json:"task"`
}

// AddTaskConditionReq AddTaskConditionReq
type AddTaskConditionReq struct {
	UserID     string `json:"userID"`
	Assignee   string `json:"assignee"`
	TaskID     string `json:"taskID"`
	InstanceID string `json:"instanceID"`
	DueTime    string `json:"dueTime"`
}

// CompleteTaskReq CompleteTaskReq
type CompleteTaskReq struct {
	InstanceID     string                 `json:"instanceID"`
	TaskID         string                 `json:"taskID"`
	UserID         string                 `json:"userID"`
	NextNodeDefKey string                 `json:"nextNodeDefKey"`
	Params         map[string]interface{} `json:"params"`
	Comments       string                 `json:"comments"`
}

// DeployReq req
type DeployReq struct {
	Model     string `json:"model"`
	CreatorID string `json:"creatorId"`
}

// DeployResp resp
type DeployResp struct {
	ID string `json:"ID"`
}

// GetTasksCountResp resp
type GetTasksCountResp struct {
	Total int64 `json:"total"`
}

// StartProcessReq request param
type StartProcessReq struct {
	ProcessID string                 `json:"processID"`
	UserID    string                 `json:"userID"`
	Params    map[string]interface{} `json:"params"`
}

// StartProcessResp response param
type StartProcessResp struct {
	InstanceID string `json:"instanceID"`
}

// GetTasksReq query task request params
type GetTasksReq struct {
	Page       int          `json:"page"`
	Limit      int          `json:"limit"`
	Desc       []string     `json:"desc"`
	Name       string       `json:"taskName"`
	Order      []QueryOrder `json:"orders"`
	ProcessID  []string     `json:"processID"`
	InstanceID []string     `json:"instanceID"`
	Assignee   string       `json:"assignee"`
	DueTime    string       `json:"dueTime"`
	NodeDefKey string       `json:"nodeDefKey"`
	TaskID     []string     `json:"taskID"`
	Status     string       `json:"status"`
}

// GetTasksResp resp
type GetTasksResp struct {
	PageSize    int            `json:"-"`
	TotalCount  int64          `json:"total"`
	TotalPage   int            `json:"-"`
	CurrentPage int            `json:"-"`
	StartIndex  int            `json:"-"`
	Data        []*ProcessTask `json:"data"`
}

// QueryOrder QueryOrder
type QueryOrder struct {
	OrderType string `json:"orderType"`
	Column    string `json:"column"`
}

// ProcessTask process task model
type ProcessTask struct {
	ID              string `json:"id"`
	ProcID          string `json:"procId"`
	ProcInstanceID  string `json:"procInstanceId"`
	ExecutionID     string `json:"executionId"`
	NodeID          string `json:"nodeId"`
	NodeDefKey      string `json:"nodeDefKey"`
	NextNodeDefKey  string `json:"nextNodeDefKey"`
	Name            string `json:"name"`
	Desc            string `json:"desc"`     // 例如：审批、填写、抄送、阅示
	TaskType        string `json:"taskType"` // Model模型任务、TempModel临时模型任务、NonModel非模型任务
	Assignee        string `json:"assignee"`
	Status          string `json:"status"` // COMPLETED, ACTIVE
	DueTime         string `json:"dueTime"`
	EndTime         string `json:"endTime"`
	CreatorID       string `json:"creatorId"`
	CreateTime      string `json:"createTime"`
	ModifierID      string `json:"modifierId"`
	ModifyTime      string `json:"modifyTime"`
	Comments        string `json:"comments"`
	TenantID        string `json:"tenantId"`
	NodeInstanceID  string `json:"tenantId"`
	NodeInstancePid string `json:"tenantId"`
}

// ProcessInstance info
type ProcessInstance struct {
	ID         string `json:"id"`
	ProcID     string `json:"procId"`
	Name       string `json:"name"`
	PID        string `json:"pId"`
	Status     string `json:"status"`    // COMPLETED, ACTIVE
	AppStatus  string `json:"appStatus"` // ACTIVE,SUSPEND
	EndTime    string `json:"endTime"`
	CreatorID  string `json:"creatorId"`
	CreateTime string `json:"createTime"`
	ModifierID string `json:"modifierId"`
	ModifyTime string `json:"modifyTime"`
	TenantID   string `json:"tenantId"`
}

// GetInstancesReq req
type GetInstancesReq struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	// Desc       []string     `json:"desc"`
	Name       string       `json:"taskName"`
	Order      []QueryOrder `json:"orders"`
	ProcessID  []string     `json:"processID"`
	InstanceID []string     `json:"instanceID"`
	Assignee   string       `json:"assignee"`
	// DueTime    int64        `json:"dueTime"`
	// NodeDefKey string       `json:"nodeDefKey"`
	// TaskID     []string     `json:"taskID"`
	ProcessStatus string `json:"processStatus"`
}

// GetInstancesResp resp
type GetInstancesResp struct {
	PageSize    int                `json:"-"`
	TotalCount  int64              `json:"total"`
	TotalPage   int                `json:"-"`
	CurrentPage int                `json:"-"`
	StartIndex  int                `json:"-"`
	Data        []*ProcessInstance `json:"data"`
}

// Node info
type Node struct {
	ID             string `json:"id"`
	ProcID         string `json:"procId"`
	ProcInstanceID string `json:"procInstanceId"`
	Name           string `json:"name"`
	DefKey         string `json:"defKey"`
	NodeType       string `json:"nodeType"`   // Start、End、User、MultiUser、Service、Script、ParallelGateway、InclusiveGateway、SubProcess
	SubProcID      string `json:"subProcId"`  // Type is SubProcess
	PairDefKey     string `json:"pairDefKey"` // ParallelGateway < == > InclusiveGateway
	Desc           string `json:"desc"`
	CreatorID      string `json:"creatorId"`
	CreateTime     string `json:"createTime"`
	ModifierID     string `json:"modifierId"`
	ModifyTime     string `json:"modifyTime"`
	TenantID       string `json:"tenantId"`
}

// NodeInfo NodeInfo
type NodeInfo struct {
	Name       string `json:"name"`
	NodeDefKey string `json:"nodeDefKey"`
	NodeType   string `json:"nodeType"`
}

// TaskPreNodeResp TaskPreNodeResp
type TaskPreNodeResp struct {
	Nodes []*NodeInfo `json:"nodes"`
}

// ParentExecutionReq ParentExecutionReq
type ParentExecutionReq struct {
	TaskID string `json:"taskID" binding:"required"`
	DefKey string `json:"defKey"`
}

// ParentExecutionResp ParentExecutionResp
type ParentExecutionResp struct {
	ExecutionID string `json:"executionID"`
}

// CompleteNodeReq process延时节点回调请求
type CompleteNodeReq struct {
	ProcessID   string                 `json:"processID" binding:"required"`
	InstanceID  string                 `json:"instanceID" binding:"required"`
	NodeDefKey  string                 `json:"nodeDefKey" binding:"required"`
	ExecutionID string                 `json:"executionID" binding:"required"`
	NextNodes   string                 `json:"nextNodes"`
	Params      map[string]interface{} `json:"params"`
	UserID      string                 `json:"userID"`
}

// CompleteNodeResp 延时节点响应
type CompleteNodeResp struct {
}

// NodeInstanceListReq req
type NodeInstanceListReq struct {
	ProcInstanceID string `json:"procInstanceId"`
}

// ProcessNodeInstance entity
type ProcessNodeInstance struct {
	ID             string `json:"id"`
	ProcID         string `json:"procId"`
	ProcInstanceID string `json:"procInstanceId"`
	PID            string `json:"pId"`
	ExecutionID    string `json:"executionId"`
	NodeDefKey     string `json:"nodeDefKey"`
	NodeName       string `json:"nodeName"`
	NodeType       string `json:"nodeType"`
	TaskID         string `json:"taskId"`
	Comments       string `json:"comments"`
	CreatorID      string `json:"creatorId"`
	CreateTime     string `json:"createTime"`
	ModifierID     string `json:"modifierId"`
	ModifyTime     string `json:"modifyTime"`
	TenantID       string `json:"tenantId"`
}

// ProcessNodeInstanceVO vo
type ProcessNodeInstanceVO struct {
	*ProcessNodeInstance
	Assignee string `json:"assignee"`
}
