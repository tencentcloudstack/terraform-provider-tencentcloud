## Context

TencentCloud TEO (EdgeOne) 产品在回源 IP 网段发生变更时会向用户推送通知。用户需要确认已将最新回源 IP 网段更新至源站防火墙，以停止推送变更通知。云 API 提供 `ConfirmOriginACLUpdate` 接口用于此确认操作。

当前 Terraform Provider 中缺少该操作对应的资源，用户无法通过 Terraform 自动化完成回源 IP 网段更新确认。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_confirm_origin_acl_update_operation` 一次性操作资源
- Create 阶段调用 `ConfirmOriginACLUpdate` API，传入 `ZoneId` 参数
- Read / Delete 为 no-op，无 Update 操作
- 在 provider.go 中注册资源
- 提供单元测试（使用 gomonkey mock）
- 生成 .md 文档

**Non-Goals:**
- 不实现 Read / Update / Delete 的实际逻辑（一次性操作无持久状态）
- 不支持 Import（一次性操作无需导入）
- 不处理异步轮询（API 为同步接口）

## Decisions

1. **资源类型选择 RESOURCE_KIND_OPERATION**：该接口仅为确认操作，不创建持久化资源，符合一次性操作资源的定义。Create 调用 API 后设置 ID 为 `helper.BuildToken()`。

2. **Schema 设计**：仅暴露 `zone_id` 作为 Required + ForceNew 参数，与 API 入参 `ZoneId` 一一对应。无需 Computed 字段，因为 API 响应仅返回 `RequestId`。

3. **重试策略**：Create 阶段使用 `tccommon.WriteRetryTimeout` 和 `resource.Retry` 包装 API 调用，失败时使用 `tccommon.RetryError()` 包装错误。

4. **代码风格参考**：参考同产品下的 `resource_tc_teo_identify_zone_operation.go` 和 config 产品的 `resource_tc_config_start_config_rule_evaluation_operation.go`。

## Risks / Trade-offs

- **[API 调用幂等性]** → `ConfirmOriginACLUpdate` 多次调用对同一 ZoneId 是安全的（仅确认操作），不存在副作用风险
- **[无状态追踪]** → 由于是操作型资源，Terraform state 中仅存储 token ID，无法通过 Read 刷新状态，这是 OPERATION 类型的预期行为
