# Quickstart: VS Code Download Hub

**Date**: 2026-08-28 | **Spec**: [spec.md](./spec.md) | **Plan**: [plan.md](./plan.md)

本指南用于里程碑结束后回归"端到端是否可用"。每个场景均可独立运行。

## Prerequisites

- Docker ≥ 24 + docker compose v2
- 端口 8080(前端)+ 8081(后端,可由 docker-compose 暴露)空闲
- 公网可达 `update.code.visualstudio.com` 与 `marketplace.visualstudio.com`
- GeoIP mmdb 缺失时,服务仍可启动,所有事件 `countryCode = UNKNOWN`

## Bootstrap

```bash
git clone <repo>
cd vscodemarket
docker compose up -d --build
# 等约 30s,后端 /healthz 通过即就绪
curl -fsS http://localhost:8081/api/v1/healthz
```

`docker compose` 会拉起:
- `backend`:Go 单二进制,监听 :8081
- `frontend`:nginx 静态托管,监听 :8080,反代 `/api/` → backend:8081
- 数据卷 `data/` 挂载到 backend 容器内,SQLite 文件 `data/vscodemarket.db` 持久化

## Scenarios

### S1: Stable 版本主路径(同屏 client + server)

```bash
curl -fsS -X POST http://localhost:8081/api/v1/versions/lookup \
  -H 'Content-Type: application/json' \
  -d '{"channel":"stable","version":"1.94.2","platform":"darwin","architecture":"arm64"}' | jq
```

**期望**:`client.downloadUrl` 与 `server.downloadUrl` 均存在,`server.commitHash` 非空。
浏览器打开 `http://localhost:8080/`,输入 `1.94.2`,勾选 darwin/arm64,提交 →
页面同时渲染两块 DownloadLinkCard,server 卡片显式展示 commit hash。

### S2: Insider 路径(无 server)

```bash
curl -fsS -X POST http://localhost:8081/api/v1/versions/lookup \
  -H 'Content-Type: application/json' \
  -d '{"channel":"insider","version":"...","platform":"darwin","architecture":"arm64"}' | jq
```

**期望**:`server` 字段缺失,HTTP 200。

### S3 / S4: 扩展搜索

```bash
# 列表
curl -fsS 'http://localhost:8081/api/v1/extensions/search?q=python&page=1&pageSize=10' | jq
# 详情(全部版本)
curl -fsS http://localhost:8081/api/v1/extensions/ms-python/python/versions | jq
# 直达单版本
curl -fsS http://localhost:8081/api/v1/extensions/ms-python/python/versions/2024.20.0 | jq
```

**期望**:`enginesVscode` 字段非空。

### S5: UA 推断

```bash
curl -fsS -X POST http://localhost:8081/api/v1/ua/infer \
  -H 'Content-Type: application/json' \
  -d '{"userAgent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"}' | jq
```

**期望**:`{"platform":"darwin","architecture":"arm64","confidence":"HIGH"}`。

### S6: 热榜

```bash
# 触发埋点(可重复多次以提高计数)
for i in $(seq 1 5); do
  curl -fsS -X POST http://localhost:8081/api/v1/events \
    -H 'Content-Type: application/json' \
    -d '{"eventType":"SEARCH","targetType":"EXTENSION","targetIdentifier":"ms-python.python"}'
done
# 查询
curl -fsS 'http://localhost:8081/api/v1/trending?targetType=EXTENSION&window=24h' | jq
```

**期望**:`items` 包含 `ms-python.python`,count ≥ 5。
切换 `window=7d` / `window=30d` 应返回相同或聚合结果。

### S7: i18n 与辅助脚本

- 浏览器打开首页右上角语言切换器,切到 EN → 所有按钮、表单 label、文档区文案切换为英文。
- 切回 zh-CN → 全中文。
- DevTools Console 执行:`fetch('/').then(r=>r.text()).then(t=>console.log(/wget|curl/.test(t)))`
  → 期望输出 `false`(源码与文档不含下载命令字面量)。

## Validation Tests

- **后端单测**:`cd backend && go test ./...` — 覆盖
  - 版本号解析 / UA 推断 / 扩展搜索透传(用 `httptest` mock 官方)/
    热榜 SQL 聚合 / 埋点入库(含 countryCode fallback)
- **前端 E2E**:`cd frontend && pnpm playwright test` — 覆盖 S1 / S3 / S4 / S6
- **静态扫描**:`grep -RE 'wget|curl|Invoke-WebRequest' frontend/src backend/internal`
  → 期望无匹配(`frontend/src/components/ScriptSnippet.vue` 的注释与占位提示
  内出现的"下载"必须为纯文字描述,不可执行)。

## Cleanup

```bash
docker compose down            # 保留数据卷
docker compose down -v         # 彻底清理(SQLite 数据卷一并删除)
```

## Troubleshooting

| 现象 | 排查 |
|------|------|
| `/healthz` 502 | `docker compose logs backend`,确认 Go 进程未 panic;数据卷权限 |
| 扩展搜索 502 / 超时 | 公网到 `marketplace.visualstudio.com` 是否可达 |
| 热榜空 | 确认先触发埋点;SQLite 文件是否被 sweep 误删(检查 created_at) |
| GeoIP 总是 UNKNOWN | 启动日志是否含 `geoip: mmdb loaded`;否则 GeoIP 文件缺失,行为合规 |
| 前端 404 静态资源 | `docker compose logs frontend`,确认 nginx 反代 `/api/` 正确 |