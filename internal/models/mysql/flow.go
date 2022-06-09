package mysql

import (
	"git.internal.yunify.com/qxp/misc/id2"
	"git.internal.yunify.com/qxp/misc/time2"
	"github.com/quanxiang-cloud/flow/internal/models"
	"gorm.io/gorm"
)

// app status
const (
	// 预删除
	AppPreDelete = "preDelete"
	// 删除
	AppDelete = "delete"
	// 恢复
	AppRecovery = "recovery"

	AppActiveStatus  = "ACTIVE"
	AppSuspendStatus = "SUSPEND"
	AppDeleteStatus  = "DELETE"
)

type flowRepo struct{}

// NewFlowRepo new repo
func NewFlowRepo() models.FlowRepo {
	return &flowRepo{}
}

// TableName db table name
func (r *flowRepo) TableName() string {
	return "flow"
}

// Create create model
func (r *flowRepo) Create(db *gorm.DB, entity *models.Flow) error {
	entity.ID = id2.GenID()
	entity.CreateTime = time2.Now()
	err := db.Table(r.TableName()).
		Create(entity).
		Error
	return err
}

func (r *flowRepo) Create2(db *gorm.DB, entity *models.Flow) error {
	entity.CreateTime = time2.Now()
	err := db.Table(r.TableName()).
		Create(entity).
		Error
	return err
}

// UpdateAppStatus update app status
func (r *flowRepo) UpdateAppStatus(db *gorm.DB, appID string, appStatus string) error {
	updateMap := map[string]interface{}{
		"app_status": appStatus,
	}
	err := db.Table(r.TableName()).Where("app_id=?", appID).Updates(updateMap).Error
	return err
}

// Update update model
func (r *flowRepo) Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error {
	updateMap["modify_time"] = time2.Now()
	err := db.Table(r.TableName()).Where("id=?", ID).Updates(updateMap).Error

	return err
}

// UpdateFlow update flow
func (r *flowRepo) UpdateFlow(db *gorm.DB, model *models.Flow) error {
	m := make(map[string]interface{})
	m["app_id"] = model.AppID
	m["source_id"] = model.SourceID
	m["name"] = model.Name
	m["trigger_mode"] = model.TriggerMode
	m["form_id"] = model.FormID
	m["bpmn_text"] = model.BpmnText
	m["process_key"] = model.ProcessKey
	m["status"] = model.Status
	m["can_cancel"] = model.CanCancel
	m["can_urge"] = model.CanUrge
	m["can_view_status_msg"] = model.CanViewStatusMsg
	m["can_msg"] = model.CanMsg
	m["can_cancel_type"] = model.CanCancelType
	m["can_cancel_nodes"] = model.CanCancelNodes
	m["instance_name"] = model.InstanceName
	m["key_fields"] = model.KeyFields
	m["modifier_id"] = model.ModifierID
	m["modify_time"] = time2.Now()
	err := db.Table(r.TableName()).Where("id=?", model.ID).Updates(m).Error

	return err
}

// DeleteByID delete model
func (r *flowRepo) Delete(db *gorm.DB, ID string) error {
	entity := &models.Flow{
		BaseModel: models.BaseModel{
			ID: ID,
		},
	}
	err := db.Table(r.TableName()).Delete(entity).Error
	return err
}

// FindByID find model by ID
func (r *flowRepo) FindByID(db *gorm.DB, ID string) (*models.Flow, error) {
	entity := new(models.Flow)
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

func (r *flowRepo) FindByIDs(db *gorm.DB, IDs []string) ([]*models.Flow, error) {
	entity := make([]*models.Flow, 0)
	err := db.Table(r.TableName()).
		Where(IDs).
		Find(&entity).
		Error
	if err != nil {
		return nil, err
	}

	return entity, nil
}

// FindFlows find flow list
func (r *flowRepo) FindFlows(db *gorm.DB, condition map[string]interface{}) ([]*models.Flow, error) {
	flows := make([]*models.Flow, 0)
	err := db.Table(r.TableName()).
		Where(condition).
		Where("status != ?", models.DELETED).
		Find(&flows).
		Error
	if err != nil {
		return nil, err
	}
	return flows, nil
}

// GetFlows find flow list
func (r *flowRepo) GetFlows(db *gorm.DB, condition map[string]interface{}) ([]*models.Flow, error) {
	flows := make([]*models.Flow, 0)
	err := db.Table(r.TableName()).
		Where(condition).
		Find(&flows).
		Error
	if err != nil {
		return nil, err
	}
	return flows, nil
}

// FindFlows find flow list
func (r *flowRepo) FindFlowList(db *gorm.DB, condition map[string]interface{}) ([]*models.Flow, error) {
	flows := make([]*models.Flow, 0)
	err := db.Table(r.TableName()).Select([]string{"name", "id", "status"}).
		Where(condition).
		Where("status != ?", models.DELETED).
		Find(&flows).
		Error
	if err != nil {
		return nil, err
	}
	return flows, nil
}

// UpdateFlows update flows
func (r *flowRepo) UpdateFlows(db *gorm.DB, condition map[string]interface{}, updateMap map[string]interface{}) error {
	updateMap["modify_time"] = time2.Now()
	err := db.Table(r.TableName()).Where(condition).Updates(updateMap).Error

	return err
}

func (r *flowRepo) FindPageFlows(db *gorm.DB, condition map[string]interface{}, page, size int) ([]*models.Flow, int64) {
	flows := make([]*models.Flow, 0)
	db = db.Table(r.TableName()).Where(condition).Where("status != ?", models.DELETED)
	var num int64
	if page > 0 && size > 0 {
		db = db.Limit(size).Offset((page - 1) * size)
	}
	db = db.Order("create_time" + " " + "desc")
	err := db.Find(&flows).Count(&num).Error
	if err != nil {
		return flows, 0
	}
	return flows, num
}

func (r *flowRepo) FindByProcessID(db *gorm.DB, processID string) (*models.Flow, error) {
	entity := new(models.Flow)
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

func (r *flowRepo) DeleteByIDs(db *gorm.DB, IDs []string) error {
	err := db.Table(r.TableName()).Where("id in (?)", IDs).Delete(&models.Flow{}).Error
	return err
}
