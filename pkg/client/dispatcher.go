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

package client

import (
	"context"
	"fmt"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"net/http"
)

// Dispatcher client
type Dispatcher interface {
	TakePost(c context.Context, req TaskPostReq) error
	UpdateState(c context.Context, req UpdateTaskStateReq) error
}

type dispatcher struct {
	conf   *config.Configs
	client http.Client
}

// NewDispatcher new
func NewDispatcher(conf *config.Configs) Dispatcher {
	return &dispatcher{
		conf:   conf,
		client: NewClient(conf.InternalNet),
	}
}

func (d *dispatcher) TakePost(ctx context.Context, req TaskPostReq) error {
	var resp interface{}
	url := fmt.Sprintf("%s%s", d.conf.APIHost.DispatcherHost, "api/dispatcher/task/post")
	err := POST(ctx, &d.client, url, req, resp)
	return err
}

// TaskPostReq request params
type TaskPostReq struct {
	Code       string `json:"code"`
	Title      string `json:"title"`
	Describe   string `json:"describe"`
	Type       int    `json:"type"`
	TimeBar    string `json:"timeBar"`
	State      int    `json:"state"`
	Retry      int    `json:"retry"`
	RetryDelay int    `json:"retryDelay"`
}

// UpdateTaskStateReq 修改任务状态[参数]
type UpdateTaskStateReq struct {
	TaskID  string `json:"taskID"`
	Code    string `json:"code"`
	State   int    `json:"state"`
	TimeBar string `json:"timeBar"`
}

func (d *dispatcher) UpdateState(ctx context.Context, req UpdateTaskStateReq) error {
	var resp interface{}
	url := fmt.Sprintf("%s%s", d.conf.APIHost.DispatcherHost, "api/dispatcher/task/state/put")
	err := POST(ctx, &d.client, url, req, resp)
	return err
}
