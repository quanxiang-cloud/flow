package models

import "gorm.io/gorm"

// Comment info
type Comment struct {
	BaseModel

	FlowInstanceID string `json:"flowInstanceId"`
	CommentUserID  string `json:"commentUserId"`
	Content        string `json:"content"`
}

// CommentRepo interface
type CommentRepo interface {
	Create(db *gorm.DB, model *Comment) error
	Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error
	Delete(db *gorm.DB, ID string) error
	FindByID(db *gorm.DB, ID string) (*Comment, error)
	FindComments(db *gorm.DB, condition map[string]interface{}, order string) ([]*Comment, error)
	DeleteByInstanceIDs(db *gorm.DB, InstanceIDs []string) error
}
