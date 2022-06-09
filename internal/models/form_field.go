package models

import "gorm.io/gorm"

// FormField info
type FormField struct {
	BaseModel

	FlowID         string `json:"flowId"`
	FormID         string `json:"formId"`
	FieldName      string `json:"fieldName"`
	FieldValuePath string `json:"fieldValuePath"`
}

// FormFieldRepo interface
type FormFieldRepo interface {
	Create(db *gorm.DB, model *FormField) error
	Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error
	Delete(db *gorm.DB, ID string) error
	DeleteByFlowID(db *gorm.DB, flowID string) error
	FindByFlowID(db *gorm.DB, ID string) ([]*FormField, error)
	DeleteByFlowIDs(db *gorm.DB, flowIDs []string) error
}
