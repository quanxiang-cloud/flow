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

type urgeRepo struct{}

// NewUrgeRepo new repo
func NewUrgeRepo() models.UrgeRepo {
	return &urgeRepo{}
}

// TableName db table name
func (r *urgeRepo) TableName() string {
	return "flow_urge"
}

// Create create model
func (r *urgeRepo) Create(db *gorm.DB, entity *models.Urge) error {
	entity.ID = id2.GenID()
	entity.CreateTime = time2.Now()
	err := db.Table(r.TableName()).
		Create(entity).
		Error
	return err
}

// Update update model
func (r *urgeRepo) Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error {
	updateMap["modify_time"] = time2.Now()
	err := db.Table(r.TableName()).Where("id=?", ID).Updates(updateMap).Error

	return err
}

// DeleteByID delete model
func (r *urgeRepo) Delete(db *gorm.DB, ID string) error {
	entity := &models.Urge{BaseModel: models.BaseModel{
		ID: ID,
	}}
	err := db.Table(r.TableName()).Delete(entity).Error
	return err
}

// FindByID find model by ID
func (r *urgeRepo) FindByID(db *gorm.DB, ID string) (*models.Urge, error) {
	entity := new(models.Urge)
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

func (r *urgeRepo) FindTaskIDs(db *gorm.DB) ([]string, error) {
	var result []string
	err := db.Table(r.TableName()).Select("task_id").Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *urgeRepo) FindByTaskID(db *gorm.DB, taskID string) ([]*models.Urge, error) {
	entity := make([]*models.Urge, 0)
	err := db.Table(r.TableName()).
		Where("task_id = ?", taskID).
		Find(&entity).
		Error
	if err != nil {
		return nil, err
	}

	return entity, nil
}

// GetUrgeNums get urge nums
func (r *urgeRepo) GetUrgeNums(db *gorm.DB, taskIDs []string) (map[string]int64, error) {
	// select fu.task_id,COUNT(fu.task_id) as urgeNum from flow_urge fu where fu.task_id in ('2778a110-c139-47e7-b7d6-e54e75f4a3b4')  group by fu.task_id

	entity := make([]*models.UrgeNumModel, 0)
	err := db.Table(r.TableName()).
		Select("task_id,COUNT(task_id) as urge_num").
		Where("task_id in (?)", taskIDs).
		Group("task_id").
		Find(&entity).
		Error
	if err != nil {
		return nil, err
	}

	urgeNums := make(map[string]int64)
	if len(entity) > 0 {
		for _, value := range entity {
			urgeNums[value.TaskID] = value.UrgeNum
		}
	}

	return urgeNums, nil
}

func (r *urgeRepo) DeleteByProcessInstanceIDs(db *gorm.DB, processInstanceIDs []string) error {
	err := db.Table(r.TableName()).Where("process_instance_id in (?)", processInstanceIDs).Delete(&models.Urge{}).Error
	return err
}
