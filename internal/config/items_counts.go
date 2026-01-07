package config

import (
	"fmt"

	"github.com/AmbitiousJun/go-emby2openlist/v2/internal/util/logs"
)

// ItemsCounts Items/Counts 接口配置
type ItemsCounts struct {
	// Enable 是否启用自定义媒体库数量统计
	Enable bool `yaml:"enable"`
	// MovieCount 电影数量
	MovieCount int `yaml:"movie-count"`
	// SeriesCount 剧集数量
	SeriesCount int `yaml:"series-count"`
	// EpisodeCount 分集数量
	EpisodeCount int `yaml:"episode-count"`
	// GameCount 游戏数量
	GameCount int `yaml:"game-count"`
	// ArtistCount 艺术家数量
	ArtistCount int `yaml:"artist-count"`
	// ProgramCount 节目数量
	ProgramCount int `yaml:"program-count"`
	// GameSystemCount 游戏系统数量
	GameSystemCount int `yaml:"game-system-count"`
	// TrailerCount 预告片数量
	TrailerCount int `yaml:"trailer-count"`
	// SongCount 歌曲数量
	SongCount int `yaml:"song-count"`
	// AlbumCount 专辑数量
	AlbumCount int `yaml:"album-count"`
	// MusicVideoCount 音乐视频数量
	MusicVideoCount int `yaml:"music-video-count"`
	// BoxSetCount 合集数量
	BoxSetCount int `yaml:"box-set-count"`
	// BookCount 书籍数量
	BookCount int `yaml:"book-count"`
	// ItemCount 总项目数量（通常是所有类型的总和）
	ItemCount int `yaml:"item-count"`
}

// Init 配置初始化
func (ic *ItemsCounts) Init() error {
	// 输出启用状态
	if !ic.Enable {
		logs.Info("Items/Counts 自定义统计: 未启用，将回源透传")
		return nil
	}

	logs.Info("Items/Counts 自定义统计: 已启用")

	// 验证数量配置不能为负数
	if ic.MovieCount < 0 {
		return fmt.Errorf("items-counts.movie-count 不能为负数: %d", ic.MovieCount)
	}
	if ic.SeriesCount < 0 {
		return fmt.Errorf("items-counts.series-count 不能为负数: %d", ic.SeriesCount)
	}
	if ic.EpisodeCount < 0 {
		return fmt.Errorf("items-counts.episode-count 不能为负数: %d", ic.EpisodeCount)
	}
	if ic.ItemCount < 0 {
		return fmt.Errorf("items-counts.item-count 不能为负数: %d", ic.ItemCount)
	}

	// 输出主要媒体类型配置
	logs.Info("Items/Counts 主要配置: 电影=%d, 剧集=%d, 分集=%d",
		ic.MovieCount, ic.SeriesCount, ic.EpisodeCount)

	// 如果配置了其他类型的媒体，也输出日志
	if ic.SongCount > 0 || ic.AlbumCount > 0 || ic.ArtistCount > 0 {
		logs.Info("Items/Counts 音乐配置: 歌曲=%d, 专辑=%d, 艺术家=%d",
			ic.SongCount, ic.AlbumCount, ic.ArtistCount)
	}

	if ic.GameCount > 0 {
		logs.Info("Items/Counts 游戏配置: 游戏=%d, 游戏系统=%d",
			ic.GameCount, ic.GameSystemCount)
	}

	if ic.TrailerCount > 0 || ic.ProgramCount > 0 || ic.BookCount > 0 || ic.BoxSetCount > 0 || ic.MusicVideoCount > 0 {
		logs.Info("Items/Counts 其他配置: 预告片=%d, 节目=%d, 书籍=%d, 合集=%d, 音乐视频=%d",
			ic.TrailerCount, ic.ProgramCount, ic.BookCount, ic.BoxSetCount, ic.MusicVideoCount)
	}

	// 如果 ItemCount 为 0，自动计算为所有类型的总和
	if ic.ItemCount == 0 {
		ic.ItemCount = ic.MovieCount + ic.SeriesCount + ic.EpisodeCount +
			ic.GameCount + ic.ArtistCount + ic.ProgramCount +
			ic.GameSystemCount + ic.TrailerCount + ic.SongCount +
			ic.AlbumCount + ic.MusicVideoCount + ic.BoxSetCount +
			ic.BookCount
		logs.Info("Items/Counts 总数: 自动计算为 %d（所有类型总和）", ic.ItemCount)
	} else {
		logs.Info("Items/Counts 总数: 手动配置为 %d", ic.ItemCount)
	}

	return nil
}

// ToJSON 生成 Emby 格式的 JSON 响应
func (ic *ItemsCounts) ToJSON() map[string]int {
	return map[string]int{
		"MovieCount":      ic.MovieCount,
		"SeriesCount":     ic.SeriesCount,
		"EpisodeCount":    ic.EpisodeCount,
		"GameCount":       ic.GameCount,
		"ArtistCount":     ic.ArtistCount,
		"ProgramCount":    ic.ProgramCount,
		"GameSystemCount": ic.GameSystemCount,
		"TrailerCount":    ic.TrailerCount,
		"SongCount":       ic.SongCount,
		"AlbumCount":      ic.AlbumCount,
		"MusicVideoCount": ic.MusicVideoCount,
		"BoxSetCount":     ic.BoxSetCount,
		"BookCount":       ic.BookCount,
		"ItemCount":       ic.ItemCount,
	}
}
