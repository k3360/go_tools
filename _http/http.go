package _http

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// BuffToRequest 通过buff，获取HTTP请求：http.Request
func BuffToRequest(buff []byte) (*http.Request, error) {
	// 创建字符串读取器
	newReader := strings.NewReader(string(buff))
	reader := bufio.NewReader(newReader)
	// 从字符串读取器中解析请求
	return http.ReadRequest(reader)
}

func Get(url string, param map[string]string) ([]byte, error) {
	uri := ""
	for key, val := range param {
		if len(uri) == 0 {
			uri += key + "=" + val
		} else {
			uri += "&" + key + "=" + val
		}
	}
	resp, err := http.Get(fmt.Sprintf("%s?%s", url, uri))
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return body, err2
	}
	return body, nil
}

//
//import (
//	"bytes"
//	"encoding/json"
//	"io/ioutil"
//	"net/_http"
//)
//
//// 定义一个接口
//type Any interface{}
//
//func postJson(url string) (string, error) {
//	//url := "http://localhost:8003/proxy/getProxyIPByUserPass"
//	//params := CommonParams{
//	//	Platform: "socks_proxy",
//	//	Version:  "1.0.0",
//	//	Params: UserPassParams{
//	//		Username: "k3360@qq.com",
//	//		Password: "King3360",
//	//	},
//	//}
//	// 将数据编码为JSON格式
//	jsonData, err := json.Marshal("params")
//	if err != nil {
//		return "", err
//	}
//	body := bytes.NewBuffer(jsonData)
//
//	resp, err := _http.Post(url, "application/json", body)
//	if err != nil {
//		return "", err
//	}
//	defer resp.Body.Close()
//
//	// 读取响应正文
//	content, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return "", err
//	}
//	return string(content), nil
//}
