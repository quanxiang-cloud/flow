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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/quanxiang-cloud/flow/internal/convert"
	"github.com/quanxiang-cloud/flow/internal/models"
	"github.com/quanxiang-cloud/flow/pkg/client"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/misc/logger"
	"github.com/quanxiang-cloud/flow/pkg/redis"
	"github.com/quanxiang-cloud/flow/pkg/utils"
	"github.com/quanxiang-cloud/flow/rpc/pb"
	"net/url"
	"reflect"
	"strings"
	"time"
)

// WebHook struct
type WebHook struct {
	*Node
	PolyAPI client.PolyAPI
}

// NewWebHook new
func NewWebHook(conf *config.Configs, node *Node) *WebHook {
	return &WebHook{
		Node:    node,
		PolyAPI: client.NewPolyAPI(conf),
	}
}

// InitBegin event
func (n *WebHook) InitBegin(ctx context.Context, eventData *EventData) (*pb.NodeEventRespData, error) {
	fmt.Println("================执行b")
	//if !n.CheckRefuse(ctx, n.Db, eventData.ProcessID) {
	//	return nil, nil
	//}
	//flow, err := n.FlowRepo.FindByProcessID(n.Db, eventData.ProcessID)
	//if err != nil {
	//	return nil, err
	//}
	//// do req
	//queryFlag := false
	//requestBody := make(map[string]interface{})
	//requestHeader := make(map[string]string) // req header
	//path := ""
	//bd := eventData.Shape.Data.BusinessData
	//if bd == nil {
	//	return nil, nil
	//}
	//hookType := utils.Strval(bd["type"])
	//var conf map[string]interface{}
	//if v := bd["config"]; v != nil {
	//	conf = v.(map[string]interface{})
	//}
	//var inputs []convert.Input
	//if v := conf["inputs"]; v != nil {
	//	arr := v.([]interface{})
	//	for _, e := range arr {
	//		marshal, err := json.Marshal(e)
	//		if err != nil {
	//			return nil, err
	//		}
	//
	//		input := convert.Input{}
	//		err = json.Unmarshal(marshal, &input)
	//		if err != nil {
	//			return nil, err
	//		}
	//		inputs = append(inputs, input)
	//	}
	//}
	//v := url.Values{}
	//if flow.TriggerMode == "FORM_DATA" {
	//	formShape, err := convert.GetShapeByChartType(flow.BpmnText, convert.FormData)
	//	if err != nil {
	//		return nil, err
	//	}
	//	formDefKey := formShape.ID
	//
	//	instance, err := n.InstanceRepo.GetEntityByProcessInstanceID(n.Db, eventData.ProcessInstanceID)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	variables, err := n.Instance.GetInstanceVariableValues(ctx, instance)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	dataReq := client.FormDataConditionModel{
	//		AppID:   instance.AppID,
	//		TableID: instance.FormID,
	//		DataID:  instance.FormInstanceID,
	//	}
	//	dataResp, err := n.FormAPI.GetFormData(ctx, dataReq)
	//	if err != nil {
	//		return nil, err
	//	}
	//	if dataResp == nil {
	//		return nil, nil
	//	}
	//
	//	for k, v := range dataResp {
	//		variables[k] = v
	//	}
	//
	//	// gen req
	//
	//	if hookType == "request" {
	//		var api map[string]interface{}
	//		if v := conf["api"]; v != nil {
	//			api = v.(map[string]interface{})
	//		}
	//		path = utils.Strval(api["value"])
	//	} else {
	//		path = utils.Strval(conf["sendUrl"])
	//	}
	//
	//	for _, e := range inputs {
	//		if e.In == convert.Header {
	//			val, err := n.webHookCal(ctx, e, variables, formDefKey)
	//			if err != nil {
	//				return nil, err
	//			}
	//			if e.Name != "" {
	//				requestHeader[e.Name] = utils.Strval(val)
	//			}
	//		} else if e.In == convert.Path {
	//			val, err := n.webHookCal(ctx, e, variables, formDefKey)
	//			if err != nil {
	//				return nil, err
	//			}
	//			path = strings.Replace(path, e.Name, utils.Strval(val), 1)
	//		} else if e.In == convert.Query {
	//			val, err := n.webHookCal(ctx, e, variables, formDefKey)
	//			if err != nil {
	//				return nil, err
	//			}
	//			queryFlag = true
	//			if e.Name != "" {
	//				v.Add(e.Name, utils.Strval(val))
	//			}
	//			// path += e.Name + "=" + utils.Strval(val) + "&"
	//		} else if e.In == convert.Body {
	//			if e.Type == "object" {
	//				obj := make(map[string]interface{})
	//				if v := e.Data; v != nil {
	//					if reflect.TypeOf(v).Kind() == reflect.Slice {
	//						arr := v.([]interface{})
	//						bdata := make(map[string]interface{})
	//						for _, e1 := range arr {
	//							marshal, err := json.Marshal(e1)
	//							if err != nil {
	//								return nil, err
	//							}
	//
	//							input := convert.Input{}
	//							err = json.Unmarshal(marshal, &input)
	//							if err != nil {
	//								return nil, err
	//							}
	//							val, err := n.exchangeParam(ctx, input, nil, "")
	//							if err != nil {
	//								return nil, err
	//							}
	//							if input.Name != "" {
	//								bdata[input.Name] = val
	//							}
	//
	//						}
	//						obj[e.Name] = bdata
	//						requestBody = obj
	//					} else {
	//						// requestBody[e.Name] = e.Data
	//
	//						val, err := n.webHookCal(ctx, e, nil, "")
	//						if err != nil {
	//							return nil, err
	//						}
	//						if e.Name != "" {
	//							requestBody[e.Name] = utils.Strval(val)
	//						}
	//
	//					}
	//
	//				}
	//			} else {
	//				if v := e.Data; v != nil {
	//					if reflect.TypeOf(v).Kind() == reflect.Slice {
	//						arr := v.([]interface{})
	//						for _, e := range arr {
	//							marshal, err := json.Marshal(e)
	//							if err != nil {
	//								return nil, err
	//							}
	//
	//							input := convert.Input{}
	//							err = json.Unmarshal(marshal, &input)
	//							if err != nil {
	//								return nil, err
	//							}
	//							val, err := n.exchangeParam(ctx, input, variables, formDefKey)
	//							if err != nil {
	//								return nil, err
	//							}
	//							if input.Name != "" {
	//								requestBody[input.Name] = val
	//							}
	//
	//						}
	//					} else {
	//						// requestBody[e.Name] = e.Data
	//
	//						val, err := n.webHookCal(ctx, e, variables, formDefKey)
	//						if err != nil {
	//							return nil, err
	//						}
	//						if e.Name != "" {
	//							requestBody[e.Name] = utils.Strval(val)
	//						}
	//
	//					}
	//
	//				}
	//			}
	//
	//		}
	//
	//	}
	//} else {
	//	mp := conf["api"].(map[string]interface{})
	//	path = mp["value"].(string)
	//	for _, e := range inputs {
	//		if e.In == convert.Header {
	//			val, err := n.webHookCal(ctx, e, nil, "")
	//			if err != nil {
	//				return nil, err
	//			}
	//			if e.Name != "" {
	//				requestHeader[e.Name] = utils.Strval(val)
	//			}
	//		} else if e.In == convert.Path {
	//			val, err := n.webHookCal(ctx, e, nil, "")
	//			if err != nil {
	//				return nil, err
	//			}
	//			path = strings.Replace(path, e.Name, utils.Strval(val), 1)
	//		} else if e.In == convert.Query {
	//			val, err := n.webHookCal(ctx, e, nil, "")
	//			if err != nil {
	//				return nil, err
	//			}
	//			queryFlag = true
	//			if e.Name != "" {
	//				v.Add(e.Name, utils.Strval(val))
	//			}
	//			// path += e.Name + "=" + utils.Strval(val) + "&"
	//		} else if e.In == convert.Body {
	//			if e.Type == "object" {
	//				if v := e.Data; v != nil {
	//					if reflect.TypeOf(v).Kind() == reflect.Slice {
	//						arr := v.([]interface{})
	//						bdata := make(map[string]interface{})
	//						for _, e1 := range arr {
	//							marshal, err := json.Marshal(e1)
	//							if err != nil {
	//								return nil, err
	//							}
	//
	//							input := convert.Input{}
	//							err = json.Unmarshal(marshal, &input)
	//							if err != nil {
	//								return nil, err
	//							}
	//							val, err := n.exchangeParam(ctx, input, nil, "")
	//							if err != nil {
	//								return nil, err
	//							}
	//							if input.Name != "" {
	//								bdata[input.Name] = val
	//							}
	//
	//						}
	//						requestBody[e.Name] = bdata
	//					} else {
	//						// requestBody[e.Name] = e.Data
	//
	//						val, err := n.webHookCal(ctx, e, nil, "")
	//						if err != nil {
	//							return nil, err
	//						}
	//						if e.Name != "" {
	//							requestBody[e.Name] = utils.FormatValue(utils.Strval(val), e.Type)
	//						}
	//
	//					}
	//
	//				}
	//			} else {
	//				if v := e.Data; v != nil {
	//					if reflect.TypeOf(v).Kind() == reflect.Slice {
	//						arr := v.([]interface{})
	//						for _, e := range arr {
	//							marshal, err := json.Marshal(e)
	//							if err != nil {
	//								return nil, err
	//							}
	//
	//							input := convert.Input{}
	//							err = json.Unmarshal(marshal, &input)
	//							if err != nil {
	//								return nil, err
	//							}
	//							val, err := n.exchangeParam(ctx, input, nil, "")
	//							if err != nil {
	//								return nil, err
	//							}
	//							if input.Name != "" {
	//								if val != nil {
	//									requestBody[input.Name] = utils.Strval(val)
	//								}
	//							}
	//
	//						}
	//					} else {
	//						// requestBody[e.Name] = e.Data
	//
	//						val, err := n.webHookCal(ctx, e, nil, "")
	//						if err != nil {
	//							return nil, err
	//						}
	//						if e.Name != "" {
	//							if val != nil {
	//								requestBody[e.Name] = utils.FormatValue(utils.Strval(val), e.Type)
	//							}
	//
	//						}
	//
	//					}
	//
	//				}
	//
	//			}
	//
	//		}
	//
	//	}
	//}
	//
	//if queryFlag {
	//	if strings.Contains(path, "?") {
	//		if v.Encode() != "" {
	//			path += "&"
	//			path += v.Encode()
	//		}
	//
	//	} else {
	//		if v.Encode() != "" {
	//			path += "?"
	//			path += v.Encode()
	//		}
	//
	//	}
	//
	//}
	//var method = ""
	//if hookType == "request" {
	//	_, ok := requestHeader["Content-Type"]
	//	if !ok {
	//		requestHeader["Content-Type"] = "application/json"
	//	}
	//	method = utils.Strval(conf["method"])
	//	resp, err := n.PolyAPI.InnerRequest(ctx, path, requestBody, requestHeader, method)
	//	if err != nil {
	//		return nil, err
	//	}
	//	if resp == nil {
	//		return nil, nil
	//	}
	//	apiMap := n.apiMap(resp)
	//	// save resp
	//	for k, v := range apiMap {
	//		code := "$" + eventData.Shape.ID + "." + k
	//		variable, err := n.InstanceVariablesRepo.FindVariablesByCode(n.Db, eventData.ProcessInstanceID, code)
	//		if err != nil {
	//			return nil, err
	//		}
	//
	//		fieldValue, fieldType := utils.StrvalAndType(v)
	//		if variable.ID != "" {
	//			if err := n.InstanceVariablesRepo.UpdateTypeAndValue(n.Db, eventData.ProcessInstanceID, code, fieldType, fieldValue); err != nil {
	//				return nil, err
	//			}
	//		} else {
	//			variable := models.InstanceVariables{
	//				ProcessInstanceID: eventData.ProcessInstanceID,
	//				Code:              code,
	//				FieldType:         fieldType,
	//				Value:             fieldValue,
	//			}
	//			if err := n.InstanceVariablesRepo.Create(n.Db, &variable); err != nil {
	//				return nil, err
	//			}
	//		}
	//	}
	//} else if hookType == "send" {
	//	requestHeader["Content-Type"] = utils.Strval(conf["contentType"])
	//	_, ok := requestHeader["Content-Type"]
	//	if !ok {
	//		requestHeader["Content-Type"] = "application/json"
	//	}
	//	method = utils.Strval(conf["sendMethod"])
	//	_, err := n.PolyAPI.SendRequest(ctx, path, requestBody, requestHeader, method)
	//	if err != nil {
	//		return nil, err
	//	}
	//}
	return nil, nil
}

func (n *WebHook) apiMap(apiMap map[string]interface{}) map[string]interface{} {
	ret := make(map[string]interface{})
	for k, v := range apiMap {
		if v != nil {
			typeOf := reflect.TypeOf(v).Kind()
			if typeOf == reflect.Map {
				m := n.apiMap(v.(map[string]interface{}))
				for k1, v1 := range m {
					ret[k+"."+k1] = v1
				}
			} else {
				ret[k] = v
			}
		}

	}
	return ret
}

// InitEnd event
func (n *WebHook) InitEnd(ctx context.Context, eventData *EventData) (*pb.NodeEventRespData, error) {
	fmt.Println("================执行e")
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
			return nil, errors.New("webhook node not match flow")
		}
	}
	preNodeKey := CheckPreNode(flow.BpmnText, eventData.NodeDefKey)
	if preNodeKey != "" {
		if preNodeKey == "processBranchTarget" {
			time.Sleep(6 * time.Second)
		} else {
			var i = 0
			for {
				get, err := redis.ClusterClient.Get(ctx, "flow:node:"+eventData.ProcessInstanceID+":"+preNodeKey).Result()
				if err != nil {
					fmt.Println(err)
				}
				if get == "over" {
					break
				}
				logger.Logger.Info("等待上个节点执行完成---", preNodeKey)
				i++
				if i >= 13 {
					break
				}
				time.Sleep(1 * time.Second)
			}
		}
	}

	// do req
	queryFlag := false
	requestBody := make(map[string]interface{})
	requestHeader := make(map[string]string) // req header
	path := ""
	bd := eventData.Shape.Data.BusinessData
	if bd == nil {
		return nil, nil
	}
	hookType := utils.Strval(bd["type"])
	var conf map[string]interface{}
	if v := bd["config"]; v != nil {
		conf = v.(map[string]interface{})
	}
	var inputs []convert.Input
	if v := conf["inputs"]; v != nil {
		arr := v.([]interface{})
		for _, e := range arr {
			marshal, err := json.Marshal(e)
			if err != nil {
				return nil, err
			}

			input := convert.Input{}
			err = json.Unmarshal(marshal, &input)
			if err != nil {
				return nil, err
			}
			inputs = append(inputs, input)
		}
	}
	if hookType == "request" {
		var api map[string]interface{}
		if v := conf["api"]; v != nil {
			api = v.(map[string]interface{})
		}
		path = utils.Strval(api["value"])
	} else {
		path = utils.Strval(conf["sendUrl"])
	}
	v := url.Values{}
	if flow.TriggerMode == "FORM_DATA" {
		formShape, err := convert.GetShapeByChartType(flow.BpmnText, convert.FormData)
		if err != nil {
			return nil, err
		}
		formDefKey := formShape.ID

		instance, err := n.InstanceRepo.GetEntityByProcessInstanceID(n.Db, eventData.ProcessInstanceID)
		if err != nil {
			return nil, err
		}

		variables, err := n.Instance.GetInstanceVariableValues(ctx, instance)
		if err != nil {
			return nil, err
		}
		var ref = make(map[string]map[string]string)
		for k := range inputs {
			if inputs[k].FieldType == "AssociatedRecords" {
				m := make(map[string]string)
				m["appID"] = instance.AppID
				m["tableID"] = inputs[k].TableID
				m["type"] = "associated_records"
				ref[inputs[k].FieldName] = m
			}
			if inputs[k].FieldType == "ForeignTable" {
				m := make(map[string]string)
				m["appID"] = instance.AppID
				m["tableID"] = inputs[k].TableID
				m["type"] = "foreign_table"
				ref[inputs[k].FieldName] = m
			}
			if inputs[k].FieldType == "SubTable" {
				m := make(map[string]string)
				m["appID"] = instance.AppID
				m["tableID"] = inputs[k].TableID
				m["type"] = "sub_table"
				ref[inputs[k].FieldName] = m
			}
			if inputs[k].FieldType == "AggregationRecords" {
				m := make(map[string]string)
				m["appID"] = instance.AppID
				m["tableID"] = inputs[k].TableID
				m["type"] = "aggregation"
				ref[inputs[k].FieldName] = m
			}
		}

		dataReq := client.FormDataConditionModel{
			AppID:   instance.AppID,
			TableID: instance.FormID,
			DataID:  instance.FormInstanceID,
			Ref:     ref,
		}
		dataResp, err := n.FormAPI.GetFormData(ctx, dataReq)
		if err != nil {
			return nil, err
		}
		if dataResp == nil {
			return nil, nil
		}

		for k, v := range dataResp {
			variables[k] = v
		}

		// gen req

		for _, e := range inputs {
			if e.Name == "" {
				continue
			}
			if e.In == convert.Header {
				val, err := n.webHookCal(ctx, e, variables, formDefKey)
				if err != nil {
					return nil, err
				}
				if e.Name != "" {
					requestHeader[e.Name] = utils.Strval(val)
				}
			} else if e.In == convert.Path {
				val, err := n.webHookCal(ctx, e, variables, formDefKey)
				if err != nil {
					return nil, err
				}
				path = strings.Replace(path, e.Name, utils.Strval(val), 1)
			} else if e.In == convert.Query {
				val, err := n.webHookCal(ctx, e, variables, formDefKey)
				if err != nil {
					return nil, err
				}
				queryFlag = true
				if e.Name != "" {
					v.Add(e.Name, utils.Strval(val))
				}
				// path += e.Name + "=" + utils.Strval(val) + "&"
			} else if e.In == convert.Body {
				if e.Type == "object" {
					obj := make(map[string]interface{})
					if v := e.Data; v != nil {
						if reflect.TypeOf(v).Kind() == reflect.Slice {
							arr := v.([]interface{})
							bdata := make(map[string]interface{})
							for _, e1 := range arr {
								marshal, err := json.Marshal(e1)
								if err != nil {
									return nil, err
								}

								input := convert.Input{}
								err = json.Unmarshal(marshal, &input)
								if err != nil {
									return nil, err
								}
								val, err := n.exchangeParam(ctx, input, nil, "")
								if err != nil {
									return nil, err
								}
								if input.Name != "" {
									bdata[input.Name] = val
								}

							}
							obj[e.Name] = bdata
							requestBody = obj
						} else {
							// requestBody[e.Name] = e.Data

							val, err := n.webHookCal(ctx, e, nil, "")
							if err != nil {
								return nil, err
							}
							if e.Name != "" {
								requestBody[e.Name] = utils.Strval(val)
							}

						}

					}
				} else {
					if v := e.Data; v != nil {
						if reflect.TypeOf(v).Kind() == reflect.Slice {
							arr := v.([]interface{})
							for _, e := range arr {
								marshal, err := json.Marshal(e)
								if err != nil {
									return nil, err
								}

								input := convert.Input{}
								err = json.Unmarshal(marshal, &input)
								if err != nil {
									return nil, err
								}
								val, err := n.exchangeParam(ctx, input, variables, formDefKey)
								if err != nil {
									return nil, err
								}
								if input.Name != "" {
									requestBody[input.Name] = val
								}

							}
						} else {
							// requestBody[e.Name] = e.Data

							val, err := n.webHookCal(ctx, e, variables, formDefKey)
							if err != nil {
								return nil, err
							}
							if e.Name != "" {
								requestBody[e.Name] = utils.Strval(val)
							}

						}

					}
				}

			}

		}
	} else {
		for _, e := range inputs {
			if e.In == convert.Header {
				val, err := n.webHookCal(ctx, e, nil, "")
				if err != nil {
					return nil, err
				}
				if e.Name != "" {
					requestHeader[e.Name] = utils.Strval(val)
				}
			} else if e.In == convert.Path {
				val, err := n.webHookCal(ctx, e, nil, "")
				if err != nil {
					return nil, err
				}
				path = strings.Replace(path, e.Name, utils.Strval(val), 1)
			} else if e.In == convert.Query {
				val, err := n.webHookCal(ctx, e, nil, "")
				if err != nil {
					return nil, err
				}
				queryFlag = true
				if e.Name != "" {
					v.Add(e.Name, utils.Strval(val))
				}
				// path += e.Name + "=" + utils.Strval(val) + "&"
			} else if e.In == convert.Body {
				if e.Type == "object" {
					if v := e.Data; v != nil {
						if reflect.TypeOf(v).Kind() == reflect.Slice {
							arr := v.([]interface{})
							bdata := make(map[string]interface{})
							for _, e1 := range arr {
								marshal, err := json.Marshal(e1)
								if err != nil {
									return nil, err
								}

								input := convert.Input{}
								err = json.Unmarshal(marshal, &input)
								if err != nil {
									return nil, err
								}
								val, err := n.exchangeParam(ctx, input, nil, "")
								if err != nil {
									return nil, err
								}
								if input.Name != "" {
									bdata[input.Name] = val
								}

							}
							requestBody[e.Name] = bdata
						} else {
							// requestBody[e.Name] = e.Data

							val, err := n.webHookCal(ctx, e, nil, "")
							if err != nil {
								return nil, err
							}
							if e.Name != "" {
								requestBody[e.Name] = utils.FormatValue(utils.Strval(val), e.Type)
							}

						}

					}
				} else {
					if v := e.Data; v != nil {
						if reflect.TypeOf(v).Kind() == reflect.Slice {
							arr := v.([]interface{})
							for _, e := range arr {
								marshal, err := json.Marshal(e)
								if err != nil {
									return nil, err
								}

								input := convert.Input{}
								err = json.Unmarshal(marshal, &input)
								if err != nil {
									return nil, err
								}
								val, err := n.exchangeParam(ctx, input, nil, "")
								if err != nil {
									return nil, err
								}
								if input.Name != "" {
									if val != nil {
										requestBody[input.Name] = utils.Strval(val)
									}
								}

							}
						} else {
							// requestBody[e.Name] = e.Data

							val, err := n.webHookCal(ctx, e, nil, "")
							if err != nil {
								return nil, err
							}
							if e.Name != "" {
								if val != nil {
									requestBody[e.Name] = utils.FormatValue(utils.Strval(val), e.Type)
								}

							}

						}

					}

				}

			}

		}
	}

	if queryFlag {
		if strings.Contains(path, "?") {
			if v.Encode() != "" {
				path += "&"
				path += v.Encode()
			}

		} else {
			if v.Encode() != "" {
				path += "?"
				path += v.Encode()
			}

		}

	}
	var method = ""
	if hookType == "request" {
		_, ok := requestHeader["Content-Type"]
		if !ok {
			requestHeader["Content-Type"] = "application/json"
		}
		method = utils.Strval(conf["method"])
		resp, err := n.PolyAPI.InnerRequest(ctx, path, requestBody, requestHeader, method)
		if err != nil {
			return nil, err
		}
		if resp == nil {
			return nil, nil
		}
		apiMap := n.apiMap(resp)
		// save resp
		for k, v := range apiMap {
			code := "$" + eventData.Shape.ID + "." + k
			variable, err := n.InstanceVariablesRepo.FindVariablesByCode(n.Db, eventData.ProcessInstanceID, code)
			if err != nil {
				return nil, err
			}

			fieldValue, fieldType := utils.StrvalAndType(v)
			if variable.ID != "" {
				if err := n.InstanceVariablesRepo.UpdateTypeAndValue(n.Db, eventData.ProcessInstanceID, code, fieldType, fieldValue); err != nil {
					return nil, err
				}
			} else {
				variable := models.InstanceVariables{
					ProcessInstanceID: eventData.ProcessInstanceID,
					Code:              code,
					FieldType:         fieldType,
					Value:             fieldValue,
				}
				if err := n.InstanceVariablesRepo.Create(n.Db, &variable); err != nil {
					return nil, err
				}
			}
		}
	} else if hookType == "send" {
		requestHeader["Content-Type"] = utils.Strval(conf["contentType"])
		_, ok := requestHeader["Content-Type"]
		if !ok {
			requestHeader["Content-Type"] = "application/json"
		}
		method = utils.Strval(conf["sendMethod"])
		_, err := n.PolyAPI.SendRequest(ctx, path, requestBody, requestHeader, method)
		if err != nil {
			return nil, err
		}
	}
	fmt.Println("================执行e over")
	redis.ClusterClient.SetEX(ctx, "flow:node:"+eventData.ProcessInstanceID+":"+eventData.NodeDefKey, "over", 20*time.Second)
	return nil, nil
}

// exchangeParam exchange body val
func (n *WebHook) exchangeParam(ctx context.Context, input convert.Input, variables map[string]interface{}, formDefKey string) (interface{}, error) {
	if input.Type == "object" {
		var inputs []convert.Input
		if v := input.Data; v != nil {
			arr := v.([]interface{})
			for _, e := range arr {
				marshal, err := json.Marshal(e)
				if err != nil {
					return nil, err
				}

				input := convert.Input{}
				err = json.Unmarshal(marshal, &input)
				if err != nil {
					return nil, err
				}
				inputs = append(inputs, input)
			}
		}
		param0 := make(map[string]interface{})
		for _, e := range inputs {
			val, err := n.exchangeParam(ctx, e, variables, formDefKey)
			if err != nil {
				return nil, err
			}
			param0[e.Name] = val
		}
		return param0, nil
	}
	val, err := n.webHookCal(ctx, input, variables, formDefKey)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (n *WebHook) webHookCal(ctx context.Context, input convert.Input, variables map[string]interface{}, formDefKey string) (interface{}, error) {
	if input.Type == "direct_expr" {
		// 判断有没有运算符
		if strings.Contains(utils.Strval(input.Data), "$") {

			expression := strings.TrimSpace(utils.Strval(input.Data))
			var flag = false
			if strings.Contains(expression, "+") || strings.Contains(expression, "-") || strings.Contains(expression, "*") || strings.Contains(expression, "/") {
				flag = true
			}
			expression = strings.Replace(expression, "$variable.", "", -1)
			expression = strings.Replace(expression, "$"+formDefKey+".", "", -1)

			for k, v := range variables {
				// if strings.HasPrefix(k, "$") && strings.Contains(expression, k) {
				if strings.Contains(expression, k) {
					expression = strings.Replace(expression, k, utils.Strval(v), -1)
					break
				}
			}
			if !flag {
				return expression, nil
			}

			expression = strings.Replace(expression, "$", "", -1)
			ret, err := n.StructorAPI.CalExpression(ctx, map[string]interface{}{
				"expression": expression,
				"parameter":  variables,
			})
			if err != nil { // 公式计算，不能计算字符串，只能计算数值
				return expression, nil
			}
			return ret, nil
		}
		return input.Data, nil
	}
	return input.Data, nil
}
