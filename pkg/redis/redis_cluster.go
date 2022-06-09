package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	internal "github.com/quanxiang-cloud/flow/internal"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/misc/redis2"
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
	entityBytes, err := ClusterClient.Get(ctx, key).Bytes()
	if err != nil {
		return ""
	}

	return string(entityBytes)
}

// SetValueToRedis put value to redis
func SetValueToRedis(ctx context.Context, key string, value interface{}) {
	key = internal.RedisPreKey + key
	ClusterClient.Set(ctx, key, value, time.Minute*5)
}
