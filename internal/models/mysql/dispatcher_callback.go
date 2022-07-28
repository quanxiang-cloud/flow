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
	"github.com/quanxiang-cloud/flow/pkg/misc/time2"
	"gorm.io/gorm"
)

type dispatcherCallbackRepo struct{}

// NewDispatcherCallbackRepo new repo
func NewDispatcherCallbackRepo() models.DispatcherCallbackRepo {
	return &dispatcherCallbackRepo{}
}

// TableName db table name
func (r *dispatcherCallbackRepo) TableName() string {
	return "dispatcher_callback"
}

// Create create model
func (r *dispatcherCallbackRepo) Create(db *gorm.DB, entity *models.DispatcherCallback) error {
	entity.CreateTime = time2.Now()
	err := db.Table(r.TableName()).
		Create(entity).
		Error
	return err
}

// Update update model
func (r *dispatcherCallbackRepo) Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error {
	updateMap["modify_time"] = time2.Now()
	err := db.Table(r.TableName()).Where("id=?", ID).Updates(updateMap).Error

	return err
}

// Delete delete model
func (r *dispatcherCallbackRepo) Delete(db *gorm.DB, ID string) error {
	entity := &models.DispatcherCallback{BaseModel: models.BaseModel{
		ID: ID,
	}}
	err := db.Table(r.TableName()).Delete(entity).Error
	return err
}

// FindByID find model by ID
func (r *dispatcherCallbackRepo) FindByID(db *gorm.DB, ID string) (*models.DispatcherCallback, error) {
	entity := new(models.DispatcherCallback)
	err := db.Table(r.TableName()).
		Where("id = ?", ID).
		Find(entity).
		Error
	if err != nil {
		return nil, err
	}
	if entity.ID == "" {
		return nil, nil
	}
	return entity, nil
}
