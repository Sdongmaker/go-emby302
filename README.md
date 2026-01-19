<div align="center">
  <img height="150" src="./assets/logo.png" />
  <h1>go-emby302</h1>
  <p>Emby 反向代理：将 STRM 中的本地路径映射为 CDN 直链并 302 重定向，同时改写 PlaybackInfo 强制直放（可选缓存）。</p>
</div>

## 功能

- **STRM 直链播放（302）**：解析 Emby 媒体本地路径，按 `emby.strm` 规则映射到 CDN URL 并 302。
- **防转码 / 直放优先**：改写 `/Items/*/PlaybackInfo`，启用 DirectPlay/DirectStream，移除转码相关字段。
- **PlaybackInfo 缓存（可选）**：开启 `cache.enable` 后减少对 Emby 的重复请求。
- **自定义统计（可选）**：拦截 `/Items/Counts`，按 `items-counts` 返回自定义数量。
- **OpenList 本地目录树（可选）**：按 `openlist.local-tree-gen` 将 OpenList 目录生成到本地（`strm`/虚拟媒体/或下载源文件）。

## 部署

### Docker Hub 镜像

镜像已发布到 Docker Hub：

```bash
docker pull tdck/go-emby302
```

如果你使用本仓库自带的 `docker-compose.yml`，建议将 `build:` 改为：

```yaml
image: tdck/go-emby302:latest
```

然后启动：

```bash
docker compose up -d
```

### Docker Compose 本地构建

1. 复制配置：

   ```bash
   cp config.example.yml config.yml
   ```

2. 编辑 `config.yml`（至少需要配置 `emby.host`、`emby.strm.cdns`）。
3. 启动：

   ```bash
   docker compose up -d --build
   ```

4. 访问：`http://<服务器IP>:8095`

查看日志：

```bash
docker compose logs -f
```

## 配置要点

最小可用配置示例（更多见 `config.example.yml` / `config_multi_cdn.example.yml`）：

```yaml
emby:
  host: http://YOUR_EMBY:8096
  strm:
    cdns:
      - name: "main"
        type: none # none/goedge/tencent
        base: https://cdn.example.com
        path-mappings:
          - local-prefix: /mnt/media/剧集
            remote-prefix: /series

cache:
  enable: true
```

常用开关：

- `emby.proxy-error-strategy`: `origin`（回源透传）/ `reject`（直接报错）
- `ssl.enable`: 是否启用 HTTPS（证书放在 `ssl/` 并在配置里写文件名）
- `items-counts.enable`: 是否启用 `/Items/Counts` 自定义统计
- `openlist.local-tree-gen.enable`: 是否生成 OpenList 本地目录树（需要配置 `openlist.host`、`openlist.token`）

命令行参数：

- `-p`: HTTP 监听端口（默认 `8095`）
- `-ps`: HTTPS 监听端口（默认 `8094`）
- `-dr`: 数据根目录（默认当前目录，`config.yml` 也从这里读取）
- `-version`: 打印版本号

## 相关文档

- `CONFIG_README.md`：配置文件使用说明
- `FEATURES.md`：功能清单
- `CDN_AUTH_GUIDE.md`：CDN 鉴权说明
- `MULTI_CDN_GUIDE.md`：多 CDN 配置示例
- `ITEMS_COUNTS_GUIDE.md`：媒体库数量统计说明

## License

见 `LICENSE`。
