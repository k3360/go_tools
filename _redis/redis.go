package _redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type RedisServer struct {
	host     string
	port     int
	password string
	slot     int
	*redis.Client
}

func New(host string, port int, password string, slot int) (*RedisServer, error) {
	server := RedisServer{host: host, port: port, password: password, slot: slot}
	return server.connect()
}

func (s *RedisServer) connect() (*RedisServer, error) {
	redisAddr := fmt.Sprintf("%s:%d", s.host, s.port)
	// 创建Redis客户端
	s.Client = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: s.password, // 验证密码
		DB:       s.slot,     // Redis槽
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

func (s *RedisServer) SCard(key string) int64 {
	result, err := s.Client.SCard(context.Background(), key).Result()
	if err != nil {
		log.Fatal(err)
		return 0
	}
	return result
}

func (s *RedisServer) SPopInt64(key string) int64 {
	i, err := s.Client.SPop(context.Background(), key).Int64()
	if err != nil {
		log.Fatal(err)
		return 0
	}
	return i
}

// 获取Set集合中的所有值，以字符串数组返回
func (s *RedisServer) SMembers(key string) []string {
	i, err := s.Client.SMembers(context.Background(), key).Result()
	if err != nil {
		log.Fatal(err)
		return []string{}
	}
	return i
}

// 向Set集合添加新值
func (s *RedisServer) SAdd(key string, members ...interface{}) bool {
	err := s.Client.SAdd(context.Background(), key, members...).Err()
	if err != nil {
		return false
	}
	return true
}

// 移除Set集合中一个或多个成员
func (s *RedisServer) SRem(key string, members ...interface{}) bool {
	err := s.Client.SRem(context.Background(), key, members...).Err()
	if err != nil {
		return false
	}
	return true
}

// 移除并返回Set集合中的一个随机元素
func (s *RedisServer) SPop(key string) string {
	val, err := s.Client.SPop(context.Background(), key).Result()
	if err != nil {
		return ""
	}
	return val
}

// 设置过期时间
func (s *RedisServer) Expire(key string, expiration time.Duration) bool {
	err := s.Client.Expire(context.Background(), key, expiration).Err()
	return err == nil
}

// 将一个或多个值插入到List列表头部
func (s *RedisServer) LPush(key string, values ...interface{}) bool {
	err := s.Client.LPush(context.Background(), key, values).Err()
	return err == nil
}

// 移除List列表的最后一个元素，返回值为移除的元素。
func (s *RedisServer) RPop(key string) string {
	val, err := s.Client.RPop(context.Background(), key).Result()
	if err != nil {
		return ""
	}
	return val
}

// 移除List列表元素
func (s *RedisServer) LRem(key string, count int64, value interface{}) bool {
	err := s.Client.LRem(context.Background(), key, count, value).Err()
	return err == nil
}
