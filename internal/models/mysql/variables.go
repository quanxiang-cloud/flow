package mysql

import (
	"github.com/quanxiang-cloud/flow/internal/models"
	"github.com/quanxiang-cloud/flow/pkg/misc/id2"
	"github.com/quanxiang-cloud/flow/pkg/misc/time2"
	"gorm.io/gorm"
)

type variablesRepo struct{}

// NewVariablesRepo new repo
func NewVariablesRepo() models.VariablesRepo {
	return &variablesRepo{}
}

// TableName db table name
func (r *variablesRepo) TableName() string {
	return "flow_variables"
}

// Create create model
func (r *variablesRepo) Create(db *gorm.DB, entity *models.Variables) error {
	entity.ID = id2.GenID()
	entity.CreateTime = time2.Now()
	err := db.Table(r.TableName()).
		Create(entity).
		Error
	return err
}

// Create create model
func (r *variablesRepo) Create2(db *gorm.DB, entity *models.Variables) error {
	entity.CreateTime = time2.Now()
	err := db.Table(r.TableName()).
		Create(entity).
		Error
	return err
}

// Update update model
func (r *variablesRepo) Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error {
	updateMap["modify_time"] = time2.Now()
	err := db.Table(r.TableName()).Where("id=?", ID).Updates(updateMap).Error

	return err
}

// DeleteByID delete model
func (r *variablesRepo) Delete(db *gorm.DB, ID string) error {
	entity := &models.Variables{BaseModel: models.BaseModel{
		ID: ID,
	}}
	err := db.Table(r.TableName()).Delete(entity).Error
	return err
}

// DeleteByFlowID delete by flow id
func (r *variablesRepo) DeleteByFlowID(db *gorm.DB, flowID string) error {
	entity := &models.Variables{}
	err := db.Table(r.TableName()).Where("flow_id = ?", flowID).Delete(entity).Error
	return err
}

// FindByID find model by ID
func (r *variablesRepo) FindByID(db *gorm.DB, ID string) (*models.Variables, error) {
	entity := new(models.Variables)
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

func (r *variablesRepo) FindVariablesByFlowID(db *gorm.DB, ID string) ([]*models.Variables, error) {
	variables := make([]*models.Variables, 0)
	db = db.Table(r.TableName()).Where("flow_id = '0'")
	if ID != "" {
		db = db.Or("flow_id=?", ID)
	}
	err := db.Find(&variables).Error
	if err != nil {
		return nil, err
	}
	return variables, nil
}

func (r *variablesRepo) FindVariables(db *gorm.DB, conditionMap map[string]interface{}) ([]*models.Variables, error) {
	variables := make([]*models.Variables, 0)
	err := db.Table(r.TableName()).Where(conditionMap).Find(&variables).Error
	if err != nil {
		return nil, err
	}
	return variables, nil
}

func (r *variablesRepo) DeleteByFlowIDs(db *gorm.DB, flowIDs []string) error {
	err := db.Table(r.TableName()).Where("flow_id in (?)", flowIDs).Delete(&models.Variables{}).Error
	return err
}
