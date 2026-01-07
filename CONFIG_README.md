# 配置文件使用指南

## 📁 配置文件说明

项目包含两个配置文件模板：

### 1. config.example.yml（主配置模板）

**用途**: 标准配置模板，适合大多数使用场景。

**特点**:
- ✅ 包含所有功能的配置项
- ✅ 详细的中文注释
- ✅ 适合快速上手
- ✅ 推荐日常使用

**适用场景**:
- 新手快速配置
- 单 CDN 配置
- 基础功能使用

---

### 2. config_multi_cdn.example.yml（多 CDN 示例）

**用途**: 多 CDN 配置的详细示例。

**特点**:
- ✅ 多个 CDN 并存的完整配置
- ✅ 各种使用场景示例
- ✅ 详细的配置说明和注释
- ✅ 适合参考学习

**适用场景**:
- 需要配置多个 CDN
- 不同类型媒体使用不同 CDN
- 主备 CDN 配置
- 复杂场景参考

---

## 🚀 快速开始

### 步骤 1: 复制配置文件

根据你的需求选择一个模板：

**单 CDN 或基础配置**:
```bash
cp config.example.yml config.yml
```

**多 CDN 配置**:
```bash
cp config_multi_cdn.example.yml config.yml
```

### 步骤 2: 编辑配置

```bash
vim config.yml
# 或
nano config.yml
```

### 步骤 3: 必改项

无论使用哪个模板，以下配置**必须修改**：

#### ① Emby 服务器地址
```yaml
emby:
  host: http://YOUR_EMBY_SERVER:8096  # 改成你的 Emby 地址
```

#### ② CDN 配置（如果使用 STRM）
```yaml
emby:
  strm:
    cdns:
      - name: "你的CDN名称"
        type: goedge  # 或 tencent 或 none
        base: https://YOUR_CDN_DOMAIN  # 改成你的 CDN 域名
        private-key: "YOUR_SECRET_KEY"  # 改成你的鉴权密钥
        path-mappings:
          - local-prefix: /YOUR/LOCAL/PATH  # 改成你的本地路径
            remote-prefix: /YOUR/REMOTE/PATH  # 改成你的 CDN 路径
```

### 步骤 4: 启动服务

```bash
docker compose up -d
```

---

## 📝 配置对比

| 配置项 | config.example.yml | config_multi_cdn.example.yml |
|--------|-------------------|------------------------------|
| **CDN 数量** | 1 个示例 | 4 个示例（详细） |
| **配置场景** | 基础场景 | 多种复杂场景 |
| **注释详细度** | 中等 | 非常详细 |
| **文件大小** | 较小（4.9K） | 较大（5.3K） |
| **推荐用户** | 新手、单 CDN | 高级用户、多 CDN |

---

## 🔧 主要配置项说明

### Emby 配置

```yaml
emby:
  host: http://localhost:8096      # Emby 服务器地址
  images-quality: 70               # 图片质量 (1-100)
  proxy-error-strategy: origin     # 错误策略: origin/reject
  episodes-unplay-prior: true      # 未播剧集优先
  resort-random-items: true        # 随机列表重排序
```

### CDN 配置

```yaml
emby:
  strm:
    cdns:
      - name: "CDN名称"
        type: goedge              # 鉴权类型: none/goedge/tencent
        base: https://cdn.com     # CDN 域名
        private-key: "secret"     # 鉴权密钥
        rand-length: 16           # 随机字符串长度
        path-mappings:            # 路径映射
          - local-prefix: /local
            remote-prefix: /remote
```

### 媒体库数量统计

```yaml
items-counts:
  enable: true           # 是否启用
  movie-count: 1000      # 电影数量
  series-count: 500      # 剧集数量
  episode-count: 5000    # 分集数量
  item-count: 0          # 总数（0=自动计算）
```

### 缓存配置

```yaml
cache:
  enable: true  # 是否启用缓存
```

---

## 📚 相关文档

- **CDN 鉴权**: [CDN_AUTH_GUIDE.md](CDN_AUTH_GUIDE.md)
- **多 CDN 配置**: [MULTI_CDN_GUIDE.md](MULTI_CDN_GUIDE.md)
- **媒体库统计**: [ITEMS_COUNTS_GUIDE.md](ITEMS_COUNTS_GUIDE.md)
- **功能列表**: [FEATURES.md](FEATURES.md)

---

## ⚙️ 配置文件位置

### Docker 部署

配置文件挂载位置（docker-compose.yml）:
```yaml
volumes:
  - ./config.yml:/app/config.yml
```

### 本地运行

```bash
./go-emby302 -dr /path/to/config/dir
```

配置文件应放在 `/path/to/config/dir/config.yml`

---

## 🔍 验证配置

### 1. 检查配置文件格式

```bash
# 使用 yamllint 检查（如果安装了）
yamllint config.yml

# 或者直接启动查看日志
docker compose up
```

### 2. 查看启动日志

```bash
docker compose logs -f
```

**正确启动**:
```
[INFO] 正在初始化本地目录树模块...
[INFO] 正在初始化路由规则...
[INFO] 路由规则初始化完成
[INFO] 正在启动服务...
```

**配置错误**:
```
[ERROR] 初始化配置文件失败: xxx
```

---

## ⚠️ 常见错误

### 1. YAML 格式错误

```
错误: yaml: line 10: mapping values are not allowed in this context
```

**原因**: YAML 缩进错误或格式不正确

**解决**: 使用支持 YAML 的编辑器（如 VS Code）编辑

---

### 2. 路径配置错误

```
错误: strm.cdns[0].base 不能为空
```

**原因**: 必填配置项为空

**解决**: 检查所有必填项是否已填写

---

### 3. 鉴权密钥错误

```
错误: strm.cdns[0].private-key 不能为空（鉴权类型: goedge）
```

**原因**: 启用了鉴权但没有配置密钥

**解决**: 填写 `private-key` 或将 `type` 改为 `none`

---

## 💡 最佳实践

### 1. 配置文件管理

```bash
# 备份配置
cp config.yml config.yml.backup

# 使用版本控制（注意不要提交密钥！）
git add config.yml
# 或将配置文件加入 .gitignore
echo "config.yml" >> .gitignore
```

### 2. 敏感信息保护

```yaml
# ❌ 不要提交到公开仓库
private-key: "my_secret_key_123"

# ✅ 使用环境变量或密钥管理工具
```

### 3. 配置注释

```yaml
# ✅ 好的做法 - 添加注释说明用途
emby:
  # 主 Emby 服务器（内网地址）
  host: http://192.168.1.100:8096

# ❌ 不好的做法 - 没有注释
emby:
  host: http://192.168.1.100:8096
```

---

## 🎯 选择指南

### 使用 config.example.yml 如果:

- ✅ 你是新手用户
- ✅ 只使用一个 CDN
- ✅ 需要快速上手
- ✅ 配置比较简单

### 使用 config_multi_cdn.example.yml 如果:

- ✅ 需要配置多个 CDN
- ✅ 不同媒体使用不同 CDN
- ✅ 需要参考复杂配置示例
- ✅ 想了解所有配置选项

---

**推荐**: 新手先使用 `config.example.yml`，熟悉后再参考 `config_multi_cdn.example.yml` 进行高级配置。

---

**最后更新**: 2026-01-08
