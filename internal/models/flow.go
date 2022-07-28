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

// Flow info
type Flow struct {
	BaseModel

	AppID       string `json:"appId" binding:"required"`
	AppStatus   string `json:"appStatus"`
	SourceID    string `json:"sourceId"` // flow initial model id
	Name        string `json:"name" binding:"required"`
	TriggerMode string `json:"triggerMode" binding:"required"` // FORM_DATA|FORM_TIME
	FormID      string `json:"formId"`
	BpmnText    string `json:"bpmnText"`   // flow model json
	Cron        string `json:"cron"`       // if TriggerMode eq FORM_TIME required
	ProcessKey  string `json:"processKey"` // Process key, used to start the process by the id
	Status      string `json:"status"`
	CanCancel   int8   `json:"canCancel"`
	/**
	1:It can only be cancel when the next event is not processed
	2:Any event can be cancel
	3:Cancel under the specified event
	*/
	CanCancelType    int8   `json:"canCancelType"`
	CanCancelNodes   string `json:"canCancelNodes"` // taskDefKey array
	CanUrge          int8   `json:"canUrge"`
	CanViewStatusMsg int8   `json:"canViewStatusMsg"`
	CanMsg           int8   `json:"canMsg"`
	InstanceName     string `json:"instanceName"` // Instance name template
	KeyFields        string `json:"keyFields"`    // Flow key fields
	ProcessID        string `json:"processID"`    // Process id

	Variables []*Variables `json:"variables" gorm:"-"`
}

const (
	// ENABLE status
	ENABLE = "ENABLE"
	// DISABLE status
	DISABLE = "DISABLE"
	// DELETED status
	DELETED = "DELETED"
)

// FlowRepo interface
type FlowRepo interface {
	Create(db *gorm.DB, model *Flow) error
	Create2(db *gorm.DB, model *Flow) error
	Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error
	Delete(db *gorm.DB, ID string) error
	FindByID(db *gorm.DB, ID string) (*Flow, error)
	FindByIDs(db *gorm.DB, IDs []string) ([]*Flow, error)
	UpdateFlow(db *gorm.DB, model *Flow) error
	FindFlows(db *gorm.DB, condition map[string]interface{}) ([]*Flow, error)
	GetFlows(db *gorm.DB, condition map[string]interface{}) ([]*Flow, error)
	UpdateFlows(db *gorm.DB, condition map[string]interface{}, updateMap map[string]interface{}) error
	FindPageFlows(db *gorm.DB, condition map[string]interface{}, page, limit int) ([]*Flow, int64)
	FindByProcessID(db *gorm.DB, processID string) (*Flow, error)
	FindFlowList(db *gorm.DB, condition map[string]interface{}) ([]*Flow, error)
	UpdateAppStatus(db *gorm.DB, appID string, appStatus string) error
	DeleteByIDs(db *gorm.DB, IDs []string) error

	FindPublishIDs(db *gorm.DB, flowID string) ([]string, error)
}
