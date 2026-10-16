package database

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

var RedisDB = InitRedis()

func InitRedis() *redis.Client {
	addr := "127.0.0.1:6379"
	client := redis.NewClient(
		&redis.Options{
			Addr: addr,
		})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	// 检查 Redis 是否可用
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		fmt.Errorf("Redis 连接失败: %w", err)
		return nil
	}
	log.Println("redis 连接成功")
	return client
}
