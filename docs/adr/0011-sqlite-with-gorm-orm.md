# ADR-0011: 持久化层选型 SQLite + GORM(预留 Postgres 切换)

**状态**: Accepted
**日期**: 2026-08-28
**决策者**: 后端

## 背景

ADR-0010 引入持久化需求,需要选型。

候选:
- (a) SQLite + GORM
- (b) Postgres + GORM
- (c) ClickHouse
- (d) DuckDB

## 决策

**采用 (a) SQLite + GORM。但代码层只用 GORM 接口,不绑定 SQLite 方言,后续可换 Postgres。**

## 理由

- **零运维**:SQLite 是单文件,docker volume 挂载即可,无独立服务。
- **GORM 抽象**:不写原生 SQL,后期切换 Postgres 只需改 DSN。
- **容量评估**:每条事件约 200B,假设 1000 DAU × 20 事件/人/天 = 20K 行/天 ≈ 4MB/天,90 天保留 360MB,SQLite 轻松应付。
- **多维聚合**:SQLite + 索引够用;真到千万级行再迁。

## 后果

- 新增 GORM v2 依赖
- 仓库 `data/` 目录挂载,git 忽略
- 模型设计:
  - `events`(id, type, target_id, query, ip, CountryCode, created_at)  // GORM tag 内 snake_case 为 `country_code`
  - 索引:(type, target_id, created_at) 用于物+时间维度
  - 索引:(created_at) 用于全局时间窗口
  - 索引:(type, created_at, country_code) 用于地区维度
- 90 天保留策略:定时清理或启动时 sweep
- 禁止:业务代码中写 `gorm:"type:..."` 等方言特定 tag

## 切换路径(写入 README)

当 SQLite 不够时:
1. `docker run postgres`
2. 改 DSN 为 Postgres URL
3. `gorm.AutoMigrate` 重建表
4. 改用 `pgx` 驱动

## 替代方案

- **(b) 直接 Postgres**:否决。单机部署,Postgres 运维过重。
- **(c) ClickHouse**:否决。当前数据量级杀鸡用牛刀。
- **(d) DuckDB**:否决。DuckDB 偏 OLAP batch,不适合高频 INSERT。