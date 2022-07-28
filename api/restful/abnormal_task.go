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

package restful

import (
	"github.com/quanxiang-cloud/flow/internal/flow"
	"github.com/quanxiang-cloud/flow/internal/models"
	"github.com/quanxiang-cloud/flow/internal/server/options"
	"github.com/quanxiang-cloud/flow/pkg"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/misc/error2"
	"github.com/quanxiang-cloud/flow/pkg/misc/resp"

	"github.com/gin-gonic/gin"
	"net/http"
)

const noStatus = -1

// AbnormalTask abnormal task
type AbnormalTask struct {
	abnormalTask flow.AbnormalTask
}

// NewAbnormalTask new abnormalTask
func NewAbnormalTask(c *config.Configs, opts ...options.Options) (*AbnormalTask, error) {
	at, err := flow.NewAbnormalTask(c, opts...)
	if err != nil {
		return nil, err
	}

	return &AbnormalTask{
		abnormalTask: at,
	}, nil
}

func (a *AbnormalTask) list(ctx *gin.Context) {
	req := &models.AbnormalTaskReq{
		Status: noStatus,
	}
	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(a.abnormalTask.List(pkg.CTXTransfer(ctx), req)).Context(ctx)
}

func (a *AbnormalTask) adminStepBack(ctx *gin.Context) {
	adminTaskReq, ok := a.getURIParams(ctx)
	if !ok {
		return
	}

	req := &models.HandleTaskModel{}
	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(a.abnormalTask.AdminStepBack(pkg.CTXTransfer(ctx), adminTaskReq, req)).Context(ctx)
}

func (a *AbnormalTask) adminSendBack(ctx *gin.Context) {
	adminTaskReq, ok := a.getURIParams(ctx)
	if !ok {
		return
	}
	resp.Format(a.abnormalTask.AdminSendBack(pkg.CTXTransfer(ctx), adminTaskReq)).Context(ctx)
}

func (a *AbnormalTask) adminAbandon(ctx *gin.Context) {
	adminTaskReq, ok := a.getURIParams(ctx)
	if !ok {
		return
	}

	resp.Format(a.abnormalTask.AdminAbandon(pkg.CTXTransfer(ctx), adminTaskReq)).Context(ctx)
}

func (a *AbnormalTask) adminDeliverTask(ctx *gin.Context) {
	adminTaskReq, ok := a.getURIParams(ctx)
	if !ok {
		return
	}

	req := &models.HandleTaskModel{}
	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(a.abnormalTask.AdminDeliverTask(pkg.CTXTransfer(ctx), adminTaskReq, req)).Context(ctx)
}

func (a *AbnormalTask) adminGetTaskForm(ctx *gin.Context) {
	adminTaskReq, ok := a.getURIParams(ctx)
	if !ok {
		return
	}

	resp.Format(a.abnormalTask.AdminGetTaskForm(pkg.CTXTransfer(ctx), adminTaskReq)).Context(ctx)
}

func (a *AbnormalTask) getURIParams(ctx *gin.Context) (*flow.AdminTaskReq, bool) {
	processInstanceID, ok1 := ctx.Params.Get("processInstanceID")
	taskID, ok2 := ctx.Params.Get("taskID")

	ok := ok1 && ok2
	if !ok {
		ctx.AbortWithError(http.StatusInternalServerError, error2.NewErrorWithString(error2.Internal, "invalid URI"))
	}

	return &flow.AdminTaskReq{
		ProcessInstanceID: processInstanceID,
		TaskID:            taskID,
	}, ok
}
