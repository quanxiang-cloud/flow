package models

import "gorm.io/gorm"

// DispatcherCallback info
type DispatcherCallback struct {
	BaseModel

	Type              string `json:"type"`
	OtherInfo         string `json:"otherInfo"`
	ProcessInstanceID string `json:"processInstanceId"`
	TaskDefKey        string `json:"taskDefKey"`
}

// DispatcherCallbackRepo interface
type DispatcherCallbackRepo interface {
	Create(db *gorm.DB, model *DispatcherCallback) error
	Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error
	Delete(db *gorm.DB, ID string) error
	FindByID(db *gorm.DB, ID string) (*DispatcherCallback, error)
}
