# STRM 路径映射功能修改说明

## 修改概述

将 STRM 文件从 URL 重定向模式改为本地路径映射模式，支持多个 CDN 域名映射。

---

## 主要改动

### 1. 配置文件格式变更

#### 旧格式（已移除）
```yaml
emby:
  strm:
    path-map:
      - https://origin.com => https://cdn.example.com
      - token=old => token=new
    internal-redirect-enable: false
```

#### 新格式
```yaml
emby:
  strm:
    path-mappings:
      - local-prefix: /mnt/media/剧集
        cdn-base: https://tx-cdn.example.com
        remote-prefix: /剧集

      - local-prefix: /mnt/media/电影
        cdn-base: https://tx-cdn.example.com
        remote-prefix: /电影

      - local-prefix: /mnt/media/近期更新
        cdn-base: https://ali-cdn.example.com
        remote-prefix: /近期更新
```

---

## 功能说明

### 路径映射逻辑

**STRM 文件内容示例：**
```
/mnt/media/剧集/国产剧/101次抢婚 (2023)/Season 1/101次抢婚 - S01E01 - 第 1 集.mp4
```

**映射配置：**
```yaml
- local-prefix: /mnt/media/剧集
  cdn-base: https://tx-cdn.example.com
  remote-prefix: /剧集
```

**转换过程：**
1. 检测到路径以 `/mnt/media/剧集` 开头
2. 去除本地前缀：`/国产剧/101次抢婚 (2023)/Season 1/101次抢婚 - S01E01 - 第 1 集.mp4`
3. 拼接 CDN 地址：`https://tx-cdn.example.com/剧集/国产剧/101次抢婚 (2023)/Season 1/101次抢婚 - S01E01 - 第 1 集.mp4`
4. 302 重定向到该地址

---

## 代码修改详情

### 文件：`internal/config/emby.go`

**变更：**
- 新增 `PathMapping` 结构体，包含 3 个字段：
  - `LocalPrefix`：本地路径前缀
  - `CdnBase`：CDN 域名
  - `RemotePrefix`：CDN 上的路径前缀

- 修改 `Strm` 结构体：
  - 移除：`PathMap []string`
  - 移除：`InternalRedirectEnable bool`
  - 新增：`PathMappings []PathMapping`

- 重写 `MapPath(localPath string) (string, error)` 方法：
  - 旧逻辑：简单的字符串替换
  - 新逻辑：前缀匹配 + URL 拼接
  - 返回值：增加 error 返回

### 文件：`internal/service/emby/redirect.go`

**变更：**
- 移除 `getFinalRedirectLink()` 函数（内部重定向跟随）
- 简化 `Redirect2OpenlistLink()` 函数：
  - 直接调用 `MapPath()` 获取 CDN URL
  - 移除内部重定向跟随逻辑
  - 直接 302 重定向到 CDN 地址

### 文件：`config.example.yml`

**变更：**
- 更新配置格式为结构化的 `path-mappings`
- 移除 `internal-redirect-enable` 配置项
- 添加详细的配置注释和示例

---

## 使用示例

### 场景：多CDN配置

```yaml
emby:
  host: http://192.168.1.100:8096
  strm:
    path-mappings:
      # 腾讯云 CDN - 剧集和电影
      - local-prefix: /mnt/media/剧集
        cdn-base: https://tx-cdn.example.com
        remote-prefix: /剧集

      - local-prefix: /mnt/media/电影
        cdn-base: https://tx-cdn.example.com
        remote-prefix: /电影

      # 阿里云 CDN - 近期更新
      - local-prefix: /mnt/media/近期更新
        cdn-base: https://ali-cdn.example.com
        remote-prefix: /近期更新

      # 其他存储 - 音乐
      - local-prefix: /mnt/media/音乐
        cdn-base: https://music-cdn.example.com
        remote-prefix: /music
```

### STRM 文件示例

#### 剧集（腾讯云）
**文件：** `/mnt/media/剧集/国产剧/狂飙/Season 1/狂飙 - S01E01.mp4`
**重定向到：** `https://tx-cdn.example.com/剧集/国产剧/狂飙/Season 1/狂飙 - S01E01.mp4`

#### 近期更新（阿里云）
**文件：** `/mnt/media/近期更新/2024年1月/新剧.mp4`
**重定向到：** `https://ali-cdn.example.com/近期更新/2024年1月/新剧.mp4`

---

## 执行流程

```
客户端请求播放
    ↓
解析 ItemId, ApiKey
    ↓
请求 Emby PlaybackInfo 接口
    ↓
获取 STRM 文件中的本地路径
例如: /mnt/media/剧集/国产剧/101次抢婚 (2023)/Season 1/101次抢婚 - S01E01 - 第 1 集.mp4
    ↓
遍历 path-mappings 配置
    ↓
匹配到 local-prefix: /mnt/media/剧集
    ↓
转换为 CDN URL
https://tx-cdn.example.com/剧集/国产剧/101次抢婚 (2023)/Season 1/101次抢婚 - S01E01 - 第 1 集.mp4
    ↓
返回 307 Redirect
    ↓
客户端直接从 CDN 拉流
```

---

## 配置校验

程序启动时会自动校验配置：

1. ✅ `path-mappings` 不能为空
2. ✅ 每个映射的 `local-prefix` 不能为空
3. ✅ 每个映射的 `cdn-base` 不能为空
4. ✅ 每个映射的 `remote-prefix` 不能为空
5. ✅ 自动去除路径末尾的 `/` 字符

---

## 注意事项

1. **匹配顺序**：从上到下匹配第一个符合的规则
2. **路径格式**：STRM 文件中必须是绝对路径（如 `/mnt/media/...`）
3. **URL 编码**：路径中的中文和特殊字符会自动进行 URL 编码
4. **缓存时间**：302 重定向响应缓存 10 分钟
5. **错误处理**：未匹配到规则时，根据 `proxy-error-strategy` 配置处理（回源或拒绝）

---

## 移除的功能

- ❌ 内部重定向跟随（`internal-redirect-enable`）
- ❌ 简单字符串替换映射（`path-map`）
- ❌ OpenList API 调用
- ❌ 路径转换服务（`path.Emby2Openlist`）
- ❌ 本地媒体根目录检测

---

## 性能优化

- ✅ 直接 302 重定向，无需额外的 HTTP 请求
- ✅ 配置在启动时预处理，运行时直接使用
- ✅ 前缀匹配效率高（O(n) 复杂度）
- ✅ 异步触发格式解析，不阻塞主流程

---

## 升级指南

### 从旧版本升级

1. **备份旧配置文件**
   ```bash
   cp config.yml config.yml.backup
   ```

2. **更新配置格式**
   ```yaml
   # 旧格式
   strm:
     path-map:
       - /old/path => /new/path

   # 新格式
   strm:
     path-mappings:
       - local-prefix: /old/path
         cdn-base: https://your-cdn.com
         remote-prefix: /new/path
   ```

3. **更新 STRM 文件内容**
   - 旧格式：`https://cdn.example.com/video.mp4`
   - 新格式：`/mnt/media/剧集/video.mp4`

4. **重启程序**
   ```bash
   ./gemby
   ```

---

## 故障排查

### 问题：302 重定向失败

**检查清单：**
1. 查看日志：`未找到匹配的路径映射规则`
2. 确认 STRM 文件路径格式正确
3. 确认配置中的 `local-prefix` 与 STRM 路径前缀一致

### 问题：CDN 地址拼接错误

**检查清单：**
1. 确认 `cdn-base` 不包含末尾的 `/`
2. 确认 `remote-prefix` 格式正确
3. 查看日志中的 "路径映射" 信息

---

生成时间：2024-01-08
