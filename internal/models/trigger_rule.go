package models

import "gorm.io/gorm"

// TriggerRule info
type TriggerRule struct {
	BaseModel
	FlowID string `json:"flowId"`
	FormID string `json:"formId"`
	Rule   string `json:"rule"`
}

// TriggerRuleRepo interface
type TriggerRuleRepo interface {
	Create(db *gorm.DB, model *TriggerRule) error
	Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error
	Delete(db *gorm.DB, ID string) error
	FindByID(db *gorm.DB, ID string) (*TriggerRule, error)
	DeleteByFlowIDS(db *gorm.DB, flowIDS []string) error
	FindTriggerRules(db *gorm.DB, condition map[string]interface{}) ([]*TriggerRule, error)
	FindTriggerRulesByFormID(db *gorm.DB, formID string) ([]*TriggerRule, error)
}
