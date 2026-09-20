## Context

Terraform Provider for TencentCloud 当前已支持 TEO（EdgeOne）的多种数据源（如 `tencentcloud_teo_config_group_versions`、`tencentcloud_teo_multi_path_gateway_origin_acl` 等），但尚未提供查询推理服务（Inference Service）部署历史的数据源。云 API `DescribeInferenceServiceDeploymentRecords`（位于 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`）已可用于查询指定站点和推理服务下的部署历史列表，返回每次部署的操作类型、状态、耗时、配置快照及是否为当前生效配置。

该 API 为同步查询接口（非异步接口），支持分页（`Offset` / `Limit`，`Limit` 最大值 100）与排序（`SortBy` / `SortOrder`）。响应结构 `RecordSet` 为 `InferenceServiceDeploymentRecord` 数组，每个记录内嵌 `InferenceServiceConfig`，其中又包含容器配置、资源配置、亲和性配置等深层嵌套结构。

现有可参考的模式包括：`data_source_tc_teo_config_group_versions.go`（列表型数据源 + 服务层分页）和 `data_source_tc_teo_multi_path_gateway_origin_acl.go`（深层嵌套结构 nil 安全）。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_inference_service_deployment_records` 数据源，调用 `DescribeInferenceServiceDeploymentRecords` API
- 遵循现有 TEO 数据源模式（如 `data_source_tc_teo_config_group_versions.go`）
- 将 `RecordSet` 列表展开为 `record_set`（TypeList），每个元素平铺所有字段，不额外嵌套一层 "列表型数据" schema
- 正确处理深层嵌套响应结构（`InferenceServiceConfig` → `Containers` / `ResourceConfig` / `AffinityConfig`），每一层均做 nil 安全检查
- 服务层内部实现自动分页（`Limit` 取云 API 最大值 100），不向用户暴露 `limit`/`offset` 参数
- 在 `provider.go` 和 `provider.md` 中注册数据源
- 使用 gomonkey mock 方式编写单元测试（新增数据源不使用 TF 测试套件）
- 新增文档 `.md` 文件

**Non-Goals:**
- 不提供 CRUD 操作（这是只读数据源）
- 不修改已有的任何 TEO 资源或数据源
- 不向用户暴露分页参数 `limit` / `offset`

## Decisions

1. **Schema 设计遵循云 API 结构并将列表展开**：`record_set` 定义为 `TypeList`（Computed），其 `Elem` 为 `schema.Resource`，内部平铺 `record_id`、`operation`、`status`、`duration`、`create_time`、`active_status` 以及 `inference_service_config` 等字段。遵循规则要求，不创建把所有字段再嵌套一层的 `xxx_set`/`xxx_list` 结构（`record_set` 本身即为展开后的列表）。

2. **深层嵌套字段使用 TypeList + MaxItems:1 表示单对象**：`inference_service_config`、`containers`（数组，无 MaxItems）、`tcr_repository_config`、`resource_config`、`hardware_config`、`auto_scaling_config`、`scaling_policies`（数组）、`scheduled_scaling_policy`、`manual_instance_config`、`affinity_config`、`session_id_affinity_config` 等结构按其本身是否为数组分别使用 `TypeList`。其中单对象使用 `TypeList` 且不限制 MaxItems（与参考数据源风格一致），数组类型如 `containers`、`scaling_policies`、`environment_variables`、`request_paths` 为多元素列表。

3. **数值类型映射**：`duration`、`listen_port`、`gpu_num`、`cpu_num`、`mem_size`、`disk_size`、`min_instance_count`、`fixed_instance_count`、`concurrency` 为数值，使用 `schema.TypeInt`；其中 `gpu_num`、`cpu_num` 在 SDK 中为 `*float64`，读取时通过类型转换处理。`request_paths` 为字符串列表使用 `TypeList` + `TypeString` 元素。

4. **复合 ID**：数据源 ID 使用 `zone_id + FILED_SP + service_id` 作为复合标识，便于定位与排障。

5. **服务方法分页**：新增 `DescribeTeoInferenceServiceDeploymentRecordsByFilter` 到 `service_tencentcloud_teo.go`，使用 `paramMap` 传入 `ZoneId`、`ServiceId`、`SortBy`、`SortOrder`，内部循环分页（`Offset` 递增、`Limit` 固定 100），合并所有 `RecordSet` 后返回。分页不暴露给用户。

6. **Read 函数 retry 与空响应处理**：在 `dataSourceTencentCloudTeoInferenceServiceDeploymentRecordsRead` 的 retry 块内，调用服务方法；若返回为空（`respData == nil` 或 `len == 0`），按 DATASOURCE 规则返回 `NonRetryableError`，避免直接 `d.SetId("")` 清空 state，并在失败路径保留 `log.Printf("[DATASOURCE] read empty, skip SetId")` 日志。

7. **nil 安全**：在 set 每个字段前判断对应 SDK 字段是否为 nil；嵌套结构（如 `InferenceServiceConfig`、`Containers[i]`、`ResourceConfig`、`AffinityConfig` 等）在访问子字段前先判断父对象是否为 nil。

## Risks / Trade-offs

- **深层嵌套 nil 安全**：响应结构嵌套较深（`RecordSet` → `InferenceServiceConfig` → `Containers` → `TcrRepositoryConfig`/`EnvironmentVariables`，`ResourceConfig` → `HardwareConfig`/`AutoScalingConfig` → `ScalingPolicies` → `ScheduledScalingPolicy`，`AffinityConfig` → `SessionIdAffinityConfig`），任意一层可能为 nil。→ 缓解：每一层访问前均做 nil 判断，参考 `data_source_tc_teo_multi_path_gateway_origin_acl.go` 的模式。
- **float64 数值字段**：`gpu_num`、`cpu_num` 为 `*float64`，需转换为 int 后再 set，避免类型不匹配。→ 缓解：使用 `int(*field)` 转换。
- **数据源空响应**：云 API 可能因短暂波动返回空列表，直接清空 ID 会导致 state 数据丢失。→ 缓解：按规则返回 `NonRetryableError` 让外层重试，并保留日志。
- **分页边界**：需正确判断分页终止条件（返回条数小于 limit 时停止），避免无限循环或漏取数据。→ 缓解：沿用 `DescribeTeoConfigGroupVersionsByFilter` 的分页循环模式。