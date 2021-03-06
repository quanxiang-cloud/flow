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

type triggerRuleRepo struct{}

// NewTriggerRuleRepo new repo
func NewTriggerRuleRepo() models.TriggerRuleRepo {
	return &triggerRuleRepo{}
}

// TableName db table name
func (r *triggerRuleRepo) TableName() string {
	return "flow_trigger_rule"
}

// Create create model
func (r *triggerRuleRepo) Create(db *gorm.DB, entity *models.TriggerRule) error {
	entity.ID = id2.GenID()
	entity.CreateTime = time2.Now()
	err := db.Table(r.TableName()).
		Create(entity).
		Error
	return err
}

// Update update model
func (r *triggerRuleRepo) Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error {
	updateMap["modify_time"] = time2.Now()
	err := db.Table(r.TableName()).Where("id=?", ID).Updates(updateMap).Error

	return err
}

// DeleteByID delete model
func (r *triggerRuleRepo) Delete(db *gorm.DB, ID string) error {
	entity := &models.TriggerRule{BaseModel: models.BaseModel{
		ID: ID,
	}}
	err := db.Table(r.TableName()).Delete(entity).Error
	return err
}

// FindByID find model by ID
func (r *triggerRuleRepo) FindByID(db *gorm.DB, ID string) (*models.TriggerRule, error) {
	entity := new(models.TriggerRule)
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

// FindByFormIDAndDFlowID find model by ID
func (r *triggerRuleRepo) FindByFormIDAndDFlowID(db *gorm.DB, formID, flowID string) (*models.TriggerRule, error) {
	entity := new(models.TriggerRule)
	err := db.Table(r.TableName()).
		Where("form_id=? and flow_id=?", formID, flowID).
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

// Deletes delete by condition
func (r *triggerRuleRepo) DeleteByFlowIDS(db *gorm.DB, flowIDS []string) error {
	err := db.Table(r.TableName()).Where("flow_id in (?)", flowIDS).Delete(&models.TriggerRule{}).Error
	return err
}

// FindTriggerRules find trigger rule list
func (r *triggerRuleRepo) FindTriggerRules(db *gorm.DB, condition map[string]interface{}) ([]*models.TriggerRule, error) {
	rules := make([]*models.TriggerRule, 0)
	err := db.Table(r.TableName()).
		Where(condition).
		Find(&rules).
		Error
	if err != nil {
		return nil, err
	}
	return rules, nil
}

// FindTriggerRulesByFormID find tigger rules by formID
func (r *triggerRuleRepo) FindTriggerRulesByFormID(db *gorm.DB, formID string) ([]*models.TriggerRule, error) {
	rules := make([]*models.TriggerRule, 0)
	err := db.Table(r.TableName()).
		Where("form_id = ?", formID).
		Find(&rules).
		Error
	if err != nil {
		return nil, err
	}
	return rules, nil
}
