# VS Code Download Hub

公网辅助下载站,辅助内网用户获取与本地 VSCode 客户端配套的官方下载物。
自身**不持有二进制、不代理下载流量**;只做"找直链 + 聚合 + 埋点 + 热榜"。

## Quick start

```bash
docker compose up -d --build
# 浏览器: http://localhost:8080
# 后端:   http://localhost:8081/api/v1/healthz
```

数据卷 `data/`(由 docker compose 创建的命名卷)持久化 SQLite 文件
`vscodemarket.db`;删除容器不会丢事件。GeoIP mmdb 可选放入数据卷根目录,
命名为 `GeoLite2-Country.mmdb`,缺失时所有事件 `countryCode = UNKNOWN`。

## 验证

```bash
bash scripts/smoke.sh   # 7 个端到端场景
cd frontend && pnpm test # i18n key 相等 + 零网络脚本静态扫描
```

## 端点速览

- `POST /api/v1/versions/lookup` — 主路径,同屏返回 client + vscode-server 直链
- `POST /api/v1/ua/infer` — UA 推断
- `GET  /api/v1/extensions/search?q=` — Marketplace 搜索(默认排序)
- `GET  /api/v1/extensions/{pub}/{name}/versions[/{ver}]` — 扩展版本列表 / 单版本
- `GET  /api/v1/releases?channel=&platform=&architecture=&page=&pageSize=` — 版本浏览
- `POST /api/v1/events` — 行为埋点,失败降级(仍 202)
- `GET  /api/v1/trending?targetType=&window=` — 热榜,Top 10,空态 OK

详细契约见 `specs/001-vscodemarket-hub/contracts/`。

## 隐私声明 / GDPR / 个保法

本服务面向公网访客部署,数据采集遵循以下原则:

- **不收集 PII**:不记录原始 IP、不存 UA、不写 Cookie、不设追踪像素、不引入第三方分析 SDK。
- **最小化**:行为事件仅记录 `eventType / targetType / targetIdentifier / platform /
  architecture / channel / countryCode(国家码)/ createdAt`,无任何可识别自然人的字段。
- **GeoIP 折中**:若 mmdb 缺失或解析失败,事件写入 `countryCode = "UNKNOWN"`,不阻塞请求。
- **保留 90 天**:进程内 sweep 每 24h 删除 `created_at < now() - 90d` 的事件行,符合
  FR-014 / SC-006。
- **第三方资源**:本站仅返回 Microsoft 官方域名的直链,不代理、不缓存、不镜像。
  访客点击下载后由浏览器与 Microsoft 建立连接,本站不参与。
- **不向境外传输**:除调用 Microsoft / VSCode 官方 API 之外,本站不与任何第三方
  服务通信,无 CDN、无分析、无社交分享追踪。

如对数据处理有疑问,请通过仓库 Issue 联系我们;我们会在 30 天内回复。

## 范围外(明确不做)

- 鉴权 / 用户系统 / 登录 / UID
- 限频 / 反作弊 / 验证码
- VSCodium / Cursor / Windsurf / Coder / Gitpod 等替代客户端
- 下载代理 / 打包脚本 / 批量下载
- 独立文档站 / Wiki 后台

详细 ADR 见 `specs/001-vscodemarket-hub/plan.md` 与 `specs/001-vscodemarket-hub/research.md`。