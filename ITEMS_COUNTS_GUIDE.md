# 媒体库数量统计自定义功能指南

## 📖 功能概述

此功能允许你自定义 Emby 客户端显示的媒体库统计数量，用于美化展示效果或隐藏真实的媒体库数量。

### 适用场景

- 🎨 **美化展示**: 显示整数或特定数字，让界面更美观
- 🔒 **隐私保护**: 不希望暴露真实的媒体库数量
- 🎯 **营销展示**: 公开服务器时显示更吸引人的数量
- 📊 **测试调试**: 快速测试不同数量下的界面表现

---

## 🎯 工作原理

### 请求流程

```
Emby 客户端请求: GET /Items/Counts
    ↓
检查配置: items-counts.enable
    ↓
├─ false → 回源透传（使用真实数据）
└─ true
    ↓
    检查 Query 参数: ParentId
    ↓
    ├─ 有 ParentId → 回源透传（查询特定媒体库）
    └─ 无 ParentId → 返回自定义数据（全局统计）
```

### 拦截规则

| 请求 | ParentId | 行为 | 返回数据 |
|------|---------|------|---------|
| `/Items/Counts` | 无 | 拦截 | 自定义统计 |
| `/Items/Counts?ParentId=123` | 有 | 透传 | Emby 真实数据 |
| `enable=false` | - | 透传 | Emby 真实数据 |

---

## ⚙️ 配置说明

### 完整配置示例

```yaml
items-counts:
  # 是否启用自定义统计
  enable: true

  # 电影数量
  movie-count: 1000

  # 剧集数量
  series-count: 500

  # 分集数量
  episode-count: 5000

  # 游戏数量
  game-count: 0

  # 艺术家数量
  artist-count: 0

  # 节目数量
  program-count: 0

  # 游戏系统数量
  game-system-count: 0

  # 预告片数量
  trailer-count: 0

  # 歌曲数量
  song-count: 0

  # 专辑数量
  album-count: 0

  # 音乐视频数量
  music-video-count: 0

  # 合集数量
  box-set-count: 0

  # 书籍数量
  book-count: 0

  # 总项目数量（0 = 自动计算）
  item-count: 0
```

### 字段说明

| 字段 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| `enable` | bool | 是 | `false` | 是否启用自定义统计 |
| `movie-count` | int | 否 | `0` | 电影数量 |
| `series-count` | int | 否 | `0` | 剧集数量 |
| `episode-count` | int | 否 | `0` | 分集数量 |
| `game-count` | int | 否 | `0` | 游戏数量 |
| `artist-count` | int | 否 | `0` | 艺术家数量 |
| `program-count` | int | 否 | `0` | 节目数量 |
| `game-system-count` | int | 否 | `0` | 游戏系统数量 |
| `trailer-count` | int | 否 | `0` | 预告片数量 |
| `song-count` | int | 否 | `0` | 歌曲数量 |
| `album-count` | int | 否 | `0` | 专辑数量 |
| `music-video-count` | int | 否 | `0` | 音乐视频数量 |
| `box-set-count` | int | 否 | `0` | 合集数量 |
| `book-count` | int | 否 | `0` | 书籍数量 |
| `item-count` | int | 否 | `0` | 总数量（0 = 自动计算） |

### 自动计算 item-count

如果将 `item-count` 设为 `0`，程序会自动计算为所有类型的总和：

```yaml
items-counts:
  enable: true
  movie-count: 1000
  series-count: 500
  episode-count: 5000
  item-count: 0  # 自动计算为 6500
```

---

## 📋 返回格式

### JSON 响应示例

```json
{
  "MovieCount": 1000,
  "SeriesCount": 500,
  "EpisodeCount": 5000,
  "GameCount": 0,
  "ArtistCount": 0,
  "ProgramCount": 0,
  "GameSystemCount": 0,
  "TrailerCount": 0,
  "SongCount": 0,
  "AlbumCount": 0,
  "MusicVideoCount": 0,
  "BoxSetCount": 0,
  "BookCount": 0,
  "ItemCount": 6500
}
```

### 响应头

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
```

---

## 🚀 使用示例

### 示例 1: 基础配置

**需求**: 显示 1000 部电影、500 部剧集

```yaml
items-counts:
  enable: true
  movie-count: 1000
  series-count: 500
  episode-count: 5000
  # 其他类型不使用，设为 0 或省略
```

**效果**: Emby 客户端首页显示 "电影: 1000, 剧集: 500"

---

### 示例 2: 整数美化

**需求**: 显示整千、整百的数字

```yaml
items-counts:
  enable: true
  movie-count: 5000
  series-count: 2000
  episode-count: 10000
  item-count: 17000
```

**效果**: 所有数字都是整数，界面更美观

---

### 示例 3: 隐藏真实数量

**需求**: 不想暴露真实的 30000+ 部电影

```yaml
items-counts:
  enable: true
  movie-count: 1000  # 真实有 30000，只显示 1000
  series-count: 500
```

**效果**: 客户端只看到 1000 部电影

---

### 示例 4: 音乐库配置

**需求**: 主要是音乐库

```yaml
items-counts:
  enable: true
  movie-count: 100
  series-count: 50
  artist-count: 500
  album-count: 2000
  song-count: 10000
```

**效果**: 显示音乐相关的统计数据

---

## 🔍 调试验证

### 1. 查看日志

启用后，日志会输出：

```
[INFO] Items/Counts: 返回自定义统计数据
```

未启用或有 ParentId 时：

```
[INFO] Items/Counts: 未启用自定义统计，回源透传
[INFO] Items/Counts: 请求包含 ParentId=123，回源透传
```

### 2. 浏览器测试

直接在浏览器访问：

```
http://your-server:8095/Items/Counts
```

应该返回配置的 JSON 数据。

### 3. 客户端验证

1. 打开 Emby 客户端
2. 查看首页的媒体库统计数字
3. 确认显示的是配置的自定义数量

---

## ⚠️ 注意事项

### 1. ParentId 参数

当请求包含 `ParentId` 时会回源透传：

```
/Items/Counts                    → 返回自定义数据 ✅
/Items/Counts?ParentId=123       → 回源透传 ⚠️
```

**原因**: `ParentId` 表示查询特定媒体库的统计，必须返回真实数据。

### 2. 数量限制

- ✅ 所有数量必须 >= 0
- ❌ 不能设置负数
- ⚠️ 设置过大的数字可能看起来不真实

### 3. 客户端缓存

某些 Emby 客户端会缓存统计数据，修改配置后可能需要：
- 清除客户端缓存
- 重启客户端
- 等待缓存过期

### 4. 配置验证

程序启动时会验证配置：

```
✅ 正确: movie-count: 1000
❌ 错误: movie-count: -100 (负数)
```

如果配置错误，程序会输出错误信息并启动失败。

---

## 🎨 最佳实践

### 1. 合理的数字

```yaml
# ❌ 不推荐 - 数字太假
movie-count: 999999

# ✅ 推荐 - 整数但合理
movie-count: 5000
```

### 2. 保持比例

```yaml
# ✅ 符合常理的比例
series-count: 500      # 剧集数量
episode-count: 5000    # 每部剧平均 10 集

# ❌ 不合理的比例
series-count: 500
episode-count: 100     # 平均每部剧 0.2 集？
```

### 3. 只配置需要的类型

```yaml
# ✅ 只设置视频类型
movie-count: 1000
series-count: 500
episode-count: 5000
# 音乐、游戏等不用的设为 0 或省略

# ❌ 设置了不使用的类型
game-count: 100  # 但实际没有游戏
```

### 4. 使用自动计算

```yaml
# ✅ 让程序自动计算总数
item-count: 0

# ❌ 手动计算容易出错
item-count: 6500
```

---

## 🐛 常见问题

### Q1: 修改配置后不生效？

**A**:
1. 重启 go-emby302 服务
2. 清除 Emby 客户端缓存
3. 检查日志确认配置已加载

### Q2: 特定媒体库的统计还是显示真实数据？

**A**: 这是正常的。当请求包含 `ParentId` 时（查询特定媒体库），会回源返回真实数据。自定义统计只影响全局统计（无 `ParentId`）。

### Q3: 可以设置不同媒体库的不同数量吗？

**A**: 不可以。此功能仅影响全局统计（不带 `ParentId` 的请求）。特定媒体库的统计始终返回真实数据。

### Q4: 数量设为 0 会有问题吗？

**A**: 不会。设为 0 表示该类型没有媒体，这是合法的。

### Q5: 如何恢复显示真实数据？

**A**: 将 `enable` 设为 `false` 即可：

```yaml
items-counts:
  enable: false  # 关闭自定义统计
```

---

## 📊 技术实现

### 拦截路由

```go
// 正则匹配: /Items/Counts 或 /items/counts
Reg_ItemsCounts = `(?i)^/items/counts($|\?)`
```

### 处理逻辑

```go
func HandleItemsCounts(c *gin.Context) {
    // 1. 检查是否启用
    if !config.C.ItemsCounts.Enable {
        ProxyOrigin(c)
        return
    }

    // 2. 检查 ParentId
    if c.Query("ParentId") != "" {
        ProxyOrigin(c)
        return
    }

    // 3. 返回自定义数据
    c.JSON(200, config.C.ItemsCounts.ToJSON())
}
```

---

## 📝 总结

- ✅ **简单易用**: 修改配置文件即可
- ✅ **灵活控制**: 每种媒体类型独立配置
- ✅ **智能判断**: 自动区分全局统计和特定媒体库查询
- ✅ **安全可靠**: 配置验证，防止错误配置
- ⚠️ **仅影响全局**: 特定媒体库统计仍显示真实数据

---

**最后更新**: 2026-01-08
