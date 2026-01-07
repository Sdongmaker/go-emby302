# 多 CDN 路径映射完全指南

## 📖 目录

- [功能概述](#功能概述)
- [路径匹配逻辑](#路径匹配逻辑)
- [配置场景示例](#配置场景示例)
- [最佳实践](#最佳实践)
- [常见问题](#常见问题)

---

## 功能概述

go-emby302 支持**多个 CDN 同时生效**，会根据文件的本地路径自动选择对应的 CDN。

### 核心特性

- ✅ **多 CDN 并存**: 配置多个 CDN，每个 CDN 独立配置鉴权
- ✅ **自动路径匹配**: 根据本地路径前缀自动选择对应的 CDN
- ✅ **灵活分配**: 不同类型的媒体可以使用不同的 CDN
- ✅ **独立鉴权**: 每个 CDN 可以使用不同的鉴权方式和密钥

### 典型应用场景

| 场景 | 说明 |
|------|------|
| **按媒体类型分配** | 剧集用 CDN A，电影用 CDN B |
| **按地域分配** | 国内资源用 CDN A，海外资源用 CDN B |
| **按画质分配** | 高清资源用 CDN A，标清资源用 CDN B |
| **主备 CDN** | 主要资源用主 CDN，备份资源用备用 CDN |
| **按存储位置分配** | 不同存储路径使用不同 CDN |

---

## 路径匹配逻辑

### 匹配流程图

```
客户端请求播放
    ↓
Emby 返回 STRM 文件路径: /mnt/media/剧集/国产剧/xxx.mp4
    ↓
go-emby302 开始匹配 CDN
    ↓
┌─────────────────────────────────────────┐
│ 遍历配置文件中的所有 CDN（按顺序）        │
└─────────────────────────────────────────┘
    ↓
┌─────────────────────────────────────────┐
│ CDN 1: GoEdge-剧集专用                   │
│   └─ 匹配路径 1: /mnt/media/剧集          │  ← ✅ 匹配成功！
│       ├─ 匹配检查: /mnt/media/剧集/...    │
│       └─ 匹配结果: 是                     │
└─────────────────────────────────────────┘
    ↓
使用 GoEdge CDN 生成签名
    ↓
返回 302 重定向: https://cdn1.goedge.com/...?sign=...
```

### 核心代码逻辑

```go
func (s *Strm) MapPath(localPath string) (string, error) {
    // 遍历所有 CDN 配置
    for _, cdn := range s.Cdns {
        // 遍历该 CDN 下的所有路径映射
        for _, mapping := range cdn.PathMappings {
            // 严格匹配路径前缀
            if !matchPathPrefix(localPath, mapping.LocalPrefix) {
                continue  // 不匹配，继续下一个
            }

            // 找到匹配，生成 CDN URL
            relativePath := strings.TrimPrefix(localPath, mapping.LocalPrefix)
            cdnPath := mapping.RemotePrefix + relativePath

            // 根据 CDN 类型生成鉴权 URL
            return generateAuthUrl(cdn, cdnPath)
        }
    }

    // 没有找到匹配的映射
    return "", fmt.Errorf("未找到匹配的路径映射规则: %s", localPath)
}
```

### 严格前缀匹配

程序使用**严格前缀匹配**，避免误匹配：

```go
// ✅ 正确匹配
local-prefix: /mnt/media
路径: /mnt/media/file.mp4         → 匹配成功

// ❌ 不会误匹配
local-prefix: /mnt/media
路径: /mnt/media2/file.mp4        → 匹配失败（避免误匹配）

// ✅ 完全匹配
local-prefix: /mnt/media
路径: /mnt/media                  → 匹配成功
```

---

## 配置场景示例

### 场景 1: 按媒体类型分配 CDN

**需求**: 剧集和动漫用 GoEdge，电影用腾讯云

```yaml
emby:
  strm:
    cdns:
      # GoEdge CDN - 用于剧集和动漫
      - name: "GoEdge-剧集"
        type: goedge
        base: https://cdn1.goedge.com
        private-key: "goedge_secret"
        path-mappings:
          - local-prefix: /mnt/media/剧集
            remote-prefix: /series
          - local-prefix: /mnt/media/动漫
            remote-prefix: /anime

      # 腾讯云 CDN - 用于电影
      - name: "腾讯云-电影"
        type: tencent
        base: https://cdn2.tencent.com
        private-key: "tencent_secret"
        path-mappings:
          - local-prefix: /mnt/media/电影
            remote-prefix: /movies
```

**匹配结果**:
```
/mnt/media/剧集/国产剧/xxx.mp4  → GoEdge CDN
/mnt/media/动漫/xxx.mp4         → GoEdge CDN
/mnt/media/电影/xxx.mp4         → 腾讯云 CDN
```

---

### 场景 2: 按存储位置分配 CDN

**需求**: 主存储用 CDN A，备份存储用 CDN B

```yaml
emby:
  strm:
    cdns:
      # 主存储 CDN
      - name: "主CDN"
        type: goedge
        base: https://main-cdn.com
        private-key: "main_secret"
        path-mappings:
          - local-prefix: /mnt/storage1
            remote-prefix: /main

      # 备份存储 CDN
      - name: "备份CDN"
        type: tencent
        base: https://backup-cdn.com
        private-key: "backup_secret"
        path-mappings:
          - local-prefix: /mnt/backup
            remote-prefix: /backup
```

**匹配结果**:
```
/mnt/storage1/剧集/xxx.mp4  → 主 CDN
/mnt/backup/剧集/xxx.mp4    → 备份 CDN
```

---

### 场景 3: 按画质分配 CDN

**需求**: 4K 资源用高速 CDN，普通资源用标准 CDN

```yaml
emby:
  strm:
    cdns:
      # 高速 CDN - 4K 专用
      - name: "高速CDN-4K"
        type: goedge
        base: https://fast-cdn.com
        private-key: "fast_secret"
        path-mappings:
          - local-prefix: /mnt/media/4K剧集
            remote-prefix: /4k-series
          - local-prefix: /mnt/media/4K电影
            remote-prefix: /4k-movies

      # 标准 CDN - 普通资源
      - name: "标准CDN"
        type: tencent
        base: https://standard-cdn.com
        private-key: "standard_secret"
        path-mappings:
          - local-prefix: /mnt/media/剧集
            remote-prefix: /series
          - local-prefix: /mnt/media/电影
            remote-prefix: /movies
```

**匹配结果**:
```
/mnt/media/4K剧集/xxx.mp4  → 高速 CDN
/mnt/media/剧集/xxx.mp4    → 标准 CDN
```

---

### 场景 4: 混合鉴权模式

**需求**: 付费内容用鉴权 CDN，免费内容用公共 CDN

```yaml
emby:
  strm:
    cdns:
      # 鉴权 CDN - 付费内容
      - name: "付费CDN"
        type: goedge
        base: https://premium-cdn.com
        private-key: "premium_secret"
        path-mappings:
          - local-prefix: /mnt/media/付费剧集
            remote-prefix: /premium

      # 无鉴权 CDN - 免费内容
      - name: "免费CDN"
        type: none
        base: https://public-cdn.com
        path-mappings:
          - local-prefix: /mnt/media/免费剧集
            remote-prefix: /public
```

**匹配结果**:
```
/mnt/media/付费剧集/xxx.mp4  → 带鉴权 URL
/mnt/media/免费剧集/xxx.mp4  → 无鉴权 URL
```

---

### 场景 5: 复杂多层级配置

**需求**: 多种媒体类型，多个存储位置，多个 CDN

```yaml
emby:
  strm:
    cdns:
      # CDN 1: GoEdge - 剧集（存储1）
      - name: "GoEdge-剧集-存储1"
        type: goedge
        base: https://cdn1.goedge.com
        private-key: "goedge1_secret"
        path-mappings:
          - local-prefix: /mnt/storage1/剧集
            remote-prefix: /series

      # CDN 2: GoEdge - 剧集（存储2）
      - name: "GoEdge-剧集-存储2"
        type: goedge
        base: https://cdn2.goedge.com
        private-key: "goedge2_secret"
        path-mappings:
          - local-prefix: /mnt/storage2/剧集
            remote-prefix: /series

      # CDN 3: 腾讯云 - 电影
      - name: "腾讯云-电影"
        type: tencent
        base: https://cdn.tencent.com
        private-key: "tencent_secret"
        path-mappings:
          - local-prefix: /mnt/storage1/电影
            remote-prefix: /movies
          - local-prefix: /mnt/storage2/电影
            remote-prefix: /movies

      # CDN 4: 公共 CDN - 公开资源
      - name: "公共CDN"
        type: none
        base: https://public.cdn.com
        path-mappings:
          - local-prefix: /mnt/public
            remote-prefix: /public
```

**匹配结果**:
```
/mnt/storage1/剧集/xxx.mp4  → GoEdge CDN 1
/mnt/storage2/剧集/xxx.mp4  → GoEdge CDN 2
/mnt/storage1/电影/xxx.mp4  → 腾讯云 CDN
/mnt/storage2/电影/xxx.mp4  → 腾讯云 CDN
/mnt/public/xxx.mp4         → 公共 CDN
```

---

## 最佳实践

### 1. 配置顺序很重要 ⚠️

**错误示例** - 通用规则在前，具体规则被忽略：
```yaml
cdns:
  - name: "通用CDN"
    path-mappings:
      - local-prefix: /mnt/media        # ❌ 太宽泛，会匹配所有

  - name: "剧集CDN"
    path-mappings:
      - local-prefix: /mnt/media/剧集   # ⚠️ 永远不会被匹配到
```

**正确示例** - 具体规则在前：
```yaml
cdns:
  - name: "剧集CDN"
    path-mappings:
      - local-prefix: /mnt/media/剧集   # ✅ 先匹配具体路径

  - name: "通用CDN"
    path-mappings:
      - local-prefix: /mnt/media        # ✅ 兜底匹配
```

### 2. 合并相同 CDN 配置

**不推荐** - 重复配置：
```yaml
cdns:
  - name: "CDN-剧集"
    type: goedge
    base: https://cdn.goedge.com       # 相同域名
    private-key: "same_secret"          # 相同密钥
    path-mappings:
      - local-prefix: /mnt/media/剧集
        remote-prefix: /series

  - name: "CDN-电影"
    type: goedge
    base: https://cdn.goedge.com       # 重复！
    private-key: "same_secret"          # 重复！
    path-mappings:
      - local-prefix: /mnt/media/电影
        remote-prefix: /movies
```

**推荐** - 合并配置：
```yaml
cdns:
  - name: "GoEdge主CDN"
    type: goedge
    base: https://cdn.goedge.com
    private-key: "same_secret"
    path-mappings:
      - local-prefix: /mnt/media/剧集
        remote-prefix: /series
      - local-prefix: /mnt/media/电影
        remote-prefix: /movies
```

### 3. 使用有意义的 CDN 名称

```yaml
# ❌ 不好的命名
- name: "cdn1"
- name: "cdn2"

# ✅ 好的命名
- name: "GoEdge-剧集专用"
- name: "腾讯云-电影备份"
- name: "公共CDN-免费内容"
```

### 4. 添加注释说明

```yaml
cdns:
  # ============================================
  # 主力 CDN - 用于所有剧集和动漫
  # 鉴权: GoEdge
  # 带宽: 100Mbps
  # ============================================
  - name: "GoEdge-主力"
    type: goedge
    base: https://main-cdn.com
    private-key: "xxx"
    path-mappings:
      - local-prefix: /mnt/media/剧集
        remote-prefix: /series
```

### 5. 测试配置

**验证步骤**:
1. 查看启动日志，确认配置加载成功
2. 播放不同路径的文件，检查日志输出
3. 确认使用了正确的 CDN

**日志示例**:
```
[INFO] 路径映射 [GoEdge-剧集]: [/mnt/media/剧集/xxx.mp4] -> [https://cdn1.goedge.com/...]
[INFO] 路径映射 [腾讯云-电影]: [/mnt/media/电影/xxx.mp4] -> [https://cdn2.tencent.com/...]
```

---

## 常见问题

### Q1: 如何查看文件匹配到了哪个 CDN？

**A**: 查看程序日志，会输出类似：
```
[INFO] 路径映射 [CDN名称]: [本地路径] -> [CDN URL]
```

### Q2: 路径匹配不生效怎么办？

**A**: 检查以下几点：
1. ✅ `local-prefix` 是否正确（严格匹配，区分大小写）
2. ✅ 配置顺序是否正确（具体规则在前）
3. ✅ 路径前缀不要以 `/` 结尾
4. ✅ 查看日志中的错误信息

### Q3: 可以为同一个路径配置多个 CDN 作为备份吗？

**A**: 当前实现找到第一个匹配就返回。如需备份，建议：
- 在 CDN 层面配置主备切换
- 或使用负载均衡

### Q4: 支持正则表达式匹配吗？

**A**: 当前仅支持前缀匹配，不支持正则。如有需求可以提 Issue。

### Q5: 如何处理中文路径？

**A**: 程序会自动进行 URL 编码，无需手动处理。

### Q6: 一个 CDN 可以配置多少个路径映射？

**A**: 理论上无限制，但建议合理组织，避免配置过于复杂。

---

## 调试技巧

### 1. 启用详细日志

查看日志输出：
```bash
docker compose logs -f | grep "路径映射"
```

### 2. 测试配置

创建测试文件，观察匹配结果：
```bash
# 创建测试路径
mkdir -p /mnt/media/剧集/测试
touch /mnt/media/剧集/测试/test.mp4

# 在 Emby 中播放，观察日志
```

### 3. 验证 URL

复制日志中生成的 URL，在浏览器中访问，验证是否能正常访问。

---

## 总结

- ✅ **多 CDN 并存**: 支持配置多个 CDN，自动路径匹配
- ✅ **灵活配置**: 按媒体类型、存储位置、画质等任意维度分配
- ✅ **独立鉴权**: 每个 CDN 独立配置鉴权方式和密钥
- ✅ **严格匹配**: 避免路径误匹配
- ⚠️ **顺序重要**: 具体规则配置在前，通用规则配置在后
- 📝 **日志查看**: 通过日志确认路径匹配结果

---

**最后更新**: 2026-01-08
