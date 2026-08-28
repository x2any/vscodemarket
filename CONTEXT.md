---
name: vscodemarket-glossary
description: 精炼的领域词汇表,无实现细节,作为单一事实源供 ADRs/spec/plan 引用
---

# VSCodeMarket 领域词汇表

本文档只定义概念与术语,不涉及技术选型、API、模块、存储等实现细节。
实现细节见 `docs/adr/`。

---

## 一、产品级概念

| 术语 | 定义 |
|------|------|
| **VSCodeMarket** | 本项目——一个面向公网部署的辅助下载站,帮助内网用户获取与本地 VSCode 客户端配套的下载物。 |
| **辅助下载物** | VSCodeMarket 提供的全部产出物,包含三类:客户端二进制、服务端二进制、扩展包。 |
| **内网用户** | VSCodeMarket 的最终用户。在公网下载辅助下载物后,通过 U 盘/光盘/中转等方式拷入内网机器执行。 |
| **公网部署** | VSCodeMarket 自身运行于公网,服务对象为公网访客,不做鉴权。 |
| **公网应用属性** | 本项目是公网应用(非面向内网员工的内部工具),会记录匿名 IP + 行为事件用于热门榜单与地区维度统计。 |

---

## 二、客户端域

| 术语 | 定义 |
|------|------|
| **客户端 (Client)** | 用户本地的 VSCode 编辑器二进制。本项目仅支持 Microsoft 官方的两个频道。 |
| **Channel** | 客户端发布渠道。本项目仅含:Stable、Insider。 |
| **Stable 版本号** | 形如 `1.94.1` 的语义化版本号,指向某个稳定发布。 |
| **Insider 版本号** | 形如 `20240814` 的日期版本号,指向每日构建。 |
| **Platform** | 客户端目标操作系统家族:windows、linux、darwin。 |
| **Architecture** | 客户端目标 CPU 架构:x86_64、arm64、armv7。 |
| **PlatformBuild** | Platform × Architecture 的笛卡尔积,例如 `linux-x64`、`darwin-arm64`。 |

---

## 三、服务端域

| 术语 | 定义 |
|------|------|
| **服务端 (Server)** | 用户内网机器上运行的远程 VSCode 后端,Microsoft 官方名为 `vscode-server`,与 Stable 客户端的 commit hash 严格一一对应。 |
| **服务端版本号** | 服务端没有自己的版本号;它**总是与某个 Stable 客户端版本号绑定**。 |
| **Commit Hash** | 服务端产物内嵌的 git commit 标识,与客户端完全一致,是匹配两者的最终依据。 |

> **关键约束**:Stable 客户端与服务端必须 commit hash 一致才能成功建立 Remote-SSH 连接。

---

## 四、扩展域

| 术语 | 定义 |
|------|------|
| **扩展 (Extension)** | VSCode Marketplace 上的插件,以 `extensionId`(形如 `publisher.name`)唯一标识。 |
| **扩展版本 (ExtensionVersion)** | 某个扩展的某个具体版本,包含 vsix 文件 URL、engines.vscode 声明、发布时间等元数据。 |
| **VSIX** | 扩展的离线分发包格式,`.vsix` 文件,可在内网手动安装。 |
| **engines.vscode** | 扩展版本声明的最低兼容 VSCode 版本(语义化版本范围)。 |
| **Publisher** | 扩展发布者,可能为个人或组织。 |

---

## 五、平台与适配

| 术语 | 定义 |
|------|------|
| **Platform Build(扩展语境)** | 扩展安装目标平台(通常与客户端平台一致);扩展版本可不区分平台,但展示时可按平台过滤。 |
| **User-Agent 推断** | 后端根据请求者浏览器 UA,默认填入合适的 PlatformBuild,允许用户手动修改。 |
| **官方直链** | 来自 `update.code.visualstudio.com`、`code.visualstudio.com`、`marketplace.visualstudio.com` 等域名的官方下载地址。 |

---

## 六、角色与归属

| 术语 | 定义 |
|------|------|
| **匿名访客** | 唯一角色,无需登录/注册。 |
| **归属单元** | 以请求者 IP 作为唯一最小归属单元;不引入 UID、Cookie 跨站追踪。 |

---

## 七、UI 域

| 术语 | 定义 |
|------|------|
| **版本查询** | 用户输入一个版本号 → 页面同时展示该版本的客户端与服务端下载链接。 |
| **扩展搜索** | 用户输入插件名(可选版本号) → 返回匹配的扩展或单一版本下载页。 |
| **版本列表** | 客户端/服务端的历史发布版本清单,支持筛选 Channel、PlatformBuild。 |
| **适配版本列表** | 同一扩展下,所有 ExtensionVersion 的清单,展示各自的 engines.vscode,用户据此判断能否在本地 VSCode 上安装。 |
| **辅助脚本** | 内联在页面中的 shell 片段,帮助用户在内网机器上完成启动 vscode-server、安装扩展等操作。脚本不下载任何东西,只在内网机器上执行。 |
| **Tooltip / 内嵌文档区** | 站点不设独立文档站;文档以 Tooltip + 嵌入式文档区形式融入主页面。 |

---

## 八、文档主题(以 Tooltip/嵌入式文档区承载)

| 主题 | 内容 |
|------|------|
| **客户端安装** | Windows / macOS / Linux 三个平台安装客户端步骤与验证方法。 |
| **服务端安装** | 在内网机器上获取并运行 vscode-server 的步骤。 |
| **服务端启动与连接** | 启动 vscode-server、客户端 Remote-SSH 连接流程。 |
| **离线导入扩展** | 在内网机器上从 .vsix 文件手动安装扩展的步骤。 |
| **排错 FAQ** | 常见连接失败、commit hash 不匹配、端口冲突等问题的排错指南。 |

---

## 九、事件域(行为埋点)

| 术语 | 定义 |
|------|------|
| **事件 (Event)** | 用户一次有意义的操作,本项目只埋两类:搜索、下载点击。 |
| **事件类型 (EventType)** | `search` / `download`。 |
| **事件目标 (Target)** | `search` 事件携带原始关键词;`download` 事件携带 extensionId 或 client/server 版本号。 |
| **事件归属 (Attribution)** | 以请求者 IP 作为唯一归属单元;不引入 UID、不跨站追踪。 |
| **CountryCode** | 由 IP 解析得到的 ISO 3166-1 alpha-2 国家码,用于地区维度聚合;解析失败时为 `UNKNOWN`。 |
| **事件保留期** | 90 天,过期清理(产品决策,非实现细节)。 |

---

## 十、热榜域

| 术语 | 定义 |
|------|------|
| **热榜 (PopularRanking)** | 基于事件聚合而来的榜单,展示"最近被搜/下载最多的对象"。 |
| **热榜维度 (RankingDimension)** | 物维度 × 时间维度 的笛卡尔积:物维度∈{客户端,服务端,扩展},时间维度∈{24h,7d,30d}。 |
| **榜单条目 (RankingEntry)** | 单条榜单结果,含目标 ID + 计数 + 排名。 |
| **TopN** | 榜单取前 N 条(N 在 spec 阶段定为 10)。 |

---

## 十一、术语对照(中英)

| 中文 | English |
|------|---------|
| 客户端 | Client |
| 服务端 | Server |
| 扩展 | Extension |
| 扩展版本 | Extension Version |
| 渠道 | Channel |
| 平台 | Platform |
| 架构 | Architecture |
| 发布者 | Publisher |
| 适配版本 | Compatible VSCode Version |
| 辅助脚本 | Helper Script |
| 事件 | Event |
| 热榜 | Popular Ranking |
| 榜单维度 | Ranking Dimension |

---

## 十二、不在词汇表范围内的内容

- 具体 API 路径、域名、字段名(由 ADR 决定,不属于术语)
- 任何 HTTP 框架、ORM、缓存方案、CI 工具(纯实现)
- 前端组件库、数据流方案(纯实现)