package mysql

import (
	"git.internal.yunify.com/qxp/misc/id2"
	"git.internal.yunify.com/qxp/misc/time2"
	"github.com/quanxiang-cloud/flow/internal/models"
	"github.com/quanxiang-cloud/flow/pkg/utils"
	"gorm.io/gorm"
)

type instanceRepo struct{}

// NewInstanceRepo new repo
func NewInstanceRepo() models.InstanceRepo {
	return &instanceRepo{}
}

// TableName db table name
func (r *instanceRepo) TableName() string {
	return "flow_instance"
}

// Create create model
func (r *instanceRepo) Create(db *gorm.DB, entity *models.Instance) error {
	entity.ID = id2.GenID()
	if len(entity.CreateTime) == 0 {
		entity.CreateTime = time2.Now()
	}

	err := db.Table(r.TableName()).
		Create(entity).
		Error
	return err
}

// UpdateAppStatus update app status
func (r *instanceRepo) UpdateAppStatus(db *gorm.DB, appID string, appStatus string) error {
	updateMap := map[string]interface{}{
		"app_status": appStatus,
	}
	err := db.Table(r.TableName()).Where("app_id=?", appID).Updates(updateMap).Error
	return err
}

// Update update model
func (r *instanceRepo) Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error {
	updateMap["modify_time"] = time2.Now()
	err := db.Table(r.TableName()).Where("id=?", ID).Updates(updateMap).Error

	return err
}

// DeleteByID delete model
func (r *instanceRepo) Delete(db *gorm.DB, ID string) error {
	entity := &models.Instance{BaseModel: models.BaseModel{
		ID: ID,
	}}
	err := db.Table(r.TableName()).Delete(entity).Error
	return err
}

// FindByID find model by ID
func (r *instanceRepo) FindByID(db *gorm.DB, ID string) (*models.Instance, error) {
	entity := new(models.Instance)
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

func (r *instanceRepo) FindByIDs(db *gorm.DB, IDs []string) ([]*models.Instance, error) {
	entity := make([]*models.Instance, 0)
	err := db.Table(r.TableName()).
		Where(IDs).
		Find(&entity).
		Error
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *instanceRepo) FindByProcessInstanceIDs(db *gorm.DB, processInstanceIDs []string) ([]*models.Instance, error) {
	entity := make([]*models.Instance, 0)
	err := db.Table(r.TableName()).
		Where("process_instance_id in (?)", processInstanceIDs).
		Find(&entity).
		Error
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *instanceRepo) GetEntityByProcessInstanceID(db *gorm.DB, processInstanceID string) (*models.Instance, error) {
	entity := new(models.Instance)
	err := db.Table(r.TableName()).
		Where("process_instance_id = ?", processInstanceID).
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

func (r *instanceRepo) FindInstances(db *gorm.DB, formInstanceID string, status []string) ([]*models.Instance, error) {
	instances := make([]*models.Instance, 0)
	err := db.Table(r.TableName()).
		Where("form_instance_id=? and status in (?)", formInstanceID, status).
		Find(&instances).
		Error
	if err != nil {
		return nil, err
	}
	return instances, nil
}

func (r *instanceRepo) PageInstances(db *gorm.DB, req *models.PageInstancesReq) ([]*models.Instance, int64, error) {
	instances := make([]*models.Instance, 0)
	var count int64

	tx := db.Table(r.TableName())
	tx = tx.Where("app_status=?", AppActiveStatus)
	if len(req.ApplyUserID) > 0 {
		tx = tx.Where("apply_user_id=?", req.ApplyUserID)
	}
	if len(req.AppID) > 0 {
		tx = tx.Where("app_id=?", req.AppID)
	}
	if len(req.Status) > 0 {
		tx = tx.Where("status=?", req.Status)
	}
	if len(req.CreateTimeBegin) > 0 {
		tx = tx.Where("create_time >= ?", utils.ChangeBjTimeToISO8601(req.CreateTimeBegin+" 00:00:00"))
	}
	if len(req.CreateTimeEnd) > 0 {
		tx = tx.Where("create_time <= ?", utils.ChangeBjTimeToISO8601(req.CreateTimeEnd+" 23:59:59"))
	}
	if len(req.Keyword) > 0 {
		tx = tx.Where("name LIKE ? or app_name LIKE ? or apply_user_name LIKE ? ", "%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	if len(req.Orders) > 0 {
		for _, value := range req.Orders {
			tx = tx.Order(value.Column + " " + value.Direction)
		}
	}
	if req.Page > 0 && req.Size > 0 {
		tx = tx.Limit(req.Size).Offset((req.Page - 1) * req.Size)
	}

	err := tx.Find(&instances).Offset(-1).Limit(-1).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	return instances, count, nil
}

func (r *instanceRepo) GetInstances(db *gorm.DB, condition map[string]interface{}) ([]*models.Instance, error) {
	instances := make([]*models.Instance, 0)
	err := db.Table(r.TableName()).
		Where(condition).
		Find(&instances).
		Error
	if err != nil {
		return nil, err
	}
	return instances, nil
}

func (r *instanceRepo) DeleteByFlowIDs(db *gorm.DB, flowIDs []string) error {
	err := db.Table(r.TableName()).Where("flow_id in (?)", flowIDs).Delete(&models.Instance{}).Error
	return err
}
