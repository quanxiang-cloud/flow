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
	"github.com/quanxiang-cloud/flow/internal/eval"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/utils"
	"net/http"
	"regexp"
)

// Structor structor client
type Structor interface {
	CalExpression(c context.Context, req map[string]interface{}) (interface{}, error)

	GetCalExpressionFields(c context.Context, expr string) ([]string, error)
}

type structor struct {
	conf   *config.Configs
	client http.Client
}

// NewStructor new
func NewStructor(conf *config.Configs) Structor {
	return &structor{
		conf:   conf,
		client: NewClient(conf.InternalNet),
	}
}

func (s *structor) CalExpression(c context.Context, req map[string]interface{}) (interface{}, error) {
	// resp := &calExpressionResp{}
	// url := fmt.Sprintf("%s%s", s.conf.APIHost.StructorHost, "api/v1/structor/formula/home/eval")
	// err := POST(c, &s.client, url, req, resp)
	// if err != nil {
	// 	return false, err
	// }
	//
	// return resp.Result, nil

	params := &eval.FormulaReq{
		Expression: req["expression"].(string),
		Parameter:  utils.ChangeObjectToMap(req["parameter"]),
	}
	return eval.Handler(c, params)
}

// GetCalExpressionFields split fields from cal expression
func (s *structor) GetCalExpressionFields(c context.Context, expr string) ([]string, error) {
	reg := regexp.MustCompile(`\$[a-zA-Z0-9_.\[\]]+`)
	params := reg.FindAllString(expr, -1)

	for _, param := range params {
		fmt.Println(param)
	}

	return params, nil
}

type calExpressionResp struct {
	Result interface{} `json:"result"`
}
