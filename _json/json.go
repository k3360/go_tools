package _json

import (
	"encoding/json"
	"errors"
	"fmt"
)

func ToJson(obj any) []byte {
	value, err := json.Marshal(obj)
	if err != nil {
		fmt.Println("Json解析为字符串异常：" + err.Error())
		return []byte{}
	}
	return value
}

func ToObject(value []byte, obj any) error {
	err := json.Unmarshal(value, obj)
	if err != nil {
		return errors.New("Json解析为对象异常：" + err.Error())
	}
	return nil
}

func AnyToObject(value any, obj any) error {
	return ToObject(ToJson(value), obj)
}
