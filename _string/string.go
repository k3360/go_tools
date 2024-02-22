package _string

import (
	"math/rand"
	"strings"
	"time"
)

const (
	MixedChars  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	UpperChars  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LowerChars  = "abcdefghijklmnopqrstuvwxyz"
	NumberChars = "0123456789"
)

func createRandString(length int, parseCode string) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune(parseCode)
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

// Random 随机混合字母
func Random(length int) string {
	return createRandString(length, MixedChars)
}

// RandomUpper 随机大写字母
func RandomUpper(length int) string {
	return createRandString(length, UpperChars)
}

// RandomLower 随机小写字母
func RandomLower(length int) string {
	return createRandString(length, LowerChars)
}

// RandomNumber 随机数字
func RandomNumber(length int) string {
	return createRandString(length, NumberChars)
}
