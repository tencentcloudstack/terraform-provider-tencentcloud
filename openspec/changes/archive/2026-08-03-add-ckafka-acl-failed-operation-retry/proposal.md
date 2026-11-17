## Why

`tencentcloud_ckafka_acl` 资源在创建 ACL（调用云 API `CreateAcl`）时，可能遇到 `Code=FailedOperation` 的瞬时错误（例如实例正在进行变更、集群内部任务进行中等场景）。当前 `CkafkaService.CreateAcl` 通过 `tccommon.RetryError(err)` 处理错误，而 `FailedOperation` 不在 `retryableErrorCode` 列表中，会被直接判定为**不可重试**错误并立即返回，导致 `terraform apply` 创建 ACL 失败，需要用户手动重试。

## What Changes

- 修改 `tencentcloud/services/ckafka/service_tencentcloud_ckafka.go` 中 `CkafkaService.CreateAcl` 方法：
  - 当调用 `CreateAcl` 云 API 返回 `Code == "FailedOperation"` 时，进行**最多 3 次重试**，每次**间隔 5 秒**。
  - 其余错误保持原有 `resource.Retry` + `tccommon.RetryError` 处理逻辑不变（网络错误、限流等仍按原策略处理）。
- 在 `tencentcloud/services/ckafka/extension_ckafka.go` 中新增重试次数与间隔常量（`3` 次、`5 * time.Second`），便于维护与单测引用。
- 在 `resource_tc_ckafka_acl_test.go` 中补充单元测试，使用 gomonkey mock `CreateAcl` 云 API，验证 FailedOperation 重试行为（重试次数、间隔、最终成功/失败路径）。
- 更新 `resource_tc_ckafka_acl.md` 说明创建时的重试行为（保持向后兼容，不改动 schema 与 ID 格式）。

## Capabilities

### New Capabilities
- `ckafka-acl-failed-operation-retry`: 定义 `tencentcloud_ckafka_acl` 资源创建时（`CreateAcl` 云 API）对 `FailedOperation` 错误的重试行为：最多重试 3 次、每次间隔 5 秒。

### Modified Capabilities
（无。`openspec/specs/` 下不存在 ckafka acl 相关的既有 spec，属于新增能力。）

## Impact

- **代码文件**：`tencentcloud/services/ckafka/service_tencentcloud_ckafka.go`（修改 `CreateAcl` 方法）
- **常量文件**：`tencentcloud/services/ckafka/extension_ckafka.go`（新增重试常量）
- **测试文件**：`tencentcloud/services/ckafka/resource_tc_ckafka_acl_test.go`（补充单元测试）
- **文档文件**：`tencentcloud/services/ckafka/resource_tc_ckafka_acl.md`（补充重试行为说明）
- **云 API**：`CreateAcl`（ckafka/v20190819）。已确认该接口可能返回的错误码包含 `FAILEDOPERATION = "FailedOperation"`（见 vendor 中 `client.go` 的 `CreateAclWithContext` 注释），需求可行。
- **依赖**：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors`（已引入，用于 `TencentCloudSDKError` 类型断言）、标准库 `time`（需新增 import）。
- **兼容性**：不改变资源 Schema、资源 ID 格式、`CreateAcl` 方法签名；成功路径行为完全不变，仅增强失败场景的容错，向后兼容。
