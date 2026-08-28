# /spec 提示语

> 用法:将本文件正文整块喂给 `/spec` 命令。
> 引用基础:`constitution.md`、`CONTEXT.md`、`docs/adr/0001~0013`。

---

为项目 vscodemarket 编写功能规格 (spec.md)。

## 1. 概述
公网辅助下载站,辅助内网用户获取与本地 VSCode 客户端配套的下载物。自身不持有二进制、不代理下载流量。会记录匿名 IP + 行为事件用于热榜与地区统计。

**支持范围**:
- 客户端:Microsoft 官方 Stable 与 Insider 两个频道(ADR-0006)
- 服务端:Microsoft 官方 vscode-server(ADR-0001),与 Stable 客户端 commit hash 严格绑定(ADR-0001/ADR-0004)
- 客户端/服务端 PlatformBuild:Platform ∈ {windows, linux, darwin},Architecture ∈ {x86_64, arm64, armv7},笛卡尔积见 ADR 列表(ADR-0006)
- 扩展:VSCode Marketplace 全部扩展(ADR-0009 采用其默认排序)

**外部数据源**:客户端/服务端/扩展元数据均从 VSCode 官方 API 实时拉取,本站不缓存。

## 2. 角色与场景
单一匿名访客角色。场景:
- S1 输入 Stable 版本号 → 同屏展示 client + vscode-server 下载链接(含 commit hash)
- S2 输入 Insider 版本号 → 展示 client 下载链接
- S3 输入扩展名 → 多匹配列表 → 进入适配版本列表与下载
- S4 输入扩展名 + 版本号 → 直接单一版本下载页
- S5 浏览版本列表(支持 Channel / Platform / Architecture 筛选)
- S6 浏览热榜(客户端 / 服务端 / 扩展 × 24h/7d/30d,Top 10)
- S7 文档(Tooltip + 内嵌区,中英双语)

## 3. 功能需求
- F1 版本查询:Stable/Insider 区分;版本不存在时明确错误。
- F2 Platform/Architecture 表单:默认按 UA 推断,可手动改。
- F3 vscode-server 下载:必须展示对应 Stable 客户端 commit hash。
- F4 扩展搜索:Marketplace 默认排序;多结果列表,单结果详情。
- F5 扩展版本过滤:指定版本号只返该版本;否则返回所有版本(含 engines.vscode / 发布时间)。
- F6 适配版本展示:每个 ExtensionVersion 显式 engines.vscode。
- F7 内嵌文档区:五个主题(客户端安装 / 服务端安装 / 服务端启动与连接 / 离线导入扩展 / 排错 FAQ),双语。
- F8 辅助脚本:Tooltip/页面内嵌,仅含运行/安装命令,严禁 wget/curl。
- F9 行为埋点:搜索提交、下载点击两类事件入库,字段见 ADR-0010。
- F10 热榜:客户端 / 服务端 / 扩展 三物维度,各支持 24h/7d/30d 时间维度;Top 10;数据源为本地事件表(聚合查询)。
- F11 地区维度(预留字段,不本期暴露 UI):CountryCode 入库,用于后期"按地区看热榜",本期 UI 不渲染。
- F12 文档形态:不设独立文档站,文档以 Tooltip + 内嵌文档区承载(ADR-0008)。

## 4. 接口契约(语义级,无技术栈)
- 客户端版本清单查询
- 服务端版本清单查询
- 扩展搜索
- 扩展版本列表
- UA 推断
- 事件埋点接收(搜索 / 下载点击)
- 热榜查询(物维度 + 时间维度)

## 5. 非功能
- NFR1 不引入用户系统、不引入登录
- NFR2 行为事件保存 90 天,过期清理
- NFR3 GeoIP 解析失败时 CountryCode=UNKNOWN,不阻塞请求
- NFR4 所有下载链接必须官方域名
- NFR5 UI 中英双语
- NFR6 响应式(桌面/平板/移动)
- NFR7 持久化层抽象由 ADR 决定(spec 不绑定具体方言)
- NFR8 不引入 ORM/DB driver 之外的"重型基础设施":无 service worker、无 cron daemon、无 sidecar(与 ADR-0012 一致的极简立场)

## 6. Done 定义
- F1~F10 均有最小可演示路径(F11/F12 字段入库即可,UI 不要求)
- 核心契约单测覆盖(版本号解析、UA 推断、扩展搜索透传、热榜聚合)
- E2E 覆盖 S1/S3/S4/S6
- i18n 完整
- 提供 Docker 一键部署 与 本地直接运行 两条等价路径(任一可启动完整功能,见 plan M10)
- README 含 GDPR/隐私声明

## 7. 范围外(明确不做)
- 鉴权、用户系统、UID、Cookie 跨站追踪
- 反作弊/限频/验证码(ADR-0012)
- 独立文档站、Wiki 后台
- VSCodium/Cursor/Windsurf、Coder/Gitpod(ADR-0001/0006)
- 下载代理/打包脚本/批量下载(ADR-0004)
- 后端 ORM/DB driver 选型(ADR-0011,constitution 不锁定)