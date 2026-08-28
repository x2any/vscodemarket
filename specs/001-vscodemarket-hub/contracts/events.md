# Behavior Events Contract

**Date**: 2026-08-28 | **Spec**: [spec.md](./spec.md)

## Event Types

| EventType | 触发时机 | TargetType 集合 | 必填字段 |
|-----------|---------|---------------|---------|
| SEARCH    | 表单提交 / 搜索框回车 | CLIENT, SERVER, EXTENSION | targetIdentifier |
| DOWNLOAD  | 用户点击 DownloadLinkCard | CLIENT, SERVER, EXTENSION | targetIdentifier, platform, architecture |

## Fields

| 字段 | 类型 | 来源 | 备注 |
|------|------|------|------|
| eventType | enum | 客户端 | SEARCH \| DOWNLOAD |
| targetType | enum | 客户端 | CLIENT \| SERVER \| EXTENSION |
| targetIdentifier | string | 客户端 | 版本号 / commit hash / `publisher.name` |
| platform | enum? | 客户端 | 仅 client/server 维度 |
| architecture | enum? | 客户端 | 仅 client/server 维度 |
| channel | enum? | 客户端 | 仅 client/server 维度 |
| countryCode | string | 服务端 | GeoIP 或 `UNKNOWN`(FR-013) |
| createdAt | datetime | 服务端 | UTC |

**客户端请求体**(POST /api/v1/events):
```json
{
  "eventType": "SEARCH",
  "targetType": "CLIENT",
  "targetIdentifier": "1.94.2",
  "platform": "darwin",
  "architecture": "arm64",
  "channel": "stable"
}
```

**服务端落库**额外追加 `countryCode` 与 `createdAt`。
**服务端从不**落库明文 IP(宪法 I 公网合规边界)。

## Failure Modes

| 场景 | 行为 |
|------|------|
| DB 写入失败 | HTTP 仍返 202,内部日志告警;主请求不受影响(SC-005) |
| GeoIP mmdb 缺失 / 解析失败 | countryCode 写 `UNKNOWN`(NFR3) |
| 字段缺失 | 接受,nullable 字段写 NULL |
| 字段非法 | 422 + 不落库 |

## Retention

90 天;进程内 sweep 每 24h 删除 `created_at < now() - 90d` 的行(FR-014)。

## PII Stance

- 明文 IP **不落库**;仅进程内用于 GeoIP lookup,lookup 完丢弃。
- `countryCode` 为国家粒度,不可定位自然人。
- 不收集 UA / Referer / Cookie / 设备指纹。
- 上述约束写入 README(GDPR/隐私声明)。