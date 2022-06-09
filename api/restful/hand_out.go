package restful

import (
	"git.internal.yunify.com/qxp/misc/error2"
	"git.internal.yunify.com/qxp/misc/resp"
	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/flow/internal/flow"
	"github.com/quanxiang-cloud/flow/internal/server/options"
	"github.com/quanxiang-cloud/flow/pkg/config"
)

// HandOut info
type HandOut struct {
	urge flow.Urge
}

// NewHandOut new abnormalTask
func NewHandOut(c *config.Configs, opts ...options.Options) (*HandOut, error) {
	u, err := flow.NewUrge(c, opts...)
	if err != nil {
		return nil, err
	}

	return &HandOut{
		urge: u,
	}, nil
}

func (h *HandOut) handOut(ctx *gin.Context) {
	code, ok := ctx.Params.Get("code")
	if !ok {
		error2.NewErrorWithString(error2.Internal, "invalid URI")
		return
	}
	err := h.urge.UrgingExecute(ctx, code)
	resp.Format(nil, err).Context(ctx)
}
