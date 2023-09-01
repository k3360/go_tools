package socks5

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"golang.org/x/net/proxy"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func postJson(url string, param Any) (string, error) {
	//url := "http://localhost:8003/proxy/getProxyIPByUserPass"
	//params := CommonParams{
	//	Platform: "socks_proxy",
	//	Version:  "1.0.0",
	//	Params: UserPassParams{
	//		Username: "k3360@qq.com",
	//		Password: "King3360",
	//	},
	//}
	// 将数据编码为JSON格式
	jsonData, err := json.Marshal(param)
	if err != nil {
		return "", err
	}
	body := bytes.NewBuffer(jsonData)

	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应正文
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func Start() {
	// 开启Socks服务
	address := fmt.Sprintf(":%d", 10088)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Println("Failed to start the server. PORT: ", address)
		return
	}
	defer listener.Close()

	log.Println("Successfully start the server. PORT: ", address)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func getUserPass(conn net.Conn) (username, password string) {
	ver := make([]byte, 1)
	if _, err := io.ReadFull(conn, ver); err != nil || 1 != int(ver[0]) {
		log.Println("账密认证版本有问题")
		return
	}
	len := make([]byte, 1)
	// 账号
	if _, err := io.ReadFull(conn, len); err != nil {
		log.Println("账号长度有问题")
		return
	}
	user := make([]byte, int(len[0]))
	if _, err := io.ReadFull(conn, user); err != nil {
		log.Println("账号读取有问题")
		return
	}
	// 密码
	if _, err := io.ReadFull(conn, len); err != nil {
		log.Println("密码长度有问题")
		return
	}
	pass := make([]byte, int(len[0]))
	if _, err := io.ReadFull(conn, pass); err != nil {
		log.Println("密码读取有问题")
		return
	}
	return string(user), string(pass)
}

var mVal = make(map[string]Response)

var mutex sync.Mutex

func getProxyIP(username, password string) (*Response, error) {
	// 并发控制
	mutex.Lock()
	defer mutex.Unlock()
	// 先从缓存取
	mapVal, ok := mVal[username+password]
	if ok {
		return &mapVal, nil
	}
	url := "https://www.918ip.com/apiv1/proxy-service/proxy/getProxyIPByUserPass"
	params := CommonParams{
		Platform: "socks_proxy",
		Version:  "1.0.0",
		Params: UserPassParams{
			Username: username,
			Password: password,
		},
	}
	content, err := postJson(url, params)
	if err != nil {
		return nil, err
	}

	// 将响应内容解码为JSON格式
	var response Response
	err1 := json.Unmarshal([]byte(content), &response)
	if err1 != nil {
		return nil, err1
	}
	// 报错缓存
	mVal[username+password] = response

	// 定期销毁
	go func() {
		//// 定时器
		//timer := time.NewTimer(15 * time.Second)
		//<-timer.C
		time.Sleep(time.Second * 15)
		delete(mVal, username+password)
	}()

	return &response, nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// 检查数据源
	buff := make([]byte, 2)
	if _, err := io.ReadFull(conn, buff); err != nil {
		log.Println("握手数据有问题:", buff)
		return
	}
	// 检查是否socks5协议
	if 5 != int(buff[0]) {
		log.Println("Socks5协议有问题: ", conn.RemoteAddr())
		return
	}
	// 获取认证方法
	methods := make([]byte, int(buff[1]))
	if _, err := io.ReadFull(conn, methods); err != nil {
		log.Println("methods数据有问题:", buff)
		return
	}

	// 暂时采用账号密码认证
	var isAuth bool
	for _, m := range methods {
		// 账密认证
		if m == 0x02 {
			isAuth = true
		}
	}
	if isAuth {
		// 要求客户端通过账号密码认证
		conn.Write([]byte{0x05, 0x02})
		// 获取认证账号密码
		username, password := getUserPass(conn)
		// 获取可代理的IP
		res, err := getProxyIP(username, password)
		if err != nil {
			// 认证失败
			conn.Write([]byte{0x01, 0x01})
			return
		}
		if res.Code == 1 || (res.Code == 0 && res.Data.IpAddress == "") {
			// 认证失败
			conn.Write([]byte{0x01, 0x01})
			return
		}
		// 认证成功
		conn.Write([]byte{0x01, 0x00})
		// 开始转发
		handleRequest(conn, res)
	} else {
		conn.Write([]byte{0x05, 0xFF})
	}
}

func handleRequest(conn net.Conn, res *Response) {
	/*
		The SOCKS request is formed as follows:
		+----+-----+-------+------+----------+----------+
		|VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
		+----+-----+-------+------+----------+----------+
		| 1  |  1  | X'00' |  1   | Variable |    2     |
		+----+-----+-------+------+----------+----------+
	*/
	header := make([]byte, 3)

	if _, err := io.ReadFull(conn, header); err != nil {
		//fmt.Println("认证后的数据有问题", header)
		return
	}
	//fmt.Println("认证后的数据：", header)

	if header[1] != 1 {
		//fmt.Println("其他协议", header)
		return
	}

	addrType := make([]byte, 1)
	io.ReadFull(conn, addrType)
	// 获取目标地址信息
	var host string
	switch addrType[0] {
	case addrTypeIPv41:
		ip := make([]byte, net.IPv4len)
		if _, err := io.ReadFull(conn, ip); err != nil {
			//return fmt.Errorf("failed to read IPv4 address: %v", err)
		}
		host = net.IP(ip).String()
	case addrTypeDomain1:
		domainLen := make([]byte, 1)
		if _, err := io.ReadFull(conn, domainLen); err != nil {
			//return fmt.Errorf("failed to read domain length: %v", err)
		}
		domain := make([]byte, domainLen[0])
		if _, err := io.ReadFull(conn, domain); err != nil {
			//return fmt.Errorf("failed to read domain: %v", err)
		}
		host = string(domain)
	case addrTypeIPv61:
		//return fmt.Errorf("IPv6 address type is not supported")
	default:
		//return fmt.Errorf("unsupported address type: %d", buf[3])
	}

	// 读取目标端口号
	portBuf := make([]byte, 2)
	if _, err := io.ReadFull(conn, portBuf); err != nil {
		//return fmt.Errorf("failed to read port: %v", err)
	}
	port := binary.BigEndian.Uint16(portBuf)

	// 向客户端发送连接响应
	reply := []byte{
		5,
		0,
		0,
		1,
	}
	localAddr := conn.LocalAddr().String()
	localHost, localPort, _ := net.SplitHostPort(localAddr)
	ipBytes := net.ParseIP(localHost).To4()
	nPort, _ := strconv.Atoi(localPort)
	reply = append(reply, ipBytes...)
	portBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(portBytes, uint16(nPort))
	reply = append(reply, portBytes...)
	conn.Write(reply)

	//fmt.Println("响应：", reply)

	//if _, err := conn.Write([]byte{0x05, 0x00, 0x00, 1}); err != nil { //, 0, 0, 0, 0, 0, 0
	//	fmt.Println("回复失败")
	//	//return fmt.Errorf("failed to send connection response: %v", err)
	//}
	fmt.Println("过了", host, port)

	//runtime.GC()

	//中转
	// 创建一个Dialer，使用代理服务器连接
	socksHost := net.JoinHostPort(res.Data.IpAddress, strconv.Itoa(res.Data.IpPort))
	dialer, err := proxy.SOCKS5("tcp", socksHost, &proxy.Auth{
		User:     res.Data.AuthUser,
		Password: res.Data.AuthPass,
	}, proxy.Direct)
	if err != nil {
		fmt.Println("无法连接到代理服务器:", err)
		return
	}
	destConn, err := dialer.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		fmt.Println("请求中转服务器失败:", err)
		return
	}
	defer destConn.Close()

	//// 文件描述符
	//var rLimit syscall.Rlimit
	//err1 := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	//if err1 != nil {
	//	fmt.Println("Error getting resource limit:", err1)
	//	return
	//}
	//fmt.Println("Current resource limit:", rLimit.Cur)
	//fmt.Println("Maximum resource limit:", rLimit.Max)

	//dialer := &net.Dialer{}
	//destConn, err := dialer.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	//
	//// 连接目标服务器
	////destConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	//if err != nil {
	//	fmt.Println("连接目标服务器失败")
	//	return
	//}
	//defer destConn.Close()
	////destConn.SetReadDeadline(time.Now().Add(5 * time.Second))
	////destConn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	//
	// 进行转发
	conChk := make(chan int)
	go CopyData(conn, destConn, conChk)
	go CopyData(destConn, conn, conChk)
	//go io.Copy(conn, destConn)
	//io.Copy(destConn, conn)
	<-conChk

	//Socks5Forward(conn, destConn)
}

func CopyData(dst io.WriteCloser, src io.Reader, conChk chan int) {
	// 将 IO 流反映到另一个
	defer onDisconnect(dst, conChk)
	written, _ := io.Copy(dst, src)
	log.Println("拷贝数据：", written)
	dst.Close()
	conChk <- 1
}

// onDisconnect 函数接收一个 io.WriteCloser 类型的写入对象 dst 和一个 chan int 类型的 conChk
// 该函数在 dst 连接关闭时被调用，并向 conChk 通道发送一个信号以表示连接已关闭
func onDisconnect(dst io.WriteCloser, conChk chan int) {
	// 关闭时 -> 强制断开另一对连接
	dst.Close()
	conChk <- 1
}

func Socks5Forward(client, target net.Conn) {
	forward := func(src, dest net.Conn) {
		defer src.Close()
		defer dest.Close()
		io.Copy(src, dest)
	}
	go forward(client, target)
	go forward(target, client)
}
