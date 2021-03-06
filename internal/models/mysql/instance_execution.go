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

type instanceExecutionRepo struct{}

// NewInstanceExecutionRepo new repo
func NewInstanceExecutionRepo() models.InstanceExecutionRepo {
	return &instanceExecutionRepo{}
}

// TableName db table name
func (r *instanceExecutionRepo) TableName() string {
	return "instance_execution"
}

// Create create model
func (r *instanceExecutionRepo) Create(db *gorm.DB, entity *models.InstanceExecution) error {
	entity.ID = id2.GenID()
	entity.CreateTime = time2.Now()
	err := db.Table(r.TableName()).
		Create(entity).
		Error
	return err
}

// Update update model
func (r *instanceExecutionRepo) Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error {
	updateMap["modify_time"] = time2.Now()
	err := db.Table(r.TableName()).Where("id=?", ID).Updates(updateMap).Error

	return err
}

// DeleteByID delete model
func (r *instanceExecutionRepo) Delete(db *gorm.DB, ID string) error {
	entity := &models.Instance{BaseModel: models.BaseModel{
		ID: ID,
	}}
	err := db.Table(r.TableName()).Delete(entity).Error
	return err
}

// FindByID find model by ID
func (r *instanceExecutionRepo) FindByID(db *gorm.DB, ID string) (*models.InstanceExecution, error) {
	entity := new(models.InstanceExecution)
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

func (r *instanceExecutionRepo) FindByIDs(db *gorm.DB, IDs []string) ([]*models.InstanceExecution, error) {
	entity := make([]*models.InstanceExecution, 0)
	err := db.Table(r.TableName()).
		Where(IDs).
		Find(&entity).
		Error
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *instanceExecutionRepo) FindByExecutionIDs(db *gorm.DB, executionIDs []string) ([]*models.InstanceExecution, error) {
	entity := make([]*models.InstanceExecution, 0)
	err := db.Table(r.TableName()).
		Where("execution_id in (?)", executionIDs).
		Find(&entity).
		Error
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *instanceExecutionRepo) DeleteByProcessInstanceIDs(db *gorm.DB, processInstanceIDs []string) error {
	err := db.Table(r.TableName()).Where("process_instance_id in (?)", processInstanceIDs).Delete(&models.InstanceExecution{}).Error
	return err
}
