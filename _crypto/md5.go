package _crypto

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(value string) string {
	// 计算 MD5 哈希值
	hash := md5.Sum([]byte(value))
	// 将哈希值转换为十六进制字符串
	return hex.EncodeToString(hash[:])
}
