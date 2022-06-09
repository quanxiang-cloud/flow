package restful

import (
	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/flow/internal/flow"
	"github.com/quanxiang-cloud/flow/internal/models"
	"github.com/quanxiang-cloud/flow/internal/server/options"
	"github.com/quanxiang-cloud/flow/pkg"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/misc/error2"
	"github.com/quanxiang-cloud/flow/pkg/misc/header2"
	"github.com/quanxiang-cloud/flow/pkg/misc/logger"
	"github.com/quanxiang-cloud/flow/pkg/misc/resp"
	"net/http"
)

// Flow info
type Flow struct {
	flow flow.Flow
}

// NewFlow new flow
func NewFlow(c *config.Configs, opts ...options.Options) (*Flow, error) {
	f, err := flow.NewFlow(c, opts...)
	if err != nil {
		return nil, err
	}

	return &Flow{
		flow: f,
	}, nil
}

// saveFlow save flow
func (f *Flow) saveFlow(ctx *gin.Context) {
	profile := header2.GetProfile(ctx)
	rq := &models.Flow{}
	err := ctx.ShouldBind(rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}
	r, err := f.flow.SaveFlow(pkg.CTXTransfer(ctx), rq, profile.UserID)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}
	resp.Format(r, nil).Context(ctx)
}

func (f *Flow) info(ctx *gin.Context) {
	ID := ctx.Param("ID")
	info, err := f.flow.Info(ID)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}
	resp.Format(info, nil).Context(ctx)
}

func (f *Flow) copyFlow(ctx *gin.Context) {
	ID := ctx.Param("ID")
	info, err := f.flow.CopyFlow(pkg.CTXTransfer(ctx), ID)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}
	resp.Format(info, nil).Context(ctx)
}

func (f *Flow) deleteFlow(ctx *gin.Context) {
	flowID, ok := ctx.Params.Get("id")
	if !ok {
		error2.NewErrorWithString(error2.Internal, "invalid URI")
		return
	}

	profile := header2.GetProfile(ctx)
	r, err := f.flow.DeleteFlow(pkg.CTXTransfer(ctx), flowID, profile.UserID)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}
	resp.Format(r, nil).Context(ctx)
}

func (f *Flow) flowList(ctx *gin.Context) {
	req := &flow.QueryFlowReq{}
	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(f.flow.FlowList(pkg.CTXTransfer(ctx), req)).Context(ctx)
}

func (f *Flow) correlationFlowList(ctx *gin.Context) {
	req := &flow.CorrelationFlowReq{}
	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(f.flow.CorrelationFlowList(pkg.CTXTransfer(ctx), req)).Context(ctx)
}

func (f *Flow) deleteApp(ctx *gin.Context) {
	req := &flow.DeleteAppReq{}

	err := ctx.ShouldBind(&req)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}

	result := f.flow.DeleteApp(pkg.CTXTransfer(ctx), req)
	resp.Format(result, nil).Context(ctx)
}

func (f *Flow) getVariableList(ctx *gin.Context) {
	ID, ok := ctx.GetQuery("id")
	if !ok {
		error2.NewErrorWithString(error2.Internal, "invalid URI")
		return
	}
	resp.Format(f.flow.GetVariableList(pkg.CTXTransfer(ctx), ID)).Context(ctx)
}

func (f *Flow) getNodes(ctx *gin.Context) {
	ID, ok := ctx.Params.Get("ID")
	if !ok {
		error2.NewErrorWithString(error2.Internal, "invalid URI")
		return
	}
	resp.Format(f.flow.GetNodes(pkg.CTXTransfer(ctx), ID)).Context(ctx)
}

func (f *Flow) saveFlowVariable(ctx *gin.Context) {
	req := &flow.SaveVariablesReq{}
	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	profile := header2.GetProfile(ctx)
	resp.Format(f.flow.SaveFlowVariable(pkg.CTXTransfer(ctx), req, profile.UserID)).Context(ctx)
}

func (f *Flow) deleteFlowVariable(ctx *gin.Context) {
	ID, ok := ctx.Params.Get("ID")
	if !ok {
		error2.NewErrorWithString(error2.Internal, "invalid URI")
		return
	}
	resp.Format(f.flow.DeleteFlowVariable(pkg.CTXTransfer(ctx), ID)).Context(ctx)
}

func (f *Flow) refreshRule(ctx *gin.Context) {
	resp.Format(f.flow.RefreshRule(pkg.CTXTransfer(ctx))).Context(ctx)
}

func (f *Flow) updateFlowStatus(ctx *gin.Context) {
	req := &flow.PublishProcessReq{}
	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	profile := header2.GetProfile(ctx)
	resp.Format(f.flow.UpdateFlowStatus(pkg.CTXTransfer(ctx), req, profile.UserID)).Context(ctx)
}

func (f *Flow) appReplicationExport(ctx *gin.Context) {
	req := &flow.AppReplicationExportReq{}
	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(f.flow.AppReplicationExport(pkg.CTXTransfer(ctx), req)).Context(ctx)
}

func (f *Flow) appReplicationImport(ctx *gin.Context) {
	req := &flow.AppReplicationImportReq{}
	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	profile := header2.GetProfile(ctx)
	resp.Format(f.flow.AppReplicationImport(pkg.CTXTransfer(ctx), req, profile.UserID)).Context(ctx)
}
