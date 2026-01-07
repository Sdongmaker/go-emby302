// redirect.go - Emby 媒体资源重定向处理模块
//
// 功能：拦截 Emby 的视频播放请求，将 STRM 文件中的 URL 通过路径映射后 302 重定向到 CDN 直链
//
// 流程：
// 1. 解析请求信息 (ItemId, ApiKey, MediaSourceId)
// 2. 从 Emby 获取 STRM 文件中的 URL
// 3. 应用路径映射规则转换 URL
// 4. 302 重定向到最终的 CDN 直链
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
	"github.com/AmbitiousJun/go-emby2openlist/v2/internal/util/trys"
	"github.com/AmbitiousJun/go-emby2openlist/v2/internal/web/cache"

	"github.com/gin-gonic/gin"
)

// Redirect2OpenlistLink 重定向资源到直链
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

	// 2 请求资源在 Emby 中的 Path 参数
	embyPath, err := getEmbyFileLocalPath(itemInfo)
	if checkErr(c, err) {
		return
	}

	// 3 STRM 文件处理：应用路径映射并重定向
	finalPath := config.C.Emby.Strm.MapPath(embyPath)
	finalPath = getFinalRedirectLink(finalPath, c.Request.Header.Clone())
	logs.Success("重定向 strm: %s", finalPath)
	c.Header(cache.HeaderKeyExpired, cache.Duration(time.Minute*10))
	c.Redirect(http.StatusTemporaryRedirect, finalPath)

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

// getFinalRedirectLink 尝试对带有重定向的原始链接进行内部请求, 返回最终链接
//
// 检测到 internal-redirect-enable 配置未启用时, 直接返回原始链接
//
// 请求中途出现任何失败都会返回原始链接
func getFinalRedirectLink(originLink string, header http.Header) string {

	if !config.C.Emby.Strm.InternalRedirectEnable {
		logs.Info("internal-redirect-enable 未启用, 使用原始链接")
		return originLink
	}

	var finalLink string
	err := trys.Try(func() (err error) {
		logs.Info("正在尝试内部重定向, originLink: [%s]", originLink)
		fl, resp, e := https.Get(originLink).Header(header).DoRedirect()
		if e != nil {
			return e
		}
		defer resp.Body.Close()
		finalLink = fl
		return nil
	}, 3, time.Second*2)

	if err != nil {
		logs.Warn("内部重定向失败: %v", err)
		return originLink
	}

	return finalLink
}
