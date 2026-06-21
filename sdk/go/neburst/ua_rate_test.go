package neburst

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

const (
	testEndpoint   = "http://api.neburst.local:8000"
	testCombinedKey = "eyJrZXlfaWQiOiJuYl9rZXlfZjY3MmIzMjM1MmU1ZWIxYzMxZTA0MzRiIiwic2VjcmV0IjoibmJfc2VjcmV0XzdlYmE4NDM0NTllNjc1MWZjNTIxZTFiN2ZjODAyNzZiIn0="
)

func doRawRequest(endpoint, keyID, secret, path, ua string) (int, http.Header, string, error) {
	c := NewClient(endpoint, keyID, secret)
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	nonce := uuidV4()
	bodyHash := sha256Hex(nil)
	sortedQ := ""
	stringToSign := timestamp + "\n" + "GET" + "\n" + path + "\n" + sortedQ + "\n" + bodyHash
	signature := hmacSHA256Hex(secret, stringToSign)

	req, err := http.NewRequest("GET", endpoint+path, nil)
	if err != nil {
		return 0, nil, "", err
	}
	req.Header.Set("User-Agent", ua)
	req.Header.Set("X-Nb-Key", keyID)
	req.Header.Set("X-Nb-Timestamp", timestamp)
	req.Header.Set("X-Nb-Nonce", nonce)
	req.Header.Set("X-Nb-Signature", signature)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, nil, "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, resp.Header, string(body), nil
}

func truncStr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func TestUARateLimit_SDKClient(t *testing.T) {
	keyID, secret := parseCombinedKey(testCombinedKey)
	ua := "Neburst OpenAPI SDK/1.0"
	path := "/open/v1/user/info"

	fmt.Printf("=== 官方 SDK UA: %s (期望 2 req/s) ===\n", ua)
	for i := 0; i < 5; i++ {
		code, h, body, err := doRawRequest(testEndpoint, keyID, secret, path, ua)
		if err != nil {
			t.Fatalf("请求失败: %v", err)
		}
		fmt.Printf("[%d] status=%d  UA-Limit=%s  UA-Remaining=%s  body=%s\n",
			i+1, code, h.Get("X-Ua-Ratelimit-Limit"), h.Get("X-Ua-Ratelimit-Remaining"), truncStr(body, 80))
		time.Sleep(100 * time.Millisecond)
	}
}

func TestUARateLimit_NonSDKClient(t *testing.T) {
	keyID, secret := parseCombinedKey(testCombinedKey)
	ua := "curl/7.88.1"
	path := "/open/v1/user/info"

	fmt.Printf("=== 非 SDK UA: %s (期望 1 req/5s) ===\n", ua)
	for i := 0; i < 4; i++ {
		code, h, body, err := doRawRequest(testEndpoint, keyID, secret, path, ua)
		if err != nil {
			t.Fatalf("请求失败: %v", err)
		}
		fmt.Printf("[%d] status=%d  UA-Limit=%s  UA-Remaining=%s  body=%s\n",
			i+1, code, h.Get("X-Ua-Ratelimit-Limit"), h.Get("X-Ua-Ratelimit-Remaining"), truncStr(body, 80))
		time.Sleep(500 * time.Millisecond)
	}
}

func TestUARateLimit_SDKRecovery(t *testing.T) {
	keyID, secret := parseCombinedKey(testCombinedKey)
	ua := "Neburst OpenAPI SDK/1.0"
	path := "/open/v1/user/info"

	fmt.Printf("=== SDK UA 限速恢复测试 (2 req/s, 窗口 1s) ===\n")
	for round := 1; round <= 3; round++ {
		fmt.Printf("\n--- 第 %d 轮 ---\n", round)
		for i := 0; i < 3; i++ {
			code, h, _, err := doRawRequest(testEndpoint, keyID, secret, path, ua)
			if err != nil {
				t.Fatalf("请求失败: %v", err)
			}
			fmt.Printf("[%d-%d] status=%d  Remaining=%s\n",
				round, i+1, code, h.Get("X-Ua-Ratelimit-Remaining"))
		}
		fmt.Printf("等待 1.1s 让窗口过期...\n")
		time.Sleep(1100 * time.Millisecond)
	}
}

func TestUARateLimit_NonSDKRecovery(t *testing.T) {
	keyID, secret := parseCombinedKey(testCombinedKey)
	ua := "python-requests/2.31"
	path := "/open/v1/user/info"

	fmt.Printf("=== 非 SDK UA 限速恢复测试 (1 req/5s, 窗口 5s) ===\n")
	for round := 1; round <= 3; round++ {
		fmt.Printf("\n--- 第 %d 轮 ---\n", round)
		for i := 0; i < 2; i++ {
			code, h, _, err := doRawRequest(testEndpoint, keyID, secret, path, ua)
			if err != nil {
				t.Fatalf("请求失败: %v", err)
			}
			fmt.Printf("[%d-%d] status=%d  Remaining=%s\n",
				round, i+1, code, h.Get("X-Ua-Ratelimit-Remaining"))
		}
		fmt.Printf("等待 5.1s 让窗口过期...\n")
		time.Sleep(5100 * time.Millisecond)
	}
}

func TestSDKClient_WithNewClient(t *testing.T) {
	c := NewClient(testEndpoint, testCombinedKey)
	fmt.Println("=== 通过 SDK NewClient 调用 (自动带 Neburst UA) ===")

	var info map[string]any
	err := c.doRequest("GET", "/open/v1/user/info", nil, nil, &info)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	fmt.Printf("用户信息: %v\n", info)
}
