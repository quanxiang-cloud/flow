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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/quanxiang-cloud/flow/pkg"
	"github.com/quanxiang-cloud/flow/pkg/misc/client"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/quanxiang-cloud/flow/pkg/misc/logger"
	"github.com/quanxiang-cloud/flow/pkg/misc/resp"
)

// Config config
type Config struct {
	Timeout      time.Duration
	MaxIdleConns int
}

// HttpRetryConf retry config
type HttpRetryConf struct {
	Attempts int
	Sleep    time.Duration
	Func     func(ctx context.Context, client *http.Client, uri string, params interface{}, entity interface{}, headers map[string]string) error
}

// NewClient new a http client
func NewClient(conf client.Config) http.Client {
	return http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(conf.Timeout * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*conf.Timeout)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
			MaxIdleConns: conf.MaxIdleConns,
		},
	}
}

// MarshalInner not encode html
func MarshalInner(data interface{}) ([]byte, error) {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	if err := jsonEncoder.Encode(data); err != nil {
		return nil, err
	}
	return bf.Bytes(), nil
}

// POST http post
func POST(ctx context.Context, client *http.Client, uri string, params interface{}, entity interface{}) error {
	return POST2(ctx, client, uri, params, entity, nil)
	//return HttpRetry(ctx, client, uri, params, entity, nil, &HttpRetryConf{
	//	Attempts: 3,
	//	Sleep:    time.Second * 5,
	//	Func:     POST2,
	//})
}

// POST2 http post
func POST2(ctx context.Context, client *http.Client, uri string, params interface{}, entity interface{}, headers map[string]string) error {
	paramByte, err := MarshalInner(params)
	if err != nil {
		return err
	}

	fmt.Println("Request-Id:" + pkg.STDRequestID(ctx).String + " http request uri:" + uri + "      http request params:" + string(paramByte))
	reader := bytes.NewReader(paramByte)
	req, err := http.NewRequest("POST", uri, reader)
	if err != nil {
		logger.Logger.Errorw(err.Error(), pkg.STDRequestID(ctx))
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	// header 封装requestID,userID
	req.Header.Add(pkg.HeadRequestID, pkg.STDRequestID(ctx).String)
	req.Header.Add(pkg.HeadUserID, pkg.STDUserID(ctx))
	req.Header.Add(pkg.HeadUserName, pkg.STDUserName(ctx))
	req.Header.Add(pkg.GlobalXID, pkg.STDGlobalXID(ctx))
	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	response, err := client.Do(req)
	if err != nil {
		logger.Logger.Errorw(err.Error(), pkg.STDRequestID(ctx))
		return err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusOK {
		err = resp.DeserializationResp(ctx, response, entity)
		if err != nil {
			logger.Logger.Errorw(err.Error(), pkg.STDRequestID(ctx))
			return err
		}
	}
	return err
}

// Request http
func Request(ctx context.Context, client *http.Client, uri string, params interface{}, entity interface{}, headers map[string]string, method string) error {
	paramByte, err := json.Marshal(params)
	if err != nil {
		return err
	}

	fmt.Println("Request-Id:" + pkg.STDRequestID(ctx).String + " http request uri:" + uri + "      http request params:" + string(paramByte))
	reader := bytes.NewReader(paramByte)
	req, err := http.NewRequest(method, uri, reader)
	if err != nil {
		logger.Logger.Errorw(err.Error(), pkg.STDRequestID(ctx))
		return err
	}

	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	response, err := client.Do(req)
	if err != nil {
		logger.Logger.Errorw(err.Error(), pkg.STDRequestID(ctx))
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Logger.Errorw(err.Error(), pkg.STDRequestID(ctx))
		return err
	}
	logger.Logger.Info("Request-Id:" + pkg.STDRequestID(ctx).String + " resp:" + string(body))
	if response.StatusCode == http.StatusOK {
		if body == nil || string(body) == "null" || string(body) == "" {
			return nil
		}
		err = json.Unmarshal(body, entity)
		if err != nil {
			logger.Logger.Errorw(err.Error(), pkg.STDRequestID(ctx))
			return err
		}
	}
	// err = resp.DeserializationResp(ctx, response, entity)
	// if err != nil {
	// 	logger.Logger.Errorw(err.Error(), pkg.STDRequestID(ctx))
	// 	return err
	// }
	return nil
}

func HttpRetry(ctx context.Context, client *http.Client, uri string, params interface{}, entity interface{}, headers map[string]string, C *HttpRetryConf) error {
	err := C.Func(ctx, client, uri, params, entity, headers)
	if err != nil {
		if C.Attempts--; C.Attempts > 0 {
			logger.Logger.Errorf("http call error. retrying times: %d...\n", C.Attempts)
			// Add some randomness to prevent creating a Thundering Herd
			jitter := time.Duration(rand.Int63n(int64(C.Sleep)))
			C.Sleep = C.Sleep + jitter/2
			time.Sleep(C.Sleep)
			return HttpRetry(ctx, client, uri, params, headers, headers, C)
		}
		return err
	}
	return nil
}
