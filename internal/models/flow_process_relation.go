package models

import "gorm.io/gorm"

// FlowProcessRelation flowID,processD
type FlowProcessRelation struct {
	BaseModel
	FlowID string `json:"flowID"`

	BpmnText string `json:"bpmnText"` // flow model json

	ProcessID string `json:"processID"` // Process id

}

// FlowProcessRelationRepo interface
type FlowProcessRelationRepo interface {
	Create(db *gorm.DB, model *FlowProcessRelation) error
	DeleteByFlowID(db *gorm.DB, flowID string) error
	FindByProcessID(db *gorm.DB, processID string) (*FlowProcessRelation, error)
}
