# Research: VS Code Download Hub

**Date**: 2026-08-28 | **Plan**: [plan.md](./plan.md) | **Spec**: [spec.md](./spec.md)

## 1. Technical Decisions

### 1.1 Router(Go)

- **Decision**: chi(`github.com/go-chi/chi/v5`)
- **Rationale**: net/http 兼容、轻量(零依赖少)、中间件链直观、URL 参数解析够用;
  与宪法"极简立场"一致(不引入 gin 的 binding/validation 全套)。
- **Alternatives**:
  - gin:功能多但默认带 binding/json 校验,有过度抽象之嫌;留作备选。
  - stdlib `http.ServeMux`(Go 1.22+ 已支持 method+pattern):简单场景够用,
    但 M11 起中间件(metrics / recover / CORS)累积后,chi 更省样板。

### 1.2 前端 UI 库

- **Decision**: Element Plus
- **Rationale**: 站点以表单 + 表格 + Tooltip 为主,Element Plus 覆盖度足够;
  组件体积可接受(按需引入)。`Vue 3` 生态成熟。
- **Alternatives**:
  - Naive UI:更轻,但生态相对小,本项目不需要其主题定制能力。
  - AntDV:设计语言偏后台密集型表单,与公网下载站风格不契合。
  - 自研:违反宪法极简立场(无 abstract for later)。

### 1.3 持久化

- **Decision**: SQLite + GORM(ADR-0011)
- **Rationale**: 事件表 schema 简单,SQLite 单文件 + 数据卷即满足 NFR2 90 天保留;
  GORM 提供 migration 与查询构造,但不绑死 ORM 抽象(查询走原生 SQL)。
- **Alternatives**:
  - 纯 `database/sql` + 自写迁移:更极简,但 M11 单测与 schema 演进成本上升。
  - Postgres:违反极简立场(部署额外组件)。

### 1.4 GeoIP

- **Decision**: `oschwald/geoip2-golang` + 本地 mmdb,启动期加载
- **Rationale**: NFR3 要求 GeoIP 失败回退 `UNKNOWN` 且不阻塞;mmdb 离线,
  启动期一次性打开句柄,handler 内 lookup。
- **Alternatives**:
  - 在线 HTTP 解析服务:违反"无重型基础设施"。
  - 不做 GeoIP:违反 FR-013(本期预留字段)。

### 1.5 过期清理

- **Decision**: 进程内 `time.Ticker`(默认 24h),启动时先 sweep 一次
- **Rationale**: 宪法禁止 cron daemon / sidecar。ponytail:进程内定时器在
  长时间无流量时会因进程退出而错过执行;M11 之后若需更稳健,可考虑
  docker-compose `restart: always` 配合启动 sweep,但**不引入** cron daemon。
- **Alternatives**:
  - 系统 cron:违反 NFR8 / 极简立场。
  - GitHub Actions 周期触发外部清理:违反"自包含"。

### 1.6 路由命名

- **Decision**: 后端前缀 `/api/v1`,前端通过 nginx 反代 `/api/` 到后端
- **Rationale**: 单镜像部署 + 同源请求,避免 CORS 复杂度。
- **Alternatives**: CORS + 独立域名:违反极简立场。

## 2. Resolved NEEDS CLARIFICATION

无 — 用户输入已固化全部技术栈。

## 3. Milestone Breakdown

| ID | 范围 | 关键文件 | 验收 |
|----|------|---------|------|
| M1 | 后端 Go 模块 + `/healthz`;前端 Vite + Vue 3 + TS + Element Plus | `backend/cmd/server/main.go`、`frontend/src/main.ts` | `docker compose up` 通,前端可访问 `/healthz` |
| M2 | 客户端/服务端版本清单查询(透传官方) | `backend/internal/upstream/client.go`、`HomePage.vue` | 输入 `1.94.1` 拿到 client + server 直链 |
| M3 | UA 推断端点 + 表单默认填值 | `backend/internal/handlers/ua.go`、`VersionForm.vue` | macOS Chrome 默认 `darwin-arm64` |
| M4 | 扩展搜索 + 版本列表(透传 Marketplace) | `backend/internal/handlers/extension.go`、`ExtensionSearch.vue`、`ExtensionDetail.vue` | 搜 `python` 出列表,点 ms-python.python 见所有版本 |
| M5 | 适配版本展示 `engines.vscode` | `ExtensionDetail.vue` | UI 可直接判断某版本能否装 |
| M6 | 版本列表页 + Channel/Platform/Architecture 筛选 + 分页 | `VersionList.vue` | 组合筛选可用 |
| M7 | 五主题文档区 + 中英 i18n + 辅助脚本 | `Docs.vue`、`i18n/`、`ScriptSnippet.vue` | 中英切换无 404;脚本不含下载命令 |
| M8 | GORM + SQLite + 埋点 + GeoIP + 进程内 sweep | `backend/internal/storage/`、`backend/internal/geoip/`、`backend/internal/sweeper/` | DB 可见事件;IP 解析国家码;重启不丢 |
| M9 | 热榜查询 + 三物 × 三时间窗 Top 10 | `backend/internal/handlers/trending.go`、`Trending.vue` | UI 出现 Top 10,切换时间维度数据随之变 |
| M10 | Docker multi-stage + docker-compose + 数据卷 | `Dockerfile.backend`、`Dockerfile.frontend`、`docker-compose.yml` | 一键通,埋点数据持久化 |
| M11 | 单测 + Playwright E2E(S1/S3/S4/S6)+ README GDPR | `backend/internal/**/*_test.go`、`frontend/tests/e2e/`、`README.md` | CI 全绿,合规文本就位 |

## 4. Risks & Mitigations

| 风险 | 缓解 |
|------|------|
| 官方 API 变更 | 仅在 `backend/internal/upstream/` 适配层处理,handler 不感知 |
| 官方 CDN 抖动 | 前端 Loading + 明确错误文案(Edge Cases 已列) |
| 扩展 API 分页 | 透传,前端无限滚动 / 分页参数原样转发 |
| SQLite 文件膨胀 | 进程内 sweep 每 24h 清理 >90 天事件;M10 监控文件大小(可选) |
| GeoIP db 缺失 | 启动期 `LoadOrNil`,handler 内 nil-safe 返回 `UNKNOWN` |
| 进程内 sweep 长周期不执行 | docker-compose `restart: always` 保证每次启动先 sweep;**不引入** cron daemon |
| Marketplace 默认排序被改 | FR-006 直接透传;UI 不做二次重排,避免与官方策略偏差 |
| 行为埋点写入失败阻塞主路径 | handler 内 try/catch + 异步队列;FR-012 失败降级,主请求 200 |
| 多语言文案漂移 | i18n key 集中字典;M11 单测断言 zh-CN/en-US 同 key 集合相等 |

## 5. ADR Backlog(由 plan 阶段锁定,后续不再讨论)

- ADR-0001 服务端 commit hash 与 Stable 客户端严格绑定
- ADR-0004 直链与不代理
- ADR-0005 技术栈固化(本 plan 落地)
- ADR-0006 客户端频道与 PlatformBuild 笛卡尔积
- ADR-0008 文档形态(Tooltip + 内嵌)
- ADR-0009 扩展搜索排序(Marketplace 默认)
- ADR-0010 行为埋点字段
- ADR-0011 持久化(SQLite + GORM)
- ADR-0012 极简立场(无 cron / 无 service worker / 无 sidecar)
- ADR-0013 GeoIP(mmdb 离线)