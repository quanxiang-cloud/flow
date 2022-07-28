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
	"github.com/quanxiang-cloud/flow/pkg/client"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/misc/logger"
	"github.com/quanxiang-cloud/flow/pkg/utils"
	"github.com/quanxiang-cloud/flow/rpc/pb"
	"strconv"
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
	logger.Logger.Info("发送站内信，processID=", eventData.ProcessID, "letterDefKey=", eventData.NodeDefKey)
	var bdData emailBD
	b, err := json.Marshal(eventData.Shape.Data.BusinessData)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &bdData); err != nil {
		return nil, err
	}

	var recivers []client.Receivers
	for _, user := range bdData.ApprovePersons.Users {
		reciver := client.Receivers{
			Type: 1,
			ID:   user["id"].(string),
			Name: user["ownerName"].(string),
		}
		recivers = append(recivers, reciver)
	}
	types, err := strconv.Atoi(eventData.Shape.Data.BusinessData["sort"].(string))
	if len(recivers) > 0 {
		web := client.Web{
			IsSend: true,
			Contents: client.Contents{
				Content: bdData.Content,
			},
			Title:     utils.Strval(bdData.Title),
			Receivers: recivers,
			Types:     types,
		}
		msgReq := client.Mail{
			Web: web,
		}
		err = n.MessageCenterAPI.MessageCreateff(ctx, msgReq)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

// InitEnd event
func (n *Letter) InitEnd(ctx context.Context, eventData *EventData) (*pb.NodeEventRespData, error) {
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
	web := client.Web{
		IsSend: true,
		Title:  utils.Strval(bd["title"]),
		Contents: client.Contents{
			Content: utils.Strval(bd["content"]),
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
	return nil, nil
}
