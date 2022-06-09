package node

import (
	"context"
	"fmt"
	"github.com/quanxiang-cloud/flow/internal/convert"
	"github.com/quanxiang-cloud/flow/pkg/client"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/utils"
	"strings"
)

// UserTask struct
type UserTask struct {
	*Node
}

// NewUserTask new
func NewUserTask(conf *config.Configs, node *Node) *UserTask {
	return &UserTask{
		Node: node,
	}
}

// Init event
func (n *UserTask) Init(ctx context.Context, eventData *EventData) error {
	fmt.Println("事件：加载AssigneeList，节点为" + eventData.NodeDefKey)
	// add dynamic assignee user list
	flowInstanceEntity, err := n.InstanceRepo.GetEntityByProcessInstanceID(n.Db, eventData.ProcessInstanceID)
	if err != nil {
		return err
	}

	// assigneeList dynamic handle users
	_, assigneeList := n.Flow.GetTaskHandleUserIDs(ctx, eventData.Shape, flowInstanceEntity)

	fmt.Println("事件：加载AssigneeList，AssigneeList=" + utils.ChangeStringArrayToString(assigneeList))
	if len(assigneeList) > 0 {
		return n.ProcessAPI.SetProcessVariables(ctx, eventData.ProcessID, eventData.ProcessInstanceID, eventData.NodeDefKey, convert.AssigneeList, assigneeList)
	}

	return nil
}

// Execute event
func (n *UserTask) Execute(ctx context.Context, eventData *EventData) error {
	flow, err := n.FlowRepo.FindByProcessID(n.Db, eventData.ProcessID)
	if err != nil {
		return err
	}
	instance, err := n.InstanceRepo.GetEntityByProcessInstanceID(n.Db, eventData.ProcessInstanceID)
	if err != nil {
		return err
	}

	taskIDs := eventData.TaskID // 如果是会签则taskID是数组，逗号间隔的数组
	if len(taskIDs) > 0 {
		apiReq := client.GetTasksReq{
			TaskID: strings.Split(taskIDs, ","),
		}
		resp, err := n.ProcessAPI.GetTasks(ctx, apiReq)
		if err != nil {
			return err
		}
		if len(resp.Data) == 0 {
			return nil
		}
		for _, value := range resp.Data {
			n.Task.TaskInitHandle(ctx, flow, instance, value, eventData.UserID)
		}
	}
	return nil
}
