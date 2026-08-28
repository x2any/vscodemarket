# ADR-0005: 技术栈 Golang + Vue 3 + TypeScript + Element Plus

**状态**: Accepted
**日期**: 2026-08-28
**决策者**: 产品负责人 + 全栈

## 背景

需求已固定:后端 Golang,前端 Vue.js + TypeScript + Element Plus。

## 决策

| 层 | 选型 |
|----|------|
| 后端语言 | Go (>= 1.22) |
| HTTP 框架 | 标准库 `net/http` + 第三方路由器(chi 或 gin,后续 PR 决定) |
| 前端框架 | Vue 3 + Composition API |
| 类型系统 | TypeScript (strict) |
| UI 库 | Element Plus |
| 构建工具 | Vite |
| 包管理 | 后端 `go mod`,前端 `pnpm` |
| 容器化 | Docker (multi-stage build,镜像作为最终交付物) |

## 理由

- **需求已经指定**,此处 ADR 只固化版本基线和组合方式。
- **Vite** 是 Vue 3 的事实默认,HMR 快、构建产物小。
- **Element Plus** 是 Vue 3 对标 Element UI 的官方后续,组件库成熟度高。
- **Golang** 单二进制部署方便,适合公网小服务场景。

## 后果

- 不引入 React、Next.js、Nuxt、Ant Design Vue、Vuetify、Naive UI。
- 前端不引入 SSR、不引入 Nuxt。SPA 即可。
- **后端 ORM / DB driver 选型见 ADR-0011**,本 ADR 不重复表态。

## 不在本 ADR 范围

- 具体的 UI 组件拆分、数据流管理(pinia vs composable)、测试栈(playwright/vitest)——这些是 plan/spec 阶段细节。