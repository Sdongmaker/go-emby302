# go-emby302 功能列表

本文档列出 go-emby302 的所有主要功能。

---

## 🎬 核心功能

### 1. STRM 文件直链播放

**说明**: 将 Emby STRM 文件的本地路径映射为 CDN 直链，302 重定向播放。

**配置**: `emby.strm`

**文档**: [CDN 鉴权指南](CDN_AUTH_GUIDE.md) | [多 CDN 配置](MULTI_CDN_GUIDE.md)

**特性**:
- ✅ 支持多个 CDN 同时配置
- ✅ 支持 GoEdge 和腾讯云 CDN 鉴权
- ✅ 自动路径匹配和 URL 编码
- ✅ 支持中文路径

---

### 2. 防止客户端转码

**说明**: 修改 PlaybackInfo 响应，强制客户端直接播放，禁止转码。

**配置**: 自动启用，无需配置

**特性**:
- ✅ SupportsDirectPlay = true
- ✅ SupportsTranscoding = false
- ✅ 移除 TranscodingUrl
- ✅ 节省服务器资源

---

### 3. 自定义媒体库数量统计 🆕

**说明**: 自定义 Emby 客户端显示的媒体库统计数量。

**配置**: `items-counts`

**文档**: [使用指南](ITEMS_COUNTS_GUIDE.md)

**特性**:
- ✅ 自定义各类媒体数量
- ✅ 美化展示效果
- ✅ 隐私保护
- ✅ 自动计算总数

**使用示例**:
```yaml
items-counts:
  enable: true
  movie-count: 1000
  series-count: 500
  episode-count: 5000
```

**效果**: Emby 客户端首页显示自定义的媒体库数量。

---

### 4. PlaybackInfo 缓存

**说明**: 缓存 PlaybackInfo 响应，减少对 Emby 服务器的请求。

**配置**: `cache.enable`

**特性**:
- ✅ 12 小时缓存时长
- ✅ 按 ItemId + ApiKey 缓存
- ✅ 支持音轨、字幕切换
- ✅ 提升响应速度

---

### 5. 随机播放重排序

**说明**: 对随机播放列表进行重新 Shuffle，避免 Emby 的固定随机顺序。

**配置**: `emby.resort-random-items: true`

**特性**:
- ✅ 缓存原始列表（3 小时）
- ✅ 每次请求重新 Shuffle
- ✅ 真正的随机播放

---

### 6. 剧集未播优先

**说明**: 在获取剧集列表时，将未播放的剧集优先展示。

**配置**: `emby.episodes-unplay-prior: true`

**特性**:
- ✅ 未播放剧集排在最前
- ✅ 方便追剧
- ✅ 智能排序

---

### 7. 图片质量控制

**说明**: 统一修改所有图片请求的质量参数。

**配置**: `emby.images-quality: 70`

**特性**:
- ✅ 全局质量控制
- ✅ 减少带宽消耗
- ✅ 加快加载速度

---

### 8. 字幕直链访问

**说明**: 强制外挂字幕使用直链访问。

**特性**:
- ✅ DeliveryMethod = External
- ✅ 自动生成 DeliveryUrl
- ✅ 提高字幕加载速度

---

### 9. 媒体名称简化

**说明**: 简化 MediaSource 中的视频名称。

**示例**:
```
原始: 国产剧/101次抢婚 (2023)/Season 1/S01E01.mp4
简化: S01E01
```

---

### 10. 错误处理策略

**说明**: 配置代理失败时的处理策略。

**配置**: `emby.proxy-error-strategy`

**选项**:
- `origin`: 回源透传（默认）
- `reject`: 拒绝请求

---

## 📊 统计功能

### Items/Counts 自定义统计

| 功能 | 说明 | 配置项 |
|------|------|--------|
| 电影统计 | 自定义电影数量 | `movie-count` |
| 剧集统计 | 自定义剧集数量 | `series-count` |
| 分集统计 | 自定义分集数量 | `episode-count` |
| 音乐统计 | 歌曲、专辑、艺术家 | `song-count`, `album-count`, `artist-count` |
| 游戏统计 | 游戏、游戏系统 | `game-count`, `game-system-count` |
| 其他统计 | 预告片、节目、书籍等 | `trailer-count`, `program-count`, `book-count` |
| 总数统计 | 自动计算或手动设置 | `item-count` |

**请求示例**:
```
GET /Items/Counts           → 返回自定义统计
GET /Items/Counts?ParentId=123 → 回源返回真实统计
```

---

## 🔐 CDN 鉴权功能

### 支持的 CDN 类型

| CDN 类型 | 鉴权算法 | 特点 |
|---------|---------|------|
| GoEdge | MD5(path@ts@rand@key) | 校验前 URL 解码 |
| 腾讯云 | MD5(uri-ts-rand-uid-key) | 直接对 URI 验签 |
| 无鉴权 | - | 直接拼接 URL |

### 签名格式

**GoEdge**:
```
签名格式: ts-rand-md5
URL: /{encoded_path}?sign=1234567890-abc123-md5hash
```

**腾讯云**:
```
签名格式: ts-rand-uid-md5
URL: /{encoded_path}?sign=1234567890-abc123-0-md5hash
```

---

## 🛠️ 其他功能

### WebSocket 代理

**说明**: 透明代理 Emby WebSocket 连接。

**特性**:
- ✅ 实时消息推送
- ✅ 保持连接状态
- ✅ 完全透明

---

### 自定义脚本注入

**说明**: 在 Emby Web 界面注入自定义 JS 和 CSS。

**目录**:
- `custom-js/`: 存放 JS 文件
- `custom-css/`: 存放 CSS 文件

**特性**:
- ✅ 自动加载
- ✅ 支持多文件
- ✅ 实时生效

---

## 📝 配置文件结构

```yaml
emby:
  host: http://localhost:8096
  images-quality: 70
  proxy-error-strategy: origin
  episodes-unplay-prior: true
  resort-random-items: true
  strm:
    cdns:
      - name: "GoEdge-主CDN"
        type: goedge
        base: https://cdn.example.com
        private-key: "secret"
        path-mappings:
          - local-prefix: /mnt/media/剧集
            remote-prefix: /series

cache:
  enable: true

items-counts:
  enable: true
  movie-count: 1000
  series-count: 500
  episode-count: 5000
```

---

## 📚 相关文档

- [CDN 鉴权指南](CDN_AUTH_GUIDE.md)
- [多 CDN 配置](MULTI_CDN_GUIDE.md)
- [媒体库数量统计指南](ITEMS_COUNTS_GUIDE.md)
- [配置示例](config.example.yml)

---

## 🎯 使用场景

| 场景 | 推荐功能 |
|------|---------|
| 公共 Emby 服务器 | Items/Counts 自定义统计 + CDN 鉴权 |
| 私人媒体库 | STRM 直链播放 + 防转码 |
| 低带宽服务器 | 图片质量控制 + CDN 分流 |
| 追剧优化 | 剧集未播优先 + 随机播放重排序 |
| 隐私保护 | Items/Counts 隐藏真实数量 |

---

**最后更新**: 2026-01-08
