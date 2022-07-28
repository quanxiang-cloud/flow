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
