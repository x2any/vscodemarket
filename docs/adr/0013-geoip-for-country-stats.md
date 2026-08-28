# ADR-0013: IP→国家/地区映射用于地区维度统计

**状态**: Accepted
**日期**: 2026-08-28
**决策者**: 产品负责人 + 后端

## 背景

ADR-0010 要求"按地区维度做网站统计"。需要从 IP 解析出国家码。

候选:
- (a) 在线 API(ipapi.co / ip-api.com 等)
- (b) 本地 GeoIP2 库(MaxMind)
- (c) 不做,只存 IP,后期再说

## 决策

**采用 (b) 本地 GeoIP2(MaxMind mmdb)。可选免费 GeoLite2,精度国家级足够。**

## 理由

- **离线可用**:公网部署但避免对在线 API 的依赖(限流、付费、可用性)。
- **零延迟**:内嵌 lookup,无网络开销。
- **精度**:国家级 GeoLite2 免费、够用。
- **隐私更稳**:IP 不外传到第三方。

## 后果

- 新增依赖:`github.com/oschwald/geoip2-golang`
- 启动时加载 `./data/GeoLite2-Country.mmdb`(git 忽略,部署时下载/挂载)
- 失败兜底:db 文件缺失则 `CountryCode = "UNKNOWN"`,不阻塞请求
- 保留 IP(ADR-0010 要求);GeoIP 解析失败仍照常入库

## 替代方案

- **(a) 在线 API**:否决。限流 / 费用 / 可用性 / 隐私外传都不利。
- **(c) 不解析**:否决。产品方明确要做地区维度。