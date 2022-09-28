package mysql

import (
	"github.com/quanxiang-cloud/flow/internal/models"
	"github.com/quanxiang-cloud/flow/pkg/misc/id2"
	"github.com/quanxiang-cloud/flow/pkg/misc/time2"
	"gorm.io/gorm"
)

type flowProcessRelationRepo struct{}

// NewFlowProcessRelationRepo new repo
func NewFlowProcessRelationRepo() models.FlowProcessRelationRepo {
	return &flowProcessRelationRepo{}
}

// TableName db table name
func (r *flowProcessRelationRepo) TableName() string {
	return "flow_process_relation"
}

// Create create model
func (r *flowProcessRelationRepo) Create(db *gorm.DB, entity *models.FlowProcessRelation) error {
	entity.ID = id2.GenID()
	entity.CreateTime = time2.Now()
	err := db.Table(r.TableName()).
		Create(entity).
		Error
	return err
}

// DeleteByFlowID delete model
func (r *flowProcessRelationRepo) DeleteByFlowID(db *gorm.DB, flowID string) error {
	entity := &models.Flow{}
	err := db.Table(r.TableName()).Delete(entity).Where("flow_id=?", flowID).Error
	return err
}

// FindByProcessID find model by processID
func (r *flowProcessRelationRepo) FindByProcessID(db *gorm.DB, processID string) (*models.FlowProcessRelation, error) {
	entity := new(models.FlowProcessRelation)
	err := db.Table(r.TableName()).
		Where("process_id = ?", processID).
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
