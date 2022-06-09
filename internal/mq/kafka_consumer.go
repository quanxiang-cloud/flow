package mq

import (
	"encoding/json"
	"fmt"
	"git.internal.yunify.com/qxp/misc/logger"
	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/flow/internal/flow"
	"github.com/quanxiang-cloud/flow/pkg"
	"github.com/quanxiang-cloud/flow/pkg/redis"
	"github.com/quanxiang-cloud/flow/pkg/utils"
	"io/ioutil"
	"net/http"
)

// msg method
const (
	PostMethod   = "post"
	PutMethod    = "put"
	DeleteMethod = "delete"
)

// DaprEvent DaprEvent
type DaprEvent struct {
	Topic           string        `json:"topic"`
	Pubsubname      string        `json:"pubsubname"`
	Traceid         string        `json:"traceid"`
	ID              string        `json:"id"`
	Datacontenttype string        `json:"datacontenttype"`
	Data            *flow.FormMsg `json:"data"`
	Type            string        `json:"type"`
	Specversion     string        `json:"specversion"`
	Source          string        `json:"source"`
}

// FormData FormData
type FormData struct {
	TableID string      `json:"tableID"`
	Entity  interface{} `json:"entity"`
	Magic   string      `json:"magic"`
	Seq     string      `json:"seq"`
	Version string      `json:"version"`
	Method  string      `json:"method"`
}

// Subscription Subscription
func Subscription(triggerRule flow.TriggerRule) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			errHandle(c, err)
			return
		}
		event := new(DaprEvent)
		err = json.Unmarshal(body, event)
		if err != nil {
			errHandle(c, err)
			return
		}
		formMsg := event.Data
		logger.Logger.Infow("formMsg", "data is", formMsg)
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
			ctx := pkg.CTXWrapper(formMsg.Seq, userID, userName)
			// 过滤消息，看是否需要判断触发，用于节点上设置数据修改是否继续触发流程功能
			flag := true
			if len(formMsg.Seq) > 0 {
				seqResult := redis.GetStringValueFromRedis(ctx, formMsg.Seq)
				if len(seqResult) > 0 && "ignore" == seqResult {
					flag = false
				}
			}
			if flag {
				err = triggerRule.CheckDataModify(ctx, formMsg)
				if err != nil {
					logger.Logger.Error(err)
				}
			}
		}
		errHandle(c, err)
	}
}

func errHandle(c *gin.Context, err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
	c.JSON(http.StatusOK, nil)
}
