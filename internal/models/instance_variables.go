package models

import "gorm.io/gorm"

// InstanceVariables info
type InstanceVariables struct {
	BaseModel

	ProcessInstanceID string `json:"processInstanceID"`
	Name              string `json:"name"`
	Type              string `json:"type"`
	Code              string `json:"code"`
	FieldType         string `json:"fieldType"`
	Format            string `json:"format"`
	Value             string `json:"value"`
	Desc              string `json:"desc"`
}

// InstanceVariablesRepo interface
type InstanceVariablesRepo interface {
	Create(db *gorm.DB, model *InstanceVariables) error
	Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error
	UpdateByID(db *gorm.DB, ID string, update interface{}) error
	Delete(db *gorm.DB, ID string) error
	FindByID(db *gorm.DB, ID string) (*InstanceVariables, error)

	BatchCreate(db *gorm.DB, model []*InstanceVariables) error
	FindVariablesByProcessInstanceID(db *gorm.DB, processInstanceID string) ([]*InstanceVariables, error)
	UpdateVariable(db *gorm.DB, processInstanceID string, code string, value string) error
	UpdateTypeAndValue(db *gorm.DB, processInstanceID string, code string, fieldType string, value string) error
	DeleteByProcessInstanceIDs(db *gorm.DB, processInstanceIDs []string) error
	FindVariablesByCode(db *gorm.DB, processInstanceID string, code string) (*InstanceVariables, error)
}
