package mysql

import (
	"github.com/quanxiang-cloud/flow/internal/models"
	"github.com/quanxiang-cloud/flow/pkg/misc/id2"
	"github.com/quanxiang-cloud/flow/pkg/misc/time2"
	"gorm.io/gorm"
)

type commentAttachmentRepo struct{}

// NewCommentAttachmentRepo new repo
func NewCommentAttachmentRepo() models.CommentAttachmentRepo {
	return &commentAttachmentRepo{}
}

// TableName db table name
func (r *commentAttachmentRepo) TableName() string {
	return "flow_comment_attachment"
}

// Create create model
func (r *commentAttachmentRepo) Create(db *gorm.DB, entity *models.CommentAttachment) error {
	entity.ID = id2.GenID()
	entity.CreateTime = time2.Now()
	err := db.Table(r.TableName()).
		Create(entity).
		Error
	return err
}

// Update update model
func (r *commentAttachmentRepo) Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error {
	updateMap["modify_time"] = time2.Now()
	err := db.Table(r.TableName()).Where("id=?", ID).Updates(updateMap).Error

	return err
}

// DeleteByID delete model
func (r *commentAttachmentRepo) Delete(db *gorm.DB, ID string) error {
	entity := &models.CommentAttachment{BaseModel: models.BaseModel{
		ID: ID,
	}}
	err := db.Table(r.TableName()).Delete(entity).Error
	return err
}

// FindByID find model by ID
func (r *commentAttachmentRepo) FindByID(db *gorm.DB, ID string) (*models.CommentAttachment, error) {
	entity := new(models.CommentAttachment)
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

func (r *commentAttachmentRepo) FindAttachments(db *gorm.DB, condition map[string]interface{}, order string) ([]*models.CommentAttachment, error) {
	ca := make([]*models.CommentAttachment, 0)
	db = db.Table(r.TableName()).Where(condition)
	if order != "" {
		db = db.Order(order)
	}
	err := db.Find(&ca).Error
	if err != nil {
		return nil, err
	}
	return ca, nil
}
