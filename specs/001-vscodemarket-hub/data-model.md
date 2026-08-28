# Data Model: VS Code Download Hub

**Date**: 2026-08-28 | **Plan**: [plan.md](./plan.md)

## Entities

### 1. ClientRelease(运行时对象,不落库)

VSCode 客户端发布版本元数据。

| 字段 | 类型 | 说明 |
|------|------|------|
| Channel | enum: `stable` \| `insider` | 频道 |
| Version | string | `1.94.2` / insider 自有版本号 |
| Platform | enum: `windows` \| `linux` \| `darwin` | |
| Architecture | enum: `x86_64` \| `arm64` \| `armv7` | |
| DownloadURL | string | 官方域名直链(必填) |
| CommitHash | string? | 仅 Stable 有效,对应 vscode-server commit |

**约束**:
- `insider` 频道下 `CommitHash` 为 null。
- `Platform × Architecture` 仅允许 ADR-0006 列出的笛卡尔积。

### 2. ServerRelease(运行时对象)

`vscode-server` 发布版本。

| 字段 | 类型 | 说明 |
|------|------|------|
| CommitHash | string | 绑定 Stable 客户端 commit |
| ClientVersion | string | 对应 Stable 版本号(冗余便于查询) |
| Platform | enum | 同上 |
| Architecture | enum | 同上 |
| DownloadURL | string | 官方域名直链 |

**约束**:
- 不为 Insider 频道生成 server 记录。

### 3. Extension(运行时对象)

| 字段 | 类型 | 说明 |
|------|------|------|
| Publisher | string | `ms-python` |
| Name | string | `python` |
| DisplayName | string | "Python" |
| LatestVersion | string | 最新版本号 |

唯一标识:`Publisher.Name`(全小写)。

### 4. ExtensionVersion(运行时对象)

| 字段 | 类型 | 说明 |
|------|------|------|
| Extension | ref Extension | 父扩展 |
| Version | string | `2024.20.0` |
| PublishTime | datetime | 发布时间(UTC) |
| EnginesVSCode | string | `^1.94.0` 等 |
| DownloadURL | string | Marketplace 官方直链 |

### 5. BehaviorEvent(落库)

| 字段 | 类型 | 约束 |
|------|------|------|
| ID | uint64 PK | 自增 |
| EventType | enum: `SEARCH` \| `DOWNLOAD` | FR-012 |
| TargetType | enum: `CLIENT` \| `SERVER` \| `EXTENSION` | FR-012 |
| TargetIdentifier | string | 版本号 / 扩展 ID / commit hash |
| Platform | string? | nullable(扩展维度无) |
| Architecture | string? | nullable |
| Channel | string? | nullable(扩展维度无) |
| CountryCode | string(2 或 `UNKNOWN`) | FR-013,NFR3 |
| CreatedAt | datetime | 默认 UTC now;索引 |

**索引**:
- `idx_event_target_time (TargetType, TargetIdentifier, CreatedAt)` — 服务于热榜聚合
- `idx_event_created (CreatedAt)` — 服务于 90 天 sweep

**保留**:
- 90 天,过期由 `internal/sweeper` 每 24h 删除。
- FR-014、SC-006。

### 6. UAInference(运行时对象)

| 字段 | 类型 | 说明 |
|------|------|------|
| Platform | enum | |
| Architecture | enum | |
| Confidence | enum: `HIGH` \| `FALLBACK` | HIGH=命中 OS+arch 关键字,FALLBACK=仅命中 OS |

## Relationships

```
ClientRelease (Channel=stable)  ──commit hash──>  ServerRelease
Extension ──< ExtensionVersion
BehaviorEvent ──> ClientRelease / ServerRelease / Extension
```

## Validation Rules

- 版本号字符串必须匹配 `^\d+\.\d+\.\d+(-[\w.]+)?$`(Stable)或官方 Insider 格式。
- `CountryCode` ∈ ISO 3166-1 alpha-2 或字面量 `UNKNOWN`。
- `Platform × Architecture` 必须 ∈ ADR-0006 矩阵,否则 400。
- `TargetIdentifier` 长度 ≤ 256,避免日志膨胀。

## State Transitions

BehaviorEvent 一旦写入不可变更;过期清理为物理删除。

## Storage Choice

- `ClientRelease / ServerRelease / Extension / ExtensionVersion`:运行时对象,
  不落库,每次请求透传官方 API(NFR8 不引入缓存层)。
- `BehaviorEvent`:SQLite 单表,GORM AutoMigrate(M11 之前禁用 ddl-auto sync;
  采用显式 migration 文件)。