## Context

TencentCloud EdgeOne (TEO) 提供了 ApplyFreeCertificate 接口，允许用户为 CNAME 接入模式的域名申请免费证书。该接口是同步接口，调用后立即返回验证信息（DNS 验证或文件验证），用户需要根据返回的验证信息完成后续配置。

当前 Terraform Provider 中已有多个 TEO OPERATION 类型资源（如 `teo_identify_zone_operation`、`teo_prefetch_task_operation` 等），本资源遵循相同的模式。

## Goals / Non-Goals

**Goals:**
- 提供 `tencentcloud_teo_apply_free_certificate` 操作资源，允许用户通过 Terraform 发起免费证书申请
- 返回 DNS 验证信息（subdomain、record_type、record_value）或文件验证信息（path、content）供用户后续配置
- 遵循现有 TEO OPERATION 资源的代码风格和模式

**Non-Goals:**
- 不实现证书验证结果检查（CheckFreeCertificateVerification 是独立操作）
- 不实现证书部署（ModifyHostsCertificate 是独立操作）
- 不实现 Read/Update/Delete 逻辑（OPERATION 类型资源无需持久化状态）

## Decisions

1. **资源 ID 策略**: 使用 `zone_id + tccommon.FILED_SP + domain` 作为复合 ID。虽然 OPERATION 资源不需要 Read，但复合 ID 可以标识操作的目标，便于日志追踪。
   - 替代方案: 使用 `helper.BuildToken()` 生成随机 ID。不采用是因为复合 ID 更具可读性。

2. **代码组织**: 直接在资源文件中调用 SDK，不通过 service 层封装。
   - 理由: OPERATION 资源逻辑简单，只有一次 API 调用，无需 service 层抽象。参考 `teo_identify_zone_operation` 的模式。

3. **重试策略**: 使用 `resource.Retry` + `tccommon.ReadRetryTimeout` 包装 API 调用。
   - 理由: 遵循项目统一的重试模式，处理网络瞬时故障。

4. **验证信息存储**: 将 DnsVerification 和 FileVerification 作为 Computed 字段存入 state。
   - 理由: 用户需要获取验证信息以完成后续配置步骤。

## Risks / Trade-offs

- [风险] OPERATION 资源每次 apply 都会重新调用接口 → 这是 OPERATION 类型的预期行为，用户应了解每次 apply 都会触发新的证书申请
- [风险] 验证信息可能过期（2天有效期）→ 用户需要在有效期内完成验证，这是 API 层面的限制，不在 Terraform 层面处理
