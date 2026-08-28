# QA Checklist: VS Code Download Hub

Each item maps to a Success Criterion (SC-###) in `specs/001-vscodemarket-hub/spec.md`.

## Deployment & health

- [ ] SC-N/A: `docker compose up -d --build` brings up backend + frontend
- [ ] SC-N/A: `curl http://localhost:8081/api/v1/healthz` returns `{"status":"ok"}`
- [ ] SC-N/A: `bash scripts/smoke.sh` exits 0

## Main path

- [ ] **SC-001** Stable 1.94.2 lookup returns within 5s; both client + server URLs present
- [ ] **SC-001** Server URL contains the corresponding client commit hash
- [ ] **SC-007** Every returned `downloadUrl` URL host ∈ {update.code.visualstudio.com, marketplace.visualstudio.com, aka.ms}
- [ ] **SC-004** zh-CN ↔ en-US toggle changes all visible strings (no raw key leaks)

## Extensions

- [ ] **SC-002** `?q=python` returns results within 3s, sorted by Marketplace default
- [ ] **SC-002** Single-result keywords jump directly to detail page (US4 acceptance)
- [ ] **SC-007** Each extension `downloadUrl` host is marketplace.visualstudio.com

## Trending

- [ ] **SC-009** `/trending?targetType=EXTENSION&window=24h` returns at most 10 rows
- [ ] **SC-009** Counts strictly non-increasing across ranks
- [ ] **SC-006** Rows whose `created_at < now() - 90d` are absent after sweep

## Privacy / Constitution

- [ ] **SC-005** Stopping the backend (`docker compose stop backend`) keeps `/api/v1/events` returning 202 (degraded header) — main path still 200/202
- [ ] **SC-007** No frontend bundle contains a non-whitelisted download host (scan)
- [ ] **SC-008** `grep -RE 'wget|curl|Invoke-WebRequest|fetch' frontend/src` returns 0 hits outside comments
- [ ] **FR-013** With no GeoIP mmdb, every new event has `countryCode = "UNKNOWN"`
- [ ] **SC-006** Events older than 90 days are removed at next sweep

## Responsive

- [ ] **SC-003** No horizontal scroll at 360 / 768 / 1280 widths
- [ ] **SC-003** All pages remain usable on iPad portrait

## Documentation

- [ ] **SC-010** README has GDPR / 隐私声明 section
- [ ] **SC-010** Docs page (UI) renders all 5 topics in zh-CN and en-US
- [ ] **SC-010** No "TODO" / "placeholder" text in shipped docs

## Out-of-scope guard

- [ ] No user / / auth UI present anywhere
- [ ] No rate-limit / captcha middleware in backend
- [ ] No independent docs site / wiki backend
- [ ] No VSCodium / Cursor / Coder links anywhere