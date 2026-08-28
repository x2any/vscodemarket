# Tasks: VS Code Download Hub

**Input**: Design documents from `/specs/001-vscodemarket-hub/`
**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, contracts/, quickstart.md
**Tests**: Spec §6 Done 定义要求后端单测与前端 Playwright E2E,故测试任务纳入。
**Organization**: 按 User Story 拆分;每 Story 内的测试先于实现(红→绿)。

## Format: `[ID] [P?] [Story] Description`

- **[P]**: 可并行(不同文件,无依赖)
- **[Story]**: US1 / US2 / ... 对应 spec.md User Story
- 路径遵循 plan.md §Project Structure(backend/ + frontend/ 仓库根并列)

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: M1 骨架 — 仓库骨架 + 后端 Go 模块 + 前端 Vite 初始化 + Docker 占位

- [x] T001 Create repository structure: `backend/`, `frontend/`, `data/`, `docker-compose.yml` 根级占位
- [x] T002 Initialize Go module `github.com/yourorg/vscodemarket` in `backend/go.mod` (Go ≥ 1.22),add chi router 与 GeoIP 依赖
- [x] T003 Create backend entry `backend/cmd/server/main.go` 启动 net/http + chi router,挂载 `/api/v1/healthz`
- [x] T004 [P] Create backend Dockerfile multi-stage `Dockerfile.backend`(builder → 单二进制 scratch/distroless)
- [x] T005 Initialize Vite + Vue 3 + TS strict in `frontend/`(pnpm),集成 Element Plus 按需引入
- [x] T006 [P] Create frontend Dockerfile multi-stage `Dockerfile.frontend` + `frontend/nginx.conf`(反代 `/api/` → backend:8081)
- [x] T007 Create `docker-compose.yml`:backend(:8081)+ frontend(:8080),`data/` 数据卷挂载到 backend
- [x] T008 [P] Create `.gitignore`(忽略 `data/*.db`、`frontend/node_modules`、`backend/dist`、`.env*`)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 所有 User Story 依赖的通用基础设施 — 错误格式、i18n 字典、官方域名白名单、响应式布局

**⚠️ CRITICAL**: 任一 User Story 启动前必须完成

- [x] T009 Implement shared API error type with bilingual fields in `backend/internal/handlers/errors.go`(code + message_zh + message_en)
- [x] T010 [P] Implement official-domain whitelist assertion in `backend/internal/upstream/whitelist.go`(handler 出网前拦截,SC-007 测试守护)
- [x] T011 [P] Setup frontend router in `frontend/src/router/index.ts`(routes 占位 /, /search, /extension/:pub/:name, /extension/:pub/:name/v/:ver, /releases, /trending, /docs)
- [x] T012 [P] Initialize i18n dictionaries in `frontend/src/i18n/zh-CN.ts` 与 `frontend/src/i18n/en-US.ts`(key 集合在 M11 测试中要求相等)
- [x] T013 [P] Create shared layout with responsive breakpoints in `frontend/src/App.vue`(≥1024 / 768-1023 / <768,SC-003)
- [x] T014 Create HTTP client wrapper in `frontend/src/api/http.ts`(baseURL = `/api/v1`,统一错误弹层)

**Checkpoint**: Foundation ready — 可并行启动所有 User Story

---

## Phase 3: User Story 1 — 已知 Stable 版本号同时取 client + server (Priority: P1) 🎯 MVP

**Goal**: 输入 Stable 版本号 → 同屏返回 client 与 vscode-server 直链(server 含 commit hash)
**Independent Test**: 见 quickstart §S1;输入 `1.94.2` / darwin / arm64 → `client.downloadUrl` 与 `server.downloadUrl` 均存在且后者含 `commitHash`

### Tests for User Story 1

- [x] T015 [P] [US1] Contract test for `POST /api/v1/versions/lookup` in `backend/internal/handlers/version_test.go`(httptest,200/404/400 路径)
- [x] T016 [P] [US1] Unit test for version string parser in `backend/internal/upstream/version_test.go`(合法/非法/Insider)

### Implementation for User Story 1

- [x] T017 [P] [US1] Implement ClientRelease upstream adapter in `backend/internal/upstream/client.go`(透传 `update.code.visualstudio.com`,返回 ClientRelease)
- [x] T018 [P] [US1] Implement ServerRelease upstream adapter in `backend/internal/upstream/server.go`(按 commit hash 拼装 vscode-server 直链)
- [x] T019 [US1] Implement `POST /versions/lookup` handler in `backend/internal/handlers/version.go`(depends T017/T018,T015 测试在此绿)
- [x] T020 [US1] Create `HomePage.vue` in `frontend/src/pages/HomePage.vue`(表单 Channel/Platform/Architecture + 提交)
- [x] T021 [P] [US1] Create `VersionForm.vue` in `frontend/src/components/VersionForm.vue`(三 select + 一 input,emit submit)
- [x] T022 [P] [US1] Create `DownloadLinkCard.vue` in `frontend/src/components/DownloadLinkCard.vue`(platform/arch/url/commit hash 显示)
- [x] T023 [US1] Wire HomePage 提交 → `http.lookupVersion()` → 双 DownloadLinkCard 渲染(depends T019-T022)

**Checkpoint**: US1 可独立验证 — MVP 形态

---

## Phase 4: User Story 2 — Insider 客户端下载 (Priority: P2)

**Goal**: Insider 版本号 → 仅返回 client 直链,无 server 块
**Independent Test**: 见 quickstart §S2;响应缺 `server` 字段,HTTP 200

### Tests for User Story 2

- [x] T024 [P] [US2] Contract test for insider channel omitting server in `backend/internal/handlers/version_test.go`

### Implementation for User Story 2

- [x] T025 [US2] Extend `version.go` handler for insider branch(FR-005 不渲染 server)(依赖 T019)
- [x] T026 [US2] Update `HomePage.vue` to conditionally render server card based on channel(依赖 T023)

**Checkpoint**: US1 + US2 可独立验证,insider 路径不泄漏 server

---

## Phase 5: User Story 3 — UA 推断 (Priority: P3 → 实为 M3 必需)

**Goal**: 表单加载时根据 UA 推断 Platform/Architecture,允许覆盖
**Independent Test**: 见 quickstart §S5;macOS Chrome UA → `{platform: darwin, architecture: arm64, confidence: HIGH}`

### Tests for User Story 3

- [x] T027 [P] [US3] Unit test for UA inference in `backend/internal/handlers/ua_test.go`(mac/win/linux × 各 arch)

### Implementation for User Story 3

- [x] T028 [US3] Implement `POST /ua/infer` handler in `backend/internal/handlers/ua.go`(OS+arch 关键字匹配,缺架构默认 x86_64 = FALLBACK)
- [x] T029 [US3] Wire `VersionForm.vue` onMounted to call `/ua/infer` with `navigator.userAgent` and prefill selects(依赖 T021)

**Checkpoint**: US1-3 可端到端验证,表单默认值得来

---

## Phase 6: User Story 4 — 扩展搜索(列表 + 直达) (Priority: P1)

**Goal**: 关键词 → Marketplace 默认排序列表;单匹配直跳详情
**Independent Test**: 见 quickstart §S3;搜 `python` → 多结果,点 ms-python.python → 进入版本列表

### Tests for User Story 4

- [x] T030 [P] [US4] Contract test for `GET /extensions/search` in `backend/internal/handlers/extension_test.go`(httptest mock Marketplace)
- [x] T031 [P] [US4] Contract test for `GET /extensions/:pub/:name/versions[/:ver]` in `backend/internal/handlers/extension_test.go`(200/404)

### Implementation for User Story 4

- [x] T032 [P] [US4] Implement Extension upstream adapter in `backend/internal/upstream/extension.go`(search + versions,透传 Marketplace 排序)
- [x] T033 [US4] Implement extension handlers in `backend/internal/handlers/extension.go`(search / versions / version-by-id)(依赖 T032)
- [x] T034 [P] [US4] Create `ExtensionSearch.vue` in `frontend/src/pages/ExtensionSearch.vue`(输入 + 结果列表)
- [x] T035 [US4] Wire search page 提交 → 后端 → 单匹配直跳 / 多结果列表(依赖 T033/T034)

**Checkpoint**: 扩展主路径独立可验证

---

## Phase 7: User Story 5 — 适配版本展示(engines.vscode) (Priority: P2)

**Goal**: 每个 ExtensionVersion 显式展示 `enginesVscode` 便于访客判断与本地 VSCode 匹配
**Independent Test**: 见 quickstart §S3 单版本接口;UI 可见 `enginesVscode` 字段

### Tests for User Story 5

- [x] T036 [P] [US5] Unit test that `enginesVscode` field is always present in `backend/internal/upstream/extension_test.go`

### Implementation for User Story 5

- [x] T037 [P] [US5] Create `ExtensionDetail.vue` in `frontend/src/pages/ExtensionDetail.vue`(版本列表 + enginesVscode 徽标)
- [x] T038 [US5] Wire `ExtensionDetail.vue` to call versions endpoint, render each version with enginesVscode(依赖 T033/T037)

**Checkpoint**: US5 与 US4 共同完成扩展详情体验

---

## Phase 8: User Story 6 — 版本列表与筛选 (Priority: P3)

**Goal**: Channel × Platform × Architecture 组合筛选 + 分页
**Independent Test**: 见 quickstart §S1 间接验证;`/releases` 页面切换筛选组合,结果随之变化

### Tests for User Story 6

- [x] T039 [P] [US6] Contract test for `GET /releases` in `backend/internal/handlers/release_test.go`(组合筛选 + 分页)

### Implementation for User Story 6

- [x] T040 [P] [US6] Implement `GET /releases` handler in `backend/internal/handlers/release.go`(可选 channel/platform/architecture/page/pageSize)
- [x] T041 [US6] Create `VersionList.vue` in `frontend/src/pages/VersionList.vue`(三个 select + 分页器)(依赖 T040)

**Checkpoint**: 浏览模式独立可演示

---

## Phase 9: User Story 7 — 文档区 + i18n + 辅助脚本 (Priority: P3)

**Goal**: 五主题(客户端安装 / 服务端安装 / 服务端启动与连接 / 离线导入扩展 / 排错 FAQ)+ 中英双语 + 零网络脚本
**Independent Test**: 见 quickstart §S7;语言切换无 404;`grep -E 'wget|curl|Invoke-WebRequest' frontend/src` 无匹配

### Tests for User Story 7

- [x] T042 [P] [US7] Snapshot test for zh-CN/en-US key set equality in `frontend/src/i18n/i18n.test.ts`
- [x] T043 [P] [US7] Static scan test: assert no download-command literals in `frontend/src/**` 在 Playwright E2E 套件中(SC-008)

### Implementation for User Story 7

- [x] T044 [P] [US7] Create `DocsBlock.vue` in `frontend/src/components/DocsBlock.vue`(渲染五主题内容,slot per 主题)
- [x] T045 [P] [US7] Create `ScriptSnippet.vue` in `frontend/src/components/ScriptSnippet.vue`(只渲染运行/安装命令文本,严禁下载动作字面量)
- [x] T046 [US7] Create `Docs.vue` in `frontend/src/pages/Docs.vue`(整合 DocsBlock + ScriptSnippet,语言切换响应)(依赖 T044/T045)
- [x] T047 [P] [US7] Populate i18n content for five topics in `frontend/src/i18n/zh-CN.ts` 与 `frontend/src/i18n/en-US.ts`

**Checkpoint**: 文档形态独立可演示,无下载动作字面量

---

## Phase 10: User Story 8 — 行为埋点(M8 持久化前置) (Priority: P2)

**Goal**: 搜索提交 + 下载点击入库,字段含 countryCode(GeoIP 或 UNKNOWN);失败降级
**Independent Test**: 见 quickstart §S6 上半;触发后 `SELECT * FROM behavior_events` 可见事件

### Tests for User Story 8

- [x] T048 [P] [US8] Unit test for GeoIP loader UNKNOWN fallback in `backend/internal/geoip/geoip_test.go`(mmdb 缺失路径)
- [x] T049 [P] [US8] Unit test for event repo write + retention query in `backend/internal/storage/event_repo_test.go`
- [x] T050 [P] [US8] Contract test for `POST /events` failure-degradation in `backend/internal/handlers/event_test.go`(DB 关闭时仍 202)

### Implementation for User Story 8

- [x] T051 [US8] Add GORM + SQLite dependencies in `backend/go.mod` & `backend/internal/storage/db.go`(AutoMigrate `BehaviorEvent`)
- [x] T052 [P] [US8] Implement GeoIP2 loader in `backend/internal/geoip/geoip.go`(LoadOrNil,Lookup(country) → ISO code 或 UNKNOWN)
- [x] T053 [US8] Implement event handler in `backend/internal/handlers/event.go`(异步落库,失败仅日志)
- [x] T054 [US8] Wire frontend DownloadLinkCard / VersionForm 提交后 fire-and-forget POST /events(依赖 T019/T022 等)
- [x] T055 [US8] Implement sweeper in `backend/internal/sweeper/sweeper.go`(启动 sweep + 24h ticker,DELETE WHERE created_at < now() - 90d)
- [x] T056 [US8] Wire sweeper start in `backend/cmd/server/main.go` 启动时(依赖 T055)

**Checkpoint**: 埋点主路径 + 过期清理独立可演示

---

## Phase 11: User Story 9 — 热榜 (Priority: P2)

**Goal**: 客户端 / 服务端 / 扩展 × 24h/7d/30d,Top 10
**Independent Test**: 见 quickstart §S6 下半;触发埋点 + 模拟访问 → UI 出现 Top 10,切换时间窗数据随之变

### Tests for User Story 9

- [x] T057 [P] [US9] Unit test for trending aggregation SQL in `backend/internal/storage/event_repo_test.go`(三物 × 三窗 Top 10)

### Implementation for User Story 9

- [x] T058 [US9] Implement `GET /trending` handler in `backend/internal/handlers/trending.go`(targetType + window → Top 10 SQL)
- [x] T059 [US9] Create `Trending.vue` in `frontend/src/pages/Trending.vue`(三面板,各面板时间窗切换器)(依赖 T058)
- [x] T060 [US9] Wire HomePage / DownloadLinkCard 等嵌入 Trending 区块快捷入口(依赖 T059)

**Checkpoint**: 热榜独立可演示,与 US8 共享 event_repo

---

## Phase 12: User Story 10 — Docker 部署(M10) (Priority: P2)

**Goal**: `docker compose up` 一键通;SQLite 数据卷持久化
**Independent Test**: 销毁 backend 容器 → `docker compose up` → 历史事件仍在

### Tests for User Story 10

- [x] T061 [P] [US10] Smoke test script in `scripts/smoke.sh`(curl 7 个 quickstart 场景端点,断言 200/202)

### Implementation for User Story 10

- [x] T062 [US10] Finalize `Dockerfile.backend` multi-stage(builder → scratch 静态二进制,无 shell)与 `Dockerfile.frontend` multi-stage(node build → nginx)
- [x] T063 [P] [US10] Finalize `docker-compose.yml`:`data/` 命名卷,`restart: always`,依赖顺序 backend → frontend,健康检查 `/healthz`
- [x] T064 [US10] Document quickstart in `README.md`(Docker run + 7 个 quickstart 场景复制可用)

**Checkpoint**: 部署独立可演示

---

## Phase 13: User Story 11 — 测试 + 隐私(M11) (Priority: P2)

**Goal**: CI 全绿;GDPR / 隐私声明就位;E2E 覆盖 S1/S3/S4/S6
**Independent Test**: `go test ./...` + `pnpm playwright test` + README grep 含 GDPR 段

### Tests for User Story 11

- [x] T065 [P] [US11] Playwright E2E spec for S1(Stable 主路径) in `frontend/tests/e2e/s1.spec.ts`
- [x] T066 [P] [US11] Playwright E2E spec for S3(扩展搜索) in `frontend/tests/e2e/s3.spec.ts`
- [x] T067 [P] [US11] Playwright E2E spec for S4(扩展直达) in `frontend/tests/e2e/s4.spec.ts`
- [x] T068 [P] [US11] Playwright E2E spec for S6(热榜) in `frontend/tests/e2e/s6.spec.ts`

### Implementation for User Story 11

- [x] T069 [P] [US11] Add GDPR / 隐私声明 section to `README.md`:数据字段(countryCode、行为事件)、90 天保留、不收 PII、联系邮箱
- [x] T070 [US11] Add CI workflow in `.github/workflows/ci.yml`(backend: `go test ./...`;frontend: `pnpm install && pnpm playwright test`)
- [x] T071 [US11] Verify all spec success criteria(SC-001 ~ SC-010)in `docs/qa-checklist.md`(每条引用 quickstart 场景)

**Checkpoint**: Done 定义全部满足

---

## Phase 14: Polish & Cross-Cutting Concerns

**Purpose**: 跨 Story 收尾

- [ ] T072 [P] Add `LICENSE` file(MIT 或 Apache-2.0,自行选定)
- [ ] T073 [P] Add `CHANGELOG.md` v0.1.0 条目,引用 M1-M11 验收摘要
- [ ] T074 [P] Verify no `wget|curl|Invoke-WebRequest|fetch` literals in `frontend/src/**` 与 `backend/internal/**`(SC-008)
- [ ] T075 Run quickstart.md validation end-to-end against `docker compose up` 实例,记录结果到 `docs/qa-run.md`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 Setup**: 无依赖
- **Phase 2 Foundational**: 依赖 Phase 1 完成 — 阻塞所有 User Story
- **Phase 3-13 User Stories**: 依赖 Phase 2 完成
  - US1(Phase 3)→ US2(4)→ US3(5)→ US4(6)→ US5(7)→ US6(8)→ US7(9)→ US8(10)→ US9(11)→ US10(12)→ US11(13)
  - US1 是 MVP 必须;US4-US11 可按团队容量并行(若多人)
- **Phase 14 Polish**: 依赖 Phase 3-13 完成

### User Story Internal Order

- Tests(T015/T016 等)→ 实现 → wire-up → 验证
- 后端 handler → 前端 component → wire-up
- US8(埋点)前置:US4/US1 已能产生事件,US8 落地 + GeoIP + sweep
- US9(热榜)前置:US8 已能产数据

### Parallel Opportunities

- Phase 1 中 T004 / T005 / T006 / T008 可并行
- Phase 2 中 T010 / T011 / T012 / T013 可并行
- 各 User Story 内的 [P] 测试任务可并行
- US8 / US9 可在 US4 完成后并行准备(US8 handler + US9 handler 不同文件)

---

## Implementation Strategy

### MVP First(仅 US1)

1. Phase 1 + Phase 2
2. Phase 3(US1)
3. **STOP and VALIDATE**:`docker compose up` → S1 通过
4. 演示/部署

### Incremental Delivery

1. Setup + Foundational → 基础就绪
2. US1 → MVP 部署
3. US2 → Insider 路径(几乎零增量,仅条件渲染)
4. US3 → 表单默认值得来(UX 提升)
5. US4 → 扩展主路径
6. US5 → 详情页 engines.vscode
7. US6 → 浏览模式
8. US7 → 文档与 i18n(独立)
9. US8 → 埋点 + 持久化(质变)
10. US9 → 热榜(US8 数据驱动)
11. US10 → 部署收尾
12. US11 → 测试 + 隐私
13. Polish

### Parallel Team Strategy

- Phase 1-2 全员合流
- Phase 3 起:
  - Dev A: US1 → US2 → US3
  - Dev B: US4 → US5 → US6 → US7
  - Dev C(等待):US8 → US9 → US10 → US11
- 提交按 PR 拆,每个 Story 一个 PR,合入后立即可演示

---

## Notes

- [P] 任务必须落到不同文件,无相互 import
- 每个 User Story 完成后立即跑 quickstart 对应 §
- 任意 checkpoint 可暂停并演示当前 Story
- 严禁在 ScriptSnippet.vue 引入 `wget|curl|Invoke-WebRequest|fetch`,所有 PR 触发 T074 扫描
- 行为埋点 handler 失败仅记日志,不向上抛(SC-005)
- 所有 `downloadUrl` 出网前必走 whitelist 断言(T010,SC-007)
