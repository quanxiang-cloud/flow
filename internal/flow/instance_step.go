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

package flow

import (
	"context"
	"github.com/quanxiang-cloud/flow/internal/convert"
	"github.com/quanxiang-cloud/flow/internal/models"
	"github.com/quanxiang-cloud/flow/internal/models/mysql"
	"github.com/quanxiang-cloud/flow/internal/server/options"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"gorm.io/gorm"
)

// InstanceStep service
type InstanceStep interface {
	GetReStartFlowInstanceStep(ctx context.Context, processInstanceID string) (*models.InstanceStep, error)
	GetFlowInstanceSteps(ctx context.Context, processInstanceID string) ([]*models.InstanceStep, error)
	GetFlowInstanceStep(ctx context.Context, processInstanceID string, nodeInstanceID string) (*models.InstanceStep, error)
}

type instanceStep struct {
	db               *gorm.DB
	instanceStepRepo models.InstanceStepRepo
}

// NewInstanceStep init
func NewInstanceStep(conf *config.Configs, opts ...options.Options) (InstanceStep, error) {
	s := &instanceStep{
		instanceStepRepo: mysql.NewInstanceStepRepo(),
	}

	for _, opt := range opts {
		opt(s)
	}
	return s, nil
}

// SetDB set db
func (s *instanceStep) SetDB(db *gorm.DB) {
	s.db = db
}

func (s *instanceStep) GetReStartFlowInstanceStep(ctx context.Context, processInstanceID string) (*models.InstanceStep, error) {
	condition := &models.InstanceStep{
		ProcessInstanceID: processInstanceID,
		Status:            "",
		TaskType:          convert.Start,
	}
	steps, err := s.instanceStepRepo.FindInstanceSteps(s.db, condition)
	if err != nil {
		return nil, err
	}

	if len(steps) > 0 {
		return steps[0], nil
	}

	return nil, nil
}

func (s *instanceStep) GetFlowInstanceSteps(ctx context.Context, processInstanceID string) ([]*models.InstanceStep, error) {
	steps, err := s.instanceStepRepo.FindInstanceStepsByStatus(s.db, processInstanceID, []string{Review, InReview})
	if err != nil {
		return nil, err
	}

	return steps, nil
}

func (s *instanceStep) GetFlowInstanceStep(ctx context.Context, processInstanceID string, nodeInstanceID string) (*models.InstanceStep, error) {
	steps, err := s.instanceStepRepo.GetFlowInstanceStep(s.db, processInstanceID, nodeInstanceID, []string{Review, InReview})
	if err != nil {
		return nil, err
	}

	if len(steps) > 0 {
		return steps[0], nil
	}

	return nil, nil
}
