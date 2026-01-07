// redirect.go - Emby STRM 文件路径重定向模块
//
// 功能：拦截 Emby 的视频播放请求，将 STRM 文件中的本地路径映射为 CDN 直链，然后 302 重定向
//
// 流程：
// 1. 解析请求信息 (ItemId, ApiKey, MediaSourceId)
// 2. 从 Emby 获取 STRM 文件中的本地路径
// 3. 根据配置的路径映射规则转换为 CDN URL
// 4. 302 重定向到 CDN 直链
package emby

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/AmbitiousJun/go-emby2openlist/v2/internal/config"
	"github.com/AmbitiousJun/go-emby2openlist/v2/internal/util/https"
	"github.com/AmbitiousJun/go-emby2openlist/v2/internal/util/logs"
	"github.com/AmbitiousJun/go-emby2openlist/v2/internal/web/cache"

	"github.com/gin-gonic/gin"
)

// Redirect2OpenlistLink 重定向 STRM 文件到 CDN 直链
func Redirect2OpenlistLink(c *gin.Context) {
	// 不处理字幕接口
	if strings.Contains(strings.ToLower(c.Request.RequestURI), "subtitles") {
		ProxyOrigin(c)
		return
	}

	// 1 解析要请求的资源信息
	itemInfo, err := resolveItemInfo(c, RouteStream)
	if checkErr(c, err) {
		return
	}
	logs.Info("解析到的 itemInfo: %v", itemInfo)

	// 2 从 Emby 获取 STRM 文件中的本地路径
	localPath, err := getEmbyFileLocalPath(itemInfo)
	if checkErr(c, err) {
		return
	}
	logs.Info("STRM 文件路径: %s", localPath)

	// 3 将本地路径映射为 CDN 直链
	cdnUrl, err := config.C.Emby.Strm.MapPath(localPath)
	if checkErr(c, err) {
		return
	}

	// 4 返回 302 重定向
	logs.Success("302 重定向到: %s", cdnUrl)
	c.Header(cache.HeaderKeyExpired, cache.Duration(time.Minute*10))
	c.Redirect(http.StatusFound, cdnUrl)

	// 异步发送一个播放 Playback 请求, 触发 emby 解析 strm 视频格式
	go func() {
		originUrl, err := url.Parse(config.C.Emby.Host + itemInfo.PlaybackInfoUri)
		if err != nil {
			return
		}
		q := originUrl.Query()
		q.Set("IsPlayback", "true")
		q.Set("AutoOpenLiveStream", "true")
		originUrl.RawQuery = q.Encode()
		resp, err := https.Post(originUrl.String()).Body(io.NopCloser(bytes.NewBufferString(PlaybackCommonPayload))).Do()
		if err != nil {
			return
		}
		resp.Body.Close()
	}()
}

// ProxyOriginalResource 拦截 original 接口
func ProxyOriginalResource(c *gin.Context) {
	if strings.Contains(strings.ToLower(c.Request.RequestURI), "subtitles") {
		ProxyOrigin(c)
		return
	}
	Redirect2OpenlistLink(c)
}

// checkErr 检查 err 是否为空
// 不为空则根据错误处理策略返回响应
//
// 返回 true 表示请求已经被处理
func checkErr(c *gin.Context, err error) bool {
	if err == nil || c == nil {
		return false
	}

	// 异常接口, 不缓存
	c.Header(cache.HeaderKeyExpired, "-1")

	// 采用拒绝策略, 直接返回错误
	if config.C.Emby.ProxyErrorStrategy == config.PeStrategyReject {
		logs.Error("代理接口失败: %v", err)
		c.String(http.StatusInternalServerError, "代理接口失败, 请检查日志")
		return true
	}

	logs.Error("代理接口失败: %v, 回源处理", err)
	ProxyOrigin(c)
	return true
}
