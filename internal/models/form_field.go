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

// FormField info
type FormField struct {
	BaseModel

	FlowID         string `json:"flowId"`
	FormID         string `json:"formId"`
	FieldName      string `json:"fieldName"`
	FieldValuePath string `json:"fieldValuePath"`
}

// FormFieldRepo interface
type FormFieldRepo interface {
	Create(db *gorm.DB, model *FormField) error
	Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error
	Delete(db *gorm.DB, ID string) error
	DeleteByFlowID(db *gorm.DB, flowID string) error
	FindByFlowID(db *gorm.DB, ID string) ([]*FormField, error)
	DeleteByFlowIDs(db *gorm.DB, flowIDs []string) error
}
