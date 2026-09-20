## Why

腾讯云 EdgeOne（TEO）已推出「推理服务（Inference Service）」能力，支持用户在边缘节点上部署容器化 AI 模型并对外提供推理访问入口。当前 Terraform Provider 尚未覆盖该资源，用户只能在控制台手动创建/调整推理服务，无法纳入基础设施即代码的声明式运维流程。新增 `tencentcloud_teo_inference_service` 资源，可在 Terraform 中完整管理推理服务的创建、读取、更新、删除与导入全生命周期，使 TEO 推理服务可被版本化、可复用、可审计地交付。

## What Changes

- 新增资源 `tencentcloud_teo_inference_service`，文件 `tencentcloud/services/teo/resource_tc_teo_inference_service.go`，基于 `teo` v20220901 SDK。
- 实现 CRUD 全生命周期：
  - **Create**：调用 `CreateInferenceService`，返回 `ServiceId` 作为资源 ID（同步接口，无 TaskId）。
  - **Read**：在 `service_tencentcloud_teo.go` 新增 `DescribeTeoInferenceServiceById(ctx, zoneId, serviceId)` 帮助函数，封装 `DescribeInferenceServices`（按 `service-id` 过滤、`Limit=200` 分页），获取单个推理服务详情。
  - **Update**：调用 `ModifyInferenceService`（同步接口），按 diff 组织请求体；注意 Modify 使用 `InferenceContainerConfigForModify` / `InferenceResourceConfigForModify` 等 ForModify 变体结构体，且 Modify 接口不再接受 `HardwareSpecId`、`HardwareSpec` 字段，HardwareConfig 也变更为 `InferenceHardwareConfigForModify`（不含 GPUNum）。
  - **Delete**：调用 `OperateInferenceService`，`Operation=Delete`（同步接口，删除后不可恢复）。
- 资源 ID 采用 `zoneId` + `serviceId` 的复合 ID（使用 `tccommon.FILED_SP` 分隔），因为 Create/Describe/Modify/Delete 接口均需要 `ZoneId` 入参；导入时需使用该复合 ID。
- schema 完整映射 `CreateInferenceService` 全部入参，并补齐 `DescribeInferenceServices` 返回的只读计算字段（`status`、`scaling_status`、`current_instance_count`、`inference_url`、`create_time`、`update_time`、`service_id`、`affinity_config` 等）。
- 在 `tencentcloud/provider.go` 的 `ResourcesMap` 注册 `tencentcloud_teo_inference_service`。
- 编写 `.md` 示例文档 `resource_tc_teo_inference_service.md`（含 example HCL 与 `terraform import` 复合 ID 说明）。
- 编写单元测试 `resource_tc_teo_inference_service_test.go`，使用 gomonkey mock 云 API 进行业务逻辑测试（新增资源不使用 terraform 测试套件）。
- 代码风格严格参考 `tencentcloud_igtm_strategy` 资源。

## Capabilities

### New Capabilities
- `teo-inference-service-resource`: 通过 Terraform 管理腾讯云 EdgeOne 推理服务的全生命周期（create/read/update/delete/import），完整映射 `CreateInferenceService`/`ModifyInferenceService` 入参，并暴露 `DescribeInferenceServices` 返回的只读状态字段。

### Modified Capabilities
<!-- 本次为纯新增资源，不修改任何已有 capability 的 requirement -->

## Impact

- **新增代码**:
  - `tencentcloud/services/teo/resource_tc_teo_inference_service.go`：schema + CRUD + build/flatten helper（单文件，参考 `tencentcloud_igtm_strategy` 风格）。
  - `tencentcloud/services/teo/resource_tc_teo_inference_service.md`：资源示例文档 + import 复合 ID 说明。
  - `tencentcloud/services/teo/resource_tc_teo_inference_service_test.go`：gomonkey mock 单元测试。
- **修改代码**:
  - `tencentcloud/services/teo/service_tencentcloud_teo.go`：新增 `DescribeTeoInferenceServiceById(ctx, zoneId, serviceId)` 帮助函数。
  - `tencentcloud/provider.go`：在 `ResourcesMap` 注册 `tencentcloud_teo_inference_service`。
- **消费的 API**：`CreateInferenceService`、`DescribeInferenceServices`、`ModifyInferenceService`、`OperateInferenceService`（均已在 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/` 中确认存在）。
- **无破坏性变更**：纯新增，不修改任何已有 schema 或 state。
- **无需 SDK 升级**：所需 API 均已存在于 vendored SDK 中。
- **接口同步性说明**：经核查 SDK，`CreateInferenceService`/`ModifyInferenceService`/`OperateInferenceService` 响应均仅返回 `RequestId`，无 `TaskId`，判定为同步接口，无需任务轮询；Create/Update 后调用 Read 即可。
- **向后兼容**：纯新增资源，无影响。