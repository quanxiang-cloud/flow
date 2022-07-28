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
	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/flow/internal/flow"
	"github.com/quanxiang-cloud/flow/internal/models"
	"github.com/quanxiang-cloud/flow/internal/server/options"
	"github.com/quanxiang-cloud/flow/pkg"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/misc/logger"
	"github.com/quanxiang-cloud/flow/pkg/misc/resp"
)

// Instance info
type Instance struct {
	instance flow.Instance
}

// NewInstance new instance
func NewInstance(c *config.Configs, opts ...options.Options) (*Instance, error) {
	i, err := flow.NewInstance(c, opts...)
	if err != nil {
		return nil, err
	}

	return &Instance{
		instance: i,
	}, nil
}

// MyApplyList rest
func (i *Instance) myApplyList(ctx *gin.Context) {
	rq := &flow.MyApplyReq{}

	err := ctx.ShouldBind(rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}

	r, err := i.instance.MyApplyList(pkg.CTXTransfer(ctx), rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}
	resp.Format(r, nil).Context(ctx)
}

// WaitReviewList rest
func (i *Instance) waitReviewList(ctx *gin.Context) {
	rq := &flow.TaskListReq{}

	err := ctx.ShouldBind(rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}

	r, err := i.instance.WaitReviewList(pkg.CTXTransfer(ctx), rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}
	resp.Format(r, nil).Context(ctx)
}

// ReviewedList rest
func (i *Instance) reviewedList(ctx *gin.Context) {
	rq := &flow.TaskListReq{}

	err := ctx.ShouldBind(rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}

	r, err := i.instance.ReviewedList(pkg.CTXTransfer(ctx), rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}
	resp.Format(r, nil).Context(ctx)
}

// CcToMeList rest
func (i *Instance) ccToMeList(ctx *gin.Context) {
	rq := &flow.CcListReq{
		Status: -1,
	}

	err := ctx.ShouldBind(rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}

	r, err := i.instance.CcToMeList(pkg.CTXTransfer(ctx), rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}
	resp.Format(r, nil).Context(ctx)
}

// AllList rest
func (i *Instance) allList(ctx *gin.Context) {
	rq := &flow.TaskListReq{}

	err := ctx.ShouldBind(rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}

	r, err := i.instance.AllList(pkg.CTXTransfer(ctx), rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}
	resp.Format(r, nil).Context(ctx)
}

// Cancel rest
func (i *Instance) cancel(ctx *gin.Context) {
	processInstanceID := ctx.Param("processInstanceID")

	r, err := i.instance.Cancel(pkg.CTXTransfer(ctx), processInstanceID)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}
	resp.Format(r, nil).Context(ctx)
}

// Resubmit rest
func (i *Instance) resubmit(ctx *gin.Context) {
	processInstanceID := ctx.Param("processInstanceID")

	req := &flow.ResubmitReq{}
	err := ctx.ShouldBind(req)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}

	r, err := i.instance.Resubmit(pkg.CTXTransfer(ctx), processInstanceID, req)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}
	resp.Format(r, nil).Context(ctx)
}

// FlowInfo rest
func (i *Instance) flowInfo(ctx *gin.Context) {
	processInstanceID := ctx.Param("processInstanceID")

	r, err := i.instance.FlowInfo(pkg.CTXTransfer(ctx), processInstanceID)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}
	resp.Format(r, nil).Context(ctx)
}

// SendBack rest
func (i *Instance) sendBack(ctx *gin.Context) {
	processInstanceID := ctx.Param("processInstanceID")
	taskID := ctx.Param("taskID")

	req := &models.HandleTaskModel{}

	err := ctx.ShouldBind(req)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}

	r, err := i.instance.SendBack(pkg.CTXTransfer(ctx), processInstanceID, taskID, req)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}
	resp.Format(r, nil).Context(ctx)
}

// StepBack rest
func (i *Instance) stepBack(ctx *gin.Context) {
	processInstanceID := ctx.Param("processInstanceID")
	taskID := ctx.Param("taskID")

	req := &models.HandleTaskModel{}

	err := ctx.ShouldBind(req)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}

	r, err := i.instance.StepBack(pkg.CTXTransfer(ctx), processInstanceID, taskID, req)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}
	resp.Format(r, nil).Context(ctx)
}

// CcFlow rest
func (i *Instance) ccFlow(ctx *gin.Context) {
	processInstanceID := ctx.Param("processInstanceID")
	taskID := ctx.Param("taskID")

	req := &models.HandleTaskModel{}

	err := ctx.ShouldBind(req)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}

	r, err := i.instance.CcFlow(pkg.CTXTransfer(ctx), processInstanceID, taskID, req)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}
	resp.Format(r, nil).Context(ctx)
}

// ReadFlow rest
func (i *Instance) readFlow(ctx *gin.Context) {
	processInstanceID := ctx.Param("processInstanceID")
	taskID := ctx.Param("taskID")

	req := &models.HandleTaskModel{}

	err := ctx.ShouldBind(req)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}

	r, err := i.instance.ReadFlow(pkg.CTXTransfer(ctx), processInstanceID, taskID, req)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}
	resp.Format(r, nil).Context(ctx)
}

// HandleCc rest
func (i *Instance) handleCc(ctx *gin.Context) {
	var req []string

	err := ctx.ShouldBind(&req)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}

	r, err := i.instance.HandleCc(pkg.CTXTransfer(ctx), req)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}
	resp.Format(r, nil).Context(ctx)
}

// HandleRead rest
func (i *Instance) handleRead(ctx *gin.Context) {
	processInstanceID := ctx.Param("processInstanceID")
	taskID := ctx.Param("taskID")

	req := &models.HandleTaskModel{}

	err := ctx.ShouldBind(&req)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}

	r, err := i.instance.HandleRead(pkg.CTXTransfer(ctx), processInstanceID, taskID, req)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}
	resp.Format(r, nil).Context(ctx)
}

// AddSign rest
func (i *Instance) addSign(ctx *gin.Context) {
	processInstanceID := ctx.Param("processInstanceID")
	taskID := ctx.Param("taskID")

	req := &models.AddSignModel{}
	err := ctx.ShouldBind(req)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}

	r, err := i.instance.AddSign(pkg.CTXTransfer(ctx), processInstanceID, taskID, req)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}
	resp.Format(r, nil).Context(ctx)
}

// DeliverTask rest
func (i *Instance) deliverTask(ctx *gin.Context) {
	processInstanceID := ctx.Param("processInstanceID")
	taskID := ctx.Param("taskID")

	req := &models.HandleTaskModel{}

	err := ctx.ShouldBind(req)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}

	r, err := i.instance.DeliverTask(pkg.CTXTransfer(ctx), processInstanceID, taskID, req)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}
	resp.Format(r, nil).Context(ctx)
}

// StepBackNodes rest
func (i *Instance) stepBackNodes(ctx *gin.Context) {
	processInstanceID := ctx.Param("processInstanceID")

	r, err := i.instance.StepBackNodes(pkg.CTXTransfer(ctx), processInstanceID)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}
	resp.Format(r, nil).Context(ctx)
}

// FlowInstanceCount rest
func (i *Instance) flowInstanceCount(ctx *gin.Context) {
	r, err := i.instance.FlowInstanceCount(pkg.CTXTransfer(ctx))
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}
	resp.Format(r, nil).Context(ctx)
}

// ReviewTask rest
func (i *Instance) reviewTask(ctx *gin.Context) {
	processInstanceID := ctx.Param("processInstanceID")
	taskID := ctx.Param("taskID")

	req := &models.HandleTaskModel{}

	err := ctx.ShouldBind(req)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}

	r, err := i.instance.ReviewTask(pkg.CTXTransfer(ctx), processInstanceID, taskID, req)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}
	resp.Format(r, nil).Context(ctx)
}

// GetFlowInstanceForm rest
func (i *Instance) getFlowInstanceForm(ctx *gin.Context) {
	processInstanceID := ctx.Param("processInstanceID")

	req := &flow.TaskTypeDetailModel{}

	err := ctx.ShouldBind(req)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}

	r, err := i.instance.GetFlowInstanceForm(pkg.CTXTransfer(ctx), processInstanceID, req)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}
	resp.Format(r, nil).Context(ctx)
}

func (i *Instance) getFormData(ctx *gin.Context) {
	processInstanceID := ctx.Param("processInstanceID")
	taskID := ctx.Param("taskID")

	req := &flow.GetFormDataReq{}

	err := ctx.ShouldBind(req)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}

	r, err := i.instance.GetFormData(pkg.CTXTransfer(ctx), processInstanceID, taskID, req)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}
	resp.Format(r, nil).Context(ctx)
}

func (i *Instance) processHistories(ctx *gin.Context) {
	processInstanceID := ctx.Param("processInstanceID")

	r, err := i.instance.ProcessHistories(pkg.CTXTransfer(ctx), processInstanceID)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}
	resp.Format(r, nil).Context(ctx)
}

func (i *Instance) list(ctx *gin.Context) {
	req := &flow.ListReq{}

	err := ctx.ShouldBind(req)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}

	r, err := i.instance.List(pkg.CTXTransfer(ctx), req)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}
	resp.Format(r, nil).Context(ctx)
}

func (i *Instance) nodeList(ctx *gin.Context) {
	req := &flow.NodeListReq{}

	err := ctx.ShouldBind(req)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}

	r, err := i.instance.NodeList(pkg.CTXTransfer(ctx), req)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}
	resp.Format(r, nil).Context(ctx)
}
