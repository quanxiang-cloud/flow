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

// MessageCenter message center client
type MessageCenter interface {
	MessageCreate(c context.Context, req MsgReq) error
	MessageCreateff(c context.Context, req Mail) error
}

type messageCenter struct {
	conf   *config.Configs
	client http.Client
}

// NewMessageCenter new
func NewMessageCenter(conf *config.Configs) MessageCenter {
	return &messageCenter{
		conf:   conf,
		client: NewClient(conf.InternalNet),
	}
}

func (mc *messageCenter) MessageCreate(c context.Context, req MsgReq) error {
	var resp interface{}
	url := fmt.Sprintf("%s%s", mc.conf.APIHost.MessageCenterHost, "api/v1/message/manager/create")
	err := POST(c, &mc.client, url, req, resp)
	return err
}

func (mc *messageCenter) MessageCreateff(c context.Context, req Mail) error {
	var resp interface{}
	url := fmt.Sprintf("%s%s", mc.conf.APIHost.MessageCenterHost, "api/v1/message/manager/create")
	err := POST(c, &mc.client, url, req, resp)
	return err
}

// SendMsgModel send msg model
type SendMsgModel struct {
	ID            string                   `json:"id"`
	TemplateID    string                   `json:"template_id"`
	Title         string                   `json:"title"`
	Args          []map[string]string      `json:"args"`
	Channel       string                   `json:"channel"` // 发送渠道 ，站内信： letter  邮件：email
	Type          int                      `json:"type"`    // 1. verifycode 2、not verifycode
	Sort          int                      `json:"sort"`    // 1. 系统消息 2、 通知通告
	IsSend        bool                     `json:"is_send"`
	Recivers      []map[string]interface{} `json:"recivers"`
	MesAttachment []map[string]string      `json:"mes_attachment"`
	Source        string                   `json:"source"`
}

// Email 邮件节点
type Email struct {
	To          []string                 `json:"to"`
	Contents    Contents                 `json:"contents"`
	Title       string                   `json:"title"`
	ContentType string                   `json:"content_type"`
	Files       []map[string]interface{} `json:"files"`
}

// MsgReq req
type MsgReq struct {
	Email Email `json:"email"`
}

// Mail 站内信
type Mail struct {
	Web Web `json:"web"`
}

// Web Web
type Web struct {
	ID        string      `json:"id"`     // 如果是新建的消息有id ，不是新建的有id
	IsSend    bool        `json:"isSend"` // 是否需要发送
	Title     string      `json:"title"`  // 标题
	Contents  Contents    `json:"contents"`
	Files     []Files     `json:"files"`
	Receivers []Receivers `json:"receivers"`
	Types     int         `json:"types"`
}

// Contents 消息内容
type Contents struct {
	Content     string            `json:"content"`
	TemplateID  string            `json:"templateID"`
	KeyAndValue map[string]string `json:"keyAndValue"`
}

// Files Files
type Files struct {
	FileName string `json:"fileName"` // 文件名字
	URL      string `json:"url"`      // path
}

// Receivers Receivers
type Receivers struct {
	Type int    `json:"type"` // 1、人员  2、部门
	ID   string `json:"id"`
	Name string `json:"name"` // 1. 系统消息 2、 通知通告
}
