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

package mq

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/flow/internal/flow"
	"github.com/quanxiang-cloud/flow/pkg/misc/logger"
	"io/ioutil"
	"net/http"
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
func Subscription(trigger flow.Trigger) gin.HandlerFunc {
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
		logger.Logger.Infow("formMsg", "receive form data: ", formMsg)
		err = trigger.MessageTrigger(formMsg)
		errHandle(c, err)
	}
}

func errHandle(c *gin.Context, err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
	c.JSON(http.StatusOK, nil)
}
