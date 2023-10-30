package _json

import "encoding/json"

func ToJson(obj any) []byte {
	value, err := json.Marshal(obj)
	if err != nil {
		panic("Json解析为字符串异常：" + err.Error())
	}
	return value
}

func ToObject(value []byte, obj any) {
	err := json.Unmarshal(value, obj)
	if err != nil {
		panic("Json解析为对象异常：" + err.Error())
	}
}
