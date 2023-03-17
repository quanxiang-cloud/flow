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
	"fmt"
	"github.com/pkg/errors"
	"github.com/quanxiang-cloud/flow/pkg/client"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/misc/logger"
	"github.com/quanxiang-cloud/flow/pkg/redis"
	"github.com/quanxiang-cloud/flow/pkg/utils"
	"github.com/quanxiang-cloud/flow/rpc/pb"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Letter struct
type Letter struct {
	*Node
}

// NewLetter new
func NewLetter(conf *config.Configs, node *Node) *Letter {
	return &Letter{
		Node: node,
	}
}

// InitBegin event
func (n *Letter) InitBegin(ctx context.Context, eventData *EventData) (*pb.NodeEventRespData, error) {
	//logger.Logger.Info("发送站内信，processID=", eventData.ProcessID, "letterDefKey=", eventData.NodeDefKey)
	//var bdData emailBD
	//b, err := json.Marshal(eventData.Shape.Data.BusinessData)
	//if err != nil {
	//	return nil, err
	//}
	//if err := json.Unmarshal(b, &bdData); err != nil {
	//	return nil, err
	//}
	//
	//var recivers []client.Receivers
	//for _, user := range bdData.ApprovePersons.Users {
	//	reciver := client.Receivers{
	//		Type: 1,
	//		ID:   user["id"].(string),
	//		Name: user["ownerName"].(string),
	//	}
	//	recivers = append(recivers, reciver)
	//}
	//types, err := strconv.Atoi(eventData.Shape.Data.BusinessData["sort"].(string))
	//if len(recivers) > 0 {
	//	web := client.Web{
	//		IsSend: true,
	//		Contents: client.Contents{
	//			Content: bdData.Content,
	//		},
	//		Title:     utils.Strval(bdData.Title),
	//		Receivers: recivers,
	//		Types:     types,
	//	}
	//	msgReq := client.Mail{
	//		Web: web,
	//	}
	//	err = n.MessageCenterAPI.MessageCreateff(ctx, msgReq)
	//	if err != nil {
	//		return nil, err
	//	}
	//}

	return nil, nil
}

// InitEnd event
func (n *Letter) InitEnd(ctx context.Context, eventData *EventData) (*pb.NodeEventRespData, error) {
	logger.Logger.Info("发送站内信，processID=", eventData.ProcessID, "letterDefKey=", eventData.NodeDefKey)
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
			return nil, errors.New("send update form data not match flow")
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

	instance, err := n.InstanceRepo.GetEntityByProcessInstanceID(n.Db, eventData.ProcessInstanceID)
	if err != nil {
		return nil, err
	}

	bd := eventData.Shape.Data.BusinessData
	if bd == nil {
		return nil, nil
	}

	handleUsers := n.Flow.GetTaskHandleUsers2(ctx, bd["approvePersons"], instance)
	// gen req params
	var recivers []client.Receivers
	for _, user := range handleUsers {
		reciver := client.Receivers{
			Type: 1,
			ID:   user.ID,
			Name: user.UserName,
		}
		recivers = append(recivers, reciver)
	}
	types, err := strconv.Atoi(bd["sort"].(string))
	if err != nil {
		return nil, err
	}
	var content = ""
	// replace content
	content = utils.Strval(bd["content"])
	compile := regexp.MustCompile(`\$\{(.*?)\}`)
	allString := compile.FindAllStringSubmatch(content, -1)

	dataReq := client.FormDataConditionModel{
		AppID:   instance.AppID,
		TableID: instance.FormID,
		DataID:  instance.FormInstanceID,
	}
	for k := range allString {

		s := allString[k][1]
		split := strings.Split(s, ".")
		content = strings.Replace(content, allString[k][1], split[0], -1)
		if len(split) == 3 {
			dataReq.Ref = map[string]interface{}{
				split[0]: map[string]interface{}{
					"appID":   instance.AppID,
					"tableID": split[1],
					"type":    split[2],
				},
			}
		}
		if len(split) == 5 {
			dataReq.Ref = map[string]interface{}{
				split[0]: map[string]interface{}{
					"appID":         instance.AppID,
					"tableID":       split[1],
					"type":          split[2],
					"sourceFieldId": split[3],
					"aggType":       split[4],
				},
			}
		}
	}

	dataResp, err := n.FormAPI.GetFormData(ctx, dataReq)
	if err != nil {
		return nil, err
	}
	if dataResp == nil {
		return nil, err
	}

	value := n.Flow.FormatFormValue(instance, dataResp)
	var fieldType map[string]interface{}
	if v := bd["fieldType"]; v != nil {
		fieldType = v.(map[string]interface{})
	}
	for k, v := range value {
		t := fieldType[k]
		if t == "datepicker" {
			vt := v.(string)
			if strings.Contains(vt, ".000Z") {
				vt = strings.Replace(vt, ".000Z", "+0000", 1)
			}
			v = utils.ChangeISO8601ToBjTime(vt)
		}
		content = strings.Replace(content, "${"+k+"}", utils.Strval(v), 1)
	}

	web := client.Web{
		IsSend: true,
		Title:  utils.Strval(bd["title"]),
		Contents: client.Contents{
			Content: content,
		},
		Receivers: recivers,
		Types:     types,
	}
	m := client.Mail{
		Web: web,
	}
	// post msg
	err = n.MessageCenterAPI.MessageCreateff(ctx, m)
	if err != nil {
		return nil, err
	}
	redis.ClusterClient.SetEX(ctx, "flow:node:"+eventData.ProcessInstanceID+":"+eventData.NodeDefKey, "over", 20*time.Second)
	return nil, nil
}
