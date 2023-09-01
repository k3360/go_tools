package http

//
//import (
//	"bytes"
//	"encoding/json"
//	"io/ioutil"
//	"net/http"
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
//	resp, err := http.Post(url, "application/json", body)
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
