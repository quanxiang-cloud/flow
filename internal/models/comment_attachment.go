package models

import "gorm.io/gorm"

// CommentAttachment info
type CommentAttachment struct {
	BaseModel

	FlowCommentID  string `json:"flowCommentId"`
	AttachmentName string `json:"attachmentName"`
	AttachmentURL  string `json:"attachmentUrl"`
}

// CommentAttachmentRepo interface
type CommentAttachmentRepo interface {
	Create(db *gorm.DB, model *CommentAttachment) error
	Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error
	Delete(db *gorm.DB, ID string) error
	FindByID(db *gorm.DB, ID string) (*CommentAttachment, error)
	FindAttachments(db *gorm.DB, condition map[string]interface{}, order string) ([]*CommentAttachment, error)
}
