# 配置文件使用指南

## 配置模板

- `config.example.yml`：主配置模板（推荐从这里开始）
- `config_multi_cdn.example.yml`：多 CDN 示例（需要多规则/主备 CDN 时参考）

## 快速开始

1. 复制模板为 `config.yml`：

```bash
cp config.example.yml config.yml
```

2. 修改必填项：

- `emby.host`：你的 Emby 地址（示例：`http://YOUR_EMBY:8096`）
- `emby.strm.cdns`：至少配置一个 CDN + 一组 `path-mappings`

3. 启动（Docker Compose）：

```bash
docker compose up -d --build
```

## 常用配置

### CDN（STRM 302）

```yaml
emby:
  strm:
    cdns:
      - name: "main"
        type: none # none/goedge/tencent
        base: https://cdn.example.com
        path-mappings:
          - local-prefix: /mnt/media/剧集
            remote-prefix: /series
```

更多鉴权细节见 `CDN_AUTH_GUIDE.md`、多 CDN 场景见 `MULTI_CDN_GUIDE.md`。

### 缓存

```yaml
cache:
  enable: true
```

### /Items/Counts 自定义统计

```yaml
items-counts:
  enable: true
  movie-count: 1000
  series-count: 500
  episode-count: 5000
  item-count: 0 # 0=自动计算
```

完整说明见 `ITEMS_COUNTS_GUIDE.md`。

### HTTPS

- 证书与私钥放在 `ssl/`
- 在 `config.yml` 里配置 `ssl.enable`、`ssl.crt`、`ssl.key`

## 启动与排错

- 查看日志：`docker compose logs -f`
- 配置错误通常会在启动时直接退出并打印原因

