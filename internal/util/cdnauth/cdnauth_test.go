package cdnauth

import (
	"strings"
	"testing"
)

// TestGenerateGoEdgeSign 测试 GoEdge 签名生成
func TestGenerateGoEdgeSign(t *testing.T) {
	path := "/电影/poster.jpg"
	privateKey := "test_secret_key"
	randLength := 16

	result := GenerateGoEdgeSign(path, privateKey, randLength)

	// 检查结果格式
	if !strings.Contains(result, "?sign=") {
		t.Errorf("GoEdge 签名格式错误: 应包含 ?sign=")
	}

	// 检查是否进行了 URL 编码（中文部分）
	if !strings.Contains(result, "%E7%94%B5%E5%BD%B1") {
		t.Errorf("GoEdge 签名路径应该进行 URL 编码")
	}

	// 检查路径分隔符没有被编码
	if strings.HasPrefix(result, "%2F") {
		t.Errorf("GoEdge 签名路径不应该编码开头的 /，实际: %s", result)
	}

	// 检查签名格式: ts-rand-md5
	parts := strings.Split(result, "?sign=")
	if len(parts) != 2 {
		t.Fatalf("签名格式错误")
	}

	signParts := strings.Split(parts[1], "-")
	if len(signParts) != 3 {
		t.Errorf("GoEdge 签名应该是 ts-rand-md5 格式，实际: %s", parts[1])
	}

	// 检查 MD5 长度（32 位十六进制）
	md5Part := signParts[2]
	if len(md5Part) != 32 {
		t.Errorf("MD5 长度应该是 32，实际: %d", len(md5Part))
	}

	// 检查路径格式
	pathPart := parts[0]
	if !strings.HasPrefix(pathPart, "/") {
		t.Errorf("路径应该以 / 开头，实际: %s", pathPart)
	}

	t.Logf("GoEdge 签名结果: %s", result)
}

// TestGenerateTencentSign 测试腾讯云签名生成
func TestGenerateTencentSign(t *testing.T) {
	path := "/电影/poster.jpg"
	privateKey := "test_secret_key"
	uid := "0"
	randLength := 6

	result := GenerateTencentSign(path, privateKey, uid, randLength)

	// 检查结果格式
	if !strings.Contains(result, "?sign=") {
		t.Errorf("腾讯云签名格式错误: 应包含 ?sign=")
	}

	// 检查是否进行了 URL 编码（中文部分）
	if !strings.Contains(result, "%E7%94%B5%E5%BD%B1") {
		t.Errorf("腾讯云签名路径应该进行 URL 编码")
	}

	// 检查路径分隔符没有被编码
	if strings.HasPrefix(result, "%2F") {
		t.Errorf("腾讯云签名路径不应该编码开头的 /，实际: %s", result)
	}

	// 检查签名格式: ts-rand-uid-md5
	parts := strings.Split(result, "?sign=")
	if len(parts) != 2 {
		t.Fatalf("签名格式错误")
	}

	signParts := strings.Split(parts[1], "-")
	if len(signParts) != 4 {
		t.Errorf("腾讯云签名应该是 ts-rand-uid-md5 格式，实际: %s", parts[1])
	}

	// 检查 uid
	if signParts[2] != uid {
		t.Errorf("uid 不匹配，期望: %s, 实际: %s", uid, signParts[2])
	}

	// 检查 MD5 长度（32 位十六进制）
	md5Part := signParts[3]
	if len(md5Part) != 32 {
		t.Errorf("MD5 长度应该是 32，实际: %d", len(md5Part))
	}

	// 检查路径格式
	pathPart := parts[0]
	if !strings.HasPrefix(pathPart, "/") {
		t.Errorf("路径应该以 / 开头，实际: %s", pathPart)
	}

	t.Logf("腾讯云签名结果: %s", result)
}

// TestGenerateGoEdgeSign_WithZeroRand 测试随机字符串为 0 的情况
func TestGenerateGoEdgeSign_WithZeroRand(t *testing.T) {
	path := "/test.mp4"
	privateKey := "secret"
	randLength := 0

	result := GenerateGoEdgeSign(path, privateKey, randLength)

	// 检查随机字符串是否为 "0"
	parts := strings.Split(result, "?sign=")
	signParts := strings.Split(parts[1], "-")

	if signParts[1] != "0" {
		t.Errorf("当 randLength=0 时，随机字符串应该是 '0'，实际: %s", signParts[1])
	}

	t.Logf("零随机字符串结果: %s", result)
}

// TestGenerateRandomString 测试随机字符串生成
func TestGenerateRandomString(t *testing.T) {
	lengths := []int{6, 10, 16, 32}

	for _, length := range lengths {
		result := generateRandomString(length)

		if len(result) != length {
			t.Errorf("随机字符串长度不匹配，期望: %d, 实际: %d", length, len(result))
		}

		// 检查字符是否都是合法字符
		for _, ch := range result {
			if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9')) {
				t.Errorf("随机字符串包含非法字符: %c", ch)
			}
		}

		t.Logf("随机字符串 (长度 %d): %s", length, result)
	}
}

// BenchmarkGenerateGoEdgeSign GoEdge 签名性能测试
func BenchmarkGenerateGoEdgeSign(b *testing.B) {
	path := "/电影/国产剧/test.mp4"
	privateKey := "benchmark_secret"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateGoEdgeSign(path, privateKey, 16)
	}
}

// BenchmarkGenerateTencentSign 腾讯云签名性能测试
func BenchmarkGenerateTencentSign(b *testing.B) {
	path := "/电影/国产剧/test.mp4"
	privateKey := "benchmark_secret"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateTencentSign(path, privateKey, "0", 6)
	}
}
