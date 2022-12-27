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
	"github.com/quanxiang-cloud/flow/pkg"
	"github.com/quanxiang-cloud/flow/pkg/client"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/utils"
	"github.com/quanxiang-cloud/flow/rpc/pb"
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
func (n *DataCreate) InitBegin(ctx context.Context, eventData *EventData) (*pb.NodeEventRespData, error) {
	//flow, err := n.FlowRepo.FindByProcessID(n.Db, eventData.ProcessID)
	//if err != nil {
	//	return nil, err
	//}
	//
	//bd := eventData.Shape.Data.BusinessData
	//if bd == nil {
	//	return nil, nil
	//}
	//
	//var createRules map[string]interface{}
	//if v := bd["createRule"]; v != nil {
	//	createRules = v.(map[string]interface{})
	//}
	//if createRules == nil {
	//	return nil, nil
	//}
	//var variables = make(map[string]interface{})
	//var formID = ""
	//filedValueReq := make(map[string]interface{})
	//
	//instance, err := n.InstanceRepo.GetEntityByProcessInstanceID(n.Db, eventData.ProcessInstanceID)
	//if err != nil {
	//	return nil, err
	//}
	//if !n.CheckRefuse(ctx, n.Db, instance.ProcessInstanceID) {
	//	return nil, nil
	//}
	//
	//if flow.TriggerMode == "FORM_DATA" {
	//	formShape, err := convert.GetShapeByChartType(flow.BpmnText, convert.FormData)
	//	if err != nil {
	//		return nil, err
	//	}
	//	// master form
	//
	//	for k, v := range createRules {
	//		tmp := v.(map[string]interface{})
	//		valueFrom := utils.Strval(tmp["valueFrom"])
	//		valueOf := tmp["valueOf"]
	//
	//		variables, err := n.Instance.GetInstanceVariableValues(ctx, instance)
	//		if err != nil {
	//			return nil, err
	//		}
	//
	//		value, err := n.Instance.Cal(ctx, valueFrom, valueOf, nil, instance, variables, nil, formShape.ID)
	//		if err != nil {
	//			return nil, err
	//		}
	//		filedValueReq[k] = value
	//		formID = formShape.ID
	//	}
	//}
	//if flow.TriggerMode == "FORM_TIME" {
	//	filedValueReq = createRules
	//}
	//
	//// child form
	//refDataReq := make(map[string]client.RefData)
	//var refs map[string]interface{}
	//if v := bd["ref"]; v != nil {
	//	refs = v.(map[string]interface{})
	//}
	//if refs != nil {
	//	for k, v := range refs {
	//		tmp := v.(map[string]interface{})
	//		refData := client.RefData{
	//			AppID:   flow.AppID,
	//			TableID: utils.Strval(tmp["tableId"]),
	//			Type:    utils.Strval(tmp["type"]),
	//		}
	//
	//		var subTableCreateRules []map[string]interface{}
	//		if t := tmp["createRules"]; t != nil {
	//			arr := t.([]interface{})
	//			for _, e := range arr {
	//				subTableCreateRules = append(subTableCreateRules, e.(map[string]interface{}))
	//			}
	//		}
	//
	//		new1 := make([]client.CreateEntity, 0)
	//		if flow.TriggerMode == "FORM_DATA" {
	//			if subTableCreateRules != nil {
	//				for _, e := range subTableCreateRules {
	//					record := client.CreateEntity{}
	//					recordEntity := make(map[string]interface{})
	//					for k1, v1 := range e {
	//						valueFrom := utils.Strval(v1.(map[string]interface{})["valueFrom"])
	//						valueOf := v1.(map[string]interface{})["valueOf"]
	//						value, err := n.Instance.Cal(ctx, valueFrom, valueOf, nil, instance, variables, nil, formID)
	//						if err != nil {
	//							return nil, err
	//						}
	//						recordEntity[k1] = value
	//					}
	//					record.Entity = recordEntity
	//					new1 = append(new1, record)
	//				}
	//			}
	//		}
	//		if flow.TriggerMode == "FORM_DATA" {
	//			refData.New = new1
	//			refDataReq[k] = refData
	//		}
	//		if flow.TriggerMode == "FORM_TIME" {
	//			filedValueReq = getDataFromShape(createRules)
	//			refDataReq = getDataFromRef(flow.AppID, refs)
	//		}
	//
	//	}
	//}
	//ctx = pkg.SetRequestID2(ctx, instance.RequestID)
	//err = n.FormAPI.CreateData(ctx, flow.AppID, utils.Strval(bd["targetTableId"]), client.CreateEntity{
	//	Entity: filedValueReq,
	//	Ref:    refDataReq,
	//}, bd["silent"].(bool))
	//if err != nil {
	//	return nil, err
	//}
	return nil, nil
}

// for FORM_TIME
func getDataFromShape(data map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range data {
		if v1, ok := v.(map[string]interface{}); ok {
			if v1["valueFrom"].(string) == "fixedValue" {
				res[k] = v1["valueOf"]
			}
		}
	}
	return res
}

// for FORM_TIME
func getDataFromRef(appID string, data map[string]interface{}) map[string]client.RefData {
	res := make(map[string]client.RefData)
	for k, v := range data {
		refData := client.RefData{}
		refData.AppID = appID
		if v1, ok := v.(map[string]interface{}); ok {
			refData.TableID = v1["tableId"].(string)
			refData.Type = v1["type"].(string)
			news := make([]map[string]interface{}, 0)
			if v2, ok2 := v1["createRules"].([]interface{}); ok2 {
				for _, v3 := range v2 {
					enty := make(map[string]interface{})
					entity := make(map[string]interface{})
					if v4, ok4 := v3.(map[string]interface{}); ok4 {
						for k44, v44 := range v4 {
							if v5, ok5 := v44.(map[string]interface{}); ok5 {
								if v5["valueFrom"].(string) == "fixedValue" {
									entity[k44] = v5["valueOf"]
								}
							}
						}
					}
					enty["entity"] = entity
					news = append(news, enty)
				}
			}
			refData.New = news
		}
		res[k] = refData
	}
	return res
}

// Execute event
func (n *DataCreate) InitEnd(ctx context.Context, eventData *EventData) (*pb.NodeEventRespData, error) {
	//if !n.CheckRefuse(ctx, n.Db, eventData.ProcessInstanceID) {
	//	return nil, nil
	//}
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
			return nil, errors.New("send create form data not match flow")
		}
	}

	bd := eventData.Shape.Data.BusinessData
	if bd == nil {
		return nil, nil
	}

	var createRules map[string]interface{}
	if v := bd["createRule"]; v != nil {
		createRules = v.(map[string]interface{})
	}
	if createRules == nil {
		return nil, nil
	}
	var variables = make(map[string]interface{})
	var formID = ""
	filedValueReq := make(map[string]interface{})

	instance, err := n.InstanceRepo.GetEntityByProcessInstanceID(n.Db, eventData.ProcessInstanceID)
	if err != nil {
		return nil, err
	}

	if flow.TriggerMode == "FORM_DATA" {
		formShape, err := convert.GetShapeByChartType(flow.BpmnText, convert.FormData)
		if err != nil {
			return nil, err
		}
		// master form

		for k, v := range createRules {
			tmp := v.(map[string]interface{})
			valueFrom := utils.Strval(tmp["valueFrom"])
			valueOf := tmp["valueOf"]

			variables, err := n.Instance.GetInstanceVariableValues(ctx, instance)
			if err != nil {
				return nil, err
			}

			value, err := n.Instance.Cal(ctx, valueFrom, valueOf, nil, instance, variables, nil, formShape.ID)
			if err != nil {
				return nil, err
			}
			filedValueReq[k] = value
			formID = formShape.ID
		}
	}
	if flow.TriggerMode == "FORM_TIME" {
		filedValueReq = createRules
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
				AppID:   flow.AppID,
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

			new1 := make([]client.CreateEntity, 0)
			if flow.TriggerMode == "FORM_DATA" {
				if subTableCreateRules != nil {
					for _, e := range subTableCreateRules {
						record := client.CreateEntity{}
						recordEntity := make(map[string]interface{})
						for k1, v1 := range e {
							valueFrom := utils.Strval(v1.(map[string]interface{})["valueFrom"])
							valueOf := v1.(map[string]interface{})["valueOf"]
							value, err := n.Instance.Cal(ctx, valueFrom, valueOf, nil, instance, variables, nil, formID)
							if err != nil {
								return nil, err
							}
							recordEntity[k1] = value
						}
						record.Entity = recordEntity
						new1 = append(new1, record)
					}
				}
			}
			if flow.TriggerMode == "FORM_DATA" {
				refData.New = new1
				refDataReq[k] = refData
			}
			if flow.TriggerMode == "FORM_TIME" {
				filedValueReq = getDataFromShape(createRules)
				refDataReq = getDataFromRef(flow.AppID, refs)
			}

		}
	}
	ctx = pkg.SetRequestID2(ctx, instance.RequestID)
	err = n.FormAPI.CreateData(ctx, flow.AppID, utils.Strval(bd["targetTableId"]), client.CreateEntity{
		Entity: filedValueReq,
		Ref:    refDataReq,
	}, bd["silent"].(bool))
	if err != nil {
		return nil, err
	}
	return nil, nil

}
