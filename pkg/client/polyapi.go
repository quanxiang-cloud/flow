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
