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
