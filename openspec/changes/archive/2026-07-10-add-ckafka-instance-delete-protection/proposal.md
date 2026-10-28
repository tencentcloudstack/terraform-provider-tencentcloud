## Why

CKafka 实例在控制台支持"删除保护"开关，开启后可防止实例被误删除。当前 `tencentcloud_ckafka_instance` 资源未暴露该开关，用户无法通过 Terraform 声明式地管理删除保护状态，只能在控制台或 SDK 手动设置，导致 IaC 覆盖不完整。腾讯云 CKafka SDK 的 `ModifyInstanceAttributes` 接口已支持 `DeleteProtectionEnable` 入参（1 开启、0 关闭），且 `DescribeInstanceAttributes` 接口响应中返回该字段，具备完整支持条件。

## What Changes

- 在 `tencentcloud_ckafka_instance` 资源 schema 中新增 `delete_protection_enable` 参数（`Optional + Computed`，TypeInt，取值 `1` 开启 / `0` 关闭），保持向后兼容（未配置时不触发 plan diff）。
- 在 Create 流程末尾的 `ModifyInstanceAttributes` 调用中，当用户配置了 `delete_protection_enable` 时，填充 `request.DeleteProtectionEnable`。
- 在 Update 流程中新增 `d.HasChange("delete_protection_enable")` 分支，构造 `ModifyInstanceAttributesRequest` 并填充 `DeleteProtectionEnable`，调用 `ModifyCkafkaInstanceAttributes`。
- 在 Read 流程的 `DescribeInstanceAttributes` 响应处理中，当 `attr.DeleteProtectionEnable` 不为 nil 时回写 `d.Set("delete_protection_enable", ...)`。
- 补充 `resource_tc_ckafka_instance_test.go` 单元测试，使用 gomonkey mock 云 API，覆盖 Create/Update/Read 三个分支中 `delete_protection_enable` 的处理逻辑。
- 更新 `resource_tc_ckafka_instance.md` 文档示例，展示 `delete_protection_enable` 字段用法。

## Capabilities

### New Capabilities
- `ckafka-instance-resource`: `tencentcloud_ckafka_instance` 通用资源（RESOURCE_KIND_GENERAL）的能力定义，覆盖实例的创建、读取、更新、删除全生命周期，本次新增 `delete_protection_enable` 删除保护开关参数。

### Modified Capabilities
<!-- 无已有 spec，本次为纯新增 -->

## Impact

- 代码：
  - `tencentcloud/services/ckafka/resource_tc_ckafka_instance.go`（schema 新增字段、Create/Update/Read 逻辑扩展）
  - `tencentcloud/services/ckafka/resource_tc_ckafka_instance_test.go`（新增 `delete_protection_enable` 相关单元测试用例）
  - `tencentcloud/services/ckafka/resource_tc_ckafka_instance.md`（示例补充 `delete_protection_enable`）
- 依赖：使用已 vendored 的 `tencentcloud-sdk-go` 中 `ckafka/v20190819.ModifyInstanceAttributesRequest.DeleteProtectionEnable`（`*int64`）及 `InstanceAttributesResponse.DeleteProtectionEnable`（`*int64`），无需变更 vendor。
- 向后兼容：新增 Optional 字段，未配置时不影响已有 state 与 TF 配置，无 state 迁移需求。
- 文档：需同步更新 website docs（由 `make doc` 自动生成流程读取 `.md` 文件）。
