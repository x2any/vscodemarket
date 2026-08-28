# ADR-0006: 客户端仅支持 Stable 与 Insider 两个 Channel

**状态**: Accepted
**日期**: 2026-08-28
**决策者**: 产品负责人

## 背景

客户端 Channel 选项:

- (a) Stable + Insider
- (b) Stable + Insider + Exploration
- (c) Stable + Insider + VSCodium
- (d) + Cursor / Windsurf / Trae 等 AI 衍生品

## 决策

**仅支持 Stable 与 Insider(选项 a)。**

## 理由

- 覆盖 95%+ 实际用户。
- Stable 与 Insider 都是 Microsoft 官方源,API 一致,集成零成本。
- Insider 版本号是日期格式,与 Stable 语义化版本号结构不同,UI 需要分别处理;仅两个 channel 时实现简单。

## 后果

- 不集成 VSCodium、Exploration、Cursor、Windsurf、Trae。
- 后端只有两个 channel 的查询/过滤逻辑。
- 前端 Channel 下拉只有两个选项。

## 替代方案

- **(b) +Exploration**:否决。Exploration 极少人用,且下载源不稳定。
- **(c) +VSCodium**:否决。VSCodium 是独立发行版,有自己的版本体系,集成后增加 50% 表面积。
- **(d) +AI 衍生品**:否决。这些产品的服务端不是 vscode-server,语义上不属于本项目。