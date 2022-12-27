/*
Copyright 2022 QuanxiangCloud Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package node

import (
	"context"
	"errors"
	"github.com/quanxiang-cloud/flow/internal/convert"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/utils"
	"github.com/quanxiang-cloud/flow/rpc/pb"
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

// InitBegin event
func (n *VariableUpdate) InitBegin(ctx context.Context, eventData *EventData) (*pb.NodeEventRespData, error) {
	flow, err := n.FlowRepo.FindByProcessID(n.Db, eventData.ProcessID)
	if err != nil {
		return nil, err
	}
	if flow == nil {
		flowProcessRelation, err := n.FlowProcessRelationRepo.FindByProcessID(n.Db, eventData.ProcessID)
		if err != nil {
			return nil, err
		}
		flow, err = n.FlowRepo.FindByID(n.Db, flowProcessRelation.FlowID)
		if err != nil {
			return nil, err
		}
		if flow == nil {
			return nil, errors.New("variable update node not match flow")
		}
	}
	formShape, err := convert.GetShapeByChartType(flow.BpmnText, convert.FormData)
	if err != nil {
		return nil, err
	}

	instance, err := n.InstanceRepo.GetEntityByProcessInstanceID(n.Db, eventData.ProcessInstanceID)
	if err != nil {
		return nil, err
	}

	bd := eventData.Shape.Data.BusinessData
	if bd == nil {
		return nil, nil
	}

	var assignmentRules []map[string]interface{}
	if v := bd["assignmentRules"]; v != nil {
		arr := v.([]interface{})
		for _, e := range arr {
			assignmentRules = append(assignmentRules, e.(map[string]interface{}))
		}
	}
	if assignmentRules == nil {
		return nil, nil
	}
	variables, err := n.Flow.GetInstanceVariableValues(ctx, instance)
	if err != nil {
		return nil, err
	}
	for _, e := range assignmentRules {
		variableName := utils.Strval(e["variableName"])
		valueFrom := utils.Strval(e["valueFrom"])
		valueOf := e["valueOf"]

		value, err := n.Instance.Cal(ctx, valueFrom, valueOf, nil, instance, variables, nil, formShape.ID)
		if err != nil {
			return nil, err
		}

		// n.instanceVariablesRepo.f
		fieldValue, fieldType := utils.StrvalAndType(value)
		err = n.InstanceVariablesRepo.UpdateTypeAndValue(n.Db, eventData.ProcessInstanceID, variableName, fieldType, fieldValue)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

// InitEnd event
func (n *VariableUpdate) InitEnd(ctx context.Context, eventData *EventData) (*pb.NodeEventRespData, error) {

	return nil, nil
}
