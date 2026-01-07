package emby

import (
	"net/http"

	"github.com/AmbitiousJun/go-emby2openlist/v2/internal/config"
	"github.com/AmbitiousJun/go-emby2openlist/v2/internal/util/logs"
	"github.com/gin-gonic/gin"
)

// HandleItemsCounts 处理 /Items/Counts 请求
//
// 功能说明：
//  1. 如果 items-counts.enable=false，回源透传
//  2. 如果请求包含 ParentId 参数，回源透传（查询特定媒体库的统计）
//  3. 如果请求不包含 ParentId 参数，返回自定义统计数据（全局统计）
//
// 使用场景：
//  - Emby 客户端启动时会请求 /Items/Counts 获取媒体库统计信息
//  - 用户可以通过配置文件自定义显示的数量，美化展示效果
func HandleItemsCounts(c *gin.Context) {
	// 1. 检查是否启用自定义统计
	if !config.C.ItemsCounts.Enable {
		logs.Info("Items/Counts: 未启用自定义统计，回源透传")
		ProxyOrigin(c)
		return
	}

	// 2. 检查是否有 ParentId 参数
	parentId := c.Query("ParentId")
	if parentId != "" {
		logs.Info("Items/Counts: 请求包含 ParentId=%s，回源透传", parentId)
		ProxyOrigin(c)
		return
	}

	// 3. 返回自定义统计数据
	logs.Info("Items/Counts: 返回自定义统计数据")

	// 生成 Emby 格式的 JSON 响应
	countsData := config.C.ItemsCounts.ToJSON()

	// 返回 JSON 响应
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(http.StatusOK, countsData)
}
