## Context

当前 Terraform Provider for TencentCloud 已经支持云拨测（CAT）服务的多个数据源，包括 `tencentcloud_cat_metric_data`、`tencentcloud_cat_probe_data` 和 `tencentcloud_cat_node`。其中 `tencentcloud_cat_metric_data` 用于查询拨测指标数据，依赖维度标签值来构建过滤条件。

然而，用户在构建维度标签过滤条件时，无法预先获知可用的维度标签值列表。CAT API 提供了 `DescribeProbeMetricTagValues` 接口来查询维度标签值，但该接口尚未被 Terraform Provider 封装。

本设计基于现有的 `data_source_tc_cat_metric_data.go` 实现模式，新增 `tencentcloud_cat_probe_metric_tag_values` 数据源，封装该接口。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_cat_probe_metric_tag_values` 数据源，支持通过分析任务类型、维度标签、过滤条件和时间范围查询维度标签值
- 遵循现有的 CAT 数据源代码模式（参考 `data_source_tc_cat_metric_data.go`）
- 在 service 层新增 `DescribeCatProbeMetricTagValuesByFilter` 方法
- 在 Provider 中注册新数据源
- 生成对应的文档和测试文件

**Non-Goals:**
- 不修改现有的 `tencentcloud_cat_metric_data` 数据源
- 不涉及资源的创建、更新和删除操作（本数据源为只读查询）
- 不暴露分页参数（该接口无分页字段）

## Decisions

### 1. 数据源 Schema 设计

**决策**: 所有输入参数（`analyze_task_type`、`key`、`filter`、`filters`、`time_range`）均设为 Optional，与云 API 入参的可选性一致。`filters` 使用 `schema.TypeSet`，与 `tencentcloud_cat_metric_data` 中的 `filters` 字段保持一致。

**理由**: 云 API 的所有入参都是可选的，用户可以根据需要组合查询条件。

### 2. Service 层方法设计

**决策**: 在 `service_tencentcloud_cat.go` 中新增 `DescribeCatProbeMetricTagValuesByFilter` 方法，采用 `paramMap` 传参模式，与现有的 `DescribeCatMetricDataByFilter` 方法风格一致。

**理由**: 保持与现有 CAT 服务层代码风格的一致性，便于维护。

### 3. Read 函数的 retry 和空响应处理

**决策**: 在 Read 函数中使用 `resource.Retry` 包裹 API 调用，超时时间为 `tccommon.ReadRetryTimeout`。在 retry 块内检查响应是否为空，若为空则返回 `NonRetryableError`，避免因云 API 短暂波动导致本地 state 中的 id 被清空。

**理由**: 遵循 RESOURCE_KIND_DATASOURCE 资源代码生成的规范要求，确保数据源在查询返回空时不丢失状态。

### 4. ID 生成策略

**决策**: 使用 `helper.DataResourceIdsHash([]string{tagValueSet})` 生成数据源 ID，与 `tencentcloud_cat_metric_data` 的 ID 生成策略一致。

**理由**: 数据源没有真实的资源 ID，使用返回内容的哈希作为 ID 可以确保同一查询结果在多次执行间保持一致性。

### 5. 结果输出文件

**决策**: 支持 `result_output_file` 可选参数，用于将查询结果保存到文件。

**理由**: 与现有 CAT 数据源（`tencentcloud_cat_metric_data`）保持一致，方便用户导出查询结果。

## Risks / Trade-offs

- **[接口返回值格式]** → `TagValueSet` 是序列化后的字符串，用户需要自行解析 JSON 内容。这是云 API 的设计决定，terraform provider 不做额外解析，保持原始返回。
- **[无分页支持]** → 该接口无分页字段，不存在分页风险。
- **[所有参数可选]** → 由于所有参数均为可选，用户在不提供任何参数的情况下调用可能导致返回结果不符合预期。通过在文档中说明各参数的用途来缓解。
