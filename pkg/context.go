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

package pkg

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	// HeadRequestID request id key
	HeadRequestID = "Request-Id"
	// HeadUserID user id key
	HeadUserID = "User-Id"
	// HeadUserName user name key
	HeadUserName = "User-Name"

	requestIDName = "requestID"
	userIDName    = "userID"
	userNameName  = "userName"

	// GlobalXID global transaction id
	GlobalXID = "xid"
	// LocalXID local transaction id
	LocalXID = "localXid"
)

type ctxValue string

// CTXWrapper2 wrapper
func CTXWrapper2(requestID string, userID string, userName string, globalXID string, localXID string) context.Context {
	v := values{map[string]string{
		requestIDName: requestID,
		userIDName:    userID,
		userNameName:  userName,
		GlobalXID:     globalXID,
		LocalXID:      localXID,
	}}

	return context.WithValue(context.Background(), ctxValue("ctxValue"), v)
}

// CTXWrapper wrwpper
func CTXWrapper(requestID string, userID string, userName string) context.Context {
	return CTXWrapper2(requestID, userID, userName, "", "")
}

// CTXTransfer transfer
func CTXTransfer(ctx *gin.Context) context.Context {
	requestID := ctx.Request.Header.Get(HeadRequestID)
	userID := ctx.Request.Header.Get(HeadUserID)
	userName := ctx.Request.Header.Get(HeadUserName)
	globalXID := ctx.Request.Header.Get(GlobalXID)

	v := values{map[string]string{
		requestIDName: requestID,
		userIDName:    userID,
		userNameName:  userName,
		GlobalXID:     globalXID,
	}}

	return context.WithValue(context.Background(), ctxValue("ctxValue"), v)
}

// RPCCTXTransfer func
func RPCCTXTransfer(requestID, userID string) context.Context {
	v := values{map[string]string{
		requestIDName: requestID,
		userIDName:    userID,
	}}
	return context.WithValue(context.Background(), ctxValue("ctxValue"), v)
}

// STDRequestID get request id from context.Context
func STDRequestID(ctx context.Context) zap.Field {
	if ctx == nil {
		return zap.String(requestIDName, "")
	}

	v := ctx.Value(ctxValue("ctxValue"))
	if v != nil {
		return zap.String(requestIDName, v.(values).get(requestIDName))
	}
	return zap.String(requestIDName, "")
}

// STDRequestID2 get request id from context
func STDRequestID2(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	v := ctx.Value(ctxValue("ctxValue"))
	if v != nil {
		return v.(values).get(requestIDName)
	}

	return ""
}

// STDUserID get user id from context.Context
func STDUserID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	v := ctx.Value(ctxValue("ctxValue"))
	if v != nil {
		return v.(values).get(userIDName)
	}
	return ""
}

// STDUserName get user name from context.Context
func STDUserName(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	v := ctx.Value(ctxValue("ctxValue"))
	if v != nil {
		return v.(values).get(userNameName)
	}
	return ""
}

// STDGlobalXID global transaction id
func STDGlobalXID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	v := ctx.Value(ctxValue("ctxValue"))
	if v != nil {
		return v.(values).get(GlobalXID)
	}
	return ""
}

// STDLocalXID local transaction id
func STDLocalXID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	v := ctx.Value(ctxValue("ctxValue"))
	if v != nil {
		return v.(values).get(LocalXID)
	}
	return ""
}

// Values values
type values struct {
	m map[string]string
}

func (v values) get(key string) string {
	return v.m[key]
}
