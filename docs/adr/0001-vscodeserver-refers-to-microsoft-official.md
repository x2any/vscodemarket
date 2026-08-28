# ADR-0001: vscodeserver 仅指代 Microsoft 官方 vscode-server

**状态**: Accepted
**日期**: 2026-08-28
**决策者**: 产品负责人

## 背景

项目中 `vscodeserver` 一词存在多种可能的指代:

- (a) Microsoft 官方 `vscode-server`(CLI 工具,与 Stable 客户端 commit hash 严格对齐)
- (b) Coder 开源 `code-server`(独立版本号,Web 版 VSCode)
- (c) Gitpod `openvscode-server`(独立版本号)
- (d) 同时支持多种,UI 分 Tab

## 决策

**本项目 `vscodeserver` 仅指代 Microsoft 官方 `vscode-server`(选项 a)。**

## 理由

- **目标用户最常见**:Microsoft 官方 vscode-server 是内网 Remote-SSH 场景的事实标准。
- **与客户端版本天然对应**:Stable 客户端和服务端 commit hash 一致,本项目核心价值之一就是"输入版本号同时给客户端+服务端",这一对应关系只有 (a) 成立。
- **Coder / Gitpod 是替代品**:它们各自有独立生态,通常作为单独产品被部署,不需要"客户端/服务端版本对照"这一功能。
- **集成成本最低**:只对接一个官方源,URL 与 API 模式单一。

## 后果

- 不支持 Coder / Gitpod 生态。如未来需要,届时新立 ADR。
- UI 上无须 Tab 切换,服务端视图与 Stable 客户端视图强耦合。
- 必须展示 commit hash(或与之绑定的客户端版本号),供用户验证两端一致。

## 替代方案

- **(d) 多 Tab 支持**:否决。集成成本翻倍,且不符合"客户端对应版本"的核心叙事。