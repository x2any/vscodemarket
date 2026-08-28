# Implementation Plan: VS Code Download Hub

**Branch**: `001-vscodemarket-hub` | **Date**: 2026-08-28 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/001-vscodemarket-hub/spec.md`

## Summary

公网辅助下载站,辅助内网用户获取与本地 VSCode 客户端配套的下载物。自身**不
持有二进制、不代理下载流量**;后端只做"找直链 + 聚合 + 埋点 + 热榜",前端做
"输入 + 同屏展示 + 文档区 + 热榜"。技术栈:Go(net/http + 第三方 router)+ Vue 3
+ TS strict + Element Plus + Vite,持久化 SQLite + GORM,GeoIP 本地 mmdb,
Docker 单镜像部署。里程碑 M1–M11 顺序推进,每步可演示。

## Technical Context

**Language/Version**: Go ≥ 1.22(后端)、TypeScript strict(前端)
**Primary Dependencies**:
- 后端:net/http 标准库 + chi(默认选择,gin 备选)
- 前端:Vue 3 + Element Plus + Vite
- 持久化:GORM + SQLite(`gorm.io/gorm` + `gorm.io/driver/sqlite`)
- GeoIP:`oschwald/geoip2-golang`,本地 mmdb

**Storage**: SQLite,文件 `data/vscodemarket.db`,Docker 数据卷挂载
**Testing**:
- 后端:`testing` + `httptest`,核心契约单测
- 前端:Vitest(单元) + Playwright(E2E 覆盖 S1/S3/S4/S6)

**Target Platform**: Linux 容器(x86_64),公网部署
**Project Type**: web-service(backend + frontend,前后端分仓库根目录)
**Performance Goals**:
- 主路径 P1(Stable 双链接):端到端 ≤ 5s
- 扩展搜索:平均 ≤ 3s(由 Marketplace 响应上限决定)
- 热榜查询:≤ 500ms(本地 SQLite 聚合)

**Constraints**:
- 不引入 cron daemon / sidecar / service worker(NFR8)
- 不引入 ORM/DB driver 之外的"重型基础设施"
- 所有下载链接必须官方域名
- 辅助脚本零网络动作

**Scale/Scope**: 公网匿名访客,90 天事件保留,Top 10 热榜

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| 原则 | 满足情况 |
|------|---------|
| I. 公网应用属性 | ✅ 不收集 PII;事件表只存脱敏 IP + CountryCode + 行为聚合 |
| II. 官方直链原则 | ✅ 后端仅返回 302 至官方 URL;无预拉取 / 缓存 / 镜像 |
| III. 辅助脚本零网络原则 | ✅ 文档区与 Tooltip 脚本不含 wget/curl/Invoke-WebRequest |
| IV. 极简立场 | ✅ 无用户系统、无验证码、无 service worker / cron daemon / sidecar |

无违规项。

## Project Structure

### Documentation (this feature)

```text
specs/001-vscodemarket-hub/
├── plan.md              # 本文件
├── research.md          # Phase 0
├── data-model.md        # Phase 1
├── quickstart.md        # Phase 1
├── contracts/           # Phase 1
│   ├── api.md           # HTTP 接口契约
│   └── events.md        # 埋点事件契约
└── tasks.md             # Phase 2(由 /speckit-tasks 产出,本阶段不写)
```

### Source Code (repository root)

```text
backend/
├── cmd/server/main.go           # 入口 + 路由 + 启动期 sweep
├── internal/
│   ├── upstream/                # 官方 API 适配层(client / server / extension)
│   │   ├── client.go
│   │   ├── server.go
│   │   └── extension.go
│   ├── handlers/                # HTTP handlers
│   │   ├── version.go           # M2 客户端/服务端版本清单
│   │   ├── ua.go                # M3 UA 推断
│   │   ├── extension.go         # M4 扩展搜索 / 版本列表
│   │   ├── event.go             # M8 埋点接收
│   │   └── trending.go          # M9 热榜
│   ├── storage/                 # M8 GORM + SQLite
│   │   ├── db.go
│   │   └── event_repo.go
│   ├── geoip/                   # M8 GeoIP2 加载 + UNKNOWN 兜底
│   │   └── geoip.go
│   └── sweeper/                 # M8 90 天过期清理(进程内定时器,非 cron daemon)
│       └── sweeper.go
├── go.mod
└── go.sum

frontend/
├── src/
│   ├── main.ts
│   ├── App.vue
│   ├── router/                  # Vue Router
│   ├── pages/
│   │   ├── HomePage.vue         # 版本查询主路径(M2)
│   │   ├── ExtensionSearch.vue  # M4
│   │   ├── ExtensionDetail.vue  # M5
│   │   ├── VersionList.vue      # M6
│   │   ├── Trending.vue         # M9
│   │   └── Docs.vue             # M7 五主题文档
│   ├── components/
│   │   ├── VersionForm.vue
│   │   ├── DownloadLinkCard.vue
│   │   ├── DocsBlock.vue
│   │   └── ScriptSnippet.vue    # M7 辅助脚本片段(零网络)
│   ├── i18n/
│   │   ├── zh-CN.ts
│   │   └── en-US.ts
│   └── api/                     # 后端 HTTP 客户端封装
│       └── http.ts
├── index.html
├── vite.config.ts
├── tsconfig.json
└── package.json

data/
└── vscodemarket.db              # SQLite 数据卷挂载点(运行时生成)

docker-compose.yml
Dockerfile.backend
Dockerfile.frontend
nginx.conf                       # 前端容器内静态托管配置
README.md
```

**Structure Decision**: 采用 Web 应用分仓(backend/ + frontend/ 仓库根并列),
符合"前后端分项目 + 公网部署"的实际形态。ponytail:若后期需多服务拆分,
再做 monorepo 改造;当前 YAGNI。

## Milestones

详见 [research.md §3 Milestone Breakdown](./research.md),摘要:

- **M1 骨架**:Go 模块 + `/healthz`、Vite + Vue 3 初始化
- **M2 版本查询**:Stable/Insider × client/server 直链(M2)
- **M3 UA 推断**:表单默认值
- **M4 扩展搜索**:Marketplace 透传
- **M5 适配版本展示**:`engines.vscode` 显式标注
- **M6 版本列表与筛选**:Channel × PlatformBuild 笛卡尔积
- **M7 文档 + i18n + 辅助脚本**:五主题,中英双语,零网络脚本
- **M8 持久化 + 埋点 + GeoIP + 过期清理**
- **M9 热榜**:三物 × 三时间窗,Top 10
- **M10 Docker 部署**:multi-stage,数据卷
- **M11 测试 + 隐私**:单测 / E2E / README GDPR

## Complexity Tracking

无违规项。表格留作未来追加使用。

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| (无)      |            |                                       |