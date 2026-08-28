# Feature Specification: VS Code Download Hub

**Feature Branch**: `001-vscodemarket-hub`
**Created**: 2026-08-28
**Status**: Draft
**Input**: 公网辅助下载站,辅助内网用户获取与本地 VSCode 客户端配套的下载物;自身不持有二进制、不代理下载流量;匿名 IP + 行为事件用于热榜与地区统计。

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 已知 Stable 版本号同时取客户端与服务端 (Priority: P1)

访客输入一个 Stable 客户端版本号,期望在同一页面同时拿到客户端与配套
`vscode-server` 的官方下载链接(包含服务端对应的 commit hash),以便内网部署。

**Why this priority**: 这是产品主路径,缺它站点就无核心价值。
**Independent Test**: 输入已知 Stable 版本(如 `1.94.2`),断言页面同时渲染
客户端 Platform/Architecture 直链与 vscode-server 直链,且后者 URL 携带正确的
commit hash。

**Acceptance Scenarios**:
1. **Given** 访客位于桌面 / 平板 / 移动端,**When** 提交有效 Stable 版本号并选择
   平台,**Then** 同屏显示客户端直链与 vscode-server 直链,后者明确显示 commit hash。
2. **Given** 访客提交不存在的版本号,**When** 表单提交,**Then** 显示明确"版本不存在"
   错误,不渲染下载区。
3. **Given** 访客未手动选择平台,**When** 进入页面,**Then** 默认按 UA 推断填入,可手动覆盖。

---

### User Story 2 - 已知 Insider 版本号取客户端 (Priority: P2)

访客输入 Insider 版本号,期望拿到该 Insider 构建的客户端直链。

**Why this priority**: Insider 是次主路径,但与 Stable 一样属于"客户端下载"
单一职责,不能与 P1 同优先级。
**Independent Test**: 输入 Insider 版本,断言返回客户端直链且**不**渲染 vscode-server 区。

**Acceptance Scenarios**:
1. **Given** Insider 版本号有效,**When** 提交,**Then** 仅展示客户端下载,无 vscode-server 区块。
2. **Given** Insider 版本号无效,**When** 提交,**Then** 报错"版本不存在"。

---

### User Story 3 - 扩展搜索 → 多匹配列表 (Priority: P1)

访客输入扩展名(或显示名片段),期望得到 Marketplace 默认排序的候选列表,
再点入单个扩展的版本列表页,选择适配本地 VSCode 版本的版本下载。

**Why this priority**: 扩展是站点第二大流量来源;主路径。
**Independent Test**: 输入扩展名,断言返回多结果列表;点击单条,断言进入该扩展的
所有版本页,且每个版本都标注 `engines.vscode`。

**Acceptance Scenarios**:
1. **Given** 关键词命中多个扩展,**When** 提交,**Then** 返回多结果列表,排序与
   Marketplace 默认一致。
2. **Given** 访客点击列表中某条,**When** 进入,**Then** 展示该扩展所有版本,
   每条含发布时间与 `engines.vscode`。
3. **Given** 关键词唯一命中,**When** 提交,**Then** 直接进入该扩展版本列表页(等同 S4)。

---

### User Story 4 - 扩展名 + 版本号 → 直接下载页 (Priority: P2)

访客已知扩展名与具体扩展版本号,期望跳过列表直接进入该版本下载区。

**Why this priority**: 进阶用户的便捷路径,功能上与 P1 部分重叠但 UX 独立。
**Independent Test**: 输入 `publisher.extension@1.2.3`,断言直接渲染该版本的
单一直链区,无列表。

**Acceptance Scenarios**:
1. **Given** 扩展版本存在,**When** 提交,**Then** 渲染该扩展该版本的单一下载区。
2. **Given** 版本号不存在,**When** 提交,**Then** 报错"扩展版本不存在"。

---

### User Story 5 - 浏览版本列表并筛选 (Priority: P3)

访客浏览全部版本,按 Channel / Platform / Architecture 组合筛选。

**Why this priority**: 浏览模式是辅助入口,不阻塞主流程。
**Independent Test**: 打开浏览页,切换 Channel / Platform / Architecture,断言列表
按所选组合筛过。

**Acceptance Scenarios**:
1. **Given** 任意默认筛选,**When** 改变任一筛选维度,**Then** 列表相应变化。
2. **Given** 筛选组合无匹配,**When** 提交,**Then** 显示空态文案。

---

### User Story 6 - 浏览热榜 (Priority: P2)

访客浏览热榜,维度:客户端 / 服务端 / 扩展 × 时间窗(24h / 7d / 30d),每维度
Top 10。

**Why this priority**: 拉留存与回访的核心 UI;非 P1 但关键。
**Independent Test**: 打开热榜,切换维度与时间窗,断言每组渲染 Top 10 条目。

**Acceptance Scenarios**:
1. **Given** 任一物维度 × 任一时间窗,**When** 切换,**Then** 渲染 Top 10 条目,
   含排名、名称、计数。
2. **Given** 时间窗内无事件,**When** 切换,**Then** 显示空态,不报错。

---

### User Story 7 - 查看文档(Tooltip + 内嵌区) (Priority: P3)

访客在页面任意处通过 Tooltip 或内嵌文档区阅读"客户端安装 / 服务端安装 /
服务端启动与连接 / 离线导入扩展 / 排错 FAQ"五个主题,中英双语可切换。

**Why this priority**: 文档形态对完成率有影响,但本身不产生新功能。
**Independent Test**: 切换语言,断言五个主题均渲染对应语言;Tooltip 在 hover
控件时弹出说明。

**Acceptance Scenarios**:
1. **Given** 默认语言,**When** 切换,**Then** 所有可见文案与文档块同步切换。
2. **Given** 访客 hover 任意"有提示"控件,**When** 触发,**Then** 显示对应说明气泡。

---

### Edge Cases

- 版本号解析失败(空串、含非法字符):返回明确错误,不抛 500。
- 扩展搜索 Marketplace 超时 / 失败:返回错误并提示重试,不缓存。
- GeoIP 解析失败:埋点写入 `CountryCode = UNKNOWN`,不阻塞请求。
- 行为事件表写入失败:搜索/下载主路径仍正常返回(埋点失败降级)。
- 客户端 90 天事件过期清理:删除任务执行后,聚合查询仍正确(过期区间视为空)。
- 并发点击同一下载链接:每个点击都独立记入事件,热榜计数正确(同 IP 同日重复计入)。
- Stable 版本号指向尚未发布 vscode-server 的 commit:服务端区显示"暂未发布"占位文案,不渲染假链接。

## Requirements *(mandatory)*

### Functional Requirements

#### 版本查询
- **FR-001**: System MUST 区分 Stable / Insider 两个客户端频道,并在请求中明确当前频道。
- **FR-002**: System MUST 在版本号不存在时返回明确错误信息,且不渲染任何下载区。
- **FR-003**: System MUST 在提交版本号时,默认按 User-Agent 推断 Platform / Architecture,
  且允许访客手动覆盖。

#### vscode-server 配套
- **FR-004**: System MUST 在展示 `vscode-server` 直链时,**同时**显示该服务端构建
  对应的 Stable 客户端 commit hash。
- **FR-005**: System MUST 拒绝把 vscode-server 链接渲染到非 Stable 频道。

#### 扩展
- **FR-006**: System MUST 在扩展搜索时透传请求到 VSCode Marketplace 并以 Marketplace
  默认排序返回多结果。
- **FR-007**: System MUST 在扩展单结果关键词时直接进入版本列表;在多结果时进入候选列表。
- **FR-008**: System MUST 在查询指定扩展版本号时,仅返回该版本;未指定时返回全部版本,
  且每版本附带 `engines.vscode` 与发布时间。
- **FR-009**: System MUST 在每个 ExtensionVersion 渲染中显式标注 `engines.vscode`,
  便于访客判定与本地 VSCode 是否匹配。

#### 文档与辅助脚本
- **FR-010**: System MUST 在页面提供五个主题的内嵌文档区(客户端安装 / 服务端安装 /
  服务端启动与连接 / 离线导入扩展 / 排错 FAQ),中英双语。
- **FR-011**: System MUST 在 Tooltip / 内嵌区中提供的辅助脚本**仅含**运行/安装命令;
  不得包含 `wget`、`curl`、`Invoke-WebRequest`、`fetch` 等任何下载动作。

#### 行为埋点
- **FR-012**: System MUST 在"搜索提交"与"下载点击"两类事件发生时写入事件表,
  字段(语义级,具体落库结构由 ADR 决定):事件类型、对象类型(客户端/服务端/扩展)、
  对象标识、Platform、Architecture、Channel(客户端/服务端维度)、CountryCode、
  时间戳。
- **FR-013**: System MUST 写入 `CountryCode`;若 GeoIP 解析失败,写入 `UNKNOWN`,
  且不得阻塞主请求。
- **FR-014**: System MUST 保留行为事件 90 天,过期清理(由部署侧的清理任务执行)。

#### 热榜
- **FR-015**: System MUST 支持三物维度(客户端 / 服务端 / 扩展)× 三时间窗
  (24h / 7d / 30d)的热榜,每组 Top 10,数据源为本地产物事件表聚合查询。
- **FR-016**: System MUST 在时间窗内无事件时返回空态而非错误。

#### 地区维度(预留)
- **FR-017**: System MUST 在事件写入时落地 `CountryCode`,**本期 UI 不暴露**地区维度的
  热榜视图(仅作为未来"按地区看热榜"的数据基础)。

#### 文档形态
- **FR-018**: System MUST 不设立独立文档站 / Wiki 后台,所有文档以 Tooltip + 内嵌
  文档区承载。

#### 范围外(明确不做)
- **FR-NEG-001**: System MUST **不**引入用户系统、登录、UID、Cookie 跨站追踪。
- **FR-NEG-002**: System MUST **不**引入限频、反作弊、验证码(本期)。
- **FR-NEG-003**: System MUST **不**支持 VSCodium / Cursor / Windsurf / Coder / Gitpod 等
  替代客户端或云端方案的下载聚合。
- **FR-NEG-004**: System MUST **不**实现下载代理、批量下载、打包脚本。

### Key Entities

- **ClientRelease**:VSCode 客户端发布版本。属性:Channel(Stable / Insider)、Version、
  Platform、Architecture、DownloadURL(官方)、CommitHash(仅 Stable,对应 vscode-server)。
- **ServerRelease**:`vscode-server` 发布版本。属性:CommitHash(绑定 Stable 客户端
  commit)、Platform、Architecture、DownloadURL(官方)、对应 ClientVersion。
- **Extension**:扩展元数据。属性:Publisher、Name、DisplayName、LatestVersion。
- **ExtensionVersion**:扩展单个版本。属性:Extension、Version、PublishTime、
  `engines.vscode`、DownloadURL(官方 Marketplace 直链)。
- **BehaviorEvent**:行为埋点。属性:EventType(SEARCH / DOWNLOAD)、TargetType
  (CLIENT / SERVER / EXTENSION)、TargetIdentifier、Platform、Architecture、Channel
  (CLIENT/SERVER 维度下)、CountryCode、Timestamp。

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 访客在已知 Stable 版本号时,**5 秒内**可拿到对应客户端与 vscode-server
  直链(主路径 P1)。
- **SC-002**: 扩展搜索平均响应耗时 **≤ 3 秒**(依赖 Marketplace 官方响应)。
- **SC-003**: 任意页面在桌面 / 平板 / 移动三档断点下均可正常使用,**无横向滚动、无溢出**。
- **SC-004**: 所有可见文案(UI + 文档)在**中 / 英**两语间切换完整,**无遗漏字段**。
- **SC-005**: 行为事件写入失败时,**主请求仍 100% 返回**,失败率不入埋点。
- **SC-006**: 90 天前的事件**不存在于热榜聚合结果**中。
- **SC-007**: 站点的所有下载链接 URL **100% 命中**官方域名白名单(由测试断言)。
- **SC-008**: 任意一页源代码中**不存在** `wget` / `curl` / `Invoke-WebRequest` /
  `fetch` 等下载动作字面量(辅助脚本区由 E2E 与静态扫描双重守护)。
- **SC-009**: 热榜每个维度 × 时间窗返回**最多 10 条**;空态正确呈现。
- **SC-010**: 首次部署后,文档(GDPR/隐私声明 + 五主题)在中英两语下均**可阅读、无占位符**。

## Assumptions

- 假设 VSCode 官方 Stable / Insider / vscode-server API 在公网可访问,且响应符合
  本规格所需的字段(版本号、commit hash、Platform/Architecture 直链)。
- 假设 VSCode Marketplace 搜索 API 可被透传调用,且默认排序可由其控制。
- 假设访客的网络环境对 Microsoft 官方域名直链可达(站点不做镜像/代理)。
- 假设部署侧的"过期清理任务"每日可被调度一次;具体调度方式由 ADR 决定(spec 不绑定
  cron daemon,见 NFR8)。
- 假设 GeoIP 数据源以本地数据库形式提供;离线解析失败时回退 `UNKNOWN`,不阻塞。
- 假设 SQLite 作为事件存储足以支撑当前流量;具体持久化方言由 ADR 决定。
- 假设 i18n 通过前端双语字典 + 后端错误码双语消息实现;不引入独立翻译管线。
- 假设响应式断点为桌面 ≥ 1024px / 平板 768–1023px / 移动 < 768px。