package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/AmbitiousJun/go-emby2openlist/v2/internal/util/cdnauth"
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

// CdnAuthType CDN 鉴权类型
type CdnAuthType string

const (
	CdnAuthTypeNone     CdnAuthType = "none"     // 不使用鉴权
	CdnAuthTypeGoEdge   CdnAuthType = "goedge"   // GoEdge CDN 鉴权
	CdnAuthTypeTencent  CdnAuthType = "tencent"  // 腾讯云 CDN Type-A 鉴权
)

// validCdnAuthType 用于校验用户配置的鉴权类型是否合法
var validCdnAuthType = map[CdnAuthType]struct{}{
	CdnAuthTypeNone: {}, CdnAuthTypeGoEdge: {}, CdnAuthTypeTencent: {},
}

// PathMapping 路径映射配置
type PathMapping struct {
	// LocalPrefix 本地路径前缀
	LocalPrefix string `yaml:"local-prefix"`
	// RemotePrefix CDN 上的路径前缀
	RemotePrefix string `yaml:"remote-prefix"`
}

// CdnConfig CDN 配置
type CdnConfig struct {
	// Name CDN 名称（用于日志标识）
	Name string `yaml:"name"`
	// Type 鉴权类型 (none/goedge/tencent)
	Type CdnAuthType `yaml:"type"`
	// Base CDN 基础域名
	Base string `yaml:"base"`
	// PrivateKey 鉴权密钥
	PrivateKey string `yaml:"private-key"`
	// RandLength 随机字符串长度（GoEdge 默认 16, 腾讯云默认 6，设为 0 则使用 "0"）
	RandLength int `yaml:"rand-length"`
	// Uid 用户 ID（仅腾讯云使用，默认 "0"）
	Uid string `yaml:"uid"`
	// PathMappings 该 CDN 下的路径映射列表
	PathMappings []PathMapping `yaml:"path-mappings"`
}

// Strm strm 配置
type Strm struct {
	// Cdns CDN 配置列表
	Cdns []CdnConfig `yaml:"cdns"`
}

// Init 配置初始化
func (s *Strm) Init() error {
	if len(s.Cdns) == 0 {
		return errors.New("strm.cdns 不能为空，至少需要配置一个 CDN")
	}

	for ci, cdn := range s.Cdns {
		// 验证 CDN 名称
		if strs.AnyEmpty(cdn.Name) {
			return fmt.Errorf("strm.cdns[%d].name 不能为空", ci)
		}

		// 验证鉴权类型
		if strs.AnyEmpty(string(cdn.Type)) {
			s.Cdns[ci].Type = CdnAuthTypeNone
		}
		s.Cdns[ci].Type = CdnAuthType(strings.TrimSpace(string(cdn.Type)))
		if _, ok := validCdnAuthType[cdn.Type]; !ok {
			return fmt.Errorf("strm.cdns[%d].type 配置错误, 有效值: %v", ci, maps.Keys(validCdnAuthType))
		}

		// 验证 CDN 基础域名
		if strs.AnyEmpty(cdn.Base) {
			return fmt.Errorf("strm.cdns[%d].base 不能为空", ci)
		}
		s.Cdns[ci].Base = strings.TrimRight(cdn.Base, "/")

		// 如果启用鉴权，验证私钥
		if cdn.Type != CdnAuthTypeNone && strs.AnyEmpty(cdn.PrivateKey) {
			return fmt.Errorf("strm.cdns[%d].private-key 不能为空（鉴权类型: %s）", ci, cdn.Type)
		}

		// 设置默认随机字符串长度
		if cdn.RandLength < 0 {
			return fmt.Errorf("strm.cdns[%d].rand-length 不能为负数", ci)
		}
		if cdn.Type == CdnAuthTypeGoEdge && cdn.RandLength == 0 {
			s.Cdns[ci].RandLength = 16 // GoEdge 默认 16 位
		}
		if cdn.Type == CdnAuthTypeTencent && cdn.RandLength == 0 {
			s.Cdns[ci].RandLength = 6 // 腾讯云默认 6 位
		}

		// 腾讯云默认 uid
		if cdn.Type == CdnAuthTypeTencent && strs.AnyEmpty(cdn.Uid) {
			s.Cdns[ci].Uid = "0"
		}

		// 验证路径映射
		if len(cdn.PathMappings) == 0 {
			return fmt.Errorf("strm.cdns[%d].path-mappings 不能为空", ci)
		}

		for mi, mapping := range cdn.PathMappings {
			if strs.AnyEmpty(mapping.LocalPrefix) {
				return fmt.Errorf("strm.cdns[%d].path-mappings[%d].local-prefix 不能为空", ci, mi)
			}
			if strs.AnyEmpty(mapping.RemotePrefix) {
				return fmt.Errorf("strm.cdns[%d].path-mappings[%d].remote-prefix 不能为空", ci, mi)
			}

			// 标准化配置
			s.Cdns[ci].PathMappings[mi].LocalPrefix = strings.TrimRight(mapping.LocalPrefix, "/")
			s.Cdns[ci].PathMappings[mi].RemotePrefix = strings.TrimRight(mapping.RemotePrefix, "/")
		}
	}

	return nil
}

// MapPath 将本地路径映射为 CDN 直链（支持鉴权）
// 示例：
//
//	输入: /mnt/media/剧集/国产剧/101次抢婚 (2023)/Season 1/101次抢婚 - S01E01 - 第 1 集.mp4
//	配置:
//	  - cdn.base: https://cdn.example.com
//	  - path-mapping.local-prefix: /mnt/media/剧集
//	  - path-mapping.remote-prefix: /series
//	  - cdn.type: goedge
//	  - cdn.private-key: xxxx
//	输出: https://cdn.example.com/%E7%94%B5%E5%BD%B1/xxx.mp4?sign=1234567890-abc123-md5hash
func (s *Strm) MapPath(localPath string) (string, error) {
	// 遍历所有 CDN 配置
	for _, cdn := range s.Cdns {
		// 遍历该 CDN 下的所有路径映射
		for _, mapping := range cdn.PathMappings {
			// 检查路径是否匹配（需要严格匹配前缀）
			if !matchPathPrefix(localPath, mapping.LocalPrefix) {
				continue
			}

			// 去掉本地前缀，得到相对路径
			relativePath := strings.TrimPrefix(localPath, mapping.LocalPrefix)

			// 构造 CDN 路径（原始路径，未编码）
			cdnPath := mapping.RemotePrefix + relativePath

			// 根据鉴权类型生成最终 URL
			finalUrl, err := generateAuthUrl(cdn, cdnPath)
			if err != nil {
				return "", fmt.Errorf("生成鉴权 URL 失败: %v", err)
			}

			logs.Info("路径映射 [%s]: [%s] -> [%s]", cdn.Name, localPath, finalUrl)
			return finalUrl, nil
		}
	}

	return "", fmt.Errorf("未找到匹配的路径映射规则: %s", localPath)
}

// matchPathPrefix 严格匹配路径前缀
// 确保前缀后面紧跟 "/" 或者完全匹配，避免误匹配
// 例如：前缀 "/mnt/media" 应该匹配 "/mnt/media/file.mp4"
// 但不应该匹配 "/mnt/media2/file.mp4"
func matchPathPrefix(path, prefix string) bool {
	// 完全匹配
	if path == prefix {
		return true
	}

	// 前缀匹配且后面紧跟 /
	if strings.HasPrefix(path, prefix) {
		// 检查前缀后的下一个字符是否是 /
		if len(path) > len(prefix) && path[len(prefix)] == '/' {
			return true
		}
	}

	return false
}

// generateAuthUrl 根据 CDN 配置生成带鉴权的 URL
func generateAuthUrl(cdn CdnConfig, cdnPath string) (string, error) {
	switch cdn.Type {
	case CdnAuthTypeNone:
		// 无鉴权，直接拼接
		return cdn.Base + cdnPath, nil

	case CdnAuthTypeGoEdge:
		// GoEdge 鉴权
		signedPath := cdnauth.GenerateGoEdgeSign(cdnPath, cdn.PrivateKey, cdn.RandLength)
		return cdn.Base + signedPath, nil

	case CdnAuthTypeTencent:
		// 腾讯云鉴权
		signedPath := cdnauth.GenerateTencentSign(cdnPath, cdn.PrivateKey, cdn.Uid, cdn.RandLength)
		return cdn.Base + signedPath, nil

	default:
		return "", fmt.Errorf("不支持的鉴权类型: %s", cdn.Type)
	}
}
