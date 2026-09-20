## Context

TEO（EdgeOne）提供了推理服务能力，用户在创建推理服务时需指定硬件规格（`HardwareSpecId`）。硬件规格 ID 需通过 `DescribeInferenceHardwareSpecifications` 接口按站点（ZoneId）查询获取。当前 Terraform Provider 已有多个 TEO 数据源（如 `tencentcloud_teo_config_group_versions`、`tencentcloud_teo_zones` 等），但缺少推理硬件规格查询数据源。

该接口为同步只读接口，无分页参数，入参仅 `ZoneId`，出参为 `HardwareSpecifications` 列表。数据源类型为 `RESOURCE_KIND_DATASOURCE`，业务逻辑代码风格参考 `tencentcloud_teo_config_group_versions`（返回列表的 TEO 数据源），并遵循数据源 Read 重试与空返回校验约束。

## Goals / Non-Goals

**Goals:**
- 新增数据源 `tencentcloud_teo_inference_hardware_specifications`，支持按站点 ID 查询推理硬件规格列表
- 正确映射云 API `DescribeInferenceHardwareSpecifications` 的入参（`ZoneId`）与出参（`HardwareSpecifications` 列表及其内部字段）
- 在 provider 中注册该数据源并提供文档
- 通过 gomonkey mock 云 API 的方式编写单元测试（不依赖 Terraform 测试套件）

**Non-Goals:**
- 不实现推理服务的创建/更新/删除（属于其他资源范畴）
- 不对返回结果做排序、过滤等额外加工
- 不涉及异步轮询（接口本身为同步接口）

## Decisions

### 1. 出参列表展开方式
采用与 `tencentcloud_teo_config_group_versions` 一致的结构：将 `HardwareSpecifications` 列表作为 `TypeList` + `schema.Resource` 的 Computed 字段（`hardware_specifications`），列表内每个元素的字段平铺为该 Resource 的子 Schema（`spec`、`hardware_spec_id`、`name`、`gpu_num`、`cpu_num`、`mem_size`、`gpu_mem_size`、`disk_size`、`allowed_gpu_nums`）。

**理由**：该接口出参本身就是列表结构，遵循现有 TEO 列表数据源的既有实现风格，保持一致性且每个字段可被 Terraform 单独读取。不创建多余的"列表型数据"再嵌套一层的 schema。

### 2. 数据类型映射
- `GPUNum`、`CPUNum`（云 API 为 `float64`）→ schema.TypeFloat
- `MemSize`、`GPUMemSize`、`DiskSize`（云 API 为 `int64`）→ schema.TypeInt
- `AllowedGPUNums`（云 API 为 `[]*float64`）→ schema.TypeList + TypeFloat 元素
- `Spec`、`HardwareSpecId`、`Name`（`string`）→ schema.TypeString

**理由**：严格按 vendor 中 SDK 字段类型对应 Terraform schema 类型。

### 3. Read 重试与空返回处理
在 `dataSourceTencentCloudTeoInferenceHardwareSpecificationsRead` 中使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包装调用 service 层方法。在 retry 块内部对返回空（`response == nil` 或 `len(HardwareSpecifications) == 0`）返回 `NonRetryableError`，避免因云 API 短暂波动导致 state 被清空；并在 retry 失败路径保留 `log.Printf("[DATASOURCE] read empty, skip SetId")` 提示。

**理由**：遵循 RESOURCE_KIND_DATASOURCE 数据源 Read 的空返回校验约束。

### 4. Service 层方法
在 `service_tencentcloud_teo.go` 中新增 `DescribeTeoInferenceHardwareSpecificationsByFilter` 方法，接收 `paramMap`，构造 `DescribeInferenceHardwareSpecificationsRequest` 并设置 `ZoneId`，调用 `DescribeInferenceHardwareSpecificationsWithContext`，返回 `[]*teov20220901.InferenceHardwareSpecification`。

### 5. SetId 策略
最终一致性设置 `d.SetId(helper.DataResourceIdsHash(ids))`，`ids` 使用各规格的 `HardwareSpecId` 列表生成哈希。

### 6. 测试策略
使用 gomonkey mock `DescribeTeoInferenceHardwareSpecificationsByFilter` 返回构造数据，校验 flatten 后字段是否正确。不使用 Terraform 测试套件。

## Risks / Trade-offs

- [Spec 字段已废弃] → 保留映射以兼容云 API 返回，但在 schema Description 中注明已废弃，建议使用 `hardware_spec_id`。无需特殊处理。
- [数据源无唯一 ID] → 使用 `helper.DataResourceIdsHash(ids)` 基于返回的 HardwareSpecId 列表生成哈希作为数据源 ID，符合现有数据源惯例。
- [字段类型 float64 可能出现精度问题] → GPUNum/CPUNum/AllowedGPUNums 为 float64，Terraform schema.TypeFloat 可承载，与云 API 定义一致，无精简转换。