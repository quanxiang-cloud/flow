package mysql

import (
	"github.com/quanxiang-cloud/flow/internal/models"
	"github.com/quanxiang-cloud/flow/pkg/misc/id2"
	"github.com/quanxiang-cloud/flow/pkg/misc/time2"
	"gorm.io/gorm"
)

type abnormalTaskRepo struct{}

// NewAbnormalTaskRepo new repo
func NewAbnormalTaskRepo() models.AbnormalTaskRepo {
	return &abnormalTaskRepo{}
}

// TableName db table name
func (r *abnormalTaskRepo) TableName() string {
	return "flow_abnormal_task"
}

// Create create model
func (r *abnormalTaskRepo) Create(db *gorm.DB, entity *models.AbnormalTask) error {
	entity.ID = id2.GenID()
	entity.CreateTime = time2.Now()
	err := db.Table(r.TableName()).
		Create(entity).
		Error
	return err
}

// Update update model
func (r *abnormalTaskRepo) Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error {
	updateMap["modify_time"] = time2.Now()
	err := db.Table(r.TableName()).Where("id=?", ID).Updates(updateMap).Error

	return err
}

// Update update model
func (r *abnormalTaskRepo) UpdateByTaskID(db *gorm.DB, taskID string, updateMap map[string]interface{}) error {
	updateMap["modify_time"] = time2.Now()
	err := db.Table(r.TableName()).Where("task_id=?", taskID).Updates(updateMap).Error

	return err
}

func (r *abnormalTaskRepo) UpdateByProcessInstanceID(db *gorm.DB, processInstanceID string, updateMap map[string]interface{}) error {
	updateMap["modify_time"] = time2.Now()
	err := db.Table(r.TableName()).Where("process_instance_id=?", processInstanceID).Updates(updateMap).Error

	return err
}

// DeleteByID delete model
func (r *abnormalTaskRepo) Delete(db *gorm.DB, ID string) error {
	entity := &models.AbnormalTask{BaseModel: models.BaseModel{
		ID: ID,
	}}
	err := db.Table(r.TableName()).Delete(entity).Error
	return err
}

// FindByID find model by ID
func (r *abnormalTaskRepo) FindByID(db *gorm.DB, ID string) (*models.AbnormalTask, error) {
	entity := new(models.AbnormalTask)
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

func (r *abnormalTaskRepo) Find(db *gorm.DB, condition map[string]interface{}) ([]*models.AbnormalTask, error) {
	rs := make([]*models.AbnormalTask, 0)
	err := db.Table(r.TableName()).Where(condition).Find(&rs).Error
	if err != nil {
		return nil, err
	}
	return rs, nil
}

func (r *abnormalTaskRepo) Page(db *gorm.DB, req *models.AbnormalTaskReq) ([]*models.AbnormalTaskVo, int64, error) {
	rs := make([]*models.AbnormalTaskVo, 0)
	var count int64
	tx := db.Table(r.TableName()).Select("flow_abnormal_task.*, flow_instance.name as 'instance_name', flow_instance.app_name, flow_instance.apply_user_name, flow_instance.create_time as 'instance_create_time',flow_instance.status as 'instance_status'")
	tx = tx.Joins("left join flow_instance on flow_abnormal_task.flow_instance_id = flow_instance.id")

	if len(req.Keyword) > 0 {
		tx = tx.Where("flow_instance.app_name LIKE ? or flow_instance.name LIKE ? or flow_instance.apply_user_name LIKE ? ", "%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.Status > -1 {
		tx = tx.Where("flow_abnormal_task.status = ?", req.Status)
	}
	if len(req.InstanceStatus) > 0 {
		tx = tx.Where("flow_instance.status = ?", req.InstanceStatus)
	}
	if len(req.AppID) > 0 {
		tx = tx.Where("flow_instance.app_id = ?", req.AppID)
	}
	if len(req.AdminAppIDs) > 0 {
		tx = tx.Where("flow_instance.app_id in (?)", req.AdminAppIDs)
	}

	if len(req.Orders) > 0 {
		for _, value := range req.Orders {
			tx = tx.Order(value.Column + " " + value.Direction)
		}
	}
	if req.Page > 0 && req.Size > 0 {
		tx = tx.Limit(req.Size).Offset((req.Page - 1) * req.Size)
	}

	err := tx.Find(&rs).Offset(-1).Limit(-1).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	return rs, count, nil
}

func (r *abnormalTaskRepo) DeleteByInstanceIDs(db *gorm.DB, InstanceIDs []string) error {
	err := db.Table(r.TableName()).Where("flow_instance_id in (?)", InstanceIDs).Delete(&models.AbnormalTask{}).Error
	return err
}
