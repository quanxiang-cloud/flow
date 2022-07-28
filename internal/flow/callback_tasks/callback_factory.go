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

package callback_tasks

import (
	"errors"
	"github.com/quanxiang-cloud/flow/internal/convert"
	"github.com/quanxiang-cloud/flow/internal/server/options"
	"github.com/quanxiang-cloud/flow/pkg/config"
)

var (
	CallbackHandlers = map[string]CallBackInterface{}
	ErrCallbackType  = errors.New("no exist callback type")
)

type callbackFunc func(map[string]CallBackInterface)

// PackCallback 初始化回调对象
func PackCallback(cfs ...callbackFunc) {
	for _, cf := range cfs {
		cf(CallbackHandlers)
	}
}

// GetCallback 根据传入的类型返回回调对象
func GetCallback(Type string) (CallBackInterface, error) {
	if handler, ok := CallbackHandlers[Type]; ok {
		return handler, nil
	}
	return nil, ErrCallbackType
}

// InitCallbacks 初始化
func InitCallbacks(conf *config.Configs, opts ...options.Options) error {
	d, err := withDelay(conf, opts...)
	if err != nil {
		return err
	}

	u, err := withUrge(conf, opts...)
	if err != nil {
		return err
	}

	PackCallback(d, u)
	return nil
}

// withDelay 返回callback对象
func withDelay(conf *config.Configs, opts ...options.Options) (callbackFunc, error) {
	d, err := NewDelay(conf, opts...)
	if err != nil {
		return nil, err
	}
	return func(m map[string]CallBackInterface) {
		m[convert.CallbackOfDelay] = d
	}, nil
}

// withUrge 返回urge对象
func withUrge(conf *config.Configs, opts ...options.Options) (callbackFunc, error) {
	u, err := NewUrge(conf, opts...)
	if err != nil {
		return nil, err
	}
	return func(m map[string]CallBackInterface) {
		m[convert.CallbackOfUrge] = u

	}, nil
}
