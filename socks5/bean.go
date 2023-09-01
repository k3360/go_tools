package socks5

const (
	addrTypeIPv41   = 1
	addrTypeIPv61   = 4
	addrTypeDomain1 = 3
)

// 定义一个接口
type Any interface{}

type UserPassParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CommonParams struct {
	Platform string         `json:"platform"`
	Version  string         `json:"version"`
	Params   UserPassParams `json:"params"`
}

type ProxyIpBean struct {
	IpAddress string `json:"ipAddress"`
	IpPort    int    `json:"ipPort"`
	AuthUser  string `json:"authUser"`
	AuthPass  string `json:"authPass"`
}
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data ProxyIpBean `json:"data"`
}
