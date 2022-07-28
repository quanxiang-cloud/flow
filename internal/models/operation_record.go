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

import (
	"gorm.io/gorm"
)

// OperationRecord info
type OperationRecord struct {
	BaseModel

	ProcessInstanceID string `json:"processInstanceID"`
	InstanceStepID    string `json:"instanceStepID"` // step id
	HandleType        string `json:"handleType"`
	HandleUserID      string `json:"handleUserID"`
	HandleDesc        string `json:"handleDesc"`
	Remark            string `json:"remark"`
	Status            string `json:"status"` // COMPLETED,ACTIVE

	TaskID          string          `json:"taskId"`
	TaskName        string          `json:"taskName"`
	TaskDefKey      string          `json:"taskDefKey"`
	CorrelationData string          `json:"correlationData"`
	HandleTaskModel HandleTaskModel `gorm:"-" json:"handleTaskModel"`
	CurrentNodeType string          `gorm:"-" json:"currentNodeType"`
	RelNodeDefKey   string          `json:"RelNodeDefKey"`
}

// OperationRecordRepo interface
type OperationRecordRepo interface {
	Create(db *gorm.DB, model *OperationRecord) error
	Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error
	Delete(db *gorm.DB, ID string) error
	FindByID(db *gorm.DB, ID string) (*OperationRecord, error)
	GetAgreeUserIds(db *gorm.DB, processInstanceID string) ([]string, error)
	FindRecords(db *gorm.DB, processInstanceID, instanceStepID string, status []string, includeStatus bool) ([]*OperationRecord, error)
	GetHandleUserIDs(db *gorm.DB, processInstanceID string, taskID string, handleType string) ([]string, error)
	UpdateStatus(db *gorm.DB, IDs []string, status string, userID string) error
	UpdateStatus2(db *gorm.DB, IDs []string, status string, userID string, remark string) error
	FindRecordsByTaskIDs(db *gorm.DB, handleType string, taskIDs []string) ([]*OperationRecord, error)
	FindRecordByRelDefKey(db *gorm.DB, processInstanceID string, relDefKey string) (*OperationRecord, error)
	DeleteByProcessInstanceIDs(db *gorm.DB, processInstanceIDs []string) error
}
