package web

import (
	"github.com/AmbitiousJun/go-emby2openlist/v2/internal/constant"
	"github.com/AmbitiousJun/go-emby2openlist/v2/internal/service/emby"
	"github.com/AmbitiousJun/go-emby2openlist/v2/internal/util/logs"

	"github.com/gin-gonic/gin"
)

// rules 预定义路由拦截规则, 以及相应的处理器
//
// 每个规则为一个切片, 参数分别是: 正则表达式, 处理器
var rules [][2]any

func initRulePatterns() {
	logs.Info("正在初始化路由规则...")
	rules = compileRules([][2]any{
		// websocket
		{constant.Reg_Socket, emby.ProxySocket()},

		// PlaybackInfo 接口
		{constant.Reg_PlaybackInfo, emby.TransferPlaybackInfo},

		// 播放停止时, 辅助请求 Progress 记录进度
		{constant.Reg_PlayingStopped, emby.PlayingStoppedHelper},
		// 拦截无效的进度报告
		{constant.Reg_PlayingProgress, emby.PlayingProgressHelper},

		// 自定义媒体库数量统计
		{constant.Reg_ItemsCounts, emby.HandleItemsCounts},

		// Items 接口
		{constant.Reg_UserItems, emby.LoadCacheItems},
		// 代理 Items 并添加转码版本信息
		{constant.Reg_UserEpisodeItems, emby.ProxyAddItemsPreviewInfo},
		// 代理 Latest 接口, 解码媒体的 Path 字段
		{constant.Reg_UserLatestItems, emby.ProxyLatestItems},

		// 字幕长时间缓存
		{constant.Reg_VideoSubtitles, emby.ProxySubtitles},

		// 资源重定向到直链
		{constant.Reg_ResourceStream, emby.Redirect2OpenlistLink},
		// 处理 original 资源
		{constant.Reg_ResourceOriginal, emby.ProxyOriginalResource},

		// 资源下载, 重定向到直链
		{constant.Reg_ItemDownload, emby.Redirect2OpenlistLink},
		{constant.Reg_ItemSyncDownload, emby.HandleSyncDownload},

		// 处理图片请求
		{constant.Reg_Images, emby.HandleImages},

		// web cors 处理
		{constant.Reg_VideoModWebDefined, emby.ChangeBaseVideoModuleCorsDefined},

		// 根路径重定向到首页
		{constant.Reg_Root, emby.ProxyRoot},

		// 其余资源走重定向回源
		{constant.Reg_All, emby.ProxyOrigin},
	})
	logs.Success("路由规则初始化完成")
}

// initRoutes 初始化路由
func initRoutes(r *gin.Engine) {
	r.Any("/*vars", globalDftHandler)
}
