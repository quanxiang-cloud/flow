package svc

import (
	"github.com/quanxiang-cloud/flow/internal/callback"
)

// ServiceContext ServiceContext
type ServiceContext struct {
	Config   Config
	Callback callback.Callback
}

// NewServiceContext NewServiceContext
func NewServiceContext(c Config, callback callback.Callback) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		Callback: callback,
	}
}
