## Why

Terraform Provider for TencentCloud 目前缺少创建 KMS（Key Management Service）托管实例的能力。用户无法通过 Terraform 购买开通 KMS 专业版预付费实例，只能手动在控制台下单，导致基础设施即代码流程断裂。需要新增一个 terraform-plugin-framework **action**，封装云 API `CreateKmsInstance`，让用户在 Terraform 中一次性完成 KMS 实例的购买开通操作。

## What Changes

- 新增 framework action `tencentcloud_kms_instance`（terraform-plugin-framework action 类型），封装 KMS 云 API 接口 `CreateKmsInstance`，支持创建预付费 KMS 实例。action 命名沿用 provider 内既有 framework action 规范（`action_tc_<product>_<name>.go` / `tencentcloud_<product>_<name>`，不带 `_operation` 后缀，与 `tencentcloud_antiddos_unblock_resources`、`tencentcloud_mongodb_restore_db_instance` 等保持一致）。
  - 支持 instance_charge_type（付费类型）、product_version（版本类型）等顶层属性。
  - 不对外暴露 `client_token` 属性：provider 内部使用 `clientToken := helper.BuildToken()` 自动生成可重入客户端 ID 并透传给云 API，用户无需也无法配置该项。
  - 支持 instance_charge_prepaid 嵌套 block（SingleNestedBlock，含 period 购买时长、renew_flag 续费标记）。
  - 支持 tags 标签列表（ListNestedBlock，tag_key / tag_value）。
  - 创建成功后将 instance_id、deal_name、bill_id 通过日志输出（framework action 无法对外暴露输出属性）；若云 API 返回空 Response 或空 InstanceId，则打印现场日志并返回 error diagnostic，不报告成功。
  - action 仅实现 Invoke，无 Read/Update/Delete（一次性操作，不记录状态）。
- 在 `tencentcloud/framework/registry.go` 的 `actionFactories` 切片中注册新 action 工厂 `kms.NewKmsInstance`。
- 不在 `tencentcloud/provider.go` 中注册（action 由 framework provider 管理而非 SDKv2 provider）。
- 在 `tencentcloud/provider.md` 的 Key Management Service(KMS) 区块新增 `Action` 子区段，追加 `tencentcloud_kms_instance` 条目。
- 云 API 调用封装在既有 service 文件 `tencentcloud/services/kms/service_tencentcloud_kms.go` 的 `KmsService.CreateKmsInstance` 方法中（不新增额外 service 文件，保持 `service_tencentcloud_<product>.go` 的仓库约定），使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` + `tccommon.RetryError()`。
- 新增 action 文档文件 `tencentcloud/services/kms/action_tc_kms_instance.md`。
- 新增单元测试文件 `tencentcloud/services/kms/action_tc_kms_instance_test.go`，使用 gomonkey mock 云 API。

## Capabilities

### New Capabilities
- `kms-instance-operation`: 通过 Terraform framework action 一次性操作创建（购买开通）KMS 预付费托管实例，封装 CreateKmsInstance 云 API。

### Modified Capabilities
<!-- 无需修改现有能力 -->

## Impact

- 新增文件：
  - `tencentcloud/services/kms/action_tc_kms_instance.go`（framework action 实现）
  - `tencentcloud/services/kms/action_tc_kms_instance.md`（action 文档）
  - `tencentcloud/services/kms/action_tc_kms_instance_test.go`（单元测试）
- 修改文件：
  - `tencentcloud/services/kms/service_tencentcloud_kms.go`（新增 `CreateKmsInstance` service 方法，含 retry）
  - `tencentcloud/framework/registry.go`（注册新 action 工厂 `kms.NewKmsInstance`）
  - `tencentcloud/provider.md`（KMS 区块新增 Action 子区段）
  - `website/tencentcloud.erb`（由 `make doc` 生成，侧边栏增加 Action 条目）
  - `.changelog/4616.txt`（new-resource 条目 `tencentcloud_kms_instance`）
- 需 HashiCorp Terraform 1.14 及以上版本支持 action。
- 依赖云 API：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118` 的 `CreateKmsInstance` 接口（已在 vendor 中）。
- 不破坏向后兼容：新增 framework action，不修改已有 SDKv2 资源。
