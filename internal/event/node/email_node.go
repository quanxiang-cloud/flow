package node

import (
	"context"
	"github.com/quanxiang-cloud/flow/pkg/client"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/utils"
	"strings"
)

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

// Init event
func (n *Email) Init(ctx context.Context, eventData *EventData) error {
	return nil
}

// Execute event
func (n *Email) Execute(ctx context.Context, eventData *EventData) error {
	instance, err := n.InstanceRepo.GetEntityByProcessInstanceID(n.Db, eventData.ProcessInstanceID)
	if err != nil {
		return err
	}

	bd := eventData.Shape.Data.BusinessData
	if bd == nil {
		return nil
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

	// replace content
	content := utils.Strval(bd["content"])
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

	handleUsers := n.Flow.GetTaskHandleUsers2(ctx, bd["approvePersons"], instance)

	// gen req params
	mesAttachments := make([]map[string]interface{}, 0)
	if v := bd["mes_attachment"]; v != nil {
		arr := v.([]interface{})
		for _, e := range arr {
			tmp := e.(map[string]interface{})
			mesAttachment := make(map[string]interface{})
			mesAttachment["name"] = utils.Strval(tmp["file_name"])
			mesAttachment["path"] = utils.Strval(tmp["file_url"])
			mesAttachments = append(mesAttachments, mesAttachment)
		}
	}

	emailAddr := make([]string, 0)
	for _, user := range handleUsers {
		emailAddr = append(emailAddr, user.Email)
	}
	email := client.Email{
		To: emailAddr,
		Contents: client.Contents{
			Content: content,
		},
		Title: utils.Strval(bd["title"]),
		Files: mesAttachments,
	}
	msgReq := client.MsgReq{
		Email: email,
	}
	// post msg
	err = n.MessageCenterAPI.MessageCreate(ctx, msgReq)
	if err != nil {
		return err
	}
	return nil
}
