## Context

Terraform Provider for TencentCloud 中，TEO（EdgeOne）产品已支持多种数据源（zones、plans、origin_acl 等），但缺少多通道安全加速网关（Multi-Path Gateway）的数据源。云API已提供 `DescribeMultiPathGateways` 接口，可查询指定站点下的网关列表。

当前 TEO 数据源遵循统一模式：schema 定义 → paramMap 构建 → 服务层调用（含分页和重试）→ 结果扁平化 → 写入 state。本项目需遵循此模式实现新数据源。

关键约束：
- `DescribeMultiPathGateways` 是列表接口，返回 `[]*MultiPathGateway`，但**不返回 `Lines` 字段**（Lines 仅在详情接口 `DescribeMultiPathGateway` 中返回）
- Filter 结构使用 `Name` + `Values`（非 AdvancedFilter），支持 `gateway-type` 和 `keyword` 两种过滤条件
- 分页参数：Offset 默认 0，Limit 默认 20 最大 1000

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_multi_path_gateway` 数据源，支持查询 TEO 多通道安全加速网关列表
- 支持 `zone_id` 和 `filters` 作为查询条件
- 返回完整的网关信息（GatewayId、GatewayName、GatewayType、GatewayPort、Status、GatewayIP、RegionId、NeedConfirm）
- 在 provider.go 和 provider.md 中正确注册数据源
- 生成对应的 .md 文档

**Non-Goals:**
- 不实现 `Lines` 字段（列表接口不返回此字段）
- 不实现网关的 CRUD 资源（仅实现数据源查询）
- 不修改现有数据源的行为

## Decisions

### 1. 数据源 Schema 设计
**决策**: 遵循 TEO 现有数据源模式（参考 `data_source_tc_teo_plans.go`），使用 `filters`（TypeList + Filter 结构）作为过滤条件，`gateways`（TypeList + Computed）作为输出。

**理由**: 与现有 TEO 数据源保持一致性，降低维护成本。Filter 结构与 `data_source_tc_teo_plans.go` 中使用的 `teo.Filter`（Name + Values）完全一致。

### 2. 不包含 Lines 字段
**决策**: 在数据源的 `gateways` 输出中不包含 `Lines` 字段。

**理由**: 云API文档明确说明 `DescribeMultiPathGateways` 列表接口不返回 Lines 字段，仅在 `DescribeMultiPathGateway` 详情接口中返回。若要获取 Lines 需要额外调用详情接口，增加复杂度且不符合数据源列表查询的设计意图。

### 3. 分页策略
**决策**: 服务层分页使用 Limit=1000（云API最大值），自动循环分页直到获取所有结果。

**理由**: 按照项目规范，查询接口分页字段应给定云API注释中标注的最大值。Limit=1000 可减少 API 调用次数。

### 4. zone_id 参数为 Required
**决策**: `zone_id` 参数设为 Required。

**理由**: 云API `DescribeMultiPathGateways` 的 `ZoneId` 是必填参数，用于指定查询的站点。

### 5. 服务层方法命名
**决策**: 服务层方法命名为 `DescribeTeoMultiPathGatewaysByFilter`。

**理由**: 遵循 TEO 服务层现有命名模式（如 `DescribeTeoZonesByFilter`、`DescribeTeoPlansByFilters`）。

## Risks / Trade-offs

- **[Lines 字段缺失]** → 如果后续需要 Lines 信息，需新增调用详情接口的逻辑或新增独立数据源。当前设计优先保持简洁性。
- **[大数据量分页]** → 如果网关数量超过 1000，服务层分页循环会自动处理。但极端情况下大量数据可能导致 API 超时。→ 使用 `tccommon.ReadRetryTimeout` 进行重试。
