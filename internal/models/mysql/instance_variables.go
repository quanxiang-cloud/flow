package mysql

import (
	"git.internal.yunify.com/qxp/misc/id2"
	"git.internal.yunify.com/qxp/misc/time2"
	"github.com/quanxiang-cloud/flow/internal/models"
	"gorm.io/gorm"
)

type instanceVariablesRepo struct{}

// NewInstanceVariablesRepo new repo
func NewInstanceVariablesRepo() models.InstanceVariablesRepo {
	return &instanceVariablesRepo{}
}

// TableName db table name
func (r *instanceVariablesRepo) TableName() string {
	return "flow_instance_variables"
}

// Create create model
func (r *instanceVariablesRepo) Create(db *gorm.DB, entity *models.InstanceVariables) error {
	entity.ID = id2.GenID()
	entity.CreateTime = time2.Now()
	err := db.Table(r.TableName()).
		Create(entity).
		Error
	return err
}

// Update update model
func (r *instanceVariablesRepo) Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error {
	updateMap["modify_time"] = time2.Now()
	err := db.Table(r.TableName()).Where("id=?", ID).Updates(updateMap).Error

	return err
}

func (r *instanceVariablesRepo) UpdateByID(db *gorm.DB, ID string, update interface{}) error {
	err := db.Table(r.TableName()).Where("id=?", ID).Updates(update).Error
	return err
}

// DeleteByID delete model
func (r *instanceVariablesRepo) Delete(db *gorm.DB, ID string) error {
	entity := &models.InstanceVariables{BaseModel: models.BaseModel{
		ID: ID,
	}}
	err := db.Table(r.TableName()).Delete(entity).Error
	return err
}

// FindByID find model by ID
func (r *instanceVariablesRepo) FindByID(db *gorm.DB, ID string) (*models.InstanceVariables, error) {
	entity := new(models.InstanceVariables)
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

func (r *instanceVariablesRepo) BatchCreate(db *gorm.DB, model []*models.InstanceVariables) error {
	if len(model) > 0 {
		for _, m := range model {
			err := db.Table(r.TableName()).
				Create(m).
				Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *instanceVariablesRepo) FindVariablesByProcessInstanceID(db *gorm.DB, processInstanceID string) ([]*models.InstanceVariables, error) {
	variables := make([]*models.InstanceVariables, 0)
	db = db.Table(r.TableName()).Where("process_instance_id = ?", processInstanceID)
	err := db.Find(&variables).Error
	if err != nil {
		return nil, err
	}
	return variables, nil
}

func (r *instanceVariablesRepo) UpdateVariable(db *gorm.DB, processInstanceID string, code string, value string) error {
	updateMap := map[string]interface{}{
		"value":       value,
		"modify_time": time2.Now(),
	}
	err := db.Table(r.TableName()).Where("process_instance_id=? and code=?", processInstanceID, code).Updates(updateMap).Error

	return err
}

func (r *instanceVariablesRepo) UpdateTypeAndValue(db *gorm.DB, processInstanceID string, code string, fieldType string, value string) error {
	updateMap := map[string]interface{}{
		"field_type":  fieldType,
		"value":       value,
		"modify_time": time2.Now(),
	}
	err := db.Table(r.TableName()).Where("process_instance_id=? and code=?", processInstanceID, code).Updates(updateMap).Error

	return err
}

func (r *instanceVariablesRepo) DeleteByProcessInstanceIDs(db *gorm.DB, processInstanceIDs []string) error {
	err := db.Table(r.TableName()).Where("process_instance_id in (?)", processInstanceIDs).Delete(&models.InstanceVariables{}).Error
	return err
}

func (r *instanceVariablesRepo) FindVariablesByCode(db *gorm.DB, processInstanceID string, code string) (*models.InstanceVariables, error) {
	variable := new(models.InstanceVariables)
	err := db.Table(r.TableName()).
		Where("process_instance_id=? and code=?", processInstanceID, code).
		Find(variable).
		Error
	if err != nil {
		return nil, err
	}
	return variable, nil
}
