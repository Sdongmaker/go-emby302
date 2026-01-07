package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/AmbitiousJun/go-emby2openlist/v2/internal/util/logs"
	"github.com/AmbitiousJun/go-emby2openlist/v2/internal/util/maps"
	"github.com/AmbitiousJun/go-emby2openlist/v2/internal/util/strs"
)

// PeStrategy 代理异常策略类型
type PeStrategy string

const (
	PeStrategyOrigin PeStrategy = "origin" // 回源
	PeStrategyReject PeStrategy = "reject" // 拒绝请求
)

// DlStrategy 下载策略类型
type DlStrategy string

const (
	DlStrategyOrigin DlStrategy = "origin" // 代理到源服务器
	DlStrategyDirect DlStrategy = "direct" // 获取并重定向到直链
	DlStrategy403    DlStrategy = "403"    // 拒绝响应
)

// validPeStrategy 用于校验用户配置的策略是否合法
var validPeStrategy = map[PeStrategy]struct{}{
	PeStrategyOrigin: {}, PeStrategyReject: {},
}

// validDlStrategy 用于校验用户配置的下载策略是否合法
var validDlStrategy = map[DlStrategy]struct{}{
	DlStrategyOrigin: {}, DlStrategyDirect: {}, DlStrategy403: {},
}

// Emby 相关配置
type Emby struct {
	// Emby 源服务器地址
	Host string `yaml:"host"`
	// EpisodesUnplayPrior 在获取剧集列表时是否将未播资源优先展示
	EpisodesUnplayPrior bool `yaml:"episodes-unplay-prior"`
	// ResortRandomItems 是否对随机的 items 进行重排序
	ResortRandomItems bool `yaml:"resort-random-items"`
	// ProxyErrorStrategy 代理错误时的处理策略
	ProxyErrorStrategy PeStrategy `yaml:"proxy-error-strategy"`
	// ImagesQuality 图片质量
	ImagesQuality int `yaml:"images-quality"`
	// Strm strm 配置
	Strm *Strm `yaml:"strm"`
}

func (e *Emby) Init() error {
	if strs.AnyEmpty(e.Host) {
		return errors.New("emby.host 配置不能为空")
	}
	if strs.AnyEmpty(string(e.ProxyErrorStrategy)) {
		// 失败默认回源
		e.ProxyErrorStrategy = PeStrategyOrigin
	}

	e.ProxyErrorStrategy = PeStrategy(strings.TrimSpace(string(e.ProxyErrorStrategy)))
	if _, ok := validPeStrategy[e.ProxyErrorStrategy]; !ok {
		return fmt.Errorf("emby.proxy-error-strategy 配置错误, 有效值: %v", maps.Keys(validPeStrategy))
	}

	if e.ImagesQuality == 0 {
		// 不允许配置零值
		e.ImagesQuality = 70
	}
	if e.ImagesQuality < 0 || e.ImagesQuality > 100 {
		return fmt.Errorf("emby.images-quality 配置错误: %d, 允许配置范围: [1, 100]", e.ImagesQuality)
	}

	if e.Strm == nil {
		e.Strm = new(Strm)
	}
	if err := e.Strm.Init(); err != nil {
		return fmt.Errorf("emby.strm 配置错误: %v", err)
	}

	return nil
}

// PathMapping 路径映射配置
type PathMapping struct {
	// LocalPrefix 本地路径前缀
	LocalPrefix string `yaml:"local-prefix"`
	// CdnBase CDN 域名
	CdnBase string `yaml:"cdn-base"`
	// RemotePrefix CDN 上的路径前缀
	RemotePrefix string `yaml:"remote-prefix"`
}

// Strm strm 配置
type Strm struct {
	// PathMappings 路径映射配置列表
	PathMappings []PathMapping `yaml:"path-mappings"`
}

// Init 配置初始化
func (s *Strm) Init() error {
	if len(s.PathMappings) == 0 {
		return errors.New("strm.path-mappings 不能为空，至少需要配置一个映射规则")
	}

	for i, mapping := range s.PathMappings {
		if strs.AnyEmpty(mapping.LocalPrefix) {
			return fmt.Errorf("strm.path-mappings[%d].local-prefix 不能为空", i)
		}
		if strs.AnyEmpty(mapping.CdnBase) {
			return fmt.Errorf("strm.path-mappings[%d].cdn-base 不能为空", i)
		}
		if strs.AnyEmpty(mapping.RemotePrefix) {
			return fmt.Errorf("strm.path-mappings[%d].remote-prefix 不能为空", i)
		}

		// 标准化配置：确保 LocalPrefix 不以 / 结尾
		s.PathMappings[i].LocalPrefix = strings.TrimRight(mapping.LocalPrefix, "/")
		// 确保 CdnBase 不以 / 结尾
		s.PathMappings[i].CdnBase = strings.TrimRight(mapping.CdnBase, "/")
		// 确保 RemotePrefix 不以 / 结尾
		s.PathMappings[i].RemotePrefix = strings.TrimRight(mapping.RemotePrefix, "/")
	}

	return nil
}

// MapPath 将本地路径映射为 CDN 直链
// 示例：
//
//	输入: /mnt/media/剧集/国产剧/101次抢婚 (2023)/Season 1/101次抢婚 - S01E01 - 第 1 集.mp4
//	配置: local-prefix: /mnt/media/剧集, cdn-base: https://cdn.example.com, remote-prefix: /剧集
//	输出: https://cdn.example.com/剧集/国产剧/101次抢婚 (2023)/Season 1/101次抢婚 - S01E01 - 第 1 集.mp4
func (s *Strm) MapPath(localPath string) (string, error) {
	for _, mapping := range s.PathMappings {
		if strings.HasPrefix(localPath, mapping.LocalPrefix) {
			// 去掉本地前缀
			relativePath := strings.TrimPrefix(localPath, mapping.LocalPrefix)
			// 拼接 CDN 地址
			cdnUrl := mapping.CdnBase + mapping.RemotePrefix + relativePath
			logs.Info("路径映射: [%s] -> [%s]", localPath, cdnUrl)
			return cdnUrl, nil
		}
	}
	return "", fmt.Errorf("未找到匹配的路径映射规则: %s", localPath)
}
