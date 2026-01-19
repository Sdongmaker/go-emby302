# go-emby302 功能列表

本项目是一个放在 Emby 前面的反向代理，核心目标是 **STRM 直链播放（302）** 与 **PlaybackInfo 改写（直放优先/禁用转码）**。

## 核心功能

- **STRM 文件直链播放（302）**
  - 拦截 `/videos|audio/.../(stream|universal)`、`/items/.../download` 等播放/下载请求
  - 从 Emby 获取媒体本地路径后，按 `emby.strm.cdns[].path-mappings` 映射为 CDN URL，并 302 重定向
  - 支持 `none` / `goedge` / `tencent` 三种 CDN 方式（见 `CDN_AUTH_GUIDE.md`）

- **PlaybackInfo 改写（直放优先）**
  - 拦截 `/Items/*/PlaybackInfo`
  - 强制 DirectPlay/DirectStream，移除转码相关字段，减少服务器转码压力

- **PlaybackInfo 缓存（可选）**
  - `cache.enable: true` 后对部分接口进行缓存，降低 Emby 源站压力

## 可选增强

- **图片质量统一**：`emby.images-quality`
- **/Items/Counts 自定义统计**：`items-counts.enable: true`（见 `ITEMS_COUNTS_GUIDE.md`）
- **Web 注入自定义 JS/CSS**：把文件放到 `custom-js/`、`custom-css/`（支持在文件内写 URL 进行远程加载）
- **HTTPS 支持**：`ssl.enable: true`（证书放 `ssl/`）
- **OpenList 本地目录树生成**：`openlist.local-tree-gen.enable: true`

