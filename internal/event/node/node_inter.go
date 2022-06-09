package node

import (
	"context"
	"github.com/quanxiang-cloud/flow/internal/convert"
)

// INode inter
type INode interface {
	Init(ctx context.Context, eventData *EventData) error
	Execute(ctx context.Context, eventData *EventData) error
}

// EventData req
type EventData struct {
	ProcessID         string
	ProcessInstanceID string
	NodeDefKey        string
	RequestID         string
	UserID            string
	TaskID            string // 如果是会签则taskID是多个，逗号间隔的字符串
	Shape             *convert.ShapeModel
}
