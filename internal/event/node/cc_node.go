package node

import (
	"context"
	"github.com/quanxiang-cloud/flow/internal/convert"
	"github.com/quanxiang-cloud/flow/pkg/client"
	"github.com/quanxiang-cloud/flow/pkg/config"
)

// CC struct
type CC struct {
	*Node
}

// NewCC new
func NewCC(conf *config.Configs, node *Node) *CC {
	return &CC{
		Node: node,
	}
}

// Init event
func (n *CC) Init(ctx context.Context, eventData *EventData) error {
	return nil
}

// Execute event
func (n *CC) Execute(ctx context.Context, eventData *EventData) error {
	instance, err := n.InstanceRepo.GetEntityByProcessInstanceID(n.Db, eventData.ProcessInstanceID)
	if err != nil {
		return err
	}

	ccUsers, _ := n.Flow.GetTaskHandleUserIDs(ctx, eventData.Shape, instance)

	if len(ccUsers) == 0 {
		return nil
	}

	apiReq := &client.AddTaskReq{
		InstanceID: eventData.ProcessInstanceID,
		NodeDefKey: eventData.Shape.ID,
		Name:       eventData.Shape.Data.NodeData.Name,
		Desc:       convert.CcTask,
		Assignee:   ccUsers,
	}
	_, err = n.ProcessAPI.AddNonModelTask(ctx, apiReq)
	return err
}
