package constant

const (
	CurrentVersion = "v2.4.1"
	RepoAddr       = "https://github.com/AmbitiousJun/go-emby2openlist"
)

const (
	Reg_Socket       = `(?i)^/.*(socket|embywebsocket)`
	Reg_PlaybackInfo = `(?i)^/.*items/.*/playbackinfo\??`

	Reg_PlayingStopped  = `(?i)^/.*sessions/playing/stopped`
	Reg_PlayingProgress = `(?i)^/.*sessions/playing/progress`

	Reg_UserItems        = `(?i)^/.*users/.*/items/\d+($|\?)`
	Reg_UserEpisodeItems = `(?i)^/.*users/.*/items\?.*includeitemtypes=(episode|movie)`
	Reg_UserPlayedItems  = `(?i)^/.*users/.*/playeditems/(\d+)($|\?|/.*)?`
	Reg_UserLatestItems  = `(?i)^/.*users/.*/items/latest($|\?)`
	Reg_ItemsCounts      = `(?i)^/.*items/counts($|\?)`

	Reg_VideoSubtitles = `(?i)^/.*videos/.*/subtitles`

	Reg_ResourceStream   = `(?i)^/.*(videos|audio)/.*/(stream|universal)(\.\w+)?\??`
	Reg_ResourceOriginal = `(?i)^/.*(videos|audio)/.*/original(\.\w+)?\??`

	Reg_ItemDownload     = `(?i)^/.*items/\d+/download($|\?)`
	Reg_ItemSyncDownload = `(?i)^/.*sync/jobitems/\d+/file($|\?)`

	Reg_Images             = `(?i)^/.*images`
	Reg_VideoModWebDefined = `(?i)^/web/modules/htmlvideoplayer/plugin.js`
	Reg_Proxy2Origin       = `^/$|(?i)^.*(/web|/users|/artists|/genres|/similar|/shows|/system|/remote|/scheduledtasks)`

	Reg_Root = `(?i)^/$`

	Reg_All = `.*`
)

const (
	RouteSubMatchGinKey = "routeSubMatches" // 路由匹配成功时, 会将匹配的正则结果存放到 Gin 上下文

	CommonDlUserAgent = "libmpv" // 通用的下载 UA
)
