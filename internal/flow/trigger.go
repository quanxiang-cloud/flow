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

package flow

import (
	"github.com/quanxiang-cloud/flow/internal/server/options"
	"github.com/quanxiang-cloud/flow/pkg"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/misc/logger"
	"github.com/quanxiang-cloud/flow/pkg/redis"
	"github.com/quanxiang-cloud/flow/pkg/utils"
	"gorm.io/gorm"
)

// msg method
const (
	PostMethod   = "post"
	PutMethod    = "put"
	DeleteMethod = "delete"
)

// Trigger interface
type Trigger interface {
	MessageTrigger(formMsg *FormMsg) error
	// ApiTrigger(ctx context.Context, req *FormMsg) error
}

type trigger struct {
	db          *gorm.DB
	triggerRule TriggerRule
}

// NewTrigger new
func NewTrigger(conf *config.Configs, opts ...options.Options) (Trigger, error) {
	triggerRule, _ := NewTriggerRule(conf, opts...)
	tr := &trigger{
		triggerRule: triggerRule,
	}

	for _, opt := range opts {
		opt(tr)
	}
	return tr, nil
}

// SetDB set db
func (t *trigger) SetDB(db *gorm.DB) {
	t.db = db
}

// MessageTrigger trigger
func (t *trigger) MessageTrigger(formMsg *FormMsg) error {
	if formMsg.Entity != nil {
		entityMap := utils.ChangeObjectToMap(formMsg.Entity)
		userID := ""
		userName := ""
		if formMsg.Method == PostMethod {
			userID = entityMap["creator_id"].(string)
			if entityMap["creator_name"] != nil {
				userName = entityMap["creator_name"].(string)
			}
		} else if formMsg.Method == PutMethod {
			userID = entityMap["modifier_id"].(string)
			if entityMap["modifier_name"] != nil {
				userName = entityMap["modifier_name"].(string)
			}
		} else if formMsg.Method == DeleteMethod {
			userID = entityMap["delete_id"].(string)
		}
		ctx := pkg.CTXWrapper(formMsg.RequestID, userID, userName)
		// 过滤消息，看是否需要判断触发，用于节点上设置数据修改是否继续触发流程功能
		flag := true
		if len(formMsg.Seq) > 0 {
			seqResult := redis.GetStringValueFromRedis(ctx, formMsg.Seq)
			if len(seqResult) > 0 && "ignore" == seqResult {
				flag = false
			}
		}
		if flag {
			err := t.triggerRule.CheckDataModify(ctx, formMsg)
			if err != nil {
				logger.Logger.Error(err)
			}
			return err
		}
	}

	return nil
}

// ApiTrigger trigger
// func (t *trigger) ApiTrigger(ctx context.Context, req *FormMsg) error {
//
// }
