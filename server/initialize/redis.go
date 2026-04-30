package initialize

import (
	"context"
	"fmt"
	"go-admin/global"

	"github.com/redis/go-redis/v9"
)

func Redis() *redis.Client {
	cfg := global.Config.Redis

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 测试连接
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("连接 Redis 失败: %v", err))
	}

	return client
}
