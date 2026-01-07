package config

import (
	"testing"
)

// TestItemsCounts_Init 测试配置初始化
func TestItemsCounts_Init(t *testing.T) {
	tests := []struct {
		name      string
		config    ItemsCounts
		wantErr   bool
		expectSum int
	}{
		{
			name: "正常配置",
			config: ItemsCounts{
				Enable:       true,
				MovieCount:   1000,
				SeriesCount:  500,
				EpisodeCount: 5000,
			},
			wantErr:   false,
			expectSum: 6500,
		},
		{
			name: "未启用",
			config: ItemsCounts{
				Enable: false,
			},
			wantErr:   false,
			expectSum: 0,
		},
		{
			name: "负数电影数量",
			config: ItemsCounts{
				Enable:     true,
				MovieCount: -100,
			},
			wantErr: true,
		},
		{
			name: "负数剧集数量",
			config: ItemsCounts{
				Enable:      true,
				SeriesCount: -50,
			},
			wantErr: true,
		},
		{
			name: "自动计算总数",
			config: ItemsCounts{
				Enable:       true,
				MovieCount:   100,
				SeriesCount:  50,
				EpisodeCount: 500,
				ItemCount:    0, // 应自动计算
			},
			wantErr:   false,
			expectSum: 650,
		},
		{
			name: "手动设置总数",
			config: ItemsCounts{
				Enable:       true,
				MovieCount:   100,
				SeriesCount:  50,
				EpisodeCount: 500,
				ItemCount:    1000, // 手动设置
			},
			wantErr:   false,
			expectSum: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Init()

			// 检查错误
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 如果期望成功，检查总数
			if !tt.wantErr && tt.config.Enable {
				if tt.config.ItemCount != tt.expectSum {
					t.Errorf("ItemCount = %d, expect %d", tt.config.ItemCount, tt.expectSum)
				}
			}
		})
	}
}

// TestItemsCounts_ToJSON 测试 JSON 生成
func TestItemsCounts_ToJSON(t *testing.T) {
	ic := ItemsCounts{
		Enable:          true,
		MovieCount:      1000,
		SeriesCount:     500,
		EpisodeCount:    5000,
		GameCount:       0,
		ArtistCount:     0,
		ProgramCount:    0,
		GameSystemCount: 0,
		TrailerCount:    0,
		SongCount:       0,
		AlbumCount:      0,
		MusicVideoCount: 0,
		BoxSetCount:     0,
		BookCount:       0,
		ItemCount:       6500,
	}

	result := ic.ToJSON()

	// 验证字段存在
	expectedFields := []string{
		"MovieCount", "SeriesCount", "EpisodeCount",
		"GameCount", "ArtistCount", "ProgramCount",
		"GameSystemCount", "TrailerCount", "SongCount",
		"AlbumCount", "MusicVideoCount", "BoxSetCount",
		"BookCount", "ItemCount",
	}

	for _, field := range expectedFields {
		if _, ok := result[field]; !ok {
			t.Errorf("ToJSON() 缺少字段: %s", field)
		}
	}

	// 验证数值
	if result["MovieCount"] != 1000 {
		t.Errorf("MovieCount = %d, expect 1000", result["MovieCount"])
	}
	if result["SeriesCount"] != 500 {
		t.Errorf("SeriesCount = %d, expect 500", result["SeriesCount"])
	}
	if result["ItemCount"] != 6500 {
		t.Errorf("ItemCount = %d, expect 6500", result["ItemCount"])
	}

	t.Logf("ToJSON 结果: %+v", result)
}

// TestItemsCounts_AutoCalculation 测试自动计算
func TestItemsCounts_AutoCalculation(t *testing.T) {
	ic := ItemsCounts{
		Enable:          true,
		MovieCount:      100,
		SeriesCount:     50,
		EpisodeCount:    500,
		SongCount:       200,
		AlbumCount:      20,
		MusicVideoCount: 10,
		ItemCount:       0, // 自动计算
	}

	err := ic.Init()
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	expected := 100 + 50 + 500 + 200 + 20 + 10
	if ic.ItemCount != expected {
		t.Errorf("自动计算 ItemCount = %d, expect %d", ic.ItemCount, expected)
	}

	t.Logf("自动计算结果: %d", ic.ItemCount)
}
