## Why

Terraform Provider for TencentCloud 目前缺少创建 KMS（Key Management Service）托管实例的资源。用户无法通过 Terraform 购买开通 KMS 专业版预付费实例，只能手动在控制台下单，导致基础设施即代码流程断裂。需要新增一个 RESOURCE_KIND_OPERATION 资源，封装云 API `CreateKmsInstance`，让用户在 Terraform 中一次性完成 KMS 实例的购买开通操作。

## What Changes

- 新增资源 `tencentcloud_kms_kms_instance_operation`（RESOURCE_KIND_OPERATION 类型），封装 KMS 云 API 接口 `CreateKmsInstance`，支持创建预付费 KMS 实例。
  - 支持 instance_charge_type（付费类型）、product_version（版本类型）、client_token（可重入客户端 ID）等顶层参数。
  - 支持 instance_charge_prepaid 嵌套结构（period 购买时长、renew_flag 续费标记）。
  - 支持 tags 标签列表（tag_key / tag_value）。
  - 返回 instance_id、deal_name、bill_id 作为 computed 属性。
  - 资源只实现 Create 接口，Read/Update/Delete 为空（一次性操作，不记录状态）。
- 在 `tencentcloud/provider.go` 中注册新资源 `tencentcloud_kms_kms_instance_operation`。
- 在 `tencentcloud/provider.md` 的 KMS Resource 区块追加新资源条目。
- 新增资源文档文件 `tencentcloud/services/kms/resource_tc_kms_kms_instance_operation.md`。
- 新增单元测试文件 `tencentcloud/services/kms/resource_tc_kms_kms_instance_operation_test.go`，使用 gomonkey mock 云 API。

## Capabilities

### New Capabilities
- `kms-instance-operation`: 通过 Terraform 一次性操作创建（购买开通）KMS 预付费托管实例，封装 CreateKmsInstance 云 API。

### Modified Capabilities
<!-- 无需修改现有能力 -->

## Impact

- 新增文件：
  - `tencentcloud/services/kms/resource_tc_kms_kms_instance_operation.go`（资源 CRUD 实现）
  - `tencentcloud/services/kms/resource_tc_kms_kms_instance_operation.md`（资源文档）
  - `tencentcloud/services/kms/resource_tc_kms_kms_instance_operation_test.go`（单元测试）
- 修改文件：
  - `tencentcloud/provider.go`（注册新资源）
  - `tencentcloud/provider.md`（文档列表追加条目，最终由 `make doc` 生成）
- 依赖云 API：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118` 的 `CreateKmsInstance` 接口（已在 vendor 中）。
- 不破坏向后兼容：仅新增资源，不修改已有资源的 schema。