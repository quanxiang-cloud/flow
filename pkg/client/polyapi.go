package client

import (
	"context"
	"fmt"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"net/http"
)

// PolyAPI client
type PolyAPI interface {
	InnerRequest(c context.Context, apiPath string, req map[string]interface{}, header map[string]string, method string) (map[string]interface{}, error)
	SendRequest(c context.Context, url string, req map[string]interface{}, header map[string]string, method string) (map[string]interface{}, error)
}

type polyAPI struct {
	conf   *config.Configs
	client http.Client
}

// NewPolyAPI new
func NewPolyAPI(conf *config.Configs) PolyAPI {
	return &polyAPI{
		conf:   conf,
		client: NewClient(conf.InternalNet),
	}
}

func (p *polyAPI) InnerRequest(c context.Context, apiPath string, req map[string]interface{}, header map[string]string, method string) (map[string]interface{}, error) {
	var resp map[string]interface{}
	url := fmt.Sprintf("%s%s%s", p.conf.APIHost.PolyAPIHost, "api/v1/polyapi/inner/request", apiPath)
	err := Request(c, &p.client, url, req, &resp, header, method)
	return resp, err
}

func (p *polyAPI) SendRequest(c context.Context, url string, req map[string]interface{}, header map[string]string, method string) (map[string]interface{}, error) {
	var resp map[string]interface{}
	err := Request(c, &p.client, url, req, resp, header, method)
	return resp, err
}
