# /plan 提示语

> 用法:将本文件正文整块喂给 `/plan` 命令。
> 引用基础:`constitution.md`、`spec.md`、`CONTEXT.md`、`docs/adr/0001~0013`。

---

为项目 vscodemarket 制定实施计划 (plan.md)。

原则:极简、增量、可运行优先;每步可演示;不预先引入抽象层、interface、factory、config;后续阶段不需则不写。

## 技术栈(ADR-0005,实施期固化)

- 后端:Go ≥1.22,net/http 标准库 + 第三方 router(chi 或 gin,首 PR 决定)
- 前端:Vue 3 + TypeScript strict + Element Plus + Vite
- 包管理:后端 go mod,前端 pnpm
- 交付:Docker 单镜像(multi-stage),前端 nginx 静态托管
- 持久化:SQLite + GORM(ADR-0011),ORM/DB driver 不在本章程约束
- GeoIP:本地 GeoIP2 mmdb(ADR-0013),启动期加载,失败兜底 UNKNOWN
- 不引入:React/Next/Nuxt/SSR/AntDV/Vuetify/Naive UI,以及任何 cron daemon / sidecar / service worker

## M1 骨架
- 后端:Go 模块初始化,HTTP 路由占位,`/healthz`
- 前端:Vite + Vue 3 + TS + Element Plus 初始化
- 验收:`docker compose up`,前端可访问 `/healthz`

## M2 版本查询主路径
- 后端:客户端/服务端版本清单查询(透传官方 API)
- 前端:版本查询表单(Channel/Platform/Architecture),提交后同屏展示 client + server 链接
- 验收:`1.94.1` 能拿到对应链接

## M3 UA 推断
- 后端:UA 推断端点
- 前端:表单加载时调用一次填默认值,允许覆盖
- 验收:macOS Chrome 默认 `darwin-arm64`

## M4 扩展搜索主路径
- 后端:扩展搜索 + 版本列表(透传 Marketplace)
- 前端:搜索页,单匹配直跳详情,多匹配列表
- 验收:搜 `python` 出列表,点 ms-python.python 见所有版本

## M5 适配版本展示
- 前端:扩展详情页每个版本展示 `engines.vscode`
- 验收:UI 可直接判断某版本能否装在本地 VSCode

## M6 版本列表与筛选
- 前端:版本列表页,Channel/Platform/Architecture 组合筛选,分页
- 验收:Stable/Insider × 各 PlatformBuild 可组合筛选

## M7 文档区 + i18n + 辅助脚本
- 前端:五主题文档区;中英 i18n;辅助脚本片段嵌入
- 验收:中英切换无 404;脚本不含下载命令

## M8 持久化层 + 行为埋点
- 后端:引入 GORM + SQLite(数据卷 `data/vscodemarket.db`)
- 后端:埋点接收端点(search / download click),记录 IP + CountryCode + target + created_at
- 后端:GeoIP2 mmdb 加载(失败兜底 UNKNOWN)
- 后端:90 天过期清理(启动 sweep + cron)
- 验收:触发搜索/下载后,DB 内可见事件;IP 解析国家码;重启后数据不丢

## M9 热榜
- 后端:热榜查询端点(物维度 × 时间维度),Top 10
- 前端:热榜区块,三物维度各自一个面板;时间维度切换(24h/7d/30d)
- 验收:埋点 + 模拟访问后,UI 出现 Top 10,切换时间维度数据随之变化

## M10 Docker 部署
- 后端 multi-stage(单二进制)、前端 multi-stage(nginx)
- `docker-compose.yml`(含 SQLite 数据卷)
- 验收:`docker compose up` 一键通,埋点数据持久化在卷内

## M11 测试 + 隐私
- 后端:核心契约单测(版本号解析、UA 推断、扩展搜索透传、热榜聚合、埋点入库)
- 前端:Vitest + Playwright,E2E 覆盖 S1/S3/S4/S6
- README:GDPR/隐私声明、数据保留 90 天说明
- 验收:CI 全绿,合规文本就位

## 范围外(再次声明)
鉴权、用户系统、限频/验证码、DB 同步外部数据源、下载/打包脚本、VSCodium、Coder/Gitpod、独立文档站。

## 风险与缓解
- 官方 API 变更 → 仅在适配层处理
- 官方 CDN 抖动 → 前端 Loading + 明确错误
- 扩展 API 分页 → 透传,前端无限滚动
- SQLite 文件膨胀 → 90 天过期清理 + 监控
- GeoIP db 缺失 → 兜底 UNKNOWN,不阻塞