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

// Variables info
type Variables struct {
	BaseModel

	FlowID       string `json:"flowId"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Code         string `json:"code"`
	FieldType    string `json:"fieldType"`
	Format       string `json:"format"`
	DefaultValue string `json:"defaultValue"`
	Desc         string `json:"desc"`
}

// VariablesRepo interface
type VariablesRepo interface {
	Create(db *gorm.DB, model *Variables) error
	Create2(db *gorm.DB, entity *Variables) error
	Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error
	// UpdateByID(db *gorm.DB, ID string, update interface{}) error
	Delete(db *gorm.DB, ID string) error
	DeleteByFlowID(db *gorm.DB, flowID string) error
	DeleteByFlowIDs(db *gorm.DB, flowIDs []string) error
	FindByID(db *gorm.DB, ID string) (*Variables, error)
	FindVariablesByFlowID(db *gorm.DB, ID string) ([]*Variables, error)
	FindVariables(db *gorm.DB, conditionMap map[string]interface{}) ([]*Variables, error)
}
