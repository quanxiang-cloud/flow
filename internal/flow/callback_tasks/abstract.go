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

package callback_tasks

import (
	"context"
	flow2 "github.com/quanxiang-cloud/flow/internal/flow"
	"github.com/quanxiang-cloud/flow/internal/models"
	"github.com/quanxiang-cloud/flow/pkg/client"
	"gorm.io/gorm"
)

// CallBackInterface 回调处理接口抽象
type CallBackInterface interface {
	Execute(ctx context.Context, code *string) error
}

// C c
type C struct {
	db                     *gorm.DB
	processAPI             client.Process
	formAPI                client.Form
	dispatcherCallbackRepo models.DispatcherCallbackRepo
	dispatcherAPI          client.Dispatcher
	urgeRepo               models.UrgeRepo
	flowRepo               models.FlowRepo
	instanceRepo           models.InstanceRepo
	task                   flow2.Task
	flow                   flow2.Flow
	instance               flow2.Instance
	operationRecord        flow2.OperationRecord
}
