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
	"github.com/quanxiang-cloud/flow/internal/convert"
	"github.com/quanxiang-cloud/flow/internal/flow"
	"github.com/quanxiang-cloud/flow/internal/flow/callback_tasks"
	"github.com/quanxiang-cloud/flow/internal/server/options"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/misc/error2"
	"github.com/quanxiang-cloud/flow/pkg/misc/logger"
	"github.com/quanxiang-cloud/flow/pkg/misc/resp"
	"strings"
)

// HandOut info
type HandOut struct {
	//urge *callback_tasks.Urge
	instance flow.Instance
}

// NewHandOut new abnormalTask
func NewHandOut(c *config.Configs, opts ...options.Options) (*HandOut, error) {
	// 初始化回调对象
	if err := callback_tasks.InitCallbacks(c, opts...); err != nil {
		return nil, err
	}
	//u, err := callback_tasks.NewUrge(c, opts...)
	//if err != nil {
	//	return nil, err
	//}
	instance, err := flow.NewInstance(c, opts...)
	if err != nil {
		return nil, err
	}
	return &HandOut{
		//urge: u,
		instance: instance,
	}, nil
}

func (h *HandOut) handOut(ctx *gin.Context) {
	var H callback_tasks.CallBackInterface
	code, ok := ctx.Params.Get("code")
	if !ok {
		error2.NewErrorWithString(error2.Internal, "invalid URI")
		return
	}
	//err := h.urge.UrgingExecute(ctx, code)
	// 根据code得到不同的handler处理业务
	s := strings.Split(code, "_")
	//split := strings.Split(s[1], "_")
	// 延时节点格式
	if len(s) == 2 && s[0] == convert.CallbackOfDelay {
		t := s[0]
		h, err := callback_tasks.GetCallback(t)
		if err != nil {
			resp.Format(nil, err).Context(ctx)
			return
		}
		H = h
		// 兼容之前
		// TODO 如果添加了其他类型，这里需要增加条件判断
	} else if len(s) == 2 && s[0] == convert.CallbackOfCron {
		// 启动流程

		startFlowModel := &flow.StartFlowModel{
			FlowID:   s[1],
			FormData: nil,
			UserID:   "",
		}
		_, err := h.instance.StartFlow(ctx, startFlowModel)
		if err != nil {
			logger.Logger.Error("enable flow err,", err)
			resp.Format(nil, err).Context(ctx)
			return
		}
		resp.Format(nil, nil).Context(ctx)
		return
		//h, err := callback_tasks.GetCallback(convert.CallbackOfDelay)
		//if err != nil {
		//	resp.Format(nil, err).Context(ctx)
		//	return
		//}
		//H = h
	} else {
		h, err := callback_tasks.GetCallback(convert.CallbackOfUrge)
		if err != nil {
			resp.Format(nil, err).Context(ctx)
			return
		}
		H = h
	}
	if err := H.Execute(ctx, &code); err != nil {
		resp.Format(nil, err).Context(ctx)
		return
	}
	resp.Format(nil, nil).Context(ctx)
}
