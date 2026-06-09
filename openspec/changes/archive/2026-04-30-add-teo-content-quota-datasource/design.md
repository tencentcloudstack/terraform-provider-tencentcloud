## Context

Terraform Provider for TencentCloud 目前支持 TEO（EdgeOne）产品的多种数据源（如 `tencentcloud_teo_zones`、`tencentcloud_teo_default_certificate` 等），但缺少对内容管理配额的查询支持。云 API `DescribeContentQuota` 已提供该能力，需要将其封装为 Terraform 数据源。

当前状态：
- 云 API `teo/v20220901.DescribeContentQuota` 已在 vendor 中可用
- 请求参数：`ZoneId`（站点 ID）
- 响应包含 `PurgeQuota`（刷新配额列表）和 `PrefetchQuota`（预热配额列表），均为 `[]*Quota` 类型
- `Quota` 结构体包含 `Batch`、`Daily`、`DailyAvailable`、`Type` 四个字段
- 该接口无分页参数，单次调用返回全部数据

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_content_quota` 数据源，支持通过 `zone_id` 查询内容管理配额
- 返回 PurgeQuota 和 PrefetchQuota 两个配额列表，每个条目包含 Type、Batch、Daily、DailyAvailable
- 遵循现有 TEO 数据源的代码模式（参考 `tencentcloud_teo_zones`）
- 在 service 层封装 API 调用，支持重试逻辑
- 生成对应的单元测试和文档

**Non-Goals:**
- 不支持修改配额（配额为只读信息）
- 不支持批量查询多个站点的配额
- 不引入新的架构模式或依赖

## Decisions

### 1. 数据源 Schema 设计

**决策**：`purge_quota` 和 `prefetch_quota` 使用 `schema.TypeList` + `schema.Resource` 嵌套块，每个元素包含 `type`、`batch`、`daily`、`daily_available` 四个 Computed 字段。

**理由**：与云 API 的 `[]*Quota` 数组结构一致，且与现有 TEO 数据源模式保持统一。PurgeQuota 通常包含多个类型（purge_prefix、purge_url、purge_host、purge_all、purge_cache_tag），PrefetchQuota 包含 prefetch_url 类型。

### 2. Service 层方法设计

**决策**：在 `TeoService` 中新增 `DescribeTeoContentQuotaByFilter` 方法，接受 `map[string]interface{}` 参数，返回 `[]*teo.Quota` 列表。

**理由**：
- 遵循现有 service 层的 `Describe*ByFilter` 命名和参数模式
- 虽然 `DescribeContentQuota` 接口无分页参数，但保留统一的方法签名以维持代码一致性
- 方法内部仍需使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 处理临时性 API 错误

### 3. 数据源 ID 生成策略

**决策**：使用 `helper.DataResourceIdsHash([]string{zoneId})` 生成数据源 ID。

**理由**：与 `tencentcloud_teo_zones` 数据源的模式一致，基于查询参数生成确定性 ID，便于 Terraform 状态管理。

### 4. 接口无分页处理

**决策**：`DescribeContentQuota` 接口无分页参数，service 层方法直接单次调用返回结果，不实现分页循环。

**理由**：该接口返回的是配额信息（数据量小），API 本身不支持分页，无需做分页处理。

## Risks / Trade-offs

- **[API 返回 null]** → PurgeQuota 和 PrefetchQuota 字段可能返回 null，在 flatten 数据时需要做 nil 检查，nil 时不设置对应字段
- **[接口无分页]** → 如果未来配额类型增多，单次调用可能返回大量数据，但目前 Quota 类型有限（6种），风险极低
- **[ZoneId 可选]** → 云 API 的 ZoneId 标记为 omitnil（可选），但在 Terraform 数据源中应设置为 Required，因为查询特定站点的配额需要明确的站点标识
