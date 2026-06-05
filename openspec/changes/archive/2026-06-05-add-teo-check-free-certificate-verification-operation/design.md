## Context

TencentCloud EdgeOne (TEO) 提供免费证书申请流程，其中 `CheckFreeCertificateVerification` 接口用于验证免费证书并获取申请结果。当前 Terraform Provider 已支持 TEO 的多种资源，但缺少该验证操作的支持。

该资源为 RESOURCE_KIND_OPERATION 类型，属于一次性操作资源，执行后不需要维护状态。仅需实现 Create 方法，Read/Update/Delete 为空实现。

云 API 已在 vendor 中可用：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_teo_check_free_certificate_verification` 操作资源
- 支持传入 `zone_id` 和 `domain` 参数调用 CheckFreeCertificateVerification 接口
- 将接口返回的 `common_name`、`signature_algorithm`、`expire_time` 作为 Computed 属性暴露
- 提供完整的单元测试（使用 gomonkey mock 云 API）
- 在 provider.go 和 provider.md 中注册资源

**Non-Goals:**
- 不实现 Read/Update/Delete 方法（OPERATION 类型资源无需维护状态）
- 不处理异步轮询（该接口为同步接口）
- 不支持 Import 功能（一次性操作资源无需 import）

## Decisions

### 1. 资源类型选择：RESOURCE_KIND_OPERATION

**决策**: 使用一次性操作资源模式，仅实现 Create，RUD 为空。

**理由**: `CheckFreeCertificateVerification` 是一个验证操作，不创建或管理任何持久化资源。每次执行都是独立的验证动作，无需跟踪状态。

**替代方案**: 使用 data source —— 但 data source 不适合有副作用的操作（验证动作可能触发证书签发）。

### 2. ID 生成策略

**决策**: 使用 `zone_id#domain` 作为复合 ID（使用 tccommon.FILED_SP 分隔符）。

**理由**: 操作资源虽然不需要 Read，但 Terraform 要求每个资源有唯一 ID。使用入参组合作为 ID 可以标识每次操作的目标。

### 3. 错误处理与重试

**决策**: 使用 `resource.Retry` + `tccommon.ReadRetryTimeout` 包装 API 调用，失败时使用 `tccommon.RetryError()` 包装错误。

**理由**: 遵循项目统一的错误处理模式，确保网络抖动等临时错误可以自动重试。

### 4. 测试策略

**决策**: 使用 gomonkey mock 云 API 进行单元测试，不使用 Terraform 验收测试套件。

**理由**: 新增 RESOURCE_KIND_OPERATION 资源按要求使用 mock 方式进行单元测试，避免依赖真实云环境。

## Risks / Trade-offs

- [风险] 操作资源每次 `terraform apply` 都会重新执行 → 用户需理解 OPERATION 类型资源的行为特点，文档中需说明。
- [风险] 证书验证可能因域名配置未完成而失败 → API 错误会通过 Terraform 错误信息透传给用户，无需额外处理。
