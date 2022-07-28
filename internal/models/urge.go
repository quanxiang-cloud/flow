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

// Urge info
type Urge struct {
	BaseModel

	TaskID            string `json:"taskId"`
	ProcessInstanceID string `json:"processInstanceId"`
}

// UrgeRepo interface
type UrgeRepo interface {
	Create(db *gorm.DB, model *Urge) error
	Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error
	Delete(db *gorm.DB, ID string) error
	FindByID(db *gorm.DB, ID string) (*Urge, error)
	FindTaskIDs(db *gorm.DB) ([]string, error)
	DeleteByProcessInstanceIDs(db *gorm.DB, processInstanceIDs []string) error

	FindByTaskID(db *gorm.DB, taskID string) ([]*Urge, error)
	GetUrgeNums(db *gorm.DB, taskIDs []string) (map[string]int64, error)
}

// UrgeNumModel urge num model
type UrgeNumModel struct {
	TaskID  string
	UrgeNum int64
}
