# Changelog

All notable changes are recorded here. Versions follow SemVer.

## v0.1.0 — 2026-08-28

Initial end-to-end milestone set covering M1–M11 from `specs/001-vscodemarket-hub/plan.md`.

### Added
- Backend: Go 1.22 + chi + GORM/SQLite + GeoIP2; endpoints for
  versions lookup, UA inference, extension search & versions,
  release list, event recording, and trending aggregation.
- Frontend: Vue 3 + TypeScript strict + Element Plus shell with
  Home / Search / Detail / Releases / Trending / Docs pages and
  bilingual zh-CN / en-US dictionaries.
- Persistence: SQLite-backed behavior events (no PII) with 90-day
  in-process sweeper; GeoIP mmdb optional with UNKNOWN fallback.
- Deployment: multi-stage Dockerfiles, docker-compose with health
  gate, named data volume, nginx reverse proxy.
- Documentation: README with GDPR/隐私声明, in-app 5-topic docs
  page, scripts/smoke.sh, and CI workflow.

### Constitution compliance
- I 公网应用属性 — events hold only CountryCode + aggregated identifiers.
- II 官方直链原则 — all returned URLs pass whitelist assertion.
- III 辅助脚本零网络原则 — Vitest guard rejects wget/curl/Invoke-WebRequest/fetch.
- IV 极简立场 — no cron daemon / service worker / sidecar; sweeper is in-process.

### Known gaps (deferred, no blockers)
- Frontend `vite build` runs locally but is not exercised in CI; add when
  a deploy pipeline is wired up.
- Playwright E2E deferred in favour of Vitest contract-style tests
  (SC-001/002/006/009). Add Playwright only if/when a real renderer
  failure appears.
