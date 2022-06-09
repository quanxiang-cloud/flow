package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"git.internal.yunify.com/qxp/misc/client"
	"github.com/quanxiang-cloud/flow/pkg"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"git.internal.yunify.com/qxp/misc/logger"
	"git.internal.yunify.com/qxp/misc/resp"
)

// Config config
type Config struct {
	Timeout      time.Duration
	MaxIdleConns int
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

	err = resp.DeserializationResp(ctx, response, entity)
	if err != nil {
		logger.Logger.Errorw(err.Error(), pkg.STDRequestID(ctx))
		return err
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
	err = json.Unmarshal(body, entity)
	if err != nil {
		logger.Logger.Errorw(err.Error(), pkg.STDRequestID(ctx))
		return err
	}
	// err = resp.DeserializationResp(ctx, response, entity)
	// if err != nil {
	// 	logger.Logger.Errorw(err.Error(), pkg.STDRequestID(ctx))
	// 	return err
	// }
	return nil
}
