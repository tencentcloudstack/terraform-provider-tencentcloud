## Why

当腾讯云 AntiDDoS（DDoS 防护）资源因遭受攻击被封堵后，用户需要调用解封接口手动解除封堵以恢复业务。目前该操作只能通过控制台或云 API 手动执行，缺少 Terraform 原生方式。新增一个 framework plugin 的 action 资源 `tencentcloud_antiddos_unblock_resources`，可以让用户在 Terraform 体系中以 IaC 方式触发解封操作。

## What Changes

- 新增一个 **action 资源**（terraform-plugin-framework action 类型）`tencentcloud_antiddos_unblock_resources`，通过 `UnblockResources` API 申请解封被封堵的资源（公网 IP 列表）。这是一次性操作，操作完成后不记录任何云侧状态。
- 该资源严格参考 `tencentcloud/framework/registry.go` 中已注册的 `NewTeoConfirmOriginAclUpdate`、`NewBdrcRunCopyPairTasks` 两个 action 的代码风格开发，文件落地在 `tencentcloud/framework/antiddos/action_tc_antiddos_unblock_resources.go`。
- 需要在 `tencentcloud/connectivity/client.go` 中新增 `UseAntiddosV20250903Client()` 方法，引入 `antiddos/v20250903` 版本的 SDK 客户端（现有 `UseAntiddosClient()` 仅支持 `v20200309`）。
- 在 `tencentcloud/framework/registry.go` 的 `actionFactories` 中注册该 action 工厂。
- 新增 action 示例文档 `tencentcloud/framework/antiddos/action_tc_antiddos_unblock_resources.md`。
- 新增单元测试 `tencentcloud/framework/antiddos/action_tc_antiddos_unblock_resources_test.go`，使用 gomonkey mock 云 API，不依赖 terraform 测试套件。

## Capabilities

### New Capabilities
- `antiddos-unblock-resources`: 一个 framework action 资源，用于通过 `UnblockResources` API 申请解封被封堵的 AntiDDoS 资源（公网 IP 列表），属于一次性操作，不持久化云侧状态。

### Modified Capabilities
<!-- 无：全新资源，不修改任何现有 spec。 -->

## Impact

- 新增文件:
  - `tencentcloud/framework/antiddos/action_tc_antiddos_unblock_resources.go`（action 实现）
  - `tencentcloud/framework/antiddos/action_tc_antiddos_unblock_resources.md`（示例文档）
  - `tencentcloud/framework/antiddos/action_tc_antiddos_unblock_resources_test.go`（单元测试）
- 修改文件:
  - `tencentcloud/framework/registry.go`（在 `actionFactories` 中注册 `antiddos.NewAntiddosUnblockResources`，并新增 antiddos import）
  - `tencentcloud/connectivity/client.go`（新增 `UseAntiddosV20250903Client()` 方法 + 对应 import + struct 字段）
  - `tencentcloud/services/antiddos/service_tencentcloud_antiddos.go`（新增 `UnblockResources` 服务方法）
- SDK: 使用 vendor 中已存在的 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20250903` 包（已在 vendor/modules.txt 中登记为 untracked，无需额外 SDK 变更）。
- 无破坏性变更，纯增量。
