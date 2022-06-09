package node

import (
	"context"
	"github.com/quanxiang-cloud/flow/internal/convert"
	"github.com/quanxiang-cloud/flow/pkg/client"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/utils"
)

// DataCreate struct
type DataCreate struct {
	*Node
}

// NewDataCreate new
func NewDataCreate(conf *config.Configs, node *Node) *DataCreate {
	return &DataCreate{
		Node: node,
	}
}

// Init event
func (n *DataCreate) Init(ctx context.Context, eventData *EventData) error {
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

	var createRules map[string]interface{}
	if v := bd["createRule"]; v != nil {
		createRules = v.(map[string]interface{})
	}
	if createRules == nil {
		return nil
	}

	variables, err := n.Instance.GetInstanceVariableValues(ctx, instance)
	if err != nil {
		return err
	}

	// master form
	filedValueReq := make(map[string]interface{})
	for k, v := range createRules {
		tmp := v.(map[string]interface{})
		valueFrom := utils.Strval(tmp["valueFrom"])
		valueOf := tmp["valueOf"]

		value, err := n.Instance.Cal(ctx, valueFrom, valueOf, nil, instance, variables, nil, formShape.ID)
		if err != nil {
			return err
		}
		filedValueReq[k] = value
	}

	// child form
	refDataReq := make(map[string]client.RefData)
	var refs map[string]interface{}
	if v := bd["ref"]; v != nil {
		refs = v.(map[string]interface{})
	}
	if refs != nil {
		for k, v := range refs {
			tmp := v.(map[string]interface{})
			refData := client.RefData{
				AppID:   instance.AppID,
				TableID: utils.Strval(tmp["tableId"]),
				Type:    utils.Strval(tmp["type"]),
			}

			var subTableCreateRules []map[string]interface{}
			if t := tmp["createRules"]; t != nil {
				arr := t.([]interface{})
				for _, e := range arr {
					subTableCreateRules = append(subTableCreateRules, e.(map[string]interface{}))
				}
			}

			new := make([]client.CreateEntity, 0)
			if subTableCreateRules != nil {
				for _, e := range subTableCreateRules {
					record := client.CreateEntity{}
					recordEntity := make(map[string]interface{})
					for k1, v1 := range e {
						valueFrom := utils.Strval(v1.(map[string]interface{})["valueFrom"])
						valueOf := v1.(map[string]interface{})["valueOf"]
						value, err := n.Instance.Cal(ctx, valueFrom, valueOf, nil, instance, variables, nil, formShape.ID)
						if err != nil {
							return err
						}
						recordEntity[k1] = value
					}
					record.Entity = recordEntity
					new = append(new, record)
				}
			}
			refData.New = new
			refDataReq[k] = refData
		}
	}

	err = n.FormAPI.CreateData(ctx, instance.AppID, utils.Strval(bd["targetTableId"]), client.CreateEntity{
		Entity: filedValueReq,
		Ref:    refDataReq,
	}, bd["silent"].(bool))
	if err != nil {
		return err
	}

	return nil
}

// Execute event
func (n *DataCreate) Execute(ctx context.Context, eventData *EventData) error {
	return nil
}
