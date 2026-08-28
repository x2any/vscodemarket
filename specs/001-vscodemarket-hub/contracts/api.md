# HTTP API Contract: VS Code Download Hub

**Base URL**: `/api/v1`
**Auth**: none(FR-NEG-001)
**Error Format**: `{ "error": { "code": "VERSION_NOT_FOUND", "message_zh": "...", "message_en": "..." } }`

## Endpoints

### 1. `GET /healthz`

健康检查。

**Response 200**:
```json
{ "status": "ok" }
```

---

### 2. `POST /versions/lookup`

主路径 P1/P2。版本查询,同屏返回 client + server 直链。

**Request**:
```json
{
  "channel": "stable",          // stable | insider
  "version": "1.94.2",
  "platform": "darwin",         // 由 UA 推断或用户覆盖
  "architecture": "arm64"
}
```

**Response 200 (stable)**:
```json
{
  "channel": "stable",
  "version": "1.94.2",
  "client": {
    "downloadUrl": "https://update.code.visualstudio.com/...",
    "platform": "darwin",
    "architecture": "arm64"
  },
  "server": {
    "downloadUrl": "https://update.code.visualstudio.com/.../vscode-server-darwin-arm64.tar.gz",
    "commitHash": "e54c774e0add6f2b...",
    "clientVersion": "1.94.2",
    "platform": "darwin",
    "architecture": "arm64"
  }
}
```

**Response 200 (insider)**:`server` 字段省略。
**Response 404**:`VERSION_NOT_FOUND`
**Response 400**:`INVALID_PLATFORM_ARCH`(不在 ADR-0006 矩阵)

---

### 3. `POST /ua/infer`

UA 推断(M3)。

**Request**:
```json
{ "userAgent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/... Chrome/120.0.0.0 Safari/537.36" }
```

**Response 200**:
```json
{ "platform": "darwin", "architecture": "arm64", "confidence": "HIGH" }
```

**Confidence**:
- `HIGH`:命中 OS 与架构关键字
- `FALLBACK`:仅命中 OS,架构默认 `x86_64`

---

### 4. `GET /extensions/search?q={query}&page={n}&pageSize={n}`

扩展搜索,透传 Marketplace(FR-006)。

**Response 200**:
```json
{
  "results": [
    {
      "publisher": "ms-python",
      "name": "python",
      "displayName": "Python",
      "latestVersion": "2024.20.0"
    }
  ],
  "total": 137
}
```

排序:Marketplace 默认。

---

### 5. `GET /extensions/{publisher}/{name}/versions`

扩展所有版本(FR-008)。

**Response 200**:
```json
{
  "extension": { "publisher": "ms-python", "name": "python", "displayName": "Python" },
  "versions": [
    {
      "version": "2024.20.0",
      "publishTime": "2024-10-08T15:32:00Z",
      "enginesVscode": "^1.94.0",
      "downloadUrl": "https://marketplace.visualstudio.com/_apis/.../vsix"
    }
  ]
}
```

### 6. `GET /extensions/{publisher}/{name}/versions/{version}`

单版本(FR-008 快捷路径)。

**Response 200**:同 `versions[]` 单元素结构
**Response 404**:`EXTENSION_VERSION_NOT_FOUND`

---

### 7. `POST /events`

埋点接收(FR-012),失败降级不影响主路径。

**Request**:
```json
{
  "eventType": "SEARCH",          // SEARCH | DOWNLOAD
  "targetType": "CLIENT",         // CLIENT | SERVER | EXTENSION
  "targetIdentifier": "1.94.2",
  "platform": "darwin",           // 可选
  "architecture": "arm64",        // 可选
  "channel": "stable"             // 可选(扩展维度无)
}
```

**Response 202**:`{ "accepted": true }`(异步落库不阻塞)

服务器侧自动追加:`CountryCode`(GeoIP 解析或 `UNKNOWN`)、`CreatedAt`、来源 IP(仅进程内,不落库明文 IP — FR-013)。

---

### 8. `GET /trending?targetType={t}&window={w}`

热榜(M9)。

**Query**:
- `targetType`: `CLIENT` | `SERVER` | `EXTENSION`
- `window`: `24h` | `7d` | `30d`

**Response 200**:
```json
{
  "targetType": "EXTENSION",
  "window": "24h",
  "items": [
    { "rank": 1, "targetIdentifier": "ms-python.python", "count": 1842 },
    { "rank": 2, "targetIdentifier": "dbaeumer.vscode-eslint", "count": 1501 }
  ]
}
```

`items` 长度 ≤ 10;空态返回 `items: []`,HTTP 仍 200。

---

### 9. `GET /releases?channel={c}&platform={p}&architecture={a}&page={n}&pageSize={n}`

版本浏览(M6)。

**Response 200**:
```json
{
  "results": [
    {
      "channel": "stable",
      "version": "1.94.2",
      "platform": "darwin",
      "architecture": "arm64",
      "downloadUrl": "...",
      "commitHash": "e54c7..."
    }
  ],
  "page": 1,
  "pageSize": 20,
  "total": 132
}
```

## Cross-Cutting

- 所有 5xx 响应**不**含堆栈;仅 `{ "error": { "code": "INTERNAL", "message_zh": "...", "message_en": "..." } }`。
- 所有时间字段 ISO-8601 UTC。
- 所有 `downloadUrl` 在响应前由后端断言域名白名单(SC-007 测试守护),
  不在白名单 → 500 + 日志告警(绝不出网)。