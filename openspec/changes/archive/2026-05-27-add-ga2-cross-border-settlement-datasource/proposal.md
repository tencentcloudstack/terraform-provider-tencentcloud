## Why

需要为云产品 GA2（全球加速2.0）新增 Terraform 数据源 `tencentcloud_ga2_cross_border_settlement`，使用户能够通过 Terraform 查询跨境账单的流量用量信息，便于成本管理和监控。

## What Changes

- 新增数据源 `tencentcloud_ga2_cross_border_settlement`，调用 `DescribeCrossBorderSettlement` 接口查询跨境账单流量用量
- 入参支持：`global_accelerator_id`（全球加速实例ID）、`accelerate_region`（加速地域）、`endpoint_group_region`（终端节点组地域）、`settlement_month`（账单年月时间）
- 出参返回：`traffic`（流量用量，单位GB，精度保留小数点6位）
- 在 `provider.go` 和 `provider.md` 中注册该数据源

## Capabilities

### New Capabilities

- `cross-border-settlement-query`: 提供通过 Terraform 数据源查询 GA2 跨境账单流量用量的能力

### Modified Capabilities

（无）

## Impact

- 新增文件：`tencentcloud/services/ga2/data_source_tc_ga2_cross_border_settlement.go`
- 新增测试文件：`tencentcloud/services/ga2/data_source_tc_ga2_cross_border_settlement_test.go`
- 新增文档文件：`tencentcloud/services/ga2/data_source_tc_ga2_cross_border_settlement.md`
- 修改文件：`tencentcloud/provider.go`（注册数据源）
- 修改文件：`tencentcloud/provider.md`（添加数据源文档条目）
- 依赖：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115`（已在 vendor 中）
