package mysql

import (
	"github.com/quanxiang-cloud/flow/internal/models"
	"github.com/quanxiang-cloud/flow/pkg/misc/id2"
	"github.com/quanxiang-cloud/flow/pkg/misc/time2"
	"gorm.io/gorm"
)

type commentRepo struct{}

// NewCommentRepo new repo
func NewCommentRepo() models.CommentRepo {
	return &commentRepo{}
}

// TableName db table name
func (r *commentRepo) TableName() string {
	return "flow_comment"
}

// Create create model
func (r *commentRepo) Create(db *gorm.DB, entity *models.Comment) error {
	entity.ID = id2.GenID()
	entity.CreateTime = time2.Now()
	err := db.Table(r.TableName()).
		Create(entity).
		Error
	return err
}

// Update update model
func (r *commentRepo) Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error {
	updateMap["modify_time"] = time2.Now()
	err := db.Table(r.TableName()).Where("id=?", ID).Updates(updateMap).Error

	return err
}

// DeleteByID delete model
func (r *commentRepo) Delete(db *gorm.DB, ID string) error {
	entity := &models.Comment{BaseModel: models.BaseModel{
		ID: ID,
	}}
	err := db.Table(r.TableName()).Delete(entity).Error
	return err
}

// FindByID find model by ID
func (r *commentRepo) FindByID(db *gorm.DB, ID string) (*models.Comment, error) {
	entity := new(models.Comment)
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

func (r *commentRepo) FindComments(db *gorm.DB, condition map[string]interface{}, order string) ([]*models.Comment, error) {
	comments := make([]*models.Comment, 0)
	db = db.Table(r.TableName()).Where(condition)
	if order != "" {
		db = db.Order(order)
	}
	err := db.Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *commentRepo) DeleteByInstanceIDs(db *gorm.DB, InstanceIDs []string) error {
	err := db.Table(r.TableName()).Where("flow_instance_id in (?)", InstanceIDs).Delete(&models.Comment{}).Error
	return err
}
