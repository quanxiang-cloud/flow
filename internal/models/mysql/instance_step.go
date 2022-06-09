package mysql

import (
	"github.com/quanxiang-cloud/flow/internal/models"
	"github.com/quanxiang-cloud/flow/pkg/misc/id2"
	"github.com/quanxiang-cloud/flow/pkg/misc/time2"
	"gorm.io/gorm"
)

type instanceStepRepo struct{}

// NewInstanceStepRepo new repo
func NewInstanceStepRepo() models.InstanceStepRepo {
	return &instanceStepRepo{}
}

// TableName db table name
func (r *instanceStepRepo) TableName() string {
	return "flow_instance_step"
}

// Create create model
func (r *instanceStepRepo) Create(db *gorm.DB, entity *models.InstanceStep) error {
	entity.ID = id2.GenID()
	entity.CreateTime = time2.Now()
	err := db.Table(r.TableName()).
		Create(entity).
		Error
	return err
}

// Update update model
func (r *instanceStepRepo) Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error {
	updateMap["modify_time"] = time2.Now()
	err := db.Table(r.TableName()).Where("id=?", ID).Updates(updateMap).Error

	return err
}

// DeleteByID delete model
func (r *instanceStepRepo) Delete(db *gorm.DB, ID string) error {
	entity := &models.InstanceStep{BaseModel: models.BaseModel{
		ID: ID,
	}}
	err := db.Table(r.TableName()).Delete(entity).Error
	return err
}

// FindByID find model by ID
func (r *instanceStepRepo) FindByID(db *gorm.DB, ID string) (*models.InstanceStep, error) {
	entity := new(models.InstanceStep)
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

// FindInstanceSteps find instance step list
func (r *instanceStepRepo) FindInstanceSteps(db *gorm.DB, condition *models.InstanceStep) ([]*models.InstanceStep, error) {
	steps := make([]*models.InstanceStep, 0)
	err := db.Table(r.TableName()).
		Where(condition).
		Order("create_time" + " " + "desc").
		Find(&steps).
		Error
	if err != nil {
		return nil, err
	}
	return steps, nil
}

func (r *instanceStepRepo) FindInstanceStepsByStatus(db *gorm.DB, processInstanceID string, status []string) ([]*models.InstanceStep, error) {
	steps := make([]*models.InstanceStep, 0)
	err := db.Table(r.TableName()).
		Where("process_instance_id=? and status  IN (?)", processInstanceID, status).
		Find(&steps).
		Error
	if err != nil {
		return nil, err
	}
	return steps, nil
}

func (r *instanceStepRepo) GetFlowInstanceStep(db *gorm.DB, processInstanceID string, nodeInstanceID string, status []string) ([]*models.InstanceStep, error) {
	steps := make([]*models.InstanceStep, 0)
	err := db.Table(r.TableName()).
		Where("process_instance_id=? and node_instance_id=? and status  IN (?)", processInstanceID, nodeInstanceID, status).
		Find(&steps).
		Error
	if err != nil {
		return nil, err
	}
	return steps, nil
}

// UpdateByNodeInstanceID update model
func (r *instanceStepRepo) UpdateByNodeInstanceID(db *gorm.DB, nodeInstanceID string, updateMap map[string]interface{}) error {
	err := db.Table(r.TableName()).Where("node_instance_id=?", nodeInstanceID).Updates(updateMap).Error
	return err
}

func (r *instanceStepRepo) DeleteByProcessInstanceIDs(db *gorm.DB, processInstanceIDs []string) error {
	err := db.Table(r.TableName()).Where("process_instance_id in (?)", processInstanceIDs).Delete(&models.InstanceStep{}).Error
	return err
}
