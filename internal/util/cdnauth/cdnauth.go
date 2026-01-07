package cdnauth

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

// GenerateGoEdgeSign 生成 GoEdge CDN 鉴权签名
//
// 参数：
//   - path: 原始路径（未编码，如 "/电影/poster.jpg"）
//   - privateKey: GoEdge 专用私钥
//   - randLength: 随机字符串长度（默认 16，设为 0 则使用 "0"）
//
// 返回：
//   - 完整的带签名 URL
//
// 签名逻辑：
//  1. 原串：path + "@" + ts + "@" + rand + "@" + privateKey
//  2. MD5：md5(原串) → 16进制小写
//  3. sign：ts + "-" + rand + "-" + md5_str
//  4. URL：url_encode(path) + "?sign=" + sign
func GenerateGoEdgeSign(path, privateKey string, randLength int) string {
	// 1. 生成时间戳
	ts := fmt.Sprintf("%d", time.Now().Unix())

	// 2. 生成随机字符串
	var randStr string
	if randLength <= 0 {
		randStr = "0"
	} else {
		randStr = generateRandomString(randLength)
	}

	// 3. 构造原始字符串（注意：这里使用的是原始未编码的 path）
	raw := path + "@" + ts + "@" + randStr + "@" + privateKey

	// 4. 计算 MD5（16进制小写）
	hash := md5.Sum([]byte(raw))
	md5Str := hex.EncodeToString(hash[:])

	// 5. 生成签名
	sign := ts + "-" + randStr + "-" + md5Str

	// 6. 对路径的每个部分进行编码，但保留路径分隔符 /
	encodedPath := encodePathSegments(path)
	return encodedPath + "?sign=" + sign
}

// GenerateTencentSign 生成腾讯云 CDN Type-A 鉴权签名
//
// 参数：
//   - path: 原始路径（未编码，如 "/电影/poster.jpg"）
//   - privateKey: 腾讯云专用私钥
//   - randLength: 随机字符串长度（默认 6，设为 0 则使用 "0"）
//   - uid: 用户 ID（默认 "0"）
//
// 返回：
//   - 完整的带签名 URL
//
// 签名逻辑：
//  1. uri = url_encode(path) - 编码每个段，保留 /
//  2. 原串：uri + "-" + ts + "-" + rand + "-" + uid + "-" + privateKey
//  3. MD5：md5(原串) → 16进制小写
//  4. sign：ts + "-" + rand + "-" + uid + "-" + md5_str
//  5. URL：uri + "?sign=" + sign
func GenerateTencentSign(path, privateKey, uid string, randLength int) string {
	// 1. URL 编码路径（腾讯云使用编码后的路径参与签名）
	// 对每个路径段编码，但保留 /
	uri := encodePathSegments(path)

	// 2. 生成时间戳
	ts := fmt.Sprintf("%d", time.Now().Unix())

	// 3. 生成随机字符串
	var randStr string
	if randLength <= 0 {
		randStr = "0"
	} else {
		randStr = generateRandomString(randLength)
	}

	// 4. 设置默认 uid
	if uid == "" {
		uid = "0"
	}

	// 5. 构造原始字符串（注意：这里使用的是编码后的 uri）
	raw := uri + "-" + ts + "-" + randStr + "-" + uid + "-" + privateKey

	// 6. 计算 MD5（16进制小写）
	hash := md5.Sum([]byte(raw))
	md5Str := hex.EncodeToString(hash[:])

	// 7. 生成签名
	sign := ts + "-" + randStr + "-" + uid + "-" + md5Str

	// 8. 返回带签名的路径
	return uri + "?sign=" + sign
}

// generateRandomString 生成指定长度的随机字符串
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

// encodePathSegments 对路径的每个部分进行 URL 编码，但保留路径分隔符 /
// 例如："/剧集/国产剧/test.mp4" -> "/%E5%89%A7%E9%9B%86/%E5%9B%BD%E4%BA%A7%E5%89%A7/test.mp4"
func encodePathSegments(path string) string {
	// 分割路径
	segments := strings.Split(path, "/")

	// 对每个部分进行编码
	for i, segment := range segments {
		if segment != "" {
			segments[i] = url.PathEscape(segment)
		}
	}

	// 重新拼接
	return strings.Join(segments, "/")
}
