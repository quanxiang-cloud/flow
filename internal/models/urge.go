package models

import "gorm.io/gorm"

// Urge info
type Urge struct {
	BaseModel

	TaskID            string `json:"taskId"`
	ProcessInstanceID string `json:"processInstanceId"`
}

// UrgeRepo interface
type UrgeRepo interface {
	Create(db *gorm.DB, model *Urge) error
	Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error
	Delete(db *gorm.DB, ID string) error
	FindByID(db *gorm.DB, ID string) (*Urge, error)
	FindTaskIDs(db *gorm.DB) ([]string, error)
	DeleteByProcessInstanceIDs(db *gorm.DB, processInstanceIDs []string) error

	FindByTaskID(db *gorm.DB, taskID string) ([]*Urge, error)
	GetUrgeNums(db *gorm.DB, taskIDs []string) (map[string]int64, error)
}

// UrgeNumModel urge num model
type UrgeNumModel struct {
	TaskID  string
	UrgeNum int64
}
