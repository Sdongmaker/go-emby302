package https

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (

	// MaxRedirectDepth 重定向的最大深度
	MaxRedirectDepth = 10
)

var client *http.Client

// RedirectCodes 有重定向含义的 http 响应码
var RedirectCodes = [4]int{http.StatusMovedPermanently, http.StatusFound, http.StatusTemporaryRedirect, http.StatusPermanentRedirect}

func init() {
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			// 建立连接 1 分钟超时
			Dial: (&net.Dialer{Timeout: time.Minute}).Dial,
			// 接收数据 5 分钟超时
			ResponseHeaderTimeout: time.Minute * 5,
			// 连接池配置
			MaxIdleConns:        100,              // 最大空闲连接数
			MaxIdleConnsPerHost: 20,               // 每个 host 最大空闲连接数
			MaxConnsPerHost:     0,                // 每个 host 最大连接数（0 表示不限制）
			IdleConnTimeout:     90 * time.Second, // 空闲连接超时时间
			// 禁用压缩（避免代理时的编码问题）
			DisableCompression: false,
			// 禁用 keep-alive 超时
			DisableKeepAlives: false,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		// 总超时时间（包含连接建立、请求发送、响应读取）
		Timeout: 10 * time.Minute,
	}
}

// IsRedirectCode 判断 http code 是否是重定向
//
// 301, 302, 307, 308
func IsRedirectCode(code int) bool {
	for _, valid := range RedirectCodes {
		if code == valid {
			return true
		}
	}
	return false
}

// IsSuccessCode 判断 http code 是否为成功状态
func IsSuccessCode(code int) bool {
	codeStr := strconv.Itoa(code)
	return strings.HasPrefix(codeStr, "2")
}

// IsErrorCode 判断 http code 是否为错误状态
func IsErrorCode(code int) bool {
	codeStr := strconv.Itoa(code)
	return strings.HasPrefix(codeStr, "4") || strings.HasPrefix(codeStr, "5")
}

// MapBody 将 map 转换为 ReadCloser 流
func MapBody(body map[string]any) io.ReadCloser {
	if body == nil {
		return nil
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		log.Printf("MapBody 转换失败, body: %v, err : %v", body, err)
		return nil
	}
	return io.NopCloser(bytes.NewBuffer(bodyBytes))
}
