package restful

import (
	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/flow/internal/flow"
	"github.com/quanxiang-cloud/flow/internal/server/options"
	"github.com/quanxiang-cloud/flow/pkg"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/misc/logger"
	"github.com/quanxiang-cloud/flow/pkg/misc/resp"
)

// Urge info
type Urge struct {
	urge flow.Urge
}

// NewUrge new
func NewUrge(c *config.Configs, opts ...options.Options) (*Urge, error) {
	u, err := flow.NewUrge(c, opts...)
	if err != nil {
		return nil, err
	}

	return &Urge{
		urge: u,
	}, nil
}

func (u *Urge) taskUrge(ctx *gin.Context) {
	req := &flow.TaskUrgeModel{}

	err := ctx.ShouldBind(req)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}

	err = u.urge.TaskUrge(pkg.CTXTransfer(ctx), req)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(ctx)
		return
	}
	resp.Format(nil, nil).Context(ctx)
}
