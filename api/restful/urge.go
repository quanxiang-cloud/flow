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
	"github.com/quanxiang-cloud/flow/internal/flow/callback_tasks"
	"github.com/quanxiang-cloud/flow/internal/server/options"
	"github.com/quanxiang-cloud/flow/pkg"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/misc/logger"
	"github.com/quanxiang-cloud/flow/pkg/misc/resp"
)

// Urge info
type Urge struct {
	urge *callback_tasks.Urge
}

// NewUrge new
func NewUrge(c *config.Configs, opts ...options.Options) (*Urge, error) {
	u, err := callback_tasks.NewUrge(c, opts...)
	if err != nil {
		return nil, err
	}
	return &Urge{
		urge: u,
	}, nil
}

func (u *Urge) taskUrge(ctx *gin.Context) {
	req := &callback_tasks.TaskUrgeModel{}

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
