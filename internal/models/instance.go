package models

import (
	"github.com/quanxiang-cloud/flow/pkg/page"
	"gorm.io/gorm"
)

// Instance info
type Instance struct {
	BaseModel

	AppID             string          `json:"appId"`
	AppStatus         string          `json:"appStatus"`
	AppName           string          `json:"appName"`
	FlowID            string          `json:"flowId"`
	ProcessInstanceID string          `json:"processInstanceId"`
	FormID            string          `json:"formId"`
	FormInstanceID    string          `json:"formInstanceId"`
	Name              string          `json:"name"`
	ApplyNo           string          `json:"applyNo"`
	ApplyUserID       string          `json:"applyUserId"`
	ApplyUserName     string          `json:"applyUserName"`
	Status            string          `json:"status"`
	FormData          interface{}     `gorm:"-" json:"formData"`
	FormSchema        interface{}     `gorm:"-" json:"formSchema"`
	Tasks             []ActTaskEntity `gorm:"-" json:"tasks"`
	Nodes             []NodeModel     `gorm:"-" json:"nodes"`
}

// InstanceRepo interface
type InstanceRepo interface {
	Create(db *gorm.DB, model *Instance) error
	Update(db *gorm.DB, ID string, updateMap map[string]interface{}) error
	Delete(db *gorm.DB, ID string) error
	FindByID(db *gorm.DB, ID string) (*Instance, error)
	FindByIDs(db *gorm.DB, IDs []string) ([]*Instance, error)
	FindByProcessInstanceIDs(db *gorm.DB, processInstanceIDs []string) ([]*Instance, error)
	GetEntityByProcessInstanceID(db *gorm.DB, processInstanceID string) (*Instance, error)
	FindInstances(db *gorm.DB, formInstanceID string, status []string) ([]*Instance, error)

	PageInstances(db *gorm.DB, req *PageInstancesReq) ([]*Instance, int64, error)
	GetInstances(db *gorm.DB, condition map[string]interface{}) ([]*Instance, error)
	DeleteByFlowIDs(db *gorm.DB, flowIDs []string) error
	UpdateAppStatus(db *gorm.DB, appID string, appStatus string) error
}

// PageInstancesReq struct
type PageInstancesReq struct {
	page.ReqPage
	ApplyUserID     string
	Status          string
	CreateTimeBegin string
	CreateTimeEnd   string
	Keyword         string
	AppID           string
}
