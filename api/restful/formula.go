package restful

import (
	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/flow/internal/eval"
	"github.com/quanxiang-cloud/flow/internal/server/options"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/misc/logger"
	"github.com/quanxiang-cloud/flow/pkg/misc/resp"
	"net/http"
)

// Formula formula
type Formula struct {
	name string
}

// NewFormula 初始化
func NewFormula(c *config.Configs, opts ...options.Options) (*Formula, error) {
	return &Formula{
		name: "table formula",
	}, nil
}

// Calculation Calculation
func (f *Formula) Calculation(c *gin.Context) {
	req := &eval.FormulaReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(eval.Handler(logger.CTXTransfer(c), req)).Context(c)
}
