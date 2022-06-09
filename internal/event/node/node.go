package node

import (
	"github.com/quanxiang-cloud/flow/internal/flow"
	"github.com/quanxiang-cloud/flow/internal/models"
	"github.com/quanxiang-cloud/flow/pkg/client"
	"gorm.io/gorm"
)

// Node struct
type Node struct {
	Db                    *gorm.DB
	FlowRepo              models.FlowRepo
	InstanceRepo          models.InstanceRepo
	InstanceVariablesRepo models.InstanceVariablesRepo
	AbnormalTaskRepo      models.AbnormalTaskRepo
	InstanceExecutionRepo models.InstanceExecutionRepo
	FlowVariable          models.VariablesRepo
	Urge                  flow.Urge
	Flow                  flow.Flow
	Instance              flow.Instance
	OperationRecord       flow.OperationRecord
	Task                  flow.Task
	FormAPI               client.Form
	MessageCenterAPI      client.MessageCenter
	StructorAPI           client.Structor
	ProcessAPI            client.Process
	IdentityAPI           client.Identity
}

// SetDB set db
func (n *Node) SetDB(db *gorm.DB) {
	n.Db = db
}
