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

// Comment info
type Comment struct {
	BaseModel

	FlowInstanceID string `json:"flowInstanceId"`
	CommentUserID  string `json:"commentUserId"`
	Content        string `json:"content"`
}

// CommentRepo interface
type CommentRepo interface {
	Create(db *gorm.DB, model *Comment) error
	Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error
	Delete(db *gorm.DB, ID string) error
	FindByID(db *gorm.DB, ID string) (*Comment, error)
	FindComments(db *gorm.DB, condition map[string]interface{}, order string) ([]*Comment, error)
	DeleteByInstanceIDs(db *gorm.DB, InstanceIDs []string) error
}
