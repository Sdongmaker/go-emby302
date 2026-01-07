# CDN 鉴权功能使用指南

本文档详细说明 go-emby302 的 CDN 鉴权功能配置和使用方法。

## 目录

- [功能概述](#功能概述)
- [支持的 CDN 类型](#支持的-cdn-类型)
- [配置说明](#配置说明)
- [鉴权算法详解](#鉴权算法详解)
- [配置示例](#配置示例)
- [常见问题](#常见问题)

---

## 功能概述

go-emby302 支持将 Emby STRM 文件的本地路径映射为 CDN 直链，并支持以下三种鉴权方式：

1. **无鉴权 (none)**: 直接拼接 URL，适用于公共 CDN
2. **GoEdge 鉴权 (goedge)**: GoEdge CDN 专用鉴权算法
3. **腾讯云鉴权 (tencent)**: 腾讯云 CDN Type-A 鉴权算法

### 核心特性

- ✅ 支持多个 CDN 同时配置
- ✅ 相同域名的路径映射自动合并
- ✅ 自动生成带时间戳的鉴权签名
- ✅ 严格的路径前缀匹配，避免误匹配
- ✅ 支持中文路径自动编码

---

## 支持的 CDN 类型

### 1. 无鉴权 (none)

**适用场景**: 公共 CDN、内网访问、不需要鉴权的场景

**配置示例**:
```yaml
emby:
  strm:
    cdns:
      - name: "公共CDN"
        type: none
        base: https://cdn.example.com
        path-mappings:
          - local-prefix: /mnt/media
            remote-prefix: /public
```

**生成 URL 示例**:
```
输入: /mnt/media/电影/test.mp4
输出: https://cdn.example.com/public/电影/test.mp4
```

---

### 2. GoEdge 鉴权 (goedge)

**适用场景**: 使用 GoEdge CDN 的场景

**核心特点**:
- ⚡ 校验前先 URL 解码
- 🔐 使用 @ 作为签名原串分隔符
- 📝 签名格式: `ts-rand-md5`

**配置示例**:
```yaml
emby:
  strm:
    cdns:
      - name: "GoEdge主CDN"
        type: goedge
        base: https://cdn.goedge.com
        private-key: "your_goedge_secret_key"
        rand-length: 16  # 默认 16 位随机字符串
        path-mappings:
          - local-prefix: /mnt/media/剧集
            remote-prefix: /series
```

**签名算法**:
```
1. 原串: /series/国产剧/test.mp4 + "@" + 1234567890 + "@" + abc123def456 + "@" + your_goedge_secret_key
2. MD5: md5(原串) → 32位小写十六进制
3. 签名: 1234567890-abc123def456-{md5}
4. URL: https://cdn.goedge.com/%2Fseries%2F%E5%9B%BD%E4%BA%A7%E5%89%A7%2Ftest.mp4?sign=1234567890-abc123def456-{md5}
```

**生成 URL 示例**:
```
输入: /mnt/media/剧集/国产剧/test.mp4
输出: https://cdn.goedge.com/%2Fseries%2F%E5%9B%BD%E4%BA%A7%E5%89%A7%2Ftest.mp4?sign=1736352000-Abc123XyZ789-a1b2c3d4e5f6...
```

---

### 3. 腾讯云鉴权 (tencent)

**适用场景**: 使用腾讯云 CDN Type-A 鉴权的场景

**核心特点**:
- ⚡ 直接对 URL 编码后的路径验签
- 🔐 使用 - 作为签名原串分隔符
- 📝 签名格式: `ts-rand-uid-md5`
- 👤 支持自定义 uid

**配置示例**:
```yaml
emby:
  strm:
    cdns:
      - name: "腾讯云CDN"
        type: tencent
        base: https://cdn.tencent.com
        private-key: "your_tencent_secret_key"
        rand-length: 6   # 默认 6 位随机字符串
        uid: "0"         # 默认用户 ID
        path-mappings:
          - local-prefix: /mnt/media/电影
            remote-prefix: /movies
```

**签名算法**:
```
1. URI: URL 编码后的路径 → %2Fmovies%2F%E7%94%B5%E5%BD%B1%2Ftest.mp4
2. 原串: {URI} + "-" + 1234567890 + "-" + abc123 + "-" + 0 + "-" + your_tencent_secret_key
3. MD5: md5(原串) → 32位小写十六进制
4. 签名: 1234567890-abc123-0-{md5}
5. URL: https://cdn.tencent.com/%2Fmovies%2F%E7%94%B5%E5%BD%B1%2Ftest.mp4?sign=1234567890-abc123-0-{md5}
```

**生成 URL 示例**:
```
输入: /mnt/media/电影/国产电影/test.mp4
输出: https://cdn.tencent.com/%2Fmovies%2F%E5%9B%BD%E4%BA%A7%E7%94%B5%E5%BD%B1%2Ftest.mp4?sign=1736352000-X7yZ9a-0-b1c2d3e4f5...
```

---

## 配置说明

### 完整配置结构

```yaml
emby:
  strm:
    cdns:
      - name: "CDN名称"              # 必填，用于日志标识
        type: "goedge|tencent|none" # 必填，鉴权类型
        base: "https://cdn.com"     # 必填，CDN基础域名（不要以 / 结尾）
        private-key: "secret"       # 启用鉴权时必填
        rand-length: 16             # 可选，随机字符串长度
        uid: "0"                    # 可选，仅腾讯云使用
        path-mappings:              # 必填，路径映射列表
          - local-prefix: "/mnt/media"
            remote-prefix: "/remote"
```

### 字段说明

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | CDN 名称，用于日志输出 |
| `type` | string | 是 | 鉴权类型: `none`/`goedge`/`tencent` |
| `base` | string | 是 | CDN 基础域名，不要以 `/` 结尾 |
| `private-key` | string | 条件 | 鉴权密钥，`type` 不为 `none` 时必填 |
| `rand-length` | int | 否 | 随机字符串长度，0 表示使用 "0" |
| `uid` | string | 否 | 用户 ID，仅腾讯云使用，默认 "0" |
| `path-mappings` | array | 是 | 路径映射列表 |

### 路径映射规则

| 字段 | 说明 | 示例 |
|------|------|------|
| `local-prefix` | 本地路径前缀 | `/mnt/media/剧集` |
| `remote-prefix` | CDN 路径前缀 | `/series` |

**注意**:
- 所有路径前缀不要以 `/` 结尾（程序会自动处理）
- 路径匹配采用严格前缀匹配，避免误匹配
- 例如: `/mnt/media` 不会匹配 `/mnt/media2/file.mp4`

---

## 鉴权算法详解

### GoEdge vs 腾讯云核心差异

| 对比项 | GoEdge | 腾讯云 |
|--------|--------|--------|
| 路径形态 | 原始路径（未编码） | URL 编码路径 |
| 原串分隔符 | `@` | `-` |
| 签名格式 | `ts-rand-md5` | `ts-rand-uid-md5` |
| 随机字符串默认长度 | 16 | 6 |
| 额外字段 | 无 | uid |

### 签名安全性建议

1. ⚠️ **不要将 `rand-length` 设为 0**（除非测试）
2. 🔐 **妥善保管 `private-key`**，不要泄露
3. ⏰ **签名包含时间戳**，有效期由 CDN 服务器配置
4. 🔄 **每次请求生成新签名**，防止重放攻击

---

## 配置示例

### 示例 1: 单 CDN 多路径映射

```yaml
emby:
  strm:
    cdns:
      - name: "GoEdge主CDN"
        type: goedge
        base: https://cdn.goedge.com
        private-key: "your_secret_key"
        rand-length: 16
        path-mappings:
          - local-prefix: /mnt/media/剧集
            remote-prefix: /series
          - local-prefix: /mnt/media/电影
            remote-prefix: /movies
          - local-prefix: /mnt/media/综艺
            remote-prefix: /variety
```

### 示例 2: 多 CDN 配置

```yaml
emby:
  strm:
    cdns:
      # GoEdge CDN - 主要使用
      - name: "GoEdge主CDN"
        type: goedge
        base: https://cdn1.goedge.com
        private-key: "goedge_secret"
        path-mappings:
          - local-prefix: /mnt/media/剧集
            remote-prefix: /series

      # 腾讯云 CDN - 备用
      - name: "腾讯云备用CDN"
        type: tencent
        base: https://cdn2.tencent.com
        private-key: "tencent_secret"
        uid: "10086"
        path-mappings:
          - local-prefix: /mnt/backup/videos
            remote-prefix: /backup

      # 无鉴权 CDN - 公共资源
      - name: "公共CDN"
        type: none
        base: https://public.cdn.com
        path-mappings:
          - local-prefix: /mnt/public
            remote-prefix: /public
```

### 示例 3: 完整配置

```yaml
emby:
  host: http://localhost:8096
  images-quality: 70
  proxy-error-strategy: origin
  episodes-unplay-prior: true
  resort-random-items: true

  strm:
    cdns:
      - name: "GoEdge-主力CDN"
        type: goedge
        base: https://cdn.example.com
        private-key: "abcd1234efgh5678"
        rand-length: 16
        path-mappings:
          - local-prefix: /mnt/storage/series
            remote-prefix: /series
          - local-prefix: /mnt/storage/movies
            remote-prefix: /movies

cache:
  enable: true
```

---

## 常见问题

### Q1: 为什么要将相同域名的映射合并？

**A**: 合并配置有以下优点：
- 减少配置冗余
- 统一管理鉴权密钥
- 提高配置可读性
- 避免密钥配置错误

### Q2: GoEdge 和腾讯云鉴权有什么区别？

**A**: 主要区别：
1. **路径处理**: GoEdge 用原始路径，腾讯云用编码路径
2. **分隔符**: GoEdge 用 `@`，腾讯云用 `-`
3. **签名字段**: 腾讯云多一个 `uid` 字段

### Q3: rand-length 设为 0 安全吗？

**A**: **不推荐**！设为 0 时随机字符串固定为 "0"，安全性极低，仅适用于测试环境。

### Q4: 如何验证签名是否正确？

**A**: 可以：
1. 查看日志输出的最终 URL
2. 在浏览器中访问该 URL
3. 运行单元测试: `go test ./internal/util/cdnauth/...`

### Q5: 路径映射不生效怎么办？

**A**: 检查以下几点：
1. `local-prefix` 是否与实际路径匹配
2. 路径前缀不要以 `/` 结尾
3. 查看日志中的错误信息
4. 确认配置文件格式正确（YAML 缩进）

### Q6: 支持自定义鉴权算法吗？

**A**: 目前仅支持 GoEdge 和腾讯云。如需其他算法，可以：
1. 在 `internal/util/cdnauth/cdnauth.go` 中添加新函数
2. 在 `internal/config/emby.go` 中添加新的鉴权类型
3. 提交 Pull Request

### Q7: 鉴权 URL 的有效期是多久？

**A**: 有效期由 CDN 服务器端配置决定，通常建议：
- GoEdge: 30-60 分钟
- 腾讯云: 30-60 分钟

---

## 技术支持

如有问题，请：
1. 查看日志输出
2. 检查配置文件格式
3. 运行单元测试验证
4. 提交 GitHub Issue

---

**最后更新**: 2026-01-08
