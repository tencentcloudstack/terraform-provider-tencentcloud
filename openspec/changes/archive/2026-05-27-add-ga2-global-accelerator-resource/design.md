## Context

当前 Terraform Provider 已有 GA2 服务的基础设施：
- `tencentcloud/services/ga2/service_tencentcloud_ga2.go` 提供了 `Ga2Service`、`WaitForGa2TaskFinish` 等公共方法
- `tencentcloud/services/ga2/resource_tc_ga2_endpoint_group.go` 是同产品下已有的资源实现，可作为模式参考
- vendor 中已包含 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115` SDK

所有 GA2 写操作（Create/Modify/Delete）均为异步接口，返回 TaskId，需通过 `DescribeTaskResult` 轮询任务状态直到 SUCCESS。已有的 `WaitForGa2TaskFinish` 方法可直接复用。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_ga2_global_accelerator` 资源的完整 CRUD
- 复用已有的 `Ga2Service` 和 `WaitForGa2TaskFinish` 异步等待机制
- 支持 Import 功能（通过 global_accelerator_id 导入）
- 在 service 层新增 `DescribeGa2GlobalAcceleratorById` 方法供 Read 使用
- 使用 gomonkey mock 方式编写单元测试

**Non-Goals:**
- 不实现 tags 的独立更新（ModifyGlobalAccelerator 不支持 tags 参数）
- 不实现 instance_charge_type 的变更（仅创建时指定）
- 不实现数据源（本次仅新增 RESOURCE_KIND_GENERAL 资源）

## Decisions

### 1. 资源 ID 使用 GlobalAcceleratorId

**决策**: 使用 `CreateGlobalAccelerator` 返回的 `GlobalAcceleratorId` 作为 Terraform 资源 ID（单一 ID，无需复合 ID）。

**理由**: 该 ID 是全球加速实例的唯一标识，Modify 和 Delete 接口均只需此 ID。

### 2. 异步操作使用已有的 WaitForGa2TaskFinish

**决策**: Create/Modify/Delete 操作完成后，调用 `Ga2Service.WaitForGa2TaskFinish` 等待任务完成。

**理由**: 该方法已在 endpoint_group 资源中验证可用，使用 `DescribeTaskResult` 接口轮询，避免重复实现。

### 3. instance_charge_type 和 tags 设为 ForceNew

**决策**: `instance_charge_type` 和 `tags` 字段设为 ForceNew，因为 ModifyGlobalAccelerator 接口不支持修改这两个字段。

**理由**: 云 API 的 Modify 接口仅支持 Name、Description、CrossBorderType、CrossBorderPromiseFlag 四个字段的修改。

### 4. Read 方法通过 DescribeGlobalAccelerators + Filter 实现

**决策**: 在 service 层新增 `DescribeGa2GlobalAcceleratorById` 方法，使用 `global-accelerator-id` Filter 查询单个实例。

**理由**: DescribeGlobalAccelerators 接口支持通过 Filters 按 ID 过滤，与 endpoint_group 的模式一致。

### 5. Schema 中声明 Timeouts

**决策**: 为 Create/Update/Delete 声明 20 分钟超时，与 endpoint_group 资源保持一致。

**理由**: 异步操作需要超时控制，20 分钟是同产品其他资源的标准超时时间。

## Risks / Trade-offs

- [Tags 不可更新] → 设为 ForceNew，用户修改 tags 会触发资源重建。这是 API 限制，无法规避。
- [instance_charge_type 不可更新] → 设为 ForceNew，与 API 能力对齐。
- [异步操作可能超时] → 使用 schema Timeouts 机制，用户可自定义超时时间。
