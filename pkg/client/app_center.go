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
	"github.com/quanxiang-cloud/flow/pkg/misc/logger"
	"net/http"
)

// AppCenter app center client
type AppCenter interface {
	BatchGetAppName(c context.Context, appIDs []string) (map[string]string, error)
	GetAppName(c context.Context, appID string) (string, error)
	GetAdminApps(c context.Context) ([]AppModel, error)
	GetAdminAppIDs(c context.Context) ([]string, error)
}

type appCenter struct {
	conf   *config.Configs
	client http.Client
}

// NewAppCenter new
func NewAppCenter(conf *config.Configs) AppCenter {
	return &appCenter{
		conf:   conf,
		client: NewClient(conf.InternalNet),
	}
}

func (ac *appCenter) GetAdminAppIDs(c context.Context) ([]string, error) {
	apps, err := ac.GetAdminApps(c)
	if err != nil {
		return nil, err
	}

	appIDs := make([]string, 0)
	for _, value := range apps {
		appIDs = append(appIDs, value.ID)
	}
	return appIDs, nil
}

func (ac *appCenter) GetAdminApps(c context.Context) ([]AppModel, error) {
	req := &appPageReq{
		Page:  1,
		Limit: 99999,
	}

	resp := &appPageResp{}
	url := fmt.Sprintf("%s%s", ac.conf.APIHost.AppCenterHost, "api/v1/app-center/adminList")
	err := POST(c, &ac.client, url, req, resp)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func (ac *appCenter) BatchGetAppName(c context.Context, appIDs []string) (map[string]string, error) {
	result := make(map[string]string, 0)

	req := map[string]interface{}{
		"ids": appIDs,
	}

	resp := &batchGetAppNameResp{}
	url := fmt.Sprintf("%s%s", ac.conf.APIHost.AppCenterHost, "api/v1/app-center/apps")
	err := POST(c, &ac.client, url, req, resp)
	if err != nil {
		return nil, err
	}

	if resp != nil && len(resp.Apps) > 0 {
		for _, value := range resp.Apps {
			result[value.ID] = value.AppName
		}
	} else {
		logger.Logger.Error("Failed to batch get app info ")
	}

	return result, nil
}

func (ac *appCenter) GetAppName(c context.Context, appID string) (string, error) {
	nameMap, err := ac.BatchGetAppName(c, []string{appID})
	if err != nil {
		return "", err
	}

	if nameMap != nil {
		return nameMap[appID], nil
	}

	return "", nil
}

type batchGetAppNameResp struct {
	Apps []AppModel `json:"apps"`
}

// AppModel app model
type AppModel struct {
	ID      string `json:"id"`
	AppName string `json:"appName"`
}

type appPageReq struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type appPageResp struct {
	TotalCount int        `json:"total_count"`
	Data       []AppModel `json:"data"`
}
