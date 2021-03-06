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

package mysql

import (
	"github.com/quanxiang-cloud/flow/internal/models"
	"github.com/quanxiang-cloud/flow/pkg/misc/id2"
	"github.com/quanxiang-cloud/flow/pkg/misc/time2"
	"gorm.io/gorm"
)

type formFieldRepo struct{}

// NewFormFieldRepo new repo
func NewFormFieldRepo() models.FormFieldRepo {
	return &formFieldRepo{}
}

// TableName db table name
func (r *formFieldRepo) TableName() string {
	return "flow_form_field"
}

// Create create model
func (r *formFieldRepo) Create(db *gorm.DB, entity *models.FormField) error {
	entity.ID = id2.GenID()
	entity.CreateTime = time2.Now()
	err := db.Table(r.TableName()).
		Create(entity).
		Error
	return err
}

// Update update model
func (r *formFieldRepo) Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error {
	updateMap["modify_time"] = time2.Now()
	err := db.Table(r.TableName()).Where("id=?", ID).Updates(updateMap).Error

	return err
}

// DeleteByID delete model
func (r *formFieldRepo) Delete(db *gorm.DB, ID string) error {
	entity := &models.FormField{BaseModel: models.BaseModel{
		ID: ID,
	}}
	err := db.Table(r.TableName()).Delete(entity).Error
	return err
}

// DeleteByFlowID delete by flow id
func (r *formFieldRepo) DeleteByFlowID(db *gorm.DB, flowID string) error {
	entity := &models.FormField{}
	err := db.Table(r.TableName()).Where("flow_id = ?", flowID).Delete(entity).Error
	return err
}

func (r *formFieldRepo) FindByFlowID(db *gorm.DB, ID string) ([]*models.FormField, error) {
	entities := make([]*models.FormField, 0)
	db = db.Table(r.TableName()).Where("flow_id=?", ID)
	err := db.Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *formFieldRepo) DeleteByFlowIDs(db *gorm.DB, flowIDs []string) error {
	err := db.Table(r.TableName()).Where("flow_id in (?)", flowIDs).Delete(&models.FormField{}).Error
	return err
}
