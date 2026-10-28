## Why

Terraform Provider for TencentCloud 目前缺少创建 KMS（Key Management Service）托管实例的能力。用户无法通过 Terraform 购买开通 KMS 专业版预付费实例，只能手动在控制台下单，导致基础设施即代码流程断裂。需要新增一个 terraform-plugin-framework **action**，封装云 API `CreateKmsInstance`，让用户在 Terraform 中一次性完成 KMS 实例的购买开通操作。

## What Changes

- 新增 framework action `tencentcloud_kms_instance_operation`（terraform-plugin-framework action 类型），封装 KMS 云 API 接口 `CreateKmsInstance`，支持创建预付费 KMS 实例。
  - 支持 instance_charge_type（付费类型）、product_version（版本类型）、client_token（可重入客户端 ID）等顶层属性。
  - 支持 instance_charge_prepaid 嵌套 block（SingleNestedBlock，含 period 购买时长、renew_flag 续费标记）。
  - 支持 tags 标签列表（ListNestedBlock，tag_key / tag_value）。
  - 创建成功后将 instance_id、deal_name、bill_id 通过日志输出（framework action 无法对外暴露输出属性）。
  - action 仅实现 Invoke，无 Read/Update/Delete（一次性操作，不记录状态）。
- 在 `tencentcloud/framework/registry.go` 的 `actionFactories` 切片中注册新 action 工厂 `kms.NewKmsInstanceOperation`。
- 不在 `tencentcloud/provider.go` 中注册（action 由 framework provider 管理而非 SDKv2 provider）。
- 在 `tencentcloud/provider.md` 的 Key Management Service(KMS) 区块新增 `Action` 子区段，追加 `tencentcloud_kms_instance_operation` 条目。
- 新增资源文档文件 `tencentcloud/services/kms/action_tc_kms_instance_operation.md`。
- 新增单元测试文件 `tencentcloud/services/kms/action_tc_kms_instance_operation_test.go`，使用 gomonkey mock 云 API。

## Capabilities

### New Capabilities
- `kms-instance-operation`: 通过 Terraform framework action 一次性操作创建（购买开通）KMS 预付费托管实例，封装 CreateKmsInstance 云 API。

### Modified Capabilities
<!-- 无需修改现有能力 -->

## Impact

- 新增文件：
  - `tencentcloud/services/kms/action_tc_kms_instance_operation.go`（framework action 实现）
  - `tencentcloud/services/kms/action_tc_kms_instance_operation.md`（action 文档）
  - `tencentcloud/services/kms/action_tc_kms_instance_operation_test.go`（单元测试）
  - `tencentcloud/services/kms/service_tc_kms_instance.go`（service 方法，含 retry）
- 修改文件：
  - `tencentcloud/framework/registry.go`（注册新 action 工厂）
  - `tencentcloud/provider.md`（KMS 区块新增 Action 子区段）
  - `website/tencentcloud.erb`（侧边栏增加 Action 条目）
  - `.changelog/4616.txt`（更新变更说明）
- 需 HashiCorp Terraform 1.14 及以上版本支持 action。
- 依赖云 API：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118` 的 `CreateKmsInstance` 接口（已在 vendor 中）。
- 不破坏向后兼容：新增 framework action，不修改已有 SDKv2 资源。