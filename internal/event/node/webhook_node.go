package node

import (
	"context"
	"encoding/json"
	"github.com/quanxiang-cloud/flow/internal/convert"
	"github.com/quanxiang-cloud/flow/internal/models"
	"github.com/quanxiang-cloud/flow/pkg/client"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/utils"
	"net/url"
	"reflect"
	"strings"
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

// Init event
func (n *WebHook) Init(ctx context.Context, eventData *EventData) error {
	flow, err := n.FlowRepo.FindByProcessID(n.Db, eventData.ProcessID)
	if err != nil {
		return err
	}
	formShape, err := convert.GetShapeByChartType(flow.BpmnText, convert.FormData)
	if err != nil {
		return err
	}
	formDefKey := formShape.ID

	instance, err := n.InstanceRepo.GetEntityByProcessInstanceID(n.Db, eventData.ProcessInstanceID)
	if err != nil {
		return err
	}

	variables, err := n.Instance.GetInstanceVariableValues(ctx, instance)
	if err != nil {
		return err
	}

	dataReq := client.FormDataConditionModel{
		AppID:   instance.AppID,
		TableID: instance.FormID,
		DataID:  instance.FormInstanceID,
	}
	dataResp, err := n.FormAPI.GetFormData(ctx, dataReq)
	if err != nil {
		return err
	}
	if dataResp == nil {
		return nil
	}

	for k, v := range dataResp {
		variables[k] = v
	}

	bd := eventData.Shape.Data.BusinessData
	if bd == nil {
		return nil
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
				return err
			}

			input := convert.Input{}
			err = json.Unmarshal(marshal, &input)
			if err != nil {
				return err
			}
			inputs = append(inputs, input)
		}
	}

	// gen req
	requestBody := make(map[string]interface{})
	requestHeader := make(map[string]string) // req header
	path := ""
	method := utils.Strval(conf["method"])
	if hookType == "request" {
		var api map[string]interface{}
		if v := conf["api"]; v != nil {
			api = v.(map[string]interface{})
		}
		path = utils.Strval(api["value"])
	} else {
		path = utils.Strval(conf["sendUrl"])
	}

	// do req
	queryFlag := false
	v := url.Values{}
	for _, e := range inputs {
		if e.In == convert.Header {
			val, err := n.webHookCal(ctx, e, variables, formDefKey)
			if err != nil {
				return err
			}
			requestHeader[e.Name] = utils.Strval(val)
		} else if e.In == convert.Path {
			val, err := n.webHookCal(ctx, e, variables, formDefKey)
			if err != nil {
				return err
			}
			path = strings.Replace(path, e.Name, utils.Strval(val), 1)
		} else if e.In == convert.Query {
			val, err := n.webHookCal(ctx, e, variables, formDefKey)
			if err != nil {
				return err
			}
			queryFlag = true
			v.Add(e.Name, utils.Strval(val))
			// path += e.Name + "=" + utils.Strval(val) + "&"
		} else if e.In == convert.Body {
			if v := e.Data; v != nil {
				if reflect.TypeOf(v).Kind() == reflect.Slice {
					arr := v.([]interface{})
					for _, e := range arr {
						marshal, err := json.Marshal(e)
						if err != nil {
							return err
						}

						input := convert.Input{}
						err = json.Unmarshal(marshal, &input)
						if err != nil {
							return err
						}
						val, err := n.exchangeParam(ctx, input, variables, formDefKey)
						if err != nil {
							return err
						}
						requestBody[input.Name] = val
					}
				} else {
					// requestBody[e.Name] = e.Data

					val, err := n.webHookCal(ctx, e, variables, formDefKey)
					if err != nil {
						return err
					}
					requestBody[e.Name] = utils.Strval(val)
				}

			}

		}

	}
	if queryFlag {
		path += "?"
		path += v.Encode()
	}
	if hookType == "request" {
		_, ok := requestHeader["Content-Type"]
		if !ok {
			requestHeader["Content-Type"] = "application/json"
		}
		resp, err := n.PolyAPI.InnerRequest(ctx, path, requestBody, requestHeader, method)
		if err != nil {
			return err
		}
		if resp == nil {
			return nil
		}
		apiMap := n.apiMap(resp)
		// save resp
		for k, v := range apiMap {
			code := "$" + eventData.Shape.ID + "." + k
			variable, err := n.InstanceVariablesRepo.FindVariablesByCode(n.Db, eventData.ProcessInstanceID, code)
			if err != nil {
				return err
			}

			fieldValue, fieldType := utils.StrvalAndType(v)
			if variable.ID != "" {
				if err := n.InstanceVariablesRepo.UpdateTypeAndValue(n.Db, eventData.ProcessInstanceID, code, fieldType, fieldValue); err != nil {
					return err
				}
			} else {
				variable := models.InstanceVariables{
					ProcessInstanceID: eventData.ProcessInstanceID,
					Code:              code,
					FieldType:         fieldType,
					Value:             fieldValue,
				}
				if err := n.InstanceVariablesRepo.Create(n.Db, &variable); err != nil {
					return err
				}
			}
		}
	} else if hookType == "send" {
		requestHeader["Content-Type"] = utils.Strval(conf["contentType"])
		_, ok := requestHeader["Content-Type"]
		if !ok {
			requestHeader["Content-Type"] = "application/json"
		}

		_, err := n.PolyAPI.SendRequest(ctx, path, requestBody, requestHeader, method)
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *WebHook) apiMap(apiMap map[string]interface{}) map[string]interface{} {
	ret := make(map[string]interface{})
	for k, v := range apiMap {
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
	return ret
}

// Execute event
func (n *WebHook) Execute(ctx context.Context, eventData *EventData) error {
	return nil
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
			expression = strings.Replace(expression, "$variable.", "", -1)
			expression = strings.Replace(expression, "$"+formDefKey+".", "", -1)

			for k, v := range variables {
				// if strings.HasPrefix(k, "$") && strings.Contains(expression, k) {
				if strings.Contains(expression, k) {
					expression = strings.Replace(expression, k, utils.Strval(v), -1)
				}
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
