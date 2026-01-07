package cdnauth

import (
	"testing"
)

// TestEncodePathSegments 测试路径段编码函数
func TestEncodePathSegments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "中文路径",
			input:    "/电影/国产剧/test.mp4",
			expected: "/%E7%94%B5%E5%BD%B1/%E5%9B%BD%E4%BA%A7%E5%89%A7/test.mp4",
		},
		{
			name:     "带空格路径",
			input:    "/剧集/人之初 (2025)/Season 01/test.mp4",
			expected: "/%E5%89%A7%E9%9B%86/%E4%BA%BA%E4%B9%8B%E5%88%9D%20%282025%29/Season%2001/test.mp4",
		},
		{
			name:     "特殊字符路径",
			input:    "/剧集/test {tmdbid=123}/file.mp4",
			expected: "/%E5%89%A7%E9%9B%86/test%20%7Btmdbid=123%7D/file.mp4",
		},
		{
			name:     "纯英文路径",
			input:    "/movies/action/test.mp4",
			expected: "/movies/action/test.mp4",
		},
		{
			name:     "根路径",
			input:    "/",
			expected: "/",
		},
		{
			name:     "无开头斜杠",
			input:    "movies/test.mp4",
			expected: "movies/test.mp4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := encodePathSegments(tt.input)
			if result != tt.expected {
				t.Errorf("encodePathSegments() 结果不符:\n输入: %s\n期望: %s\n实际: %s", tt.input, tt.expected, result)
			}
		})
	}
}

// TestRealWorldPath 测试真实世界的路径
func TestRealWorldPath(t *testing.T) {
	// 根据你的日志中的真实路径
	localPath := "/mnt/media/剧集/国产剧/人之初 (2025) {tmdbid=258492}/Season 01/人之初 S01E01 .mp4"
	localPrefix := "/mnt/media/剧集"
	remotePrefix := "/剧集"

	// 模拟路径映射过程
	relativePath := localPath[len(localPrefix):] // /国产剧/人之初 (2025) {tmdbid=258492}/Season 01/人之初 S01E01 .mp4
	cdnPath := remotePrefix + relativePath        // /剧集/国产剧/人之初 (2025) {tmdbid=258492}/Season 01/人之初 S01E01 .mp4

	// 编码
	encoded := encodePathSegments(cdnPath)

	t.Logf("本地路径: %s", localPath)
	t.Logf("相对路径: %s", relativePath)
	t.Logf("CDN路径:  %s", cdnPath)
	t.Logf("编码结果: %s", encoded)

	// 验证编码结果
	// 1. 应该以 / 开头，而不是 %2F
	if encoded[0:3] == "%2F" {
		t.Errorf("编码路径不应该以 %%2F 开头，应该以 / 开头")
	}

	if encoded[0] != '/' {
		t.Errorf("编码路径应该以 / 开头，实际: %c", encoded[0])
	}

	// 2. 应该包含中文的编码
	if !contains(encoded, "%E5%89%A7%E9%9B%86") { // "剧集"
		t.Errorf("应该包含中文编码")
	}

	// 3. 路径分隔符不应该被编码
	if contains(encoded, "%2F%2F") {
		t.Errorf("不应该出现连续的 %%2F")
	}
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
