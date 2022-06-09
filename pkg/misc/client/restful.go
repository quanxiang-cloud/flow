package client

import (
	"bytes"
	"context"
	"encoding/json"
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

// New new a http client
func New(conf Config) http.Client {
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

// POST http post
func POST(ctx context.Context, client *http.Client, uri string, params interface{}, entity interface{}) error {
	paramByte, err := json.Marshal(params)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(paramByte)
	req, err := http.NewRequest("POST", uri, reader)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	// header 封装request user-id user-name department-id role
	headersCTX := logger.STDHeader(ctx)
	for HKey, HValue := range headersCTX {
		req.Header.Add(HKey, HValue)
	}
	response, err := client.Do(req)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
		return err
	}
	defer response.Body.Close()

	err = resp.DeserializationResp(ctx, response, entity)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
		return err
	}
	return err
}
