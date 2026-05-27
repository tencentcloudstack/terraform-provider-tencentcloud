## Context

GA2（全球加速2.0）已有 Terraform 资源 `tencentcloud_ga2_endpoint_group`，服务层文件 `service_tencentcloud_ga2.go` 已存在。现需新增数据源 `tencentcloud_ga2_cross_border_settlement`，用于查询跨境账单流量用量。

云API `DescribeCrossBorderSettlement` 接口已在 vendor 中可用（`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115`），该接口为同步接口，无分页，入参为查询条件，出参仅返回一个 `Traffic` 浮点数字段。

## Goals / Non-Goals

**Goals:**
- 提供 `tencentcloud_ga2_cross_border_settlement` 数据源，允许用户通过 Terraform 查询指定加速实例、加速地域、终端节点组地域和账单月份的跨境流量用量
- 遵循项目现有的数据源代码模式（参考 `tencentcloud_igtm_instance_list`）
- 包含完整的单元测试（使用 gomonkey mock）

**Non-Goals:**
- 不提供写入/修改跨境账单的能力（API 仅支持查询）
- 不实现分页逻辑（该接口无分页参数）

## Decisions

1. **数据源文件位置**：放置在 `tencentcloud/services/ga2/data_source_tc_ga2_cross_border_settlement.go`，复用已有的 ga2 服务目录。

2. **Schema 设计**：所有入参（`global_accelerator_id`、`accelerate_region`、`endpoint_group_region`、`settlement_month`）设为 Required，出参 `traffic` 设为 Computed。同时添加 `result_output_file` 用于输出结果到文件。

3. **ID 设计**：由于是数据源且无唯一资源ID，使用入参组合作为 ID（`global_accelerator_id#accelerate_region#endpoint_group_region#settlement_month`）。

4. **重试机制**：使用 `tccommon.ReadRetryTimeout` + `resource.Retry` 包装 API 调用，失败时使用 `tccommon.RetryError()` 返回。

5. **测试方式**：使用 gomonkey mock 云 API 调用进行单元测试，不依赖真实环境。

## Risks / Trade-offs

- [Risk] `SettlementMonth` 为 `uint64` 类型，Terraform Schema 中使用 `TypeInt` 表示 → 在 Go 代码中进行类型转换（`uint64(d.Get("settlement_month").(int))`）
- [Risk] `Traffic` 为 `float64` 类型，Terraform Schema 中使用 `TypeFloat` 表示 → 直接映射，无精度损失风险
