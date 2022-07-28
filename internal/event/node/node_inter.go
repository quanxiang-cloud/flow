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
	"github.com/quanxiang-cloud/flow/internal/convert"
	"github.com/quanxiang-cloud/flow/rpc/pb"
)

// INode inter
type INode interface {
	InitBegin(ctx context.Context, eventData *EventData) (*pb.NodeEventRespData, error)
	InitEnd(ctx context.Context, eventData *EventData) (*pb.NodeEventRespData, error)
}

// EventData req
type EventData struct {
	ProcessID         string
	ProcessInstanceID string
	NodeDefKey        string
	RequestID         string
	UserID            string
	ExecutionID       string
	TaskID            []string // 如果是会签则taskID是多个，逗号间隔的字符串
	Shape             *convert.ShapeModel
}
