package node

import (
	"context"
	"github.com/quanxiang-cloud/flow/pkg/client"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/utils"
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

// Init event
func (n *Letter) Init(ctx context.Context, eventData *EventData) error {
	return nil
}

// Execute event
func (n *Letter) Execute(ctx context.Context, eventData *EventData) error {
	instance, err := n.InstanceRepo.GetEntityByProcessInstanceID(n.Db, eventData.ProcessInstanceID)
	if err != nil {
		return err
	}

	bd := eventData.Shape.Data.BusinessData
	if bd == nil {
		return nil
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
		return err
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
		return err
	}
	return nil
}
