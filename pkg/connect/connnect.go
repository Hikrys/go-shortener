package connect

import (
	"net/http"
	"time"
)

// client全局的HTTP客户端
var client = &http.Client{
	Transport: &http.Transport{
		DisableKeepAlives: true,
	},
	// 强制设置超时为 500ms 或 1s。如果1秒没通，直接视为无效或放弃校验
	Timeout: 500 * time.Millisecond,
}

// Get判断url是否能请求通常
func Get(url string) bool {
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	// 只要状态码是 200 就行，不读取 Body，节省内存和时间
	return resp.StatusCode == http.StatusOK //这里我们只判断直接链接，跳转链接我们也不管
}
