package node

import (
	"context"
	"github.com/quanxiang-cloud/flow/internal/convert"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/utils"
)

// VariableUpdate struct
type VariableUpdate struct {
	*Node
}

// NewVariableUpdate new
func NewVariableUpdate(conf *config.Configs, node *Node) *VariableUpdate {
	return &VariableUpdate{
		Node: node,
	}
}

// Init event
func (n *VariableUpdate) Init(ctx context.Context, eventData *EventData) error {
	flow, err := n.FlowRepo.FindByProcessID(n.Db, eventData.ProcessID)
	if err != nil {
		return err
	}
	formShape, err := convert.GetShapeByChartType(flow.BpmnText, convert.FormData)
	if err != nil {
		return err
	}

	instance, err := n.InstanceRepo.GetEntityByProcessInstanceID(n.Db, eventData.ProcessInstanceID)
	if err != nil {
		return err
	}

	bd := eventData.Shape.Data.BusinessData
	if bd == nil {
		return nil
	}

	var assignmentRules []map[string]interface{}
	if v := bd["assignmentRules"]; v != nil {
		arr := v.([]interface{})
		for _, e := range arr {
			assignmentRules = append(assignmentRules, e.(map[string]interface{}))
		}
	}
	if assignmentRules == nil {
		return nil
	}
	variables, err := n.Flow.GetInstanceVariableValues(ctx, instance)
	if err != nil {
		return err
	}
	for _, e := range assignmentRules {
		variableName := utils.Strval(e["variableName"])
		valueFrom := utils.Strval(e["valueFrom"])
		valueOf := e["valueOf"]

		value, err := n.Instance.Cal(ctx, valueFrom, valueOf, nil, instance, variables, nil, formShape.ID)
		if err != nil {
			return err
		}

		// n.instanceVariablesRepo.f
		fieldValue, fieldType := utils.StrvalAndType(value)
		err = n.InstanceVariablesRepo.UpdateTypeAndValue(n.Db, eventData.ProcessInstanceID, variableName, fieldType, fieldValue)
		if err != nil {
			return err
		}
	}
	return nil
}

// Execute event
func (n *VariableUpdate) Execute(ctx context.Context, eventData *EventData) error {

	return nil
}
