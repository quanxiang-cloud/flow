package models

import "gorm.io/gorm"

// InstanceStep info
type InstanceStep struct {
	BaseModel

	ProcessInstanceID string             `json:"processInstanceId"`
	TaskID            string             `json:"taskId"`
	TaskType          string             `json:"taskType"` // 节点类型：或签、会签、任填、全填、开始、结束
	TaskDefKey        string             `json:"taskDefKey"`
	TaskName          string             `json:"taskName"`
	HandleUserIDs     string             `json:"handleUserIds"`
	Status            string             `json:"status"` // 步骤处理结果，通过、拒绝、完成填写、已回退、打回重填、自动跳过、自动交给管理员
	NodeInstanceID    string             `json:"nodeInstanceId"`
	OperationRecords  []*OperationRecord `gorm:"-" json:"operationRecords"`
	FlowName          string             `gorm:"-" json:"flowName"`
	Reason            string             `gorm:"-" json:"reason"`
}

// InstanceStepRepo interface
type InstanceStepRepo interface {
	Create(db *gorm.DB, model *InstanceStep) error
	Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error
	Delete(db *gorm.DB, ID string) error
	FindByID(db *gorm.DB, ID string) (*InstanceStep, error)
	FindInstanceSteps(db *gorm.DB, condition *InstanceStep) ([]*InstanceStep, error)
	FindInstanceStepsByStatus(db *gorm.DB, processInstanceID string, status []string) ([]*InstanceStep, error)
	GetFlowInstanceStep(db *gorm.DB, processInstanceID string, nodeInstanceID string, status []string) ([]*InstanceStep, error)
	UpdateByNodeInstanceID(db *gorm.DB, nodeInstanceID string, updateMap map[string]interface{}) error
	DeleteByProcessInstanceIDs(db *gorm.DB, processInstanceIDs []string) error
}
