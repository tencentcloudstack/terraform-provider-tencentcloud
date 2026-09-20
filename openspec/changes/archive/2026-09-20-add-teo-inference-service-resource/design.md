## Context

腾讯云 EdgeOne（TEO）在 `teo/v20220901` SDK 中新增了推理服务（Inference Service）相关 API。vendored SDK（`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/`）已暴露以下接口：

- `CreateInferenceServiceWithContext` → 返回 `{ ServiceId }`（同步，响应仅含 `ServiceId` + `RequestId`，无 `TaskId`）。
- `DescribeInferenceServicesWithContext` → 按 `ZoneId` + `Filters` 分页列表，`Limit` 最大值 `200`；响应 `Services` 为 `[]*InferenceService`。
- `ModifyInferenceServiceWithContext` → 同步修改；入参使用 `InferenceContainerConfigForModify` / `InferenceResourceConfigForModify` 等 ForModify 变体结构体。
- `OperateInferenceServiceWithContext` → 通用操作接口，`Operation` 取值 `Stop`/`Resume`/`Delete`；删除使用 `Operation=Delete`，删除后不可恢复。响应仅含 `RequestId`，无 `TaskId`。

provider 已具备 `UseTeoV20220901Client` 连接绑定与 `service_tencentcloud_teo.go` 服务层文件（其中已有大量 `DescribeTeoXxxById` 帮助函数，可参考其分页过滤模式）。

参考资源为 `tencentcloud_igtm_strategy`（单文件资源布局、复合 ID、`helper.String`/`helper.Int64`、`resource.Retry` 包裹、防御性 nil 检查等风格）。

## Goals / Non-Goals

**Goals:**
- 在 Terraform 中提供 TEO 推理服务的完整生命周期管理（create / read / update / delete / import）。
- schema 字段与 `CreateInferenceService` 入参一一对应（不重命名、不合并、不新增合成字段），同时暴露 `DescribeInferenceServices` 返回的只读状态字段。
- 代码风格与 `tencentcloud_igtm_strategy` 一致：单文件资源、每次 SDK 调用包裹 `resource.Retry`、响应 payload 防御性 nil 检查、复合 ID 使用 `tccommon.FILED_SP`。
- 文档与测试齐备：`.md` 示例（含 import 复合 ID 说明）、gomonkey mock 单元测试。

**Non-Goals:**
- 不实现推理服务的 datasource（`data_source_tc_teo_inference_services`），本变更仅资源。
- 不管理部署记录（`DescribeInferenceServiceDeploymentRecords`）、部署日志（`DescribeInferenceServiceDeploymentLogs`）、监控数据（`DescribeInferenceServiceMonitorData`）等周边接口。
- 不暴露 `OperateInferenceService` 的 `Stop`/`Resume` 操作为 Terraform 字段（这些属于运维动作而非声明式状态；如需后续支持，可在未来变更中以 `status` 字段驱动）。
- 不管理硬件规格查询（`DescribeInferenceHardwareSpecifications`）——`hardware_spec_id` 作为普通字符串入参由用户填写。

## Decisions

### D1. 复合资源 ID = `zoneId` + `serviceId`
**Why**: `CreateInferenceService` / `DescribeInferenceServices` / `ModifyInferenceService` / `OperateInferenceService` 均要求 `ZoneId` 作为必填入参。`ServiceId` 单独不足以支撑后续 Read/Update/Delete 的调用，必须同时持久化 `ZoneId`。采用 `tccommon.FILED_SP` 拼接的复合 ID 与 `tencentcloud_igtm_strategy`（`instanceId#strategyId`）完全一致。
**Alternative**: 仅存 `serviceId`，从 state 的 `zone_id` 字段读取。被否决——`zone_id` 在 schema 中是 `Required` 可变字段，若用户修改会导致 ID 与实际归属站点脱钩；复合 ID 更稳健，且 import 时可一次性带入两个值。

### D2. 同步接口判定，无任务轮询
**Why**: 核查 SDK 响应结构体——`CreateInferenceServiceResponseParams`、`ModifyInferenceServiceResponseParams`、`OperateInferenceServiceResponseParams` 均仅含 `RequestId`，不含 `TaskId`；client 函数注释也未标注异步。判定四个接口均为同步接口。Create/Update/Delete 成功后直接调用 Read 回填即可，无需 `DescribeTaskResult` 轮询。
**Alternative**: 预留 task 轮询逻辑。被否决——无 TaskId 可轮询，属过度设计。

### D3. schema 设计：列表型嵌套结构展开
**Why**: `containers` 在 SDK 中是 `[]*InferenceContainerConfig`（当前仅支持 1 个容器），`environment_variables` 是 `[]*InferenceEnvironmentVariable`，`request_paths` 是 `[]*string`。按项目代码生成要求第 13 条，列表型数据应展开平铺，每个元素的字段单独可 set/read。因此：
- `containers` → `TypeList`，`Elem: schema.Resource`，内含 `image_type`、`tcr_repository_config`（嵌套 `tcr_type`/`image`/`registry_id`/`region_name`）、`startup_command`、`environment_variables`（再嵌套 `key`/`value` 的 TypeList）。
- `resource_config` → `TypeList` 单元素，内含 `scaling_mode`、`hardware_spec`、`hardware_spec_id`、`hardware_config`（嵌套 `gpu_num`/`cpu_num`/`mem_size`/`disk_size`）、`auto_scaling_config`（嵌套 `min_instance_count` + `scaling_policies` 列表）、`manual_instance_config`（嵌套 `fixed_instance_count`）、`concurrency`。
- `affinity_config` → `TypeList` 单元素，内含 `switch`、`affinity_mode`、`session_id_affinity_config`（嵌套 `source`/`header_name`）。
- `scaling_policies` → `TypeList`，内含 `policy_name`、`policy_type`、`scheduled_scaling_policy`（嵌套 `scheduled_actions` 列表 + `effective_range` 单元素 + `time_zone`）；`scheduled_actions` 内含 `cron_expression`、`min_instance_count`；`effective_range` 内含 `effective_type`、`start_date`、`end_date`。

**Alternative**: 全部用 `TypeString` + jsonencode 透传。被否决——项目要求字段可被 terraform 单独 set/read，且参考 `tencentcloud_igtm_strategy` 的 `main_address_pool_set` 多层嵌套 `schema.Resource` 写法。

### D4. ForceNew 字段选择
**Why**: `zone_id` 与 `name` 设为 ForceNew。
- `zone_id`：Create/Describe/Modify/Delete 均需 ZoneId，且站点归属不可迁移 → ForceNew。
- `name`：`ModifyInferenceServiceRequest` 不含 `Name` 字段，名称创建后不可修改 → ForceNew。
其余顶层字段（`listen_port`、`containers`、`resource_config`、`affinity_config`、`request_paths`、`description`）均在 `ModifyInferenceServiceRequest` 中存在对应入参，可原地更新，不设 ForceNew。

### D5. Create 与 Modify 的结构体差异处理
**Why**: Create 用 `InferenceContainerConfig` / `InferenceResourceConfig`，Modify 用 `InferenceContainerConfigForModify` / `InferenceResourceConfigForModify`，字段集合不同：
- Modify 的 `InferenceResourceConfigForModify` 不含 `HardwareSpec` / `HardwareSpecId`（已废弃/不支持修改），`HardwareConfig` 类型为 `InferenceHardwareConfigForModify`（不含 `GPUNum`）。
- Create 的 `InferenceResourceConfig` 含 `HardwareSpec`（已废弃）、`HardwareSpecId`、`HardwareConfig`（`InferenceHardwareConfig`，含 `GPUNum`）。
因此 build helper 需分别为 Create 和 Modify 编写：`buildTeoInferenceContainers`（Create）与 `buildTeoInferenceContainersForModify`（Modify），resource_config 同理。schema 中保留 `hardware_spec`、`hardware_spec_id`、`gpu_num` 字段（Create 用），Modify 时忽略这些字段（不在 Modify 请求体中设置）。
**Alternative**: schema 仅保留 Modify 支持的字段。被否决——会导致 Create 必填的 `hardware_spec_id` 无法配置，破坏创建能力。

### D6. Read 回填与只读计算字段
**Why**: `DescribeInferenceServices` 返回的 `InferenceService` 含 Create 入参字段 + 只读状态字段。只读字段设为 `Computed`：
- `service_id`（Computed，同时存入复合 ID）
- `status`、`scaling_status`、`current_instance_count`、`inference_url`、`create_time`、`update_time`
- `affinity_config`：Create 入参含 `AffinityConfig`，但 `DescribeInferenceServices` 的 `InferenceService` 结构体**未包含** `AffinityConfig` 字段（核查确认），因此 `affinity_config` 及其子字段标记为 `Optional`（非 Computed），Read 时不回填（保持用户配置值）。这与「Create 入参 = schema 字段」原则一致。
Read 时对每个字段先判断 nil 再 `d.Set`；若 `response.Response.Services` 为空或匹配项缺失，先打印 `[CRUD] ... id=%s` 保留现场，再 `d.SetId("")`。

### D7. Describe 帮助函数 `DescribeTeoInferenceServiceById`
**Why**: 参考现有 `DescribeTeoL4ProxyById` / `DescribeTeoAccelerationDomainById` 模式：
- 在 `service_tencentcloud_teo.go` 新增 `DescribeTeoInferenceServiceById(ctx, zoneId, serviceId string) (*teo.InferenceService, error)`。
- 构建 `DescribeInferenceServicesRequest`（在 for 循环外构建），`Filters` 设为 `[{Name: "service-id", Values: [serviceId]}]`。
- 分页 `Limit=200`（SDK 注释标注最大值），循环外构建 request，循环内仅变 `Offset`。
- 对每页 SDK 调用包裹 `resource.Retry(tccommon.ReadRetryTimeout, ...)`。
- 严格相等校验 `*item.ServiceId == serviceId` 后返回。
- 未找到返回 `(nil, nil)`，由资源层 `d.SetId("")`。

### D8. 单文件资源布局
**Why**: 参考 `tencentcloud_igtm_strategy`，整个资源放在一个文件 `resource_tc_teo_inference_service.go`：package + imports → schema → Create/Read/Update/Delete → build/flatten helper。服务层 helper 放在 `service_tencentcloud_teo.go`。不拆分 `_crud.go` / `_helpers.go`。

### D9. 单元测试使用 gomonkey mock
**Why**: 项目代码生成要求第 1 条——新增 terraform 资源的 `*_test.go` 不使用 terraform 测试套件，而是用 gomonkey mock 云 API 进行业务逻辑单元测试。测试覆盖 Create/Read/Update/Delete 各路径，mock `UseTeoV20220901Client().CreateInferenceServiceWithContext` 等方法返回预设响应。

## Risks / Trade-offs

- **[Risk]** `affinity_config` 在 Describe 响应中缺失，Read 无法回填，若用户在配置外手动修改会导致 drift 不可见。→ **Mitigation**: schema 标记 `affinity_config` 为 `Optional`（非 Computed），Terraform 以用户配置为准；在 `.md` 文档中说明该字段为只写（write-only），不参与 Read 回填。
- **[Risk]** Modify 接口的 `InferenceResourceConfigForModify` 不含 `HardwareSpecId`，用户若尝试修改 `hardware_spec_id` 不会生效。→ **Mitigation**: `hardware_spec_id` 保持非 ForceNew（避免误触重建），但 Update 时不将其写入 Modify 请求体；在文档中说明硬件规格 ID 仅在创建时生效。
- **[Risk]** `containers` 当前 SDK 仅支持 1 个容器，schema 用 `TypeList` 可容纳多个，用户配置多个时 API 会报错。→ **Mitigation**: 依赖服务端校验，错误经 `tccommon.RetryError` 透传给用户；文档注明「当前仅支持设置 1 个容器」。
- **[Trade-off]** 复合 ID 使 import 需提供 `zoneId#serviceId` 两个值。→ 已在 `.md` import 示例中说明，与 `tencentcloud_igtm_strategy` 一致，用户可预期。

## Migration Plan

纯新增变更，无需 state 迁移：
1. 落地新资源 + 服务层 helper + provider 注册。
2. 发布后，用户在配置中新增 `resource "tencentcloud_teo_inference_service" "x" { ... }` 即可使用。
3. 已有 TEO 资源不受影响。

回滚：纯 revert 新增文件 + `provider.go` 注册行，无 state 需撤销。

## Open Questions

- 无需用户输入的问题。SDK 暴露了全部所需 API；复合 ID、schema 平铺、同步接口判定等决策已由 SDK 结构体核查完全确定。