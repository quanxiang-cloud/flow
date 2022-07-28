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
	"github.com/quanxiang-cloud/flow/pkg/page"
	"gorm.io/gorm"
)

// AbnormalTask info
type AbnormalTask struct {
	BaseModel

	FlowInstanceID    string `json:"flowInstanceId"`
	ProcessInstanceID string `json:"processInstanceId"`
	TaskID            string `json:"taskId"`
	TaskName          string `json:"taskName"`
	TaskDefKey        string `json:"taskDefKey"`
	Reason            string `json:"reason"`
	Remark            string `json:"remark"`
	Status            int8   `json:"status"` // 0 unhandle，1 handled，2 autoHandled
}

// AbnormalTaskRepo interface
type AbnormalTaskRepo interface {
	Create(db *gorm.DB, model *AbnormalTask) error
	Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error
	UpdateByTaskID(db *gorm.DB, taskID string, updateMap map[string]interface{}) error
	UpdateByProcessInstanceID(db *gorm.DB, processInstanceID string, updateMap map[string]interface{}) error
	Delete(db *gorm.DB, ID string) error
	FindByID(db *gorm.DB, ID string) (*AbnormalTask, error)
	Find(db *gorm.DB, condition map[string]interface{}) ([]*AbnormalTask, error)
	Page(db *gorm.DB, req *AbnormalTaskReq) ([]*AbnormalTaskVo, int64, error)
	DeleteByInstanceIDs(db *gorm.DB, InstanceIDs []string) error
}

// AbnormalTaskReq AbnormalTaskReq
type AbnormalTaskReq struct {
	page.ReqPage

	Status         int8     `json:"status"`
	InstanceStatus string   `json:"instanceStatus"`
	AppID          string   `json:"appId"`
	AdminAppIDs    []string `json:"adminAppIDs"`
	Keyword        string   `json:"keyword"` // 应用名称、流程名称、申请人名称
}

// AbnormalTaskVo vo
type AbnormalTaskVo struct {
	AbnormalTask

	AppName            string `json:"appName"`
	InstanceName       string `json:"instanceName"`
	ApplyUserName      string `json:"applyUserName"`
	InstanceCreateTime string `json:"instanceCreateTime"`
	InstanceStatus     string `json:"instanceStatus"`
}
