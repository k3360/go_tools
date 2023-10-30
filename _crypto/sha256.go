package _crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sha256(value string) string {
	// 创建一个SHA-256哈希对象
	hasher := sha256.New()
	// 写入要计算哈希的数据
	hasher.Write([]byte(value))
	// 计算SHA-256哈希值
	hashBytes := hasher.Sum(nil)
	// 将哈希值转换为十六进制字符串
	return hex.EncodeToString(hashBytes)
}
