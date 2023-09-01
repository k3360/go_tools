package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type RedisServer struct {
	Client *redis.Client
}

func (s *RedisServer) Connect() (*RedisServer, error) {
	redisAddr := fmt.Sprintf("%s:%d", host, port)
	// 创建Redis客户端
	s.Client = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: password, // 验证密码
		DB:       0,        // Redis槽
	})
	return s, nil
}

// 获取键值
func (s *RedisServer) GetValue(key string) (string, error) {
	return s.Client.Get(context.Background(), key).Result()
}

// 设置键值
func (s *RedisServer) SetValue(key string, value interface{}, expiration time.Duration) bool {
	err := s.Client.Set(context.Background(), key, value, expiration).Err()
	if err != nil {
		return false
	}
	return true
}

func (s *RedisServer) SetNXValue(key string, value interface{}, expiration time.Duration) bool {
	ok, err := s.Client.SetNX(context.Background(), key, value, expiration).Result()
	if err != nil {
		return false
	}
	return ok
}

// 删除键
func (s *RedisServer) Delete(key string) bool {
	_, err := s.Client.Del(context.Background(), key).Result()
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}
