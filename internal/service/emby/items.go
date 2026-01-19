package emby

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/AmbitiousJun/go-emby2openlist/v2/internal/config"
	"github.com/AmbitiousJun/go-emby2openlist/v2/internal/util/bytess"
	"github.com/AmbitiousJun/go-emby2openlist/v2/internal/util/https"
	"github.com/AmbitiousJun/go-emby2openlist/v2/internal/util/jsons"
	"github.com/AmbitiousJun/go-emby2openlist/v2/internal/util/logs"
	"github.com/AmbitiousJun/go-emby2openlist/v2/internal/util/urls"
	"github.com/AmbitiousJun/go-emby2openlist/v2/internal/web/cache"

	"github.com/gin-gonic/gin"
)

// ProxyAddItemsPreviewInfo 代理 Items 接口（转码功能已移除，仅处理 path 解码）
func ProxyAddItemsPreviewInfo(c *gin.Context) {
	// 代理请求
	c.Request.Header.Del("Accept-Encoding")
	resp, err := https.ProxyRequest(c.Request, config.C.Emby.Host)
	if checkErr(c, err) {
		return
	}
	defer resp.Body.Close()

	// 检查响应, 读取为 JSON
	if resp.StatusCode != http.StatusOK {
		checkErr(c, fmt.Errorf("emby 远程返回了错误的响应码: %d", resp.StatusCode))
		return
	}
	resJson, err := jsons.Read(resp.Body)
	if checkErr(c, err) {
		return
	}

	// 预响应请求
	defer func() {
		https.CloneHeader(c.Writer, resp.Header)
		jsons.OkResp(c.Writer, resJson)
		go runtime.GC()
	}()

	// 获取 Items 数组
	itemsArr, ok := resJson.Attr("Items").Done()
	if !ok || itemsArr.Empty() || itemsArr.Type() != jsons.JsonTypeArr {
		return
	}

	// 遍历每个 Item, 处理 MediaSource 信息
	itemsArr.RangeArr(func(index int, item *jsons.Item) error {
		mediaSources, ok := item.Attr("MediaSources").Done()
		if !ok || mediaSources.Empty() {
			return nil
		}

		mediaSources.RangeArr(func(_ int, ms *jsons.Item) error {
			simplifyMediaName(ms)

			// path 解码
			if path, ok := ms.Attr("Path").String(); ok {
				ms.Attr("Path").Set(urls.Unescape(path))
			}
			return nil
		})

		return nil
	})
}

// ProxyLatestItems 代理 Latest 请求
func ProxyLatestItems(c *gin.Context) {
	// 代理请求
	c.Request.Header.Del("Accept-Encoding")
	resp, err := https.ProxyRequest(c.Request, config.C.Emby.Host)
	if checkErr(c, err) {
		return
	}
	defer resp.Body.Close()

	// 检查响应, 读取为 JSON
	if resp.StatusCode != http.StatusOK {
		checkErr(c, fmt.Errorf("emby 远程返回了错误的响应码: %d", resp.StatusCode))
		return
	}
	resJson, err := jsons.Read(resp.Body)
	if checkErr(c, err) {
		return
	}

	// 预响应请求
	defer func() {
		https.CloneHeader(c.Writer, resp.Header)
		jsons.OkResp(c.Writer, resJson)
	}()

	// 遍历 MediaSources 解码 path
	if resJson.Type() != jsons.JsonTypeArr {
		return
	}
	resJson.RangeArr(func(_ int, item *jsons.Item) error {
		mediaSources, ok := item.Attr("MediaSources").Done()
		if !ok || mediaSources.Type() != jsons.JsonTypeArr || mediaSources.Empty() {
			return nil
		}
		mediaSources.RangeArr(func(_ int, ms *jsons.Item) error {
			if path, ok := ms.Attr("Path").String(); ok {
				ms.Attr("Path").Set(urls.Unescape(path))
			}
			return nil
		})
		return nil
	})

}
