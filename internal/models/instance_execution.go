package models

import (
	"gorm.io/gorm"
)

// InstanceExecution info
type InstanceExecution struct {
	BaseModel

	ProcessInstanceID string `json:"processInstanceID"`
	ExecutionID       string `json:"executionID"`
	Result            string `json:"result"`
}

// InstanceExecutionRepo interface
type InstanceExecutionRepo interface {
	Create(db *gorm.DB, model *InstanceExecution) error
	Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error

	Delete(db *gorm.DB, ID string) error
	DeleteByProcessInstanceIDs(db *gorm.DB, processInstanceIDs []string) error
	FindByID(db *gorm.DB, ID string) (*InstanceExecution, error)
	FindByIDs(db *gorm.DB, IDs []string) ([]*InstanceExecution, error)
	FindByExecutionIDs(db *gorm.DB, executionIDs []string) ([]*InstanceExecution, error)
}
