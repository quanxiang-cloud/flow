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

// ActTaskEntity task info
type ActTaskEntity struct {
	ID                 string      `json:"id"`         // task id
	TaskDefKey         string      `json:"taskDefKey"` // Task definition id
	ProcInstID         string      `json:"procInstId"` // Process instance id
	ActInstID          string      `json:"actInstId"`
	Name               string      `json:"name"`        // Task name
	Description        string      `json:"description"` // Task description
	Owner              string      `json:"owner"`
	Assignee           string      `json:"assignee"`
	StartTime          string      `json:"startTime"`
	EndTime            string      `json:"endTime"`
	Duration           int64       `json:"duration"`
	DueDate            string      `json:"dueDate"`
	FlowInstanceEntity interface{} `json:"flowInstanceEntity"` // Instance
	UrgeNum            int64       `json:"urgeNum"`
	Handled            string      `json:"handled"`
}

// NodeModel event model
type NodeModel struct {
	TaskDefKey string `json:"taskDefKey"` // Task definition id
	TaskName   string `json:"taskName"`
}

// HandleTaskModel  task handle model
type HandleTaskModel struct {
	HandleType       string                 `json:"handleType"`
	HandleDesc       string                 `json:"handleDesc"`
	Remark           string                 `json:"remark"`
	TaskDefKey       string                 `json:"taskDefKey"`
	AttachFiles      []AttachFileModel      `json:"attachFiles"`
	HandleUserIDs    []string               `json:"handleUserIds"`
	CorrelationIDs   []string               `json:"correlationIds"`
	FormData         map[string]interface{} `json:"formData"`
	AutoReviewUserID string                 `json:"autoReviewUserId"`
	RelNodeDefKey    string                 `json:"RelNodeDefKey"`
}

// AttachFileModel file model
type AttachFileModel struct {
	FileName string `json:"fileName"`
	FileURL  string `json:"fileUrl"`
}

// AddSignModel struct
type AddSignModel struct {
	Type              string        `json:"type"` // 加签方式：BEFORE前加签，AFTER后加签
	Assignee          []interface{} `json:"assignee"`
	MultiplePersonWay string        `json:"multiplePersonWay"` // 多人处理方式：and会签，or或签
}
