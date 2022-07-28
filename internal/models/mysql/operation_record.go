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
	"strings"
)

type operationRecordRepo struct{}

// NewOperationRecordRepo new repo
func NewOperationRecordRepo() models.OperationRecordRepo {
	return &operationRecordRepo{}
}

// TableName db table name
func (r *operationRecordRepo) TableName() string {
	return "flow_operation_record"
}

// Create create model
func (r *operationRecordRepo) Create(db *gorm.DB, entity *models.OperationRecord) error {
	entity.ID = id2.GenID()
	entity.CreateTime = time2.Now()
	err := db.Table(r.TableName()).
		Create(entity).
		Error
	return err
}

// Update update model
func (r *operationRecordRepo) Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error {
	updateMap["modify_time"] = time2.Now()
	err := db.Table(r.TableName()).Where("id=?", ID).Updates(updateMap).Error

	return err
}

// DeleteByID delete model
func (r *operationRecordRepo) Delete(db *gorm.DB, ID string) error {
	entity := &models.OperationRecord{BaseModel: models.BaseModel{
		ID: ID,
	}}
	err := db.Table(r.TableName()).Delete(entity).Error
	return err
}

// FindByID find model by ID
func (r *operationRecordRepo) FindByID(db *gorm.DB, ID string) (*models.OperationRecord, error) {
	entity := new(models.OperationRecord)
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

func (r *operationRecordRepo) GetAgreeUserIds(db *gorm.DB, processInstanceID string) ([]string, error) {
	operationRecords := make([]*models.OperationRecord, 0)
	err := db.Table(r.TableName()).
		Where("process_instance_id=? and handle_type='AGREE' and creator_id=''", processInstanceID).
		Find(&operationRecords).
		Error
	if err != nil {
		return nil, err
	}

	userIDs := make([]string, 0)
	for _, value := range operationRecords {
		userIDs = append(userIDs, value.CreatorID)
	}
	return userIDs, nil
}

func (r *operationRecordRepo) FindRecords(db *gorm.DB, processInstanceID, instanceStepID string, status []string, includeStatus bool) ([]*models.OperationRecord, error) {
	records := make([]*models.OperationRecord, 0)
	db = db.Table(r.TableName()).
		Where("process_instance_id=? and instance_step_id=? ", processInstanceID, instanceStepID)
	if includeStatus {
		db = db.Where(" handle_type in (?)", status)
	} else {
		db = db.Where(" handle_type not in (?)", status)
	}
	err := db.Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r *operationRecordRepo) GetHandleUserIDs(db *gorm.DB, processInstanceID string, taskID string, handleType string) ([]string, error) {
	records := make([]*models.OperationRecord, 0)
	err := db.Table(r.TableName()).
		Where("process_instance_id=? and task_id=? and handle_type=?", processInstanceID, taskID, handleType).
		Find(&records).
		Error
	if err != nil {
		return nil, err
	}

	userIDs := make([]string, 0)
	if len(records) > 0 {
		for _, value := range records {
			userIDs = append(userIDs, strings.Split(value.CorrelationData, ",")...)
		}
	}
	return userIDs, nil
}

func (r *operationRecordRepo) UpdateStatus(db *gorm.DB, IDs []string, status string, userID string) error {
	updateMap := map[string]interface{}{
		"modify_time": time2.Now(),
		"status":      status,
		"modifier_id": userID,
	}
	err := db.Table(r.TableName()).Where("id in (?)", IDs).Updates(updateMap).Error

	return err
}

func (r *operationRecordRepo) UpdateStatus2(db *gorm.DB, IDs []string, status string, userID string, remark string) error {
	updateMap := map[string]interface{}{
		"modify_time": time2.Now(),
		"status":      status,
		"remark":      remark,
		"modifier_id": userID,
	}
	err := db.Table(r.TableName()).Where("id in (?)", IDs).Updates(updateMap).Error

	return err
}

func (r *operationRecordRepo) FindRecordsByTaskIDs(db *gorm.DB, handleType string, taskIDs []string) ([]*models.OperationRecord, error) {
	records := make([]*models.OperationRecord, 0)
	err := db.Table(r.TableName()).
		Where("handle_type = ? and task_id in (?)", handleType, taskIDs).
		Find(&records).
		Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r *operationRecordRepo) FindRecordByRelDefKey(db *gorm.DB, processInstanceID string, relDefKey string) (*models.OperationRecord, error) {
	entity := new(models.OperationRecord)
	err := db.Table(r.TableName()).
		Where("process_instance_id = ? and rel_node_def_key = ?", processInstanceID, relDefKey).
		Find(&entity).
		Error
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *operationRecordRepo) DeleteByProcessInstanceIDs(db *gorm.DB, processInstanceIDs []string) error {
	err := db.Table(r.TableName()).Where("process_instance_id in (?)", processInstanceIDs).Delete(&models.OperationRecord{}).Error
	return err
}
