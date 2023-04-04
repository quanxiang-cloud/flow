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

package redis

import (
	"context"
	"fmt"
	internal "github.com/quanxiang-cloud/flow/internal"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/misc/redis2"
	"time"

	"github.com/go-redis/redis/v8"
)

// ClusterClient redis集群客户端
var ClusterClient *redis.ClusterClient

// Init 初始化
func Init() error {
	return newClient()
}

func newClient() error {
	client, err := redis2.NewClient(config.Config.Redis)
	if err != nil {
		return err
	}
	ClusterClient = client
	return nil
}

// Close 关闭
func Close() error {
	return ClusterClient.Close()
}

// GetStringValueFromRedis get string value from redis
func GetStringValueFromRedis(ctx context.Context, key string) string {
	key = internal.RedisPreKey + key
	entityBytes, err := ClusterClient.Get(ctx, key).Result()
	if err != nil {
		return ""
	}

	return entityBytes
}

// SetValueToRedis put value to redis
func SetValueToRedis(ctx context.Context, key string, value interface{}) {
	key = internal.RedisPreKey + key
	_, err := ClusterClient.SetEX(ctx, key, value, time.Minute*5).Result()
	if err != nil {
		fmt.Println("SetValueToRedis插入redis数据异常", key, "+++++++++", err.Error())
	}
	fmt.Println("SetValueToRedis插入redis数据ok", key, "+++++++++")
}
