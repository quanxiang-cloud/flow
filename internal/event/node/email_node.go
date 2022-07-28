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
	"fmt"
	"github.com/quanxiang-cloud/flow/internal/convert"
	"github.com/quanxiang-cloud/flow/internal/models"
	"github.com/quanxiang-cloud/flow/pkg/client"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/misc/logger"
	"github.com/quanxiang-cloud/flow/pkg/utils"
	"github.com/quanxiang-cloud/flow/rpc/pb"
	"strings"
)

// emailBD 邮件节点业务数据结构
type emailBD struct {
	ApprovePersons approvePersons    `json:"approvePersons"`
	Content        string            `json:"content"`
	MesAttachment  []interface{}     `json:"mes_attachment"`
	TemplateID     string            `json:"templateId"`
	Title          string            `json:"title"`
	FormulaFields  map[string]string `json:"formulaFields"`
	FieldType      map[string]string `json:"fieldType"`
}

// approvePersons ap
type approvePersons struct {
	convert.ApprovePersonsModel
	MultipleFields []string `json:"multipleFields"`
}

// Email struct
type Email struct {
	*Node
}

// NewEmail new
func NewEmail(conf *config.Configs, node *Node) *Email {
	return &Email{
		Node: node,
	}
}

// InitBegin event
func (n *Email) InitBegin(ctx context.Context, eventData *EventData) (*pb.NodeEventRespData, error) {
	logger.Logger.Info("发送邮件，processID=", eventData.ProcessID, "emailDefKey=", eventData.NodeDefKey)
	var bdData emailBD
	b, err := json.Marshal(eventData.Shape.Data.BusinessData)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &bdData); err != nil {
		return nil, err
	}

	// gen req params
	mesAttachments := make([]map[string]interface{}, 0)
	if v := bdData.MesAttachment; v != nil {
		//arr := v.([]interface{})
		for _, e := range v {
			tmp := e.(map[string]interface{})
			mesAttachment := make(map[string]interface{})
			mesAttachment["name"] = utils.Strval(tmp["file_name"])
			mesAttachment["path"] = utils.Strval(tmp["file_url"])
			mesAttachments = append(mesAttachments, mesAttachment)
		}
	}

	emailAddr := make([]string, 0)
	for _, user := range bdData.ApprovePersons.Users {
		// valid check
		s := user["email"].(string)
		if !utils.EmailAddressValid(&s) {
			logger.Logger.Warnf("value of [%s] is not valid email address", s)
			continue
		}
		emailAddr = append(emailAddr, s)
	}

	if len(emailAddr) > 0 {
		email := client.Email{
			To: emailAddr,
			Contents: client.Contents{
				Content: bdData.Content,
			},
			Title: utils.Strval(bdData.Title),
			Files: mesAttachments,
		}
		msgReq := client.MsgReq{
			Email: email,
		}
		err = n.MessageCenterAPI.MessageCreate(ctx, msgReq)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

// InitEnd event
func (n *Email) InitEnd(ctx context.Context, eventData *EventData) (*pb.NodeEventRespData, error) {
	// 业务数据
	var bdData emailBD
	var handleUsers []*client.UserInfoResp
	b, err := json.Marshal(eventData.Shape.Data.BusinessData)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &bdData); err != nil {
		return nil, err
	}
	instance, err := n.InstanceRepo.GetEntityByProcessInstanceID(n.Db, eventData.ProcessInstanceID)
	if err != nil {
		return nil, err
	}

	dataReq := client.FormDataConditionModel{
		AppID:   instance.AppID,
		TableID: instance.FormID,
		DataID:  instance.FormInstanceID,
	}
	formResp, err := n.FormAPI.GetFormData(ctx, dataReq)
	if err != nil {
		return nil, err
	}
	if formResp == nil {
		return nil, nil
	}
	// replace content
	content := utils.Strval(bdData.Content)
	value := n.Flow.FormatFormValue(instance, formResp)
	//var fieldType map[string]interface{}
	//if v := bdData.FieldType; v != nil {
	//	fieldType = v.(map[string]interface{})
	//}
	for k, v := range value {
		t := bdData.FieldType[k]
		if t == "datepicker" {
			vt := v.(string)
			if strings.Contains(vt, ".000Z") {
				vt = strings.Replace(vt, ".000Z", "+0000", 1)
			}
			v = utils.ChangeISO8601ToBjTime(vt)
		}
		content = strings.Replace(content, "${"+k+"}", utils.Strval(v), 1)
	}
	switch bdData.ApprovePersons.Type {
	case convert.EmailTypeOfField:
		handleUsers = n.Flow.GetTaskHandleUsers2(ctx, bdData.ApprovePersons, instance)
	case convert.EmailTypeOfMultipleField:
		handleUsers = n.MultipleFieldsHandle(ctx, &bdData, formResp, instance)
	default:
		return nil, fmt.Errorf("unkonw type [%s] of approvePersons.type", bdData.ApprovePersons.Type)
	}
	// gen req params
	mesAttachments := make([]map[string]interface{}, 0)
	if v := bdData.MesAttachment; v != nil {
		//arr := v.([]interface{})
		for _, e := range v {
			tmp := e.(map[string]interface{})
			mesAttachment := make(map[string]interface{})
			mesAttachment["name"] = utils.Strval(tmp["file_name"])
			mesAttachment["path"] = utils.Strval(tmp["file_url"])
			mesAttachments = append(mesAttachments, mesAttachment)
		}
	}
	emailAddr := make([]string, 0)
	for _, user := range handleUsers {
		// valid check
		if !utils.EmailAddressValid(&user.Email) {
			logger.Logger.Warnf("value of [%s] is not valid email address", user.Email)
			continue
		}
		emailAddr = append(emailAddr, user.Email)
	}
	if len(emailAddr) > 0 {
		email := client.Email{
			To: emailAddr,
			Contents: client.Contents{
				Content: content,
			},
			Title: utils.Strval(bdData.Title),
			Files: mesAttachments,
		}
		msgReq := client.MsgReq{
			Email: email,
		}
		// post msg

		err = n.MessageCenterAPI.MessageCreate(ctx, msgReq)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

// MultipleFieldsHandle 多个字段处理
// bd 业务数据
// formData 表单数据
func (n *Email) MultipleFieldsHandle(ctx context.Context, bd *emailBD, formData map[string]interface{}, instance *models.Instance) []*client.UserInfoResp {
	var userCtx []*client.UserInfoResp
	var sendEmail []string
	for i := 0; i < len(bd.ApprovePersons.Fields); i++ {
		field := bd.ApprovePersons.Fields[i]
		fieldType := bd.FieldType[field]
		// 没有找到对应类型记录日志
		if len(fieldType) == 0 {
			continue
		}
		switch fieldType {
		// 直接取值的字段
		case convert.CompOfSelect:
			fallthrough
		case convert.CompOfRadioGroup:
			fallthrough
		case convert.CompOfInput:
			fallthrough
		case convert.CompOfTextarea:
			_, value := n.FormAPI.GetValue(formData, field, formData[field])
			if value != nil {
				v := value.(string)
				// 多行文本
				if fieldType == convert.CompOfTextarea {
					values := strings.Split(v, "\n")
					sendEmail = append(sendEmail, values...)
				} else {
					sendEmail = append(sendEmail, v)
				}
			}
		// 列表值
		case convert.CompOfMultipleSelect:
			fallthrough
		case convert.CompOfCheckboxGroup:
			_, value := n.FormAPI.GetValue(formData, field, formData[field])
			if value != nil {
				v := value.([]interface{})
				for i := 0; i < len(v); i++ {
					sendEmail = append(sendEmail, v[i].(string))
				}
			}
		// 人员选择
		case convert.CompOfUserPicker:
			u := n.Flow.GetTaskHandleUsers2(ctx, bd.ApprovePersons, instance)
			userCtx = append(userCtx, u...)
		}
	}
	// 去重
	notRepeatEmail := utils.RemoveReplicaSliceString(sendEmail)
	for i := 0; i < len(notRepeatEmail); i++ {
		userCtx = append(userCtx, &client.UserInfoResp{
			Email: notRepeatEmail[i],
		})
	}
	return userCtx
}
