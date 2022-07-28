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

// InstanceExecution info
type InstanceExecution struct {
	BaseModel

	ProcessInstanceID string `json:"processInstanceID"`
	ExecutionID       string `json:"executionID"`
	Result            string `json:"result"`
}

// InstanceExecutionRepo interface
type InstanceExecutionRepo interface {
	Create(db *gorm.DB, model *InstanceExecution) error
	Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error

	Delete(db *gorm.DB, ID string) error
	DeleteByProcessInstanceIDs(db *gorm.DB, processInstanceIDs []string) error
	FindByID(db *gorm.DB, ID string) (*InstanceExecution, error)
	FindByIDs(db *gorm.DB, IDs []string) ([]*InstanceExecution, error)
	FindByExecutionIDs(db *gorm.DB, executionIDs []string) ([]*InstanceExecution, error)
}
