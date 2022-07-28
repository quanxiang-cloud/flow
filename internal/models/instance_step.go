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

package models

import "gorm.io/gorm"

// InstanceStep info
type InstanceStep struct {
	BaseModel

	ProcessInstanceID string             `json:"processInstanceId"`
	TaskID            string             `json:"taskId"`
	TaskType          string             `json:"taskType"` // 节点类型：或签、会签、任填、全填、开始、结束
	TaskDefKey        string             `json:"taskDefKey"`
	TaskName          string             `json:"taskName"`
	HandleUserIDs     string             `json:"handleUserIds"`
	Status            string             `json:"status"` // 步骤处理结果，通过、拒绝、完成填写、已回退、打回重填、自动跳过、自动交给管理员
	NodeInstanceID    string             `json:"nodeInstanceId"`
	OperationRecords  []*OperationRecord `gorm:"-" json:"operationRecords"`
	FlowName          string             `gorm:"-" json:"flowName"`
	Reason            string             `gorm:"-" json:"reason"`
}

// InstanceStepRepo interface
type InstanceStepRepo interface {
	Create(db *gorm.DB, model *InstanceStep) error
	Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error
	Delete(db *gorm.DB, ID string) error
	FindByID(db *gorm.DB, ID string) (*InstanceStep, error)
	FindInstanceSteps(db *gorm.DB, condition *InstanceStep) ([]*InstanceStep, error)
	FindInstanceStepsByStatus(db *gorm.DB, processInstanceID string, status []string) ([]*InstanceStep, error)
	GetFlowInstanceStep(db *gorm.DB, processInstanceID string, nodeInstanceID string, status []string) ([]*InstanceStep, error)
	UpdateByNodeInstanceID(db *gorm.DB, nodeInstanceID string, updateMap map[string]interface{}) error
	DeleteByProcessInstanceIDs(db *gorm.DB, processInstanceIDs []string) error
}
